package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Map struct {
	hashFunc HashFunc       // 哈希函数，用于将节点或键转换为 uint32 的哈希值
	replicas int            // 每个真实节点对应的“虚拟节点”数量（提高分布均衡性）
	keys     []int          // 所有虚拟节点的哈希值，已排序（形成一个“哈希环”）
	hashMap  map[int]string // 虚拟节点哈希值 -> 真实节点名称的映射
}

// 创建一致性哈希 Map
func NewMap(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hashFunc: fn,
		hashMap:  make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// 判断是否为空
func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// 添加真实节点（带多个虚拟节点）
func (m *Map) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < m.replicas; i++ {
			virtualKey := m.hashFunc([]byte(strconv.Itoa(i) + node))
			m.keys = append(m.keys, int(virtualKey))
			m.hashMap[int(virtualKey)] = node
		}
	}
	sort.Ints(m.keys)
}

// 获取与某 key 最近的节点
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))

	// 二分查找第一个 >= hash 的位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 如果超过最大值就回环
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}

// 删除节点
func (m *Map) Remove(node string) {
	for i := 0; i < m.replicas; i++ {
		virtualKey := int(m.hashFunc([]byte(strconv.Itoa(i) + node)))
		delete(m.hashMap, virtualKey)

		// 从 keys 切片中删除
		idx := sort.SearchInts(m.keys, virtualKey)
		if idx < len(m.keys) && m.keys[idx] == virtualKey {
			m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		}
	}
}
