package engine

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// node format:
// | 类型 | 键数量 |   指针数组    |   偏移量数组  |     键值对数据    |
// | type | nkeys |  pointers  |   offsets  | key-values
// |  2B  |   2B  | nkeys * 8B | nkeys * 2B | ...

// key-value format:
// | klen | vlen | key | val |
// |  2B  |  2B  | ... | ... |

const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(node1max <= BTREE_PAGE_SIZE)
}

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

type BNode []byte // can be dumped to the disk

type BTree struct {
	// pointer (a nonzero page number) 存储根节点的页面编号(非零值)
	root uint64

	// 用于管理磁盘页面的回调函数
	// callbacks for managing on-disk pages
	get func(uint64) []byte // dereference a pointer
	new func([]byte) uint64 // allocate a new page
	del func(uint64)        // deallocate a page
}

// header
// btype 返回节点的类型(内部节点或叶子节点)
// binary.LittleEndian.Uint16 从字节切片中读取前2个字节并转换为小端序的16位无符号整数
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// nkeys 返回节点中键的数量
// 从节点头部的第2-4字节读取键数量
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// setHeader 设置节点的头部信息(类型和键数量)
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype) // 写入节点类型
	binary.LittleEndian.PutUint16(node[2:4], nkeys) // 写入键数量
}

// pointers
// getPtr 获取指定索引位置的子节点指针
// binary.LittleEndian.Uint64 从字节切片中读取8个字节并转换为小端序的64位无符号整数
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys()) // 确保索引有效
	pos := HEADER + 8*idx      // 计算指针在节点中的位置
	return binary.LittleEndian.Uint64(node[pos:])
}

// setPtr 设置指定索引位置的子节点指针
func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx < node.nkeys())
	// assert(node.btype() == BNODE_LEAF || val != 0)
	// assert(node.btype() == BNODE_NODE || val == 0)
	pos := HEADER + 8*idx // 计算指针在节点中的位置
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// offset list 计算指定索引的偏移量在节点中的位置
func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

// getOffset 获取指定索引的键值对在数据区域的偏移量
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0 // 第一个键值对的偏移量固定为0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

// setOffset 设置指定索引的键值对偏移量
func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], offset)
}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}
func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}
func (node BNode) getVal(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos+0:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}

// node size in bytes
// nbytes 计算节点的总大小(字节数)
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}

// returns the first kid node whose range intersects the key. (kid[i] <= key)
// nodeLookupLE 查找节点中小于等于指定键的最后一个位置
// 例如：对于键数组 [a, b, d, e, f]
// 查找c: 返回b的索引(1)，因为b < c < d
// 查找d: 返回d的索引(2)，因为d = d
// TODO: bisect
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)
	// the first key is a copy from the parent node,
	// thus it's always less than or equal to the key.
	for i := uint16(1); i < nkeys; i++ {
		// bytes.Compare 返回比较结果:
		// 0: 如果两个键相等 (a==b)
		// -1: 如果第一个键小于第二个键 (a<b)
		// 1: 如果第一个键大于第二个键 (a>b)
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			// 找到完全匹配的键，返回其索引
			found = i
		}
		if cmp >= 0 {
			// 当前键大于目标键，返回前一个位置
			break
		}
	}
	return found
}

// add a new key to a leaf node
// leafInsert 向叶子节点插入一个新的键值对

func leafInsert(
	new BNode, old BNode, idx uint16,
	key []byte, val []byte,
) {
	// 设置新节点为叶子节点，键数量加1
	new.setHeader(BNODE_LEAF, old.nkeys()+1)

	// Range添加idx之前键值, 添加当前kv, Range添加idx之后键值
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// update an existing key from a leaf node
// leafUpdate 更新叶子节点中已存在的键值对
func leafUpdate(
	new BNode, old BNode, idx uint16,
	key []byte, val []byte,
) {
	// 设置新节点为叶子节点，键数量不变
	new.setHeader(BNODE_LEAF, old.nkeys())

	// Range添加idx之前键值, 修改当前kv, Range添加idx之后键值
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-(idx+1))
}

