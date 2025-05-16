package bplustree

import (
	//"log"
	"fmt"
	"sort"
)

// kv 结构体表示一个键值对
type kv struct {
	key   int    // 键（整数类型）
	value string // 值（字符串类型）
}

// kvs 类型是一个固定大小的键值对数组，大小由MaxKV常量决定
type kvs [MaxKV]kv

// 以下三个方法实现了sort.Interface接口，使kvs可以被sort包排序
func (a *kvs) Len() int           { return len(a) } // 返回数组长度
func (a *kvs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] } // 交换两个元素
func (a *kvs) Less(i, j int) bool { return a[i].key < a[j].key } // 比较两个元素

// String方法返回键值对数组的字符串表示，用于调试
func (a *kvs) String() string {
	var s string
	for _, kv := range a {
		s += fmt.Sprintf("%d\t", kv.key)
	}
	return s
}

// leafNode结构体表示B+树的叶子节点
type leafNode struct {
	kvs   kvs           // 存储键值对的数组
	count int           // 当前存储的键值对数量
	next  *leafNode     // 指向下一个叶子节点，形成链表结构，便于范围查询
	p     *interiorNode // 指向父节点（内部节点）
}

// newLeafNode创建并返回一个新的叶子节点
// 参数p指定新节点的父节点
func newLeafNode(p *interiorNode) *leafNode {
	return &leafNode{
		p: p, // 设置父节点引用
	}
}

// find方法在叶子节点中查找指定的键
// 参数:
//   - key: 要查找的键
// 返回值:
//   - int: 如果找到键，返回其索引；如果未找到，返回应插入的位置
//   - bool: 如果找到键，返回true；否则返回false
// 
// 如果键存在于节点中，返回该键的索引和true
// 如果键不存在于节点中，返回应插入该键的位置（节点中第一个大于给定键的键的索引）和false
func (l *leafNode) find(key int) (int, bool) {
	// 定义一个比较函数，用于二分查找
	c := func(i int) bool {
		return l.kvs[i].key >= key // 当数组元素的键大于等于目标键时返回true
	}

	// 使用标准库的二分查找函数
	i := sort.Search(l.count, c) // 在前l.count个元素中搜索满足c函数的最小索引

	// 判断是否找到了精确匹配
	if i < l.count && l.kvs[i].key == key {
		return i, true // 找到了键，返回索引和true
	}

	return i, false // 未找到键，返回应插入的位置和false
}

// insert方法将键值对插入叶子节点
// 参数:
//   - key: 要插入的键
//   - value: 与键关联的值
// 返回值:
//   - int: 如果节点分裂，返回新节点的第一个键；否则返回0
//   - bool: 如果节点分裂，返回true；否则返回false
func (l *leafNode) insert(key int, value string) (int, bool) {
	// 先查找键是否已存在
	i, ok := l.find(key)

	// 如果键已存在，则更新值
	if ok {
		//log.Println("insert.replace", i)
		l.kvs[i].value = value
		return 0, false // 返回0和false表示没有分裂
	}

	// 如果节点未满，直接插入
	if !l.full() {
		// 将i及之后的元素向后移动一位，为新元素腾出空间
		copy(l.kvs[i+1:], l.kvs[i:l.count])
		// 在位置i插入新的键值对
		l.kvs[i].key = key
		l.kvs[i].value = value
		l.count++ // 增加计数
		return 0, false // 返回0和false表示没有分裂
	}

	// 如果节点已满，需要分裂
	next := l.split() // 分裂当前节点

	// 根据键的大小决定将新键值对插入到哪个节点
	if key < next.kvs[0].key {
		l.insert(key, value) // 插入到原节点
	} else {
		next.insert(key, value) // 插入到新节点
	}

	// 返回新节点的第一个键和true，表示发生了分裂
	return next.kvs[0].key, true
}

// split方法将当前叶子节点分裂成两个节点
// 返回值:
//   - *leafNode: 分裂后新创建的叶子节点
func (l *leafNode) split() *leafNode {
	// 创建一个新的叶子节点，暂不设置父节点
	next := newLeafNode(nil)

	// 将当前节点后半部分的键值对复制到新节点
	copy(next.kvs[0:], l.kvs[l.count/2+1:])

	// 设置新节点的计数和next指针
	next.count = MaxKV - l.count/2 - 1 // 新节点包含原节点后半部分的键值对
	next.next = l.next // 新节点的next指向原节点的next

	// 更新当前节点的计数和next指针
	l.count = l.count/2 + 1 // 保留前半部分的键值对
	l.next = next // 当前节点的next指向新节点

	// 返回新创建的节点
	return next
}

// full方法检查叶子节点是否已满
// 返回值:
//   - bool: 如果节点已满，返回true；否则返回false
func (l *leafNode) full() bool { return l.count == MaxKV }

// parent方法返回叶子节点的父节点
// 返回值:
//   - *interiorNode: 父节点指针
func (l *leafNode) parent() *interiorNode { return l.p }

// setParent方法设置叶子节点的父节点
// 参数:
//   - p: 要设置的父节点指针
func (l *leafNode) setParent(p *interiorNode) { l.p = p }
