package byodb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"golang.org/x/sys/unix"
)

type KV struct {
	Path  string
	Fsync func(int) error // overridable; for testing
	// internals
	fd   int
	tree BTree
	free FreeList
	mmap struct {
		total  int      // mmap size, can be larger than the file size
		chunks [][]byte // multiple mmaps, can be non-continuous
	}
	page struct {
		flushed uint64            // database size in number of pages
		nappend uint64            // number of pages to be appended
		updates map[uint64][]byte // pending updates, including appended pages
	}
	failed bool // Did the last update fail?
}

// `BTree.get`, read a page.
func (db *KV) pageRead(ptr uint64) []byte {
	assert(ptr < db.page.flushed+db.page.nappend)
	if node, ok := db.page.updates[ptr]; ok {
		return node // pending update
	}
	return db.pageReadFile(ptr)
}

func (db *KV) pageReadFile(ptr uint64) []byte {
	start := uint64(0)
	for _, chunk := range db.mmap.chunks {
		end := start + uint64(len(chunk))/BTREE_PAGE_SIZE
		if ptr < end {
			offset := BTREE_PAGE_SIZE * (ptr - start)
			return chunk[offset : offset+BTREE_PAGE_SIZE]
		}
		start = end
	}
	panic("bad ptr")
}

// `BTree.new`, allocate a new page.
func (db *KV) pageAlloc(node []byte) uint64 {
	assert(len(node) == BTREE_PAGE_SIZE)
	if ptr := db.free.PopHead(); ptr != 0 { // try the free list
		assert(db.page.updates[ptr] == nil)
		db.page.updates[ptr] = node
		return ptr
	}
	return db.pageAppend(node) // append
}

// `FreeList.new`, append a new page.
func (db *KV) pageAppend(node []byte) uint64 {
	assert(len(node) == BTREE_PAGE_SIZE)
	ptr := db.page.flushed + db.page.nappend
	db.page.nappend++
	assert(db.page.updates[ptr] == nil)
	db.page.updates[ptr] = node
	return ptr
}

// `FreeList.set`, update an existing page.
func (db *KV) pageWrite(ptr uint64) []byte {
	assert(ptr < db.page.flushed+db.page.nappend)
	if node, ok := db.page.updates[ptr]; ok {
		return node // pending update
	}
	node := make([]byte, BTREE_PAGE_SIZE)
	copy(node, db.pageReadFile(ptr)) // initialized from the file
	db.page.updates[ptr] = node
	return node
}

// open or create a file and fsync the directory
// 回归到写文件问题, 写文件流程。对应着B+树的持久化步骤
// 1.创建新文件(副本)并写入数据, 2.新文件写盘,确保真的写入硬盘 fsync
// 3. 原子操作(重命名)替换掉旧文件 4.对目录做一次 fsync, 文件名是否能被找到，要靠目录也 fsync()

// 为了防止断电造成的文件状态不一致，我们不能依赖文件大小来判断页面数量，
// 而是需要在第 0 页的元页面中明确记录 使用了多少页，并控制写入顺序和落盘顺序，确保数据一致性。
// 安全地创建一个可读写的文件，并确保该文件及其目录项在系统崩溃或断电后都能被正确恢复和识别
func createFileSync(file string) (int, error) {
	// obtain the directory fd
	// 1. 打开文件所在目录（读模式 + 指定为目录）
	flags := os.O_RDONLY | syscall.O_DIRECTORY
	dirfd, err := syscall.Open(path.Dir(file), flags, 0o644)
	if err != nil {
		return -1, fmt.Errorf("open directory: %w", err)
	}
	defer syscall.Close(dirfd)

	// open or create the file
	// 2. 使用 openat 打开或创建目标文件
	flags = os.O_RDWR | os.O_CREATE
	fd, err := syscall.Openat(dirfd, path.Base(file), flags, 0o644)
	if err != nil {
		return -1, fmt.Errorf("open file: %w", err)
	}

	// 3. 对目录进行 fsync，确保目录项写入磁盘
	// fsync the directory
	// 在创建新文件时，不仅要 fsync 文件内容，还要 fsync 目录，
	// 这样才能在系统崩溃或断电时，确保文件名不会丢失。同时通过 openat 保证目录路径不会被恶意修改。
	err = syscall.Fsync(dirfd)
	if err != nil { // may leave an empty file
		_ = syscall.Close(fd)
		return -1, fmt.Errorf("fsync directory: %w", err)
	}
	// done
	return fd, nil
}

