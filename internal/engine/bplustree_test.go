package engine

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
	"unsafe"

	is "github.com/stretchr/testify/require"
)

type C struct {
	tree  BTree
	ref   map[string]string
	pages map[uint64]BNode
}

func newC() *C {
	pages := map[uint64]BNode{}
	return &C{
		tree: BTree{
			get: func(ptr uint64) []byte {
				node, ok := pages[ptr]
				assert(ok)
				return node
			},
			new: func(node []byte) uint64 {
				assert(BNode(node).nbytes() <= BTREE_PAGE_SIZE)
				ptr := uint64(uintptr(unsafe.Pointer(&node[0])))
				assert(pages[ptr] == nil)
				pages[ptr] = node
				return ptr
			},
			del: func(ptr uint64) {
				assert(pages[ptr] != nil)
				delete(pages, ptr)
			},
		},
		ref:   map[string]string{},
		pages: pages,
	}
}

func (c *C) add(key string, val string) {
	_, err := c.tree.Upsert([]byte(key), []byte(val))
	assert(err == nil)
	c.ref[key] = val
}

func (c *C) del(key string) bool {
	delete(c.ref, key)
	deleted, err := c.tree.Delete(&DeleteReq{Key: []byte(key)})
	assert(err == nil)
	return deleted
}

func (c *C) dump() ([]string, []string) {
	keys := []string{}
	vals := []string{}

	var nodeDump func(uint64)
	nodeDump = func(ptr uint64) {
		node := BNode(c.tree.get(ptr))
		nkeys := node.nkeys()
		if node.btype() == BNODE_LEAF {
			for i := uint16(0); i < nkeys; i++ {
				keys = append(keys, string(node.getKey(i)))
				vals = append(vals, string(node.getVal(i)))
			}
		} else {
			for i := uint16(0); i < nkeys; i++ {
				ptr := node.getPtr(i)
				nodeDump(ptr)
			}
		}
	}

	nodeDump(c.tree.root)
	assert(keys[0] == "")
	assert(vals[0] == "")
	return keys[1:], vals[1:]
}

type sortIF struct {
	len  int
	less func(i, j int) bool
	swap func(i, j int)
}

func (self sortIF) Len() int {
	return self.len
}
func (self sortIF) Less(i, j int) bool {
	return self.less(i, j)
}
func (self sortIF) Swap(i, j int) {
	self.swap(i, j)
}

func (c *C) verify(t *testing.T) {
	keys, vals := c.dump()

	rkeys, rvals := []string{}, []string{}
	for k, v := range c.ref {
		rkeys = append(rkeys, k)
		rvals = append(rvals, v)
	}
	is.Equal(t, len(rkeys), len(keys))
	sort.Stable(sortIF{
		len:  len(rkeys),
		less: func(i, j int) bool { return rkeys[i] < rkeys[j] },
		swap: func(i, j int) {
			k, v := rkeys[i], rvals[i]
			rkeys[i], rvals[i] = rkeys[j], rvals[j]
			rkeys[j], rvals[j] = k, v
		},
	})

	is.Equal(t, rkeys, keys)
	is.Equal(t, rvals, vals)

	var nodeVerify func(BNode)
	nodeVerify = func(node BNode) {
		nkeys := node.nkeys()
		assert(nkeys >= 1)
		if node.btype() == BNODE_LEAF {
			return
		}
		for i := uint16(0); i < nkeys; i++ {
			key := node.getKey(i)
			kid := BNode(c.tree.get(node.getPtr(i)))
			is.Equal(t, key, kid.getKey(0))
			nodeVerify(kid)
		}
	}

	nodeVerify(c.tree.get(c.tree.root))
}

func fmix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func commonTestBasic(t *testing.T, hasher func(uint32) uint32) {
	c := newC()
	c.add("k", "v")
	c.verify(t)

	// insert
	for i := 0; i < 250000; i++ {
		key := fmt.Sprintf("key%d", hasher(uint32(i)))
		val := fmt.Sprintf("vvv%d", hasher(uint32(-i)))
		c.add(key, val)
		if i < 2000 {
			c.verify(t)
		}
	}
	c.verify(t)

	// del
	for i := 2000; i < 250000; i++ {
		key := fmt.Sprintf("key%d", hasher(uint32(i)))
		is.True(t, c.del(key))
	}
	c.verify(t)

	// overwrite
	for i := 0; i < 2000; i++ {
		key := fmt.Sprintf("key%d", hasher(uint32(i)))
		val := fmt.Sprintf("vvv%d", hasher(uint32(+i)))
		c.add(key, val)
		c.verify(t)
	}

	is.False(t, c.del("kk"))

	for i := 0; i < 2000; i++ {
		key := fmt.Sprintf("key%d", hasher(uint32(i)))
		is.True(t, c.del(key))
		c.verify(t)
	}

	c.add("k", "v2")
	c.verify(t)
	c.del("k")
	c.verify(t)

	// the dummy empty key
	is.Equal(t, 1, len(c.pages))
	is.Equal(t, uint16(1), BNode(c.tree.get(c.tree.root)).nkeys())
}

func TestBTreeBasicAscending(t *testing.T) {
	commonTestBasic(t, func(h uint32) uint32 { return +h })
}

func TestBTreeBasicDescending(t *testing.T) {
	commonTestBasic(t, func(h uint32) uint32 { return -h })
}

func TestBTreeBasicRand(t *testing.T) {
	commonTestBasic(t, fmix32)
}

