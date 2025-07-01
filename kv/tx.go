package kv

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"time"

	"github.com/xiaoma03xf/sharddoc/kv/pb"
)

// 用于记录事务读取的键范围（start 到 stop
// 后续用来做读写冲突检测（Conflict Detection）
type KeyRange struct {
	start []byte
	stop  []byte
}

// 事务存在的问题:
// 1. 脏读：读取其他事务未提交的数据
//    解决方案：读取优先从本事务未提交的修改：pending BTree 中读取数据，如果没有就会在 snapshot BTree
//	  开始事务时的快照中读取，不会读取其他事务的pending数据

// 2. 不可重复读: 一个事务在读取一条数据时，由于另一个事务修改了这条数据并且提交事务，再次读取时导致数据不一致
// 问题描述：事务 A 两次读取 key = k1, 事务 B 在 A 第一次读之后修改并提交了 k1,  A 第二次读取时拿到不同结果
// 解决方案: 事故整个事务周期，同一个事务的读操作来自于自身的 pending BTree 或者 snapshot BTree 其他事务隔离开的，不会受到影响

// 3. 幻读：一个事务读取了某个范围的数据，同时另一个事务新增了这个范围的数据，再次读取发现俩次得到的结果不一致
// 事务 A 查询了一个范围 [k1, k5]， 事务 B 插入了一个新 key k3 并提交， A 再次读取该范围时，发现多了个 key（k3）→ 这就是 幻影
// KV transaction
type KVTX struct {
	snapshot BTree
	version  uint64 // 当前事务的版本号
	// captured KV updates:
	// values are prefixed by a 1-byte flag to indicate deleted keys.
	// 本地未提交的写操作（用 B+树保存）, 不是对snapshot的拷贝, 而是对写操作的增量
	pending BTree
	// a list of involved intervals of keys for detecting conflicts
	// 本事务读取过的 key 范围
	reads []KeyRange
	// should check for conflicts even if an update changes nothing

	// 是否尝试过写入操作（即使没有改变）
	updateAttempted bool

	// check misuses
	// 标记事务是否正常结束
	errors []error
	done   bool
}

// a < b
func versionBefore(a, b uint64) bool {
	return a-b > 1<<63 // this works even after wraparounds
}

// a prefix for values in KVTX.pending
const (
	FLAG_DELETED = byte(1)
	FLAG_UPDATED = byte(2)
)

// begin a transaction
// 初始化一个新的事务 tx，让它拥有当前数据库的快照版本，并准备好一个空的`待提交修改表`
func (kv *KV) Begin(tx *KVTX) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	// read-only snapshot, just the tree root and the page read callback
	tx.snapshot.root = kv.tree.root
	chunks := kv.mmap.chunks // copied to avoid updates from writers
	tx.snapshot.get = func(ptr uint64) []byte { return mmapRead(ptr, chunks) }
	tx.version = kv.version

	// in-memory tree to capture updates
	// 	设置一个 B+ 树结构，事务对数据的写入都会暂存在这个 pending 树中，实际写入磁盘前不会影响快照
	// 所有的变更（无论是写还是删）都会先写入这个临时树中，并用 1 字节的 flag 表明是删除还是更新
	pages := [][]byte(nil)
	tx.pending.get = func(ptr uint64) []byte { return pages[ptr-1] }
	tx.pending.new = func(node []byte) uint64 {
		pages = append(pages, node)
		return uint64(len(pages))
	}
	tx.pending.del = func(uint64) {}

	// keep track of concurrent TXs
	// 把该事务的版本号加入一个`正在进行的事务`列表，方便以后判断是否有并发冲突
	kv.ongoing = append(kv.ongoing, tx.version)
	tx.errors = make([]error, 0)
	// XXX: sanity check: unreliable
	// 给这个事务设置一个析构钩子，如果事务对象被 GC 了，但 tx.done = false，就会触发断言失败
	// 这是一个安全检查机制，防止忘记 Commit 或 Abort
	runtime.SetFinalizer(tx, func(tx *KVTX) { assert(tx.done) })
}

