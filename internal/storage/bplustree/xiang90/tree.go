package bplustree

import "fmt"

// BTree结构体表示B+树的整体结构
type BTree struct {
	root     *interiorNode // 指向根节点，注意是内部节点（根不是叶子）
	first    *leafNode     // 指向第一个叶子节点，用于遍历（链表结构）
	leaf     int           // 当前叶子节点数量
	interior int           // 当前内部节点数量
	height   int           // 树的高度
}

// newBTree创建并返回一个新的B+树
// 初始B+树包含一个内部节点（根节点）和一个叶子节点
// 返回值:
//   - *BTree: 新创建的B+树
func newBTree() *BTree {
	leaf := newLeafNode(nil)        // 创建一个孤立的叶子节点
	r := newInteriorNode(nil, leaf) // 创建一个内部节点作为根节点，指向叶子节点
	leaf.p = r                      // 设置叶子节点的父节点指回根节点，建立双向连接

	// 返回初始化的B+树
	return &BTree{
		root:     r,    // 设置根节点
		first:    leaf, // 设置第一个叶子节点
		leaf:     1,    // 初始有1个叶子节点
		interior: 1,    // 初始有1个内部节点
		height:   2,    // 初始树高为2（根节点+叶子节点层）
	}
}

// First方法返回B+树的第一个叶子节点
// 这个方法可用于开始顺序遍历所有叶子节点（通过叶子节点的next指针）
// 返回值:
//   - *leafNode: 第一个叶子节点的指针
func (bt *BTree) First() *leafNode {
	return bt.first
}

// Insert方法将键值对插入B+树
// 参数:
//   - key: 要插入的键
//   - value: 与键关联的值
func (bt *BTree) Insert(key int, value string) {
	// 首先查找键应该插入到哪个叶子节点
	_, oldIndex, leaf := search(bt.root, key)
	p := leaf.parent() // 获取叶子节点的父节点

	// 将键值对插入叶子节点
	mid, bump := leaf.insert(key, value)
	// 如果没有发生分裂(bump为false)，直接返回
	if !bump {
		return
	}

	// 如果发生了分裂，需要处理分裂后的情况
	var midNode node
	midNode = leaf // midNode指向原始叶子节点

	// 更新父节点中对应的子节点指针，指向分裂后的新节点
	p.kcs[oldIndex].child = leaf.next
	leaf.next.setParent(p) // 设置新节点的父节点

	// 准备向上传播分裂
	interior, interiorP := p, p.parent()

	// 循环处理内部节点的分裂，直到不需要分裂或到达根节点
	for {
		var oldIndex int
		var newNode *interiorNode

		// 检查当前处理的节点是否为根节点
		isRoot := interiorP == nil

		// 如果不是根节点，找到键在父节点中的位置
		if !isRoot {
			oldIndex, _ = interiorP.find(key)
		}

		// 将中间键和节点插入当前内部节点
		mid, newNode, bump = interior.insert(mid, midNode)
		// 如果没有发生分裂，结束处理
		if !bump {
			return
		}

		// 处理分裂后的情况
		if !isRoot {
			// 如果不是根节点，更新父节点中对应的子节点指针
			interiorP.kcs[oldIndex].child = newNode
			newNode.setParent(interiorP)

			// 准备处理下一层
			midNode = interior
		} else {
			// 如果是根节点，创建新的根节点
			bt.root = newInteriorNode(nil, newNode)
			newNode.setParent(bt.root)

			// 将原根节点作为新根节点的子节点
			bt.root.insert(mid, interior)
			return
		}

		// 移动到上一层继续处理
		interior, interiorP = interiorP, interiorP.parent()
	}
}

// Search方法在B+树中查找指定的键
// 参数:
//   - key: 要查找的键
//
// 返回值:
//   - string: 如果找到键，返回对应的值；否则返回空字符串
//   - bool: 如果找到键，返回true；否则返回false
//
// 如果键存在，返回键对应的值和true
// 如果键不存在，返回空字符串和false
func (bt *BTree) Search(key int) (string, bool) {
	// 调用内部search函数查找键
	kv, _, _ := search(bt.root, key)
	// 如果没有找到键（kv为nil），返回空字符串和false
	if kv == nil {
		return "", false
	}
	// 找到键，返回对应的值和true
	return kv.value, true
}

// search是内部查找函数，用于在B+树中查找指定的键
// 参数:
//   - n: 开始查找的节点（通常是根节点）
//   - key: 要查找的键
//
// 返回值:
//   - *kv: 如果找到键，返回指向键值对的指针；否则返回nil
//   - int: 最后一个内部节点中子节点的索引
//   - *leafNode: 包含键或应该包含键的叶子节点
func search(n node, key int) (*kv, int, *leafNode) {
	curr := n      // 当前节点，初始为传入的节点
	oldIndex := -1 // 记录最后一个内部节点中子节点的索引

	// 循环直到找到叶子节点
	for {
		// 根据节点类型进行不同处理
		switch t := curr.(type) {
		case *leafNode:
			// 如果是叶子节点，在叶子节点中查找键
			i, ok := t.find(key)
			// 如果没有找到键，返回nil、最后的索引和叶子节点
			if !ok {
				// 查找调试信息
				fmt.Printf("DEBUG: key %d not found in leaf node, i=%d, count=%d\n", key, i, t.count)
				for j := 0; j < t.count; j++ {
					fmt.Printf("DEBUG: leaf kvs[%d]=%d\n", j, t.kvs[j].key)
				}
				if key == 189 {
					fmt.Println("DEBUG: Searching for key 189")
				}
				return nil, oldIndex, t
			}
			// 找到键，返回指向键值对的指针、最后的索引和叶子节点
			return &t.kvs[i], oldIndex, t
		case *interiorNode:
			// 如果是内部节点，找到应该往下查找的子节点
			i, _ := t.find(key)
			// 更新当前节点为找到的子节点
			curr = t.kcs[i].child
			// 记录子节点的索引
			oldIndex = i
			if key == 189 {
				fmt.Printf("DEBUG: key 189 - interior node level, count=%d, i=%d\n", t.count, i)
				for j := 0; j < t.count; j++ {
					fmt.Printf("DEBUG: interior kcs[%d].key=%d\n", j, t.kcs[j].key)
				}
			}
		default:
			// 如果节点类型既不是叶子节点也不是内部节点，抛出异常
			// 这种情况不应该发生
			panic("未知的节点类型")
		}
	}
}
