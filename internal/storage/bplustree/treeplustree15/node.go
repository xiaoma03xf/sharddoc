package bplustree15

import (
	"github.com/itxiaoma0610/sharddoc/internal/lib/logger"
)

/*
* 度数、序数，阶数
UPPER_BOUND : 阶为 M, 关键字数为 M-1
UNDER_BOUND : (M-1)/2

数据节点 -- 叶子节点
索引节点 -- 非叶子节点

增加: 分裂 {数据节点, 索引节点}
删除:

	叶子节点
	借: 左兄弟, 右兄弟
	都穷: 搭伙过日子 合并

* B+树的定义:
*1.任意非叶子结点最多有M个子节点, 且M>2, M为B+树的阶数
*2.除根结点以外的非叶子结点至少有(M+1)/2个子节点:
*3.根结点至少有2个子节点;
*4.除根节点外每个结点存放至少(M-1)/2和至多M-1个关键字;(至少1个关键字)
*5.非叶子结点的子树指针比关键字多1个;

*6.非叶子节点的所有key按升序存放，假设节点的关键字分别为K[0]，K[1]…K[N-2],指向子女的指针分别为P[0]，P[1].P[N-1]。则有:

	P[0]< K[0] <= P[1]< K[1] …..< K[M-2] <= P[M-1]

*7.所有叶子结点位于同一层;
*8.为所有叶子结点增加一个链指针:
*9.所有关键字部在叶子结点出现
*/

// BPlusTree15 impl
type BPlusTree15[K Ordered, V any] struct {
	Degree     int
	UpperBound int
	UnderBound int
	Root       *leafNode[K, V]

	// 叶子节点 链表的角度
	Head *leafNode[K, V]
	Tail *leafNode[K, V]
}

func NewBPlusTree15[K Ordered, V any](degree int) *BPlusTree15[K, V] {
	if degree < 3 {
		logger.Warn("build bplustree err, degree not allowed under 3:", degree)
		return nil
	}
	b := &BPlusTree15[K, V]{}
	b.Degree = degree
	b.UpperBound = degree - 1
	b.UnderBound = b.UpperBound / 2
	b.Root = newLeafNode[K, V](true, true)
	b.Head = b.Root
	b.Tail = b.Root

	return b
}
func (tree *BPlusTree15[K, V]) Put(key K, value V) {
	var zero K
	if key == zero {
		return
	}
}
