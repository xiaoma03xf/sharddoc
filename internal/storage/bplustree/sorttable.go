package bplustree

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// SortTable 是一个有序表接口
type SortTable[K Ordered, V any] interface {
	Put(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	ContainsKey(key K) bool
	FloorKey(key K) (K, bool)
	CeilingKey(key K) (K, bool)
	FirstKey() (K, bool)
	LastKey() (K, bool)
	Size() int
}
