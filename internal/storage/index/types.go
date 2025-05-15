package index

import "github.com/itxiaoma0610/sharddoc/internal/storage"

// enum index type
type IndexType int

const (
	ExactMatch IndexType = iota
	Inverted
	Range
)

// index describtion
type IndexDefinition struct {
	Field  string
	Type   IndexType
	Nested bool
}

// common index interface
type Index interface {
	Insert(docID string, doc storage.Document) error
	Delete(docID string, doc storage.Document) error
	Search(query interface{}) ([]string, error)
	Type() string
}

type IndexManager interface {
	AddIndex(field string, index Index)
	RemoveIndex(field string)
	GetIndex(field string) (Index, bool)
	Search(field string, value interface{}) ([]string, error)
	Insert(docID string, doc storage.Document) error
	Delete(docID string, doc storage.Document) error
}

// TODO
// 主键索引 Insert / Delete / Search(接口函数)
// 单字段索引 同上
// 复合索引 同上
// 唯一索引 Insert 中需判断冲突
// TTL 索引  Insert → 添加到 min-heap Delete → 移除堆中 Search 可选
// TTL 不用于查找，而是定期清理数据
// 前缀索引	❌ 可选	做搜索提示、自动补全时才需要，暂不必优先

// | 索引类型   | 推荐数据结构                   | 原因概述          |
// | ------ | ------------------------ | ------------- |
// | 主键索引   | B+ 树 / 跳表                | 有序、高性能查找、范围查询 |
// | 单字段索引  | B+ 树 / Hash 表            | 范围查询 vs 精确匹配  |
// | 复合索引   | B+ 树（复合键）                | 支持多个字段联合查找    |
// | 唯一索引   | B+ 树 + 判断重复              | 快速检索是否存在      |
// | TTL 索引 | 小根堆 / 时间排序 B+ 树          | 快速找到最早过期的数据   |
// | 前缀索引   | Trie / Radix Tree        | 前缀匹配最优结构      |
// | 全文索引   | 倒排索引（Map + Posting List） | 文本关键词检索       |
// | 地理索引   | R-Tree / QuadTree        | 空间分区、坐标定位查询   |

// todo
// 主键索引（必须） 根据主键快速定位一条文档（如 ID, 通常用 B+ 树 来维护主键有序；查找是 O(log n)
// 单字段索引（常用） 按某个字段（如 age, name）进行查询。
// 复合索引 同时用多个字段查询，比如 (name, age) 把多个字段拼成一个复合键，比如 "小明#25" 再建 B+ 树
// TTL 索引（时效数据常见）
// 作用：让文档在某个时间后自动过期、删除（缓存系统很常见）
// 原理：
// 用 小根堆 管理所有过期时间，堆顶是最先过期的