// open or create a DB file
func (db *KV) Open() error {
	if db.Fsync == nil {
		db.Fsync = syscall.Fsync
	}
	var err error
	db.page.updates = map[uint64][]byte{}
	// B+tree callbacks
	db.tree.get = db.pageRead
	db.tree.new = db.pageAlloc
	db.tree.del = db.free.PushTail
	// free list callbacks
	db.free.get = db.pageRead
	db.free.new = db.pageAppend
	db.free.set = db.pageWrite
	// open or create the DB file
	if db.fd, err = createFileSync(db.Path); err != nil {
		return err
	}
	// get the file size
	finfo := syscall.Stat_t{}
	if err = syscall.Fstat(db.fd, &finfo); err != nil {
		goto fail
	}
	// create the initial mmap
	if err = extendMmap(db, int(finfo.Size)); err != nil {
		goto fail
	}
	// read the meta page
	if err = readRoot(db, finfo.Size); err != nil {
		goto fail
	}
	return nil
	// error
fail:
	db.Close()
	return fmt.Errorf("KV.Open: %w", err)
}

const DB_SIG = "BuildYourOwnDB07"

/*
the 1st page stores the root pointer and other auxiliary data.
| sig | root_ptr | page_used | head_page | head_seq | tail_page | tail_seq |
| 16B |    8B    |     8B    |     8B    |    8B    |     8B    |    8B    |
*/
// 把根节点页号和 已提交的最大页编号 恢复成写入前的值，让数据库逻辑回到 什么都没发生 之前的状态。
func loadMeta(db *KV, data []byte) {
	db.tree.root = binary.LittleEndian.Uint64(data[16:24])
	db.page.flushed = binary.LittleEndian.Uint64(data[24:32])
	db.free.headPage = binary.LittleEndian.Uint64(data[32:40])
	db.free.headSeq = binary.LittleEndian.Uint64(data[40:48])
	db.free.tailPage = binary.LittleEndian.Uint64(data[48:56])
	db.free.tailSeq = binary.LittleEndian.Uint64(data[56:64])
}

func saveMeta(db *KV) []byte {
	var data [64]byte
	copy(data[:16], []byte(DB_SIG))
	binary.LittleEndian.PutUint64(data[16:24], db.tree.root)
	binary.LittleEndian.PutUint64(data[24:32], db.page.flushed)
	binary.LittleEndian.PutUint64(data[32:40], db.free.headPage)
	binary.LittleEndian.PutUint64(data[40:48], db.free.headSeq)
	binary.LittleEndian.PutUint64(data[48:56], db.free.tailPage)
	binary.LittleEndian.PutUint64(data[56:64], db.free.tailSeq)
	return data[:]
}

// 读取数据库文件的元页面（meta page）信息，并做一些校验
func readRoot(db *KV, fileSize int64) error {
	if fileSize%BTREE_PAGE_SIZE != 0 {
		return errors.New("file is not a multiple of pages")
	}
	if fileSize == 0 { // empty file
		// reserve 2 pages: the meta page and a free list node
		db.page.flushed = 2
		// add an initial node to the free list so it's never empty
		db.free.headPage = 1 // the 2nd page
		db.free.tailPage = 1
		return nil // the meta page will be written in the 1st update
	}
	// read the page
	data := db.mmap.chunks[0]
	loadMeta(db, data)
	// initialize the free list
	db.free.SetMaxSeq()
	// verify the page, 校验元页面是否合法
	bad := !bytes.Equal([]byte(DB_SIG), data[:16])
	// pointers are within range?
	maxpages := uint64(fileSize / BTREE_PAGE_SIZE)
	bad = bad || !(0 < db.page.flushed && db.page.flushed <= maxpages)
	bad = bad || !(0 < db.tree.root && db.tree.root < db.page.flushed)
	bad = bad || !(0 < db.free.headPage && db.free.headPage < db.page.flushed)
	bad = bad || !(0 < db.free.tailPage && db.free.tailPage < db.page.flushed)
	if bad {
		return errors.New("bad meta page")
	}
	return nil
}

// update the meta page. it must be atomic.
func updateRoot(db *KV) error {
	// NOTE: atomic?
	if _, err := syscall.Pwrite(db.fd, saveMeta(db), 0); err != nil {
		return fmt.Errorf("write meta page: %w", err)
	}
	return nil
}

// extend the mmap by adding new mappings.
// 保证你的程序可以通过内存映射访问数据库文件中你需要的更多内容，并且映射的内存是动态增加的
// mmap 作磁盘映射时, 映射范围可以大于当前文件大小, 因为文件会增长。
// 1. 64 位系统可以寻址上百 TB 的虚拟内存, 而且映射也不会占用实际物理内存, 开销非常小
// 因此可以程序 预先申请大块连续的地址空间，方便以后文件增长时，能直接在已有映射里扩展读取或写入，而不用频繁创建新的映射

// 2. 追加新的映射块：每次文件增长时，新映射一块新的地址空间，拼接成多个块
// 每次文件扩展时添加新的映射会导致大量的映射，这会降低性能，因为操作系统必须跟踪它们。
// 指数增长可以避免这种情况，因为 mmap文件大小可能会超出指数增长的范围
// 指数增长, 每次扩展映射时，申请的大小是当前映射大小和 64MB 的最大值，再指数级翻倍，保证映射块数量不至于太多

