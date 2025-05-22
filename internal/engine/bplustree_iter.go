package engine

import (
	"bytes"
)

// B-tree iterator
type BIter struct {
	tree *BTree
	path []BNode  // from root to leaf, 从根节点到叶子节点的路径（node链）
	pos  []uint16 // indexes into nodes
}

// current KV pair
func (iter *BIter) Deref() ([]byte, []byte) {
	assert(iter.Valid())
	last := len(iter.path) - 1
	node := iter.path[last]
	pos := iter.pos[last]
	return node.getKey(pos), node.getVal(pos)
}

// 是否最左侧哨兵键（dummy key，不是真实用户数据）
// 哨兵值在树叶子节点的最左端
func iterIsFirst(iter *BIter) bool {
	for _, pos := range iter.pos {
		if pos != 0 {
			return false
		}
	}
	return true // the first key is an dummy sentry
}

// 是否超出了叶子节点末尾
func iterIsEnd(iter *BIter) bool {
	last := len(iter.path) - 1
	return last < 0 || iter.pos[last] >= iter.path[last].nkeys()
}

// 表示不在哨兵位置并且超出叶子节点末尾
func (iter *BIter) Valid() bool {
	return !(iterIsFirst(iter) || iterIsEnd(iter))
}

func iterPrev(iter *BIter, level int) {
	if iter.pos[level] > 0 {
		iter.pos[level]-- // move within this node
	} else if level > 0 {
		iterPrev(iter, level-1) // move to a sibling node
	} else {
		panic("unreachable") // dummy key
	}
	if level+1 < len(iter.pos) { // update the child node
		node := iter.path[level]
		kid := BNode(iter.tree.get(node.getPtr(iter.pos[level])))
		iter.path[level+1] = kid
		iter.pos[level+1] = kid.nkeys() - 1
	}
}

//                 [         20         ] 键数量：1（nkeys() == 1）子节点数量：2（[10] 和 [40]）
//                /                     \
//           [ 10 ]                    [ 40 ]
//          /     \                  /       \
//    [nil, 5]  [10,15]      [20,25,30]  [40,45,50]

// iter.path: [ [20], [10], [10,15] ]
// iter.pos : [   0 ,   1 ,    1    ]   // 你刚访问完 [10,15] 中的 15

// BNode.nkeys()非叶子节点存储孩子数，叶子节点存储kv键值对数
// 叶子节点和非叶子节点的后移不一样,叶子节点后移值,非叶子节点移边界
// 给一个找到节点15的下一个节点的过程:
//  1. iter.pos[level]+1 < iter.path[level].nkeys() 已经不能后移了，调用iterNext(iter, 1)
//  2. 非叶子节点[10],此时的pos已经为1,不能后移了,调用iterNext(iter, 0)
//  3. 非叶子节点[20]边界可以后移, kid := BNode(iter.tree.get(node.getPtr(iter.pos[level])))取出当前非叶子节点移动后的子节点
//     iter.path: [ [20], [10], [10,15] ] 更新为 [ [20], [40], [10,15] ], 并且pos指向新节点的 第0个位置
//     4.回溯上述步骤. iter.path:[ [20], [40], [10,15] ] 更新为 [ [20], [40], [20,25,30] ]
//
// 栈式路径回溯 B树迭代器
func iterNext(iter *BIter, level int) {
	if iter.pos[level]+1 < iter.path[level].nkeys() {
		// 情况 1：当前节点还能右移
		iter.pos[level]++ // move within this node
	} else if level > 0 {
		// 情况 2：当前节点到头了，往上层找兄弟
		iterNext(iter, level-1) // move to a sibling node
	} else {
		// 情况 3：已经在根节点的最后一个 key，彻底结束
		// leaf为叶子节点的最后一个位置, iter.pos[leaf]++
		//相当于尝试访问下一个 key，但其实已经到了最后一个 key，此时只能走到 叶子末尾之后的位置
		leaf := len(iter.pos) - 1
		iter.pos[leaf]++
		//这是个断言：现在 pos 应该正好等于叶子节点的 key 总数，说明已经超出了最后一个有效 key 的位置。
		assert(iter.pos[leaf] == iter.path[leaf].nkeys())
		return // past the last key
	}

	// 若当前层不是叶子层，需要往下更新 path 和 pos
	if level+1 < len(iter.pos) { // update the child node
		//取出当前层的节点
		node := iter.path[level]
		//从当前节点中取出第 pos[level] 个指针（指向子节点）
		kid := BNode(iter.tree.get(node.getPtr(iter.pos[level])))
		iter.path[level+1] = kid
		iter.pos[level+1] = 0
	}
}

func (iter *BIter) Prev() {
	if !iterIsFirst(iter) {
		iterPrev(iter, len(iter.path)-1)
	}
}

func (iter *BIter) Next() {
	if !iterIsEnd(iter) { // 不是最后一个 key 才能前进
		iterNext(iter, len(iter.path)-1) // 从最底层叶子节点开始尝试前进
	}
}

// find the closest position that is less or equal to the input key
func (tree *BTree) SeekLE(key []byte) *BIter {
	iter := &BIter{tree: tree}
	for ptr := tree.root; ptr != 0; {
		node := BNode(tree.get(ptr))
		idx := nodeLookupLE(node, key)
		iter.path = append(iter.path, node)
		iter.pos = append(iter.pos, idx)

		ptr = node.getPtr(idx)
	}
	return iter
}

const (
	CMP_GE = +3 // >=
	CMP_GT = +2 // >
	CMP_LT = -2 // <
	CMP_LE = -3 // <=
)

// key cmp ref
func cmpOK(key []byte, cmp int, ref []byte) bool {
	r := bytes.Compare(key, ref)
	switch cmp {
	case CMP_GE:
		return r >= 0
	case CMP_GT:
		return r > 0
	case CMP_LT:
		return r < 0
	case CMP_LE:
		return r <= 0
	default:
		panic("what?")
	}
}

// find the closest position to a key with respect to the `cmp` relation
func (tree *BTree) Seek(key []byte, cmp int) *BIter {
	iter := tree.SeekLE(key)
	assert(iterIsFirst(iter) || !iterIsEnd(iter))
	if cmp != CMP_LE {
		cur := []byte(nil) // dummy key
		if !iterIsFirst(iter) {
			cur, _ = iter.Deref()
		}
		if len(key) == 0 || !cmpOK(cur, cmp, key) {
			// off by one
			if cmp > 0 {
				iter.Next()
			} else {
				iter.Prev()
			}
		}
	}
	if iter.Valid() {
		cur, _ := iter.Deref()
		assert(cmpOK(cur, cmp, key))
	}
	return iter
}