// 判断两个有序区间是否相交
func sortedRangesOverlap(s1, s2 []KeyRange) bool {
	for len(s1) > 0 && len(s2) > 0 {
		if bytes.Compare(s1[0].stop, s2[0].start) < 0 {
			s1 = s1[1:]
		} else if bytes.Compare(s2[0].stop, s1[0].start) < 0 {
			s2 = s2[1:]
		} else {
			return true
		}
	}
	return false
}

// 检测当前事务 tx 的读取范围是否与任何 更新过的历史事务 的写入范围冲突，以确保事务的隔离性
// （事务）在开始工作后，不能读到后来被别人改过的数据，否则可能产生脏读/幻读
// 解决 写-后-读 冲突（Write-After-Read）
func detectConflicts(kv *KV, tx *KVTX) bool {
	// sort the dependency ranges for easier overlap detection
	// 先把当前事务读取区间排序，方便后续快速判断重叠
	slices.SortFunc(tx.reads, func(r1, r2 KeyRange) int {
		return bytes.Compare(r1.start, r2.start)
	})
	// do they overlap with newer versions?
	// 从最新的历史事务往回遍历
	for i := len(kv.history) - 1; i >= 0; i-- {
		// 如果历史事务版本比当前事务版本新，则继续比较，否则跳出循环（因为history是按版本排序的）
		// 因为当前事务是写操作, 当且仅当历史事务在前面读, 我在后面写他要的数据才报错
		// tx.version < kv.history[i].version? 表示从事务末尾到tx.version了, 再往前不用看了
		if !versionBefore(tx.version, kv.history[i].version) {
			break // sorted
		}
		if sortedRangesOverlap(tx.reads, kv.history[i].writes) {
			return true
		}
	}
	return false
}

var ErrorConflict = errors.New("cannot commit due to conflict")

func (kv *KV) Commit(tx *KVTX) error {
	//保证只提交一次，进入`关键区`，并确保无论成功与否，最后会执行 txFinalize 清理工作
	assert(!tx.done)
	tx.done = true
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	defer txFinalize(kv, tx)

	// check conflicts, 当前写事务后面有事务读(写后读冲突/脏读)
	if tx.updateAttempted && detectConflicts(kv, tx) {
		return ErrorConflict
	}
	// 检查事务是否已失效
	if len(tx.errors) > 0 {
		return fmt.Errorf("cannot commit: transaction has %d errors, first: %v", len(tx.errors), tx.errors[0])
	}

	// save the meta page, 保存元数据并记录当前B+树根节点
	meta, root := saveMeta(kv), kv.tree.root
	// transfer updates to the current tree, 更新空闲列表版本号
	kv.free.curVer = kv.version + 1 // version in the free list
	writes := []KeyRange(nil)       // collect updated ranges in this TX

	// 保存pending对key, val的修改
	incrals := &pb.SnapshotBatch{}
	// 从头遍历pending增量树, iter迭代器遍历此树, 写入增量树的数据头部会带一个1byte的标识 FLAG_DELETED|FLAG_UPDATED
	for iter := tx.pending.Seek(nil, CMP_GT); iter.Valid(); iter.Next() {
		modified := false
		key, val := iter.Deref()
		oldVal, isOld := tx.snapshot.Get(key) // 从本树中取旧值

		switch val[0] {
		case FLAG_DELETED:
			// 操作为删除, 标记modified为isOld
			modified = isOld
			deleted, err := kv.tree.Delete(&DeleteReq{Key: key})
			assert(err == nil)          // can only fail by length limit
			assert(deleted == modified) // assured by conflict detection

			// 增量删除, 数据记录
			wal := &pb.IncrementalSnapshot{Key: key, Value: val,
				Timestamp: time.Now().UnixNano(), Operation: int32(OpDelete)}
			incrals.Snapshots = append(incrals.Snapshots, wal)

		case FLAG_UPDATED:
			// 不是旧值, 或者说新修改的值和旧值不相等的时候, 标记modified
			modified = !isOld || !bytes.Equal(oldVal, val[1:])
			updated, err := kv.tree.Update(&UpdateReq{Key: key, Val: val[1:]})
			assert(err == nil)          // can only fail by length limit
			assert(updated == modified) // assured by conflict detection

			// 增量更新, 数据记录
			wal := &pb.IncrementalSnapshot{Key: key, Value: val,
				Timestamp: time.Now().UnixNano(), Operation: int32(OpUpdate)}
			incrals.Snapshots = append(incrals.Snapshots, wal)
		default:
			panic("unreachable")
		}
		// 收集此次修改的键范围, 便于事务的冲突检测
		if modified && len(kv.ongoing) > 1 {
			writes = append(writes, KeyRange{key, key})
		}

	}

	// commit the update
	if root != kv.tree.root {
		kv.version++
		// 在写时复制 + MVCC 的事务模型中，Abort 不需要回滚，只需要`什么都别提交`，状态就自动回滚了
		// 如果出现写入失败, B+树自己就会回滚了
		if err := updateOrRevert(kv, meta); err != nil {
			return err
		}
	}

	//将这次事务写过的 key 区间 writes 存入 history
	// 用于将来检测后续事务是否和当前这次提交冲突
	// keep a history of updated ranges grouped by each TX
	if len(writes) > 0 {
		// sort the ranges for faster overlap detection
		slices.SortFunc(writes, func(r1, r2 KeyRange) int {
			return bytes.Compare(r1.start, r2.start)
		})
		kv.history = append(kv.history, CommittedTX{kv.version, writes})
	}

	// 导出增量快照, 只导出增删改, 查不需要导出增量快照
	if len(incrals.Snapshots) > 0 {
		if err := kv.asyncSnapshot(tx, incrals); err != nil {
			panic(err)
		}
	}
	return nil
}

