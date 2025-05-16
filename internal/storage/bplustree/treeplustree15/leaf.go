package bplustree15

import "sort"

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}
type kv[K Ordered, V any] struct {
	key   K
	value V
}

type kvs[K Ordered, V any] []kv[K, V]

func (a kvs[K, V]) Len() int           { return len(a) }
func (a kvs[K, V]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a kvs[K, V]) Less(i, j int) bool { return a[i].key < a[j].key }

type leafNode[K Ordered, V any] struct {
	IsRoot bool
	IsLeaf bool

	// key 表示关键字, value 表示叶子节点数据项
	kvs      kvs[K, V]
	Children []*leafNode[K, V] // 非叶子节点的子节点
	Next     *leafNode[K, V]   // 叶子节点的后继
	Prev     *leafNode[K, V]   // 叶子节点的前驱
}

func newLeafNode[K Ordered, V any](isroot, isleaf bool) *leafNode[K, V] {
	n := &leafNode[K, V]{}
	n.IsRoot = isroot
	n.IsLeaf = isleaf

	return n
}

// find the ceiling index of the key
func (l *leafNode[K, V]) findCeilingKeyIndex(key K) int {
	return sort.Search(l.kvs.Len(), func(i int) bool {
		return l.kvs[i].key >= key
	})
}

func (l *leafNode[K, V]) findChildIndexByCeilingKeyIndex(ceilindex int, key K) int {
	if ceilindex == l.kvs.Len() {
		return ceilindex
	}
	
}
func (l *leafNode[K, V]) insert(key K, value V) {
	// 找到大于或者等于 key 在 keys 中索引位置, 针对的是当前节点

	// 根据找到的位置推导出具体孩子节点的索引位置
}
