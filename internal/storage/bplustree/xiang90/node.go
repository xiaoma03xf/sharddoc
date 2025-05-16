package bplustree

// 定义两个重要常量：
const (
	MaxKV = 255 // 叶子节点中最大可存储的键值对数量
	MaxKC = 511 // 内部节点中最大可存储的键和子节点对数量
)

// node 接口定义了B+树中所有节点（内部节点和叶子节点）必须实现的行为
type node interface {
	// find 根据提供的键查找节点中的位置
	// 返回值: 
	// - int: 找到键的索引或应当插入的位置
	// - bool: 表示键是否已存在于节点中
	find(key int) (int, bool)
	
	// parent 返回节点的父节点（内部节点）
	parent() *interiorNode
	
	// setParent 设置节点的父节点
	setParent(*interiorNode)
	
	// full 检查节点是否已满，用于确定是否需要分裂
	full() bool
}