// common routines when exiting a transaction
// 事务生命周期结束（无论是 Commit 还是 Abort）时的清理逻辑
func txFinalize(kv *KV, tx *KVTX) {
	// remove myself from `kv.ongoing`
	idx := slices.Index(kv.ongoing, tx.version)
	last := len(kv.ongoing) - 1
	kv.ongoing[idx], kv.ongoing = kv.ongoing[last], kv.ongoing[:last]
	// find the oldest in-use version
	minVersion := kv.version
	for _, other := range kv.ongoing {
		if versionBefore(other, minVersion) {
			minVersion = other
		}
	}
	// release the free list
	kv.free.SetMaxVer(minVersion)
	// trim `kv.history` if `minVersion` has been increased
	for idx = 0; idx < len(kv.history); idx++ {
		if versionBefore(minVersion, kv.history[idx].version) {
			break // sorted
		}
	}
	kv.history = kv.history[idx:] // sorted
}

// end a transaction: rollback
func (kv *KV) Abort(tx *KVTX) {
	assert(!tx.done)
	tx.done = true
	// maintain `kv.ongoing` and `kv.history`
	kv.mutex.Lock()
	txFinalize(kv, tx)
	kv.mutex.Unlock()
}

type KVIter interface {
	Deref() (key []byte, val []byte)
	Valid() bool
	Next()
	// TODO: Prev()
}

// an iterator that combines pending updates and the snapshot
type CombinedIter struct {
	top *BIter // KVTX.pending
	bot *BIter // KVTX.snapshot
	dir int    // +1 for greater or greater-than, -1 for less or less-than
	// the end of the range
	cmp int
	end []byte
}

func (iter *CombinedIter) Deref() ([]byte, []byte) {
	var k1, v1, k2, v2 []byte
	top, bot := iter.top.Valid(), iter.bot.Valid()
	assert(top || bot) // 至少有一个迭代器是有效的（还有数据可读）
	if top {
		k1, v1 = iter.top.Deref()
	}
	if bot {
		k2, v2 = iter.bot.Deref()
	}
	// use the min/max key of the two
	// 如果两个迭代器都有效, 且他们指向的键都相等
	if top && bot && bytes.Compare(k1, k2) == +iter.dir {
		return k2, v2
	}
	if top {
		// v[1:] 表示去除标志
		return k1, v1[1:]
	} else {
		return k2, v2
	}
}

