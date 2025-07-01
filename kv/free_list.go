package kv

import "encoding/binary"

// node format:
// | next | pointers | unused |
// |  8B  |   n*8B   |   ...  |

//
// 如果更新过程中断（比如断电），元页面依然指向旧的数据，所以数据不会丢失，
// 也不需要额外的崩溃恢复步骤。而且，空闲列表页面不像元页面那样关键，
// 不需要保证写操作的原子性（即写入必须一次性完成），即使写了一半也不会导致数据损坏。

// 添加空闲页面：当某页变空闲了，就把它的编号写进链表尾巴所在的页面里（追加进去）。
// 使用空闲页面：当数据库要申请一个新页面时，从链表头部取一个空闲页面编号，表示这页被借出去了。
// 链表的下一个节点指针和追加编号字段是直接修改的，但操作安全，因为不覆盖旧数据
// 因为就算操作失败(如断电), 文件会自动加载错误回退，但是此时链表的状态还是没变，还是指向之前的位置
// 链表没坏,不会破坏数据完整性

type LNode []byte

const FREE_LIST_HEADER = 8
const FREE_LIST_CAP = (BTREE_PAGE_SIZE - FREE_LIST_HEADER) / 16

// getters & setters
func (node LNode) getNext() uint64 {
	return binary.LittleEndian.Uint64(node[0:8])
}
func (node LNode) setNext(next uint64) {
	binary.LittleEndian.PutUint64(node[0:8], next)
}
func (node LNode) getItem(idx int) (uint64, uint64) {
	offset := FREE_LIST_HEADER + 16*idx
	return binary.LittleEndian.Uint64(node[offset:]),
		binary.LittleEndian.Uint64(node[offset+8:])
}
func (node LNode) setItem(idx int, ptr uint64, version uint64) {
	assert(idx < FREE_LIST_CAP)
	offset := FREE_LIST_HEADER + 16*idx
	binary.LittleEndian.PutUint64(node[offset+0:], ptr)
	binary.LittleEndian.PutUint64(node[offset+8:], version)
}

type FreeList struct {
	// callbacks for managing on-disk pages
	get func(uint64) []byte // read a page
	new func([]byte) uint64 // append a new page
	set func(uint64) []byte // update an existing page
	// persisted data in the meta page, 持久化到元数据
	headPage uint64 // pointer to the list head node, 指向空闲列表的头节点所在的页面
	headSeq  uint64 // monotonic sequence number to index into the list head,单调递增的序列号，用来标识空闲列表头节点的位置
	tailPage uint64 // 指向空闲列表的尾节点所在的页面
	tailSeq  uint64 //标识空闲列表尾节点的位置
	// in-memory states
	maxSeq uint64 // saved `tailSeq` to prevent consuming newly added items
	maxVer uint64 // the oldest reader version
	curVer uint64 // version number when committing
}

//                      first_item
//                          ↓
// head_page -> [ next |    xxxxx ]
//                 ↓
//              [ next | xxxxxxxx ]
//                 ↓
// tail_page -> [ NULL | xxxx     ]
//                          ↑
//                      last_item

// | next (8 bytes) | ptr0 | ptr1 | ptr2 | ... | ptrN | unused |
//
//	 ↑              ↑
//	下一个节点      存储的空闲页号（uint64）

func seq2idx(seq uint64) int {
	//  | next (8 bytes) | ptr0 | ptr1 | ptr2 | ... | ptrN | unused |
	// freelist中每个节点如上表所示
	// const FREE_LIST_CAP = (BTREE_PAGE_SIZE - FREE_LIST_HEADER) / 8
	// FREE_LIST_CAP 表示从[0,511],第一个位置留给 next(8 bytes)
	return int(seq % FREE_LIST_CAP)
}

func (fl *FreeList) check() {
	assert(fl.headPage != 0 && fl.tailPage != 0)                   //指针不能同时为空
	assert(fl.headSeq != fl.tailSeq || fl.headPage == fl.tailPage) // 不能说page不相等, 但是seq相等了
}