func extendMmap(db *KV, size int) error {
	if size <= db.mmap.total {
		return nil // enough range
	}
	// 64MB是最小扩展单元, 或者当前总大小, 如果还是不够, 就按指数倍增长
	alloc := max(db.mmap.total, 64<<20) // double the current address space
	for db.mmap.total+alloc < size {
		alloc *= 2 // still not enough?
	}
	// 申请映射块
	chunk, err := syscall.Mmap(
		db.fd, int64(db.mmap.total), alloc,
		syscall.PROT_READ, syscall.MAP_SHARED, // read-only
	)
	if err != nil {
		return fmt.Errorf("mmap: %w", err)
	}
	// 	更新 db.mmap.total 和 db.mmap.chunks
	// 把新映射块追加到切片 chunks，并增加总映射大小
	db.mmap.total += alloc
	db.mmap.chunks = append(db.mmap.chunks, chunk)
	return nil
}

// 原子更新指针只能保证指针本身不出错（不会半更新），但不能保证指针指向的数据已经写好了
// 简单说就是修改指针指向不出错,但是不能确保所指的数据已经修改好了
// 如果在数据未处理结束就已经修改好指针, 指针可能指向一个空位置
// 此处采用写时复制

// 双写方案
// 1. 写到双写缓冲区（临时安全区）,  2. 复制到正式位置 3.出错了就从缓冲区恢复恢复

// 写时复制是 先改副本、再换指针
// 双写是 先写中间区，再更新原地
// 先写日志再写数据，通过日志重放或回滚数据
func updateFile(db *KV) error {
	// 1. Write new nodes.
	if err := writePages(db); err != nil {
		return err
	}
	// 2. `fsync` to enforce the order between 1 and 3.
	if err := db.Fsync(db.fd); err != nil {
		return err
	}
	// 3. Update the root pointer atomically.
	if err := updateRoot(db); err != nil {
		return err
	}
	// 4. `fsync` to make everything persistent.
	if err := db.Fsync(db.fd); err != nil {
		return err
	}
	// prepare the free list for the next update
	db.free.SetMaxSeq()
	return nil
}

// 利用写时复制 + 元页面只读策略，让失败的写入不会破坏旧状态，从而可以在失败后恢复读取甚至恢复写入
// meta 是写入前保存下来的 旧元页面 的内容，也就是写失败时要恢复的内容
func updateOrRevert(db *KV, meta []byte) error {
	// ensure the on-disk meta page matches the in-memory one after an error
	// 如果上一次写失败了，先强制恢复旧的元页面到磁盘
	// 使用 syscall.Pwrite() 把之前保存的 meta 重新写入文件开头（偏移为 0，通常是元页面的位置）
	if db.failed {
		if _, err := syscall.Pwrite(db.fd, meta, 0); err != nil {
			return fmt.Errorf("rewrite meta page: %w", err)
		}
		if err := db.Fsync(db.fd); err != nil {
			return err
		}
		db.failed = false
	}
	// 2-phase update
	err := updateFile(db)
	// revert on error
	if err != nil {
		// the on-disk meta page is in an unknown state.
		// mark it to be rewritten on later recovery.
		db.failed = true
		// in-memory states are reverted immediately to allow reads
		// 撤销当前修改，恢复数据库内存状态到写入前的安全状态
		loadMeta(db, meta)
		// discard temporaries
		db.page.nappend = 0
		db.page.updates = map[uint64][]byte{}
	}
	return err
}

func writePages(db *KV) error {
	// extend the mmap if needed
	size := (db.page.flushed + db.page.nappend) * BTREE_PAGE_SIZE
	if err := extendMmap(db, int(size)); err != nil {
		return err
	}
	// write data pages to the file
	for ptr, node := range db.page.updates {
		offset := int64(ptr * BTREE_PAGE_SIZE)
		if _, err := unix.Pwrite(db.fd, node, offset); err != nil {
			return err
		}
	}
	// discard in-memory data
	db.page.flushed += db.page.nappend
	db.page.nappend = 0
	db.page.updates = map[uint64][]byte{}
	return nil
}

// KV interfaces
func (db *KV) Get(key []byte) ([]byte, bool) {
	return db.tree.Get(key)
}
func (db *KV) Set(key []byte, val []byte) error {
	meta := saveMeta(db)
	if err := db.tree.Insert(key, val); err != nil {
		return err
	}
	return updateOrRevert(db, meta)
}
func (db *KV) Del(key []byte) (bool, error) {
	meta := saveMeta(db)
	if deleted, err := db.tree.Delete(key); !deleted {
		return false, err
	}
	err := updateOrRevert(db, meta)
	return err == nil, err
}

// cleanups
func (db *KV) Close() {
	for _, chunk := range db.mmap.chunks {
		err := syscall.Munmap(chunk)
		assert(err == nil)
	}
	_ = syscall.Close(db.fd)
}
