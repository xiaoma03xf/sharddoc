package bplustree

import (
	"fmt"
	"sort"
)

// kc 结构体表示内部节点中的键-子节点对
type kc struct {
	key   int  // 键值
	child node // 子节点指针，可以指向内部节点或叶子节点
}

// kcs 类型是一个固定大小的键-子节点对数组
// 注意数组大小为MaxKC+1，比最大容量多1，为分裂操作预留一个空位
type kcs [MaxKC + 1]kc

// 以下三个方法实现了sort.Interface接口，使kcs可以被sort包排序
func (a *kcs) Len() int { return len(a) } // 返回数组长度

func (a *kcs) Swap(i, j int) { a[i], a[j] = a[j], a[i] } // 交换两个元素

// Less方法比较两个键-子节点对的大小
// 特殊处理键为0的情况（表示空位或特殊位置）
func (a *kcs) Less(i, j int) bool {
	// 如果第一个元素的键为0，认为它更大（排序后会在后面）
	if a[i].key == 0 {
		return false
	}

	// 如果第二个元素的键为0，认为第一个元素更小（排序后会在前面）
	if a[j].key == 0 {
		return true
	}

	// 正常情况下按键的值比较大小
	return a[i].key < a[j].key
}

// String方法返回键-子节点对数组的字符串表示，用于调试
func (a *kcs) String() string {
	var s string
	for _, kc := range a {
		s += fmt.Sprintf("%d\t", kc.key)
	}
	return s
}

// interiorNode结构体表示B+树的内部节点（非叶子节点）
type interiorNode struct {
	kcs   kcs           // 存储键-子节点对的数组
	count int           // 当前存储的键-子节点对数量
	p     *interiorNode // 指向父节点（也是内部节点）
}

// newInteriorNode创建并返回一个新的内部节点
// 参数:
//   - p: 指定新节点的父节点
//   - largestChild: 指定新节点的最大子节点（最右边的子节点）
func newInteriorNode(p *interiorNode, largestChild node) *interiorNode {
	// 创建新的内部节点
	i := &interiorNode{
		p:     p,     // 设置父节点
		count: 1,     // 初始计数为1，因为有一个子节点指针
	}

	// 如果提供了最大子节点，则设置第一个子节点指针
	// 在B+树中，内部节点的子节点指针比键多一个，最左边的子节点对应的键存储在父节点中
	if largestChild != nil {
		i.kcs[0].child = largestChild
	}
	return i
}

// find方法在内部节点中查找指定的键应该对应哪个子节点
// 参数:
//   - key: 要查找的键
// 返回值:
//   - int: 子节点的索引
//   - bool: 总是返回true（与叶子节点的find方法保持接口一致）
func (in *interiorNode) find(key int) (int, bool) {
	// 定义一个比较函数，用于二分查找
	// 当节点中的键大于目标键时返回true
	c := func(i int) bool { return in.kcs[i].key > key }

	// 使用标准库的二分查找函数在前count-1个元素中搜索
	// 返回第一个大于key的键的索引
	i := sort.Search(in.count-1, c)

	// 返回找到的索引和true
	// 注意：内部节点的find总是返回true，因为它只需要确定子节点的位置，不需要精确匹配
	return i, true
}

// full方法检查内部节点是否已满
// 返回值:
//   - bool: 如果节点已满，返回true；否则返回false
func (in *interiorNode) full() bool { return in.count == MaxKC }

// parent方法返回内部节点的父节点
// 返回值:
//   - *interiorNode: 父节点指针
func (in *interiorNode) parent() *interiorNode { return in.p }

// setParent方法设置内部节点的父节点
// 参数:
//   - p: 要设置的父节点指针
func (in *interiorNode) setParent(p *interiorNode) { in.p = p }

// insert方法将键和子节点插入内部节点
// 参数:
//   - key: 要插入的键
//   - child: 与键关联的子节点
// 返回值:
//   - int: 如果节点分裂，返回中间键；否则返回0
//   - *interiorNode: 如果节点分裂，返回新创建的内部节点；否则返回nil
//   - bool: 如果节点分裂，返回true；否则返回false
func (in *interiorNode) insert(key int, child node) (int, *interiorNode, bool) {
	// 找到键应该插入的位置
	i, _ := in.find(key)

	// 如果节点未满，直接插入
	if !in.full() {
		// 将i及之后的元素向后移动一位，为新元素腾出空间
		copy(in.kcs[i+1:], in.kcs[i:in.count])

		// 在位置i插入新的键和子节点
		in.kcs[i].key = key
		in.kcs[i].child = child
		child.setParent(in) // 设置子节点的父节点为当前节点

		in.count++ // 增加计数
		return 0, nil, false // 返回0、nil和false表示没有分裂
	}

	// 如果节点已满，先将新键值对插入到预留的空位（MaxKC位置）
	in.kcs[MaxKC].key = key
	in.kcs[MaxKC].child = child
	child.setParent(in) // 设置子节点的父节点为当前节点

	// 然后分裂节点
	next, midKey := in.split()

	// 返回中间键、新节点和true，表示发生了分裂
	return midKey, next, true
}

// split方法将当前内部节点分裂成两个节点
// 返回值:
//   - *interiorNode: 分裂后新创建的内部节点
//   - int: 分裂位置的键值，将被提升到父节点
func (in *interiorNode) split() (*interiorNode, int) {
	// 先对节点中的所有键-子节点对进行排序
	// 这是因为新键值对可能被插入到了预留位置（MaxKC）
	sort.Sort(&in.kcs)

	// 获取中间位置的信息
	midIndex := MaxKC / 2 // 中间位置索引
	midChild := in.kcs[midIndex].child // 中间位置的子节点
	midKey := in.kcs[midIndex].key // 中间位置的键

	// 创建新的内部节点，暂不设置父节点
	next := newInteriorNode(nil, nil)
	// 将当前节点中间位置之后的键-子节点对复制到新节点
	copy(next.kcs[0:], in.kcs[midIndex+1:])
	next.count = MaxKC - midIndex // 新节点包含原节点后半部分的键-子节点对
	
	// 更新所有被移动的子节点的父节点指针
	for i := 0; i < next.count; i++ {
		next.kcs[i].child.setParent(next)
	}

	// 修改原节点
	in.count = midIndex + 1 // 保留前半部分的键-子节点对
	midChild.setParent(in) // 确保中间子节点的父节点正确设置

	// 返回新创建的节点和中间键
	// 中间键将被提升到父节点中
	return next, midKey
}