// get 1 item from the list head. return 0 on failure.
func (fl *FreeList) PopHead() uint64 {
	ptr, head := flPop(fl)
	if head != 0 { // the empty head node is recycled
		fl.PushTail(head)
	}
	return ptr
}

// remove 1 item from the head node, and remove the head node if empty.
// 在一个数据库系统里，存储数据的空间是通过页面（page）来管理的。
// 空闲链表（FreeList）就是管理这些 空闲页面 的一个数据结构。每当数据库需要一个新的页面来存储数据时，
// 它就从空闲链表中取一个空闲页面；当某个页面不再使用时，它又会将这个页面放回到空闲链表中
func flPop(fl *FreeList) (ptr uint64, head uint64) {
	fl.check()
	if fl.headSeq == fl.maxSeq {
		return 0, 0 // cannot advance
	}

	//  | next (8 bytes) | ptr0 | ptr1 | ptr2 | ... | ptrN | unused |
	// freelist中每个节点如上表所示
	// node -> fl.get(fl.headPage) 读取当前headPage的数据, 并转为LNode类型
	// 然后 node.getPtr(seq2idx(fl.headSeq)) 这个操作获得具体取节点上的哪个ptr

	node := LNode(fl.get(fl.headPage))
	// 顺序访问 ptr
	ptr, version := node.getItem(seq2idx(fl.headSeq))
	if versionBefore(fl.maxVer, version) {
		return 0, 0 // cannot advance; still in-use
	}

	fl.headSeq++
	// move to the next one if the head node is empty
	// 取之前headSeq = 510, headSeq%FREE_LIST_CAP = 510, 已经是ptr的最后一个位置了
	// headSeq++, headSeq=511,此时 seq2idx(fl.headSeq) == 0 表示当前空闲链表头节点已经被取满
	// headPage指向下一个, head此时指向0
	if seq2idx(fl.headSeq) == 0 {
		head, fl.headPage = fl.headPage, node.getNext()
		assert(fl.headPage != 0)
	}
	return
}

// add 1 item to the tail
func (fl *FreeList) PushTail(ptr uint64) {
	fl.check()
	// add it to the tail node
	// node -> fl.get(fl.headPage) 读取当前headPage的数据, 并转为LNode类型
	// 然后 node.setPtr(seq2idx(fl.tailSeq), ptr) 这个操作设置 ptr添加进哪个位置

	LNode(fl.set(fl.tailPage)).setItem(seq2idx(fl.tailSeq), ptr, fl.curVer)
	fl.tailSeq++

	// add a new tail node if it's full (the list is never empty)
	// seq2idx(fl.tailSeq) 表示当前已经满了
	if seq2idx(fl.tailSeq) == 0 {
		// try to reuse from the list head, 尝试从头部节点回收一个空闲节点
		// 如果next=0,fl.headSeq == fl.maxSeq, 当前节点已经满了
		// 需要新建一个节点存储, 尝试从head中取一个page，如果没有则new一个page
		next, head := flPop(fl) // may remove the head node
		if next == 0 {
			// or allocate a new node by appending
			next = fl.new(make([]byte, BTREE_PAGE_SIZE))
		}

		// link to the new tail node, next指针指向新的page
		LNode(fl.set(fl.tailPage)).setNext(next)
		fl.tailPage = next

		// also add the head node if it's removed
		// 把从head中取到的一个page加入空闲页中 ->(boltDB的freelist节点页管理问题)
		// 这页虽然 将来会被回收，但 BoltDB 会确保 现在没人再需要它 才会去用它，不然数据当然就真的会丢。
		if head != 0 {
			LNode(fl.set(fl.tailPage)).setItem(0, head, fl.curVer)
			fl.tailSeq++
		}
	}
}

// make the newly added items available for consumption
func (fl *FreeList) SetMaxVer(maxVer uint64) {
	fl.maxSeq = fl.tailSeq
	fl.maxVer = maxVer
}