// replace a link with the same key
// nodeReplaceKid1ptr 替换内部节点中的一个子节点指针(键不变)
func nodeReplaceKid1ptr(new BNode, old BNode, idx uint16, ptr uint64) {
	copy(new, old[:old.nbytes()])
	new.setPtr(idx, ptr) // only the pointer is changed
}

// replace a link with multiple links
// nodeReplaceKidN 用多个子节点替换内部节点中的一个子节点
func nodeReplaceKidN(
	tree *BTree, new BNode, old BNode, idx uint16,
	kids ...BNode,
) {
	// 如果要替换的kids数量为一, 直接修改指针引用即可
	inc := uint16(len(kids))
	if inc == 1 && bytes.Equal(kids[0].getKey(0), old.getKey(idx)) {
		// common case, only replace 1 pointer
		nodeReplaceKid1ptr(new, old, idx, tree.new(kids[0]))
		return
	}

	// Range 复制替换idx位置之前数据, 循环添加kids, Range 复制替换idx位置之后数据
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

// replace 2 adjacent links with 1
// nodeReplace2Kid 用一个子节点替换内部节点中的两个相邻子节点(合并操作)
func nodeReplace2Kid(
	new BNode, old BNode, idx uint16,
	ptr uint64, key []byte,
) {
	new.setHeader(BNODE_NODE, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, ptr, key, nil)
	nodeAppendRange(new, old, idx+1, idx+2, old.nkeys()-(idx+2))
}

// copy a KV into the position
// nodeAppendKV 在节点指定位置添加一个键值对
func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	// ptrs
	new.setPtr(idx, ptr)

	// KVs, 写入值长度(2字节), 写入键长度(2字节)
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))

	// 复制键数据, 复制值数据到新节点中
	copy(new[pos+4:], key)
	copy(new[pos+4+uint16(len(key)):], val)

	// the offset of the next key
	// 更新下一个键值对的偏移量
	// 键值对在节点中的存储格式:
	// | 键长度 | 值长度 | 键数据 | 值数据 |
	// | klen  | vlen  |  key  |  val  |
	// |  2B   |  2B   |  ...  |  ...  |
	// 新偏移量 = 旧偏移量 + 键长度(2B) + 值长度(2B) + 键值对数据长度
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

// copy multiple KVs into the position
// nodeAppendRange 从旧节点复制多个键值对到新节点
// The nodeAppendRange function copies keys from an old node to a new node.
// dstNew 表示偏移量, 从目标位置开始
// srcOld 表示旧位置的其实索引, 表示从哪儿开始复制
// n 表示复制的节点数量
func nodeAppendRange(
	new BNode, old BNode,
	dstNew uint16, srcOld uint16, n uint16,
) {
	assert(srcOld+n <= old.nkeys())
	assert(dstNew+n <= new.nkeys())
	if n == 0 {
		return
	}

	// pointers
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}
	// offsets
	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(srcOld)
	for i := uint16(1); i <= n; i++ { // NOTE: the range is [1, n]
		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
		new.setOffset(dstNew+i, offset)
	}
	// KVs
	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)
	copy(new[new.kvPos(dstNew):], old[begin:end])
}

// split a bigger-than-allowed node into two.
// the second node always fits on a page.