func (iter *CombinedIter) Valid() bool {
	if iter.top.Valid() || iter.bot.Valid() {
		key, _ := iter.Deref()
		return cmpOK(key, iter.cmp, iter.end)
	}
	return false
}

func (iter *CombinedIter) Next() {
	// which B+tree iterator to move?
	top, bot := iter.top.Valid(), iter.bot.Valid()
	if top && bot {
		k1, _ := iter.top.Deref()
		k2, _ := iter.bot.Deref()
		switch bytes.Compare(k1, k2) {
		case -iter.dir:
			top, bot = true, false
		case +iter.dir:
			top, bot = false, true
		case 0: // equal; move both
		}
	}
	assert(top || bot)
	// move B+tree iterators wrt the direction
	if top {
		if iter.dir > 0 {
			iter.top.Next()
		} else {
			iter.top.Prev()
		}
	}
	if bot {
		if iter.dir > 0 {
			iter.bot.Next()
		} else {
			iter.bot.Prev()
		}
	}
}

// point query. combines captured updates with the snapshot
// 事务读取优先从 pending 查，找不到再从 snapshot 查
func (tx *KVTX) Get(key []byte) ([]byte, bool) {
	tx.reads = append(tx.reads, KeyRange{key, key}) // dependency, 记录依赖
	val, ok := tx.pending.Get(key)                  // 优先从 pending 查
	switch {
	case ok && val[0] == FLAG_UPDATED: // updated in this TX
		return val[1:], true
	case ok && val[0] == FLAG_DELETED: // deleted in this TX
		return nil, false
	case !ok: // read from the snapshot
		return tx.snapshot.Get(key)
	default:
		panic("unreachable")
	}
}

func cmp2dir(cmp int) int {
	if cmp > 0 {
		return +1
	} else {
		return -1
	}
}

// range query. combines captured updates with the snapshot
func (tx *KVTX) Seek(key1 []byte, cmp1 int, key2 []byte, cmp2 int) KVIter {
	assert(cmp2dir(cmp1) != cmp2dir(cmp2))
	lo, hi := key1, key2
	if cmp2dir(cmp1) < 0 {
		lo, hi = hi, lo
	}
	tx.reads = append(tx.reads, KeyRange{lo, hi}) // FIXME: slightly larger
	return &CombinedIter{
		top: tx.pending.Seek(key1, cmp1),
		bot: tx.snapshot.Seek(key1, cmp1),
		dir: cmp2dir(cmp1),
		cmp: cmp2,
		end: key2,
	}
}

// capture updates
func (tx *KVTX) Update(req *UpdateReq) (bool, error) {
	tx.updateAttempted = true
	// check the existing key against the update mode
	old, exists := tx.Get(req.Key) // also add a dependency
	if req.Mode == MODE_UPDATE_ONLY && !exists {
		return false, nil
	}
	if req.Mode == MODE_INSERT_ONLY && exists {
		return false, nil
	}
	if exists && bytes.Equal(old, req.Val) {
		return false, nil
	}
	// insert the flagged KV (reduces the size limit by 1 byte)
	flaggedVal := append([]byte{FLAG_UPDATED}, req.Val...)
	_, err := tx.pending.Update(&UpdateReq{Key: req.Key, Val: flaggedVal})
	if err != nil {
		tx.errors = append(tx.errors, err)
		return false, err // length limit
	}
	req.Added = !exists
	req.Updated = true
	req.Old = old
	return true, nil
}

func (tx *KVTX) Set(key []byte, val []byte) (bool, error) {
	return tx.Update(&UpdateReq{Key: key, Val: val})
}

func (tx *KVTX) Del(req *DeleteReq) (bool, error) {
	tx.updateAttempted = true
	exists := false
	if req.Old, exists = tx.Get(req.Key); !exists {
		return false, nil
	}
	return tx.pending.Update(&UpdateReq{Key: req.Key, Val: []byte{FLAG_DELETED}})
}
