package bplustree

const (
	// 叶子节点最大存储值
	MaxKV = 4
	// 索引最大宽度
	MaxKC = 4
)

type node interface {
	find(key int) (int, bool)
	parent() *interiorNode
	setParent(*interiorNode)
	full() bool
}