// nodeSplit2 将一个过大的节点分裂为两个节点
// 确保第二个节点(right)一定能适应页面大小
// nodeSplit2 函数的目标是将一个 过大的节点（old）分裂为两个节点（left 和 right）
// 并确保 右节点 能够满足页面大小的要求，而 左节点 在某些情况下可能不满足页面大小的要求，但是会尽量接近页面大小。
func nodeSplit2(left BNode, right BNode, old BNode) {
	assert(old.nkeys() >= 2)

	// the initial guess, 初始猜测：左节点包含一半的键
	nleft := old.nkeys() / 2

	// try to fit the left half, 尝试调整左节点大小，使其适应页面大小
	left_bytes := func() uint16 {
		return HEADER + 8*nleft + 2*nleft + old.getOffset(nleft)
	}
	// 如果左节点过大，减少其键数量
	for left_bytes() > BTREE_PAGE_SIZE {
		nleft--
	}
	assert(nleft >= 1)

	// try to fit the right half
	right_bytes := func() uint16 {
		return old.nbytes() - left_bytes() + HEADER
	}
	for right_bytes() > BTREE_PAGE_SIZE {
		nleft++
	}
	assert(nleft < old.nkeys())
	nright := old.nkeys() - nleft

	left.setHeader(old.btype(), nleft)
	right.setHeader(old.btype(), nright)
	nodeAppendRange(left, old, 0, 0, nleft)
	nodeAppendRange(right, old, 0, nleft, nright)

	// the left half may be still too big
	// 注意：左节点可能仍然太大，但右节点必须适应页面大小
	assert(right.nbytes() <= BTREE_PAGE_SIZE)
}

// split a node if it's too big. the results are 1~3 nodes.
func nodeSplit3(old BNode) (uint16, [3]BNode) {
	// 如果节点未超过页面大小，无需分裂
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old} // not split
	}

	// 先尝试分裂为两个节点
	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE)) // might be split later, 左节点可能仍需再次分裂
	right := BNode(make([]byte, BTREE_PAGE_SIZE))
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left = left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right} // 2 nodes
	}

	// 左节点仍然过大，再次分裂左节点
	leftleft := BNode(make([]byte, BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, BTREE_PAGE_SIZE))
	nodeSplit2(leftleft, middle, left)
	assert(leftleft.nbytes() <= BTREE_PAGE_SIZE)
	return 3, [3]BNode{leftleft, middle, right} // 3 nodes
}

