package etcd

import (
	"fmt"
	"testing"

	"github.com/xiaoma03xf/sharddoc/kv"
)

func TestTableDef(t *testing.T) {
	rt, err := NewTableDefRegistry([]string{"118.89.66.104:2379"})
	if err != nil {
		t.Error(err)
	}
	tdef1 := &kv.TableDef{
		Name:  "StudentsInfo",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	tdef2 := &kv.TableDef{
		Name:  "UserInfo",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	err = rt.RegisterTable(tdef1)
	if err != nil {
		t.Error(err)
	}
	err = rt.RegisterTable(tdef2)
	if err != nil {
		t.Error(err)
	}
	tableDef, _ := rt.GetTable("teacher")
	fmt.Println(tableDef)

	err = rt.PutMetaKey("hello", []byte("dj"))
	if err != nil {
		t.Error(err)
	}
	val, err := rt.GetMetaKey("hello")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(val))
}