func TestBTreeRandLength(t *testing.T) {
	c := newC()
	for i := 0; i < 2000; i++ {
		klen := fmix32(uint32(2*i+0)) % BTREE_MAX_KEY_SIZE
		vlen := fmix32(uint32(2*i+1)) % BTREE_MAX_VAL_SIZE
		if klen == 0 {
			continue
		}

		key := make([]byte, klen)
		rand.Read(key)
		val := make([]byte, vlen)
		// rand.Read(val)
		c.add(string(key), string(val))
		c.verify(t)
	}
}

func TestBTreeIncLength(t *testing.T) {
	for l := 1; l < BTREE_MAX_KEY_SIZE+BTREE_MAX_VAL_SIZE; l++ {
		c := newC()

		klen := l
		if klen > BTREE_MAX_KEY_SIZE {
			klen = BTREE_MAX_KEY_SIZE
		}
		vlen := l - klen
		key := make([]byte, klen)
		val := make([]byte, vlen)

		factor := BTREE_PAGE_SIZE / l
		size := factor * factor * 2
		if size > 4000 {
			size = 4000
		}
		if size < 10 {
			size = 10
		}
		for i := 0; i < size; i++ {
			rand.Read(key)
			c.add(string(key), string(val))
		}
		c.verify(t)
	}
}

// without bisect 一千万数据
// === RUN   TestBisect
// Set time: 3m2.245314963s
// Get time: 25.479808521s
// Del time: 2m41.160816668s

// with bisect 一千万数据
// === RUN   TestBisect
// Set time: 2m47.297635704s
// Get time: 11.998730957s
// Del time: 2m29.35970286s
func TestBisect(t *testing.T) {
	c := newC()

	const n = 100000
	keys := make([][]byte, n)
	for i := 0; i < n; i++ {
		keys[i] = []byte(fmt.Sprintf("%05d", i))
	}
	// 随机打乱用于后续查询
	rand.Shuffle(n, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// 测试 Set
	start := time.Now()
	for i := 0; i < n; i++ {
		c.add(string(keys[i]), string("val-"+string(keys[i])))
	}

	fmt.Println("Set time:", time.Since(start))

	// 测试 Get
	start = time.Now()
	for i := 0; i < n; i++ {
		_, ok := c.tree.Get(keys[i])
		if !ok {
			t.Error(ok)
		}
	}
	fmt.Println("Get time:", time.Since(start))

	// 测试 Del keys[i]
	start = time.Now()
	for i := 0; i < n; i++ {
		_ = c.del(string(keys[i]))
	}
	fmt.Println("Del time:", time.Since(start))
}

func TestBTreeIter(t *testing.T) {
	{
		c := newC()
		iter := c.tree.SeekLE(nil)
		is.False(t, iter.Valid())
	}

	sizes := []int{5, 2500}
	for _, sz := range sizes {
		c := newC()

		for i := 0; i < sz; i++ {
			key := fmt.Sprintf("key%010d", i)
			val := fmt.Sprintf("vvv%d", fmix32(uint32(-i)))
			c.add(key, val)
		}
		c.verify(t)

		prevk, prevv := []byte(nil), []byte(nil)
		for i := 0; i < sz; i++ {
			key := []byte(fmt.Sprintf("key%010d", i))
			val := []byte(fmt.Sprintf("vvv%d", fmix32(uint32(-i))))
			// fmt.Println(i, string(key), val)

			iter := c.tree.SeekLE(key)
			is.True(t, iter.Valid())
			gotk, gotv := iter.Deref()
			is.Equal(t, key, gotk)
			is.Equal(t, val, gotv)

			iter.Prev()
			if i > 0 {
				is.True(t, iter.Valid())
				gotk, gotv := iter.Deref()
				is.Equal(t, prevk, gotk)
				is.Equal(t, prevv, gotv)
			} else {
				is.False(t, iter.Valid())
			}

			iter.Next()
			{
				is.True(t, iter.Valid())
				gotk, gotv := iter.Deref()
				is.Equal(t, key, gotk)
				is.Equal(t, val, gotv)
			}

			if i+1 == sz {
				iter.Next()
				is.False(t, iter.Valid())
			}

			prevk, prevv = key, val
		}
	}
}
func TestBTreeIter_LargeScaleNext(t *testing.T) {
	is := is.New(t)
	c := newC()

	sz := 100000
	// 批量插入10万个key-value
	for i := 0; i < sz; i++ {
		key := fmt.Sprintf("key%010d", i)
		val := fmt.Sprintf("val%d", i)
		c.add(key, val)
	}
	c.verify(t)

	// 测试3万个点（这里每隔3个取一个，约3.3万个）
	step := 3
	testCount := 30000
	for i := testCount; i >= 1; i-- {
		idx := i * step
		if idx >= sz-1 {
			break // 防止越界，最后一个没下一个了
		}
		key := []byte(fmt.Sprintf("key%010d", idx))
		nextKey := []byte(fmt.Sprintf("key%010d", idx+1))
		val := []byte(fmt.Sprintf("val%d", idx))
		nextVal := []byte(fmt.Sprintf("val%d", idx+1))

		// SeekLE 找到当前key
		iter := c.tree.SeekLE(key)
		is.True(iter.Valid())
		gotk, gotv := iter.Deref()
		is.Equal(key, gotk)
		is.Equal(val, gotv)

		// Next 找到下一个key
		iter.Next()
		is.True(iter.Valid())
		gotk, gotv = iter.Deref()
		is.Equal(nextKey, gotk)
		is.Equal(nextVal, gotv)
	}
}