// insert a KV into a node, the result might be split.
// the caller is responsible for deallocating the input node
// and splitting and allocating result nodes.
func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	// the result node.
	// it's allowed to be bigger than 1 page and will be split if so
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))

	// where to insert the key?
	idx := nodeLookupLE(node, key)
	// act depending on the node type
	switch node.btype() {
	case BNODE_LEAF:
		// leaf, node.getKey(idx) <= key
		if bytes.Equal(key, node.getKey(idx)) {
			// found the key, update it.
			leafUpdate(new, node, idx, key, val)
		} else {
			// insert it after the position.
			leafInsert(new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		// internal node, insert it to a kid node.
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("bad node!")
	}
	return new
}

// part of the treeInsert(): KV insertion to an internal node
func nodeInsert(
	tree *BTree, new BNode, node BNode, idx uint16,
	key []byte, val []byte,
) {
	kptr := node.getPtr(idx)
	// recursive insertion to the kid node
	knode := treeInsert(tree, tree.get(kptr), key, val)
	// split the result
	nsplit, split := nodeSplit3(knode)
	// deallocate the kid node
	tree.del(kptr)
	// update the kid links
	nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
}

// remove a key from a leaf node
func leafDelete(new BNode, old BNode, idx uint16) {
	new.setHeader(BNODE_LEAF, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-(idx+1))
}

// merge 2 nodes into 1
func nodeMerge(new BNode, left BNode, right BNode) {
	new.setHeader(left.btype(), left.nkeys()+right.nkeys())
	nodeAppendRange(new, left, 0, 0, left.nkeys())
	nodeAppendRange(new, right, left.nkeys(), 0, right.nkeys())
	assert(new.nbytes() <= BTREE_PAGE_SIZE)
}

// delete a key from the tree
func treeDelete(tree *BTree, node BNode, key []byte) BNode {
	// where to find the key?
	idx := nodeLookupLE(node, key)
	// act depending on the node type
	switch node.btype() {
	case BNODE_LEAF:
		if !bytes.Equal(key, node.getKey(idx)) {
			return BNode{} // not found
		}
		// delete the key in the leaf
		new := BNode(make([]byte, BTREE_PAGE_SIZE))
		leafDelete(new, node, idx)
		return new
	case BNODE_NODE:
		return nodeDelete(tree, node, idx, key)
	default:
		panic("bad node!")
	}
}

// part of the treeDelete()
func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {
	// recurse into the kid
	kptr := node.getPtr(idx)
	updated := treeDelete(tree, tree.get(kptr), key)
	if len(updated) == 0 {
		return BNode{} // not found
	}
	tree.del(kptr)

	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	// check for merging
	mergeDir, sibling := shouldMerge(tree, node, idx, updated)
	switch {
	case mergeDir < 0: // left
		merged := BNode(make([]byte, BTREE_PAGE_SIZE))
		nodeMerge(merged, sibling, updated)
		tree.del(node.getPtr(idx - 1))
		nodeReplace2Kid(new, node, idx-1, tree.new(merged), merged.getKey(0))
	case mergeDir > 0: // right
		merged := BNode(make([]byte, BTREE_PAGE_SIZE))
		nodeMerge(merged, updated, sibling)
		tree.del(node.getPtr(idx + 1))
		nodeReplace2Kid(new, node, idx, tree.new(merged), merged.getKey(0))
	case mergeDir == 0 && updated.nkeys() == 0:
		assert(node.nkeys() == 1 && idx == 0) // 1 empty child but no sibling
		new.setHeader(BNODE_NODE, 0)          // the parent becomes empty too
	case mergeDir == 0 && updated.nkeys() > 0: // no merge
		nodeReplaceKidN(tree, new, node, idx, updated)
	}
	return new
}

// should the updated kid be merged with a sibling?
func shouldMerge(
	tree *BTree, node BNode,
	idx uint16, updated BNode,
) (int, BNode) {
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}

	if idx > 0 {
		sibling := BNode(tree.get(node.getPtr(idx - 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling // left
		}
	}
	if idx+1 < node.nkeys() {
		sibling := BNode(tree.get(node.getPtr(idx + 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling // right
		}
	}
	return 0, BNode{}
}

func checkLimit(key []byte, val []byte) error {
	if len(key) == 0 {
		return errors.New("empty key") // used as a dummy key
	}
	if len(key) > BTREE_MAX_KEY_SIZE {
		return errors.New("key too long")
	}
	if len(val) > BTREE_MAX_VAL_SIZE {
		return errors.New("value too long")
	}
	return nil
}

// the interface
func (tree *BTree) Insert(key []byte, val []byte) error {
	// 1. check the length limit imposed by the node format
	if err := checkLimit(key, val); err != nil {
		return err // the only way for an update to fail
	}

	// 2. create the first node
	if tree.root == 0 {
		// create the first node
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeader(BNODE_LEAF, 2)
		// a dummy key, this makes the tree cover the whole key space.
		// thus a lookup can always find a containing node.
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return nil
	}
	// 3. insert the key
	node := treeInsert(tree, tree.get(tree.root), key, val)

	// 4. grow the tree if the root is split
	nsplit, split := nodeSplit3(node)
	tree.del(tree.root)
	if nsplit > 1 {
		// the root was split, add a new level.
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeader(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
	return nil
}

func (tree *BTree) Delete(key []byte) (bool, error) {
	if err := checkLimit(key, nil); err != nil {
		return false, err // the only way for an update to fail
	}

	if tree.root == 0 {
		return false, nil
	}

	updated := treeDelete(tree, tree.get(tree.root), key)
	if len(updated) == 0 {
		return false, nil // not found
	}

	tree.del(tree.root)
	if updated.btype() == BNODE_NODE && updated.nkeys() == 1 {
		// remove a level
		tree.root = updated.getPtr(0)
	} else {
		tree.root = tree.new(updated)
	}
	return true, nil
}

func nodeGetKey(tree *BTree, node BNode, key []byte) ([]byte, bool) {
	idx := nodeLookupLE(node, key)
	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			return node.getVal(idx), true
		} else {
			return nil, false
		}
	case BNODE_NODE:
		return nodeGetKey(tree, tree.get(node.getPtr(idx)), key)
	default:
		panic("bad node!")
	}
}

func (tree *BTree) Get(key []byte) ([]byte, bool) {
	if tree.root == 0 {
		return nil, false
	}
	return nodeGetKey(tree, tree.get(tree.root), key)
}
