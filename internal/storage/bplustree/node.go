package bplustree

const (
	MaxKV = 255
	MaxKC = 511
)

type Node interface {
	find(key int) (int, bool)
	parent() *interiorNode
}
