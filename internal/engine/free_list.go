package byodb

import "encoding/binary"

// node format:
// | next | pointers | unused |
// |  8B  |   n*8B   |   ...  |

// | next (8 bytes) | ptr0 | ptr1 | ptr2 | ... | ptrN | unused |
//
//	 ↑              ↑
//	下一个节点      存储的空闲页号（uint64）
//
// 如果更新过程中断（比如断电），元页面依然指向旧的数据，所以数据不会丢失，
// 也不需要额外的崩溃恢复步骤。而且，空闲列表页面不像元页面那样关键，
// 不需要保证写操作的原子性（即写入必须一次性完成），即使写了一半也不会导致数据损坏。

// 添加空闲页面：当某页变空闲了，就把它的编号写进链表尾巴所在的页面里（追加进去）。
// 使用空闲页面：当数据库要申请一个新页面时，从链表头部取一个空闲页面编号，表示这页被借出去了。
// 链表的下一个节点指针和追加编号字段是直接修改的，但操作安全，因为不覆盖旧数据
// 因为就算操作失败(如断电), 之前的数据都还在,链表没坏,不会破坏数据完整性
type LNode []byte

const FREE_LIST_HEADER = 8
const FREE_LIST_CAP = (BTREE_PAGE_SIZE - FREE_LIST_HEADER) / 8

// getters & setters
func (node LNode) getNext() uint64 {
	return binary.LittleEndian.Uint64(node[0:8])
}
func (node LNode) setNext(next uint64) {
	binary.LittleEndian.PutUint64(node[0:8], next)
}
func (node LNode) getPtr(idx int) uint64 {
	offset := FREE_LIST_HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[offset:])
}
func (node LNode) setPtr(idx int, ptr uint64) {
	assert(idx < FREE_LIST_CAP)
	offset := FREE_LIST_HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[offset:], ptr)
}

type FreeList struct {
	// callbacks for managing on-disk pages
	get func(uint64) []byte // read a page
	new func([]byte) uint64 // append a new page
	set func(uint64) []byte // update an existing page
	// persisted data in the meta page
	headPage uint64 // pointer to the list head node
	headSeq  uint64 // monotonic sequence number to index into the list head
	tailPage uint64
	tailSeq  uint64
	// in-memory states
	maxSeq uint64 // saved `tailSeq` to prevent consuming newly added items
}

func seq2idx(seq uint64) int {
	return int(seq % FREE_LIST_CAP)
}

func (fl *FreeList) check() {
	assert(fl.headPage != 0 && fl.tailPage != 0)
	assert(fl.headSeq != fl.tailSeq || fl.headPage == fl.tailPage)
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
func flPop(fl *FreeList) (ptr uint64, head uint64) {
	fl.check()
	if fl.headSeq == fl.maxSeq {
		return 0, 0 // cannot advance
	}
	node := LNode(fl.get(fl.headPage))
	ptr = node.getPtr(seq2idx(fl.headSeq))
	fl.headSeq++
	// move to the next one if the head node is empty
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
	LNode(fl.set(fl.tailPage)).setPtr(seq2idx(fl.tailSeq), ptr)
	fl.tailSeq++
	// add a new tail node if it's full (the list is never empty)
	if seq2idx(fl.tailSeq) == 0 {
		// try to reuse from the list head
		next, head := flPop(fl) // may remove the head node
		if next == 0 {
			// or allocate a new node by appending
			next = fl.new(make([]byte, BTREE_PAGE_SIZE))
		}
		// link to the new tail node
		LNode(fl.set(fl.tailPage)).setNext(next)
		fl.tailPage = next
		// also add the head node if it's removed
		if head != 0 {
			LNode(fl.set(fl.tailPage)).setPtr(0, head)
			fl.tailSeq++
		}
	}
}

// make the newly added items available for consumption
func (fl *FreeList) SetMaxSeq() {
	fl.maxSeq = fl.tailSeq
}
