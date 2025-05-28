package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"

	is "github.com/stretchr/testify/require"
)

// R 是一个测试辅助结构，封装了数据库操作和引用数据以便验证
type R struct {
	db  DB                  // 被测试的数据库实例
	ref map[string][]Record // 参考数据映射，用于验证数据库操作结果
}

// newR 创建一个新的测试辅助实例，初始化测试数据库
func newR() *R {
	os.Remove("r.db")
	r := &R{
		db:  DB{Path: "r.db"},
		ref: map[string][]Record{},
	}
	err := r.db.Open()
	assert(err == nil)
	return r
}

// dispose 关闭数据库并删除测试文件
func (r *R) dispose() {
	r.db.Close()
	os.Remove("r.db")
}

// begin 开始一个新的数据库事务
func (r *R) begin() *DBTX {
	tx := DBTX{}
	r.db.Begin(&tx)
	return &tx
}

// commit 提交一个数据库事务
func (r *R) commit(tx *DBTX) {
	err := r.db.Commit(tx)
	assert(err == nil)
}

// create 在数据库中创建一个新表
func (r *R) create(tdef *TableDef) {
	tx := r.begin()
	err := tx.TableNew(tdef)
	r.commit(tx)
	assert(err == nil)
}

// findRef 通过主键在参考映射中查找记录
// 返回记录在参考数组中的索引，如果未找到则返回-1
func (r *R) findRef(table string, rec Record) int {
	pkeys := len(r.db.tables[table].Indexes[0])
	records := r.ref[table]
	found := -1
	for i, old := range records {
		if reflect.DeepEqual(old.Vals[:pkeys], rec.Vals[:pkeys]) {
			assert(found == -1)
			found = i
		}
	}
	return found
}

// add 向数据库和参考映射中添加记录
// 返回记录是否为新增（true表示新增，false表示更新）
func (r *R) add(table string, rec Record) bool {
	tx := r.begin()
	dbreq := DBUpdateReq{Record: rec}
	_, err := tx.Set(table, &dbreq)
	assert(err == nil)
	r.commit(tx)

	records := r.ref[table]
	idx := r.findRef(table, rec)
	assert((idx < 0) == dbreq.Added)
	if idx < 0 {
		r.ref[table] = append(records, rec)
	} else {
		records[idx] = rec
	}
	return dbreq.Added
}

// del 从数据库和参考映射中删除记录
// 返回是否成功删除
func (r *R) del(table string, rec Record) bool {
	tx := r.begin()
	deleted, err := tx.Delete(table, rec)
	assert(err == nil)
	r.commit(tx)

	idx := r.findRef(table, rec)
	if deleted {
		assert(idx >= 0)
		records := r.ref[table]
		copy(records[idx:], records[idx+1:])
		r.ref[table] = records[:len(records)-1]
	} else {
		assert(idx == -1)
	}

	return deleted
}

// get 从数据库中获取记录并与参考映射进行验证
// 返回是否找到记录
func (r *R) get(table string, rec *Record) bool {
	tx := r.begin()
	ok, err := tx.Get(table, rec)
	assert(err == nil)
	r.commit(tx)
	idx := r.findRef(table, *rec)
	if ok {
		assert(idx >= 0)
		records := r.ref[table]
		assert(reflect.DeepEqual(records[idx], *rec))
	} else {
		assert(idx < 0)
	}
	return ok
}

// TestTableCreate 测试表创建功能
// 验证表结构定义是否正确存储在系统表中
func TestTableCreate(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			os.Remove("r.db")
		}
	}()
	r := newR()

	// 创建一个测试表，包含多种类型的列和索引
	tdef := &TableDef{
		Name:    "tbl_test",
		Cols:    []string{"ki1", "ks2", "s1", "i2"},
		Types:   []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES, TYPE_INT64},
		Indexes: [][]string{{"ki1", "ks2"}},
	}
	r.create(tdef)

	// 创建第二个测试表，结构更简单
	tdef = &TableDef{
		Name:    "tbl_test2",
		Cols:    []string{"ki1", "ks2"},
		Types:   []uint32{TYPE_INT64, TYPE_BYTES},
		Indexes: [][]string{{"ki1", "ks2"}},
	}
	r.create(tdef)

	tx := r.begin()
	{
		// 检查元数据表中的下一个前缀是否正确设置
		rec := (&Record{}).AddStr("key", []byte("next_prefix"))
		ok, err := tx.Get("@meta", rec)
		assert(ok && err == nil)
		is.Equal(t, []byte{102, 0, 0, 0}, rec.Get("val").Str)
	}
	{
		// 验证表定义是否正确存储在系统表中
		rec := (&Record{}).AddStr("name", []byte("tbl_test"))
		ok, err := tx.Get("@table", rec)
		assert(ok && err == nil)
		expected := `{"Name":"tbl_test","Types":[2,1,1,2],"Cols":["ki1","ks2","s1","i2"],"Indexes":[["ki1","ks2"]],"Prefixes":[100]}`
		is.Equal(t, expected, string(rec.Get("def").Str))
	}
	r.commit(tx)

	r.dispose()
}

// TestTableBasic 测试表的基本CRUD操作
// 验证记录的添加、查询、更新和删除功能
func TestTableBasic(t *testing.T) {
	r := newR()
	tdef := &TableDef{
		Name:    "tbl_test",
		Cols:    []string{"ki1", "ks2", "s1", "i2"},
		Types:   []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES, TYPE_INT64},
		Indexes: [][]string{{"ki1", "ks2"}},
	}
	r.create(tdef)

	// 测试添加记录
	rec := Record{}
	rec.AddInt64("ki1", 1).AddStr("ks2", []byte("hello"))
	rec.AddStr("s1", []byte("world")).AddInt64("i2", 2)
	added := r.add("tbl_test", rec)
	is.True(t, added)

	// 测试查询存在的记录
	{
		got := Record{}
		got.AddInt64("ki1", 1).AddStr("ks2", []byte("hello"))
		ok := r.get("tbl_test", &got)
		is.True(t, ok)
	}

	// 测试查询不存在的记录
	{
		got := Record{}
		got.AddInt64("ki1", 1).AddStr("ks2", []byte("hello2"))
		ok := r.get("tbl_test", &got)
		is.False(t, ok)
	}

	// 测试更新记录
	rec.Get("s1").Str = []byte("www")
	added = r.add("tbl_test", rec)
	is.False(t, added) // 不是新增而是更新

	// 验证更新后的记录
	{
		got := Record{}
		got.AddInt64("ki1", 1).AddStr("ks2", []byte("hello"))
		ok := r.get("tbl_test", &got)
		is.True(t, ok)
	}

	// 测试删除记录
	{
		// 尝试删除不存在的记录
		key := Record{}
		key.AddInt64("ki1", 1).AddStr("ks2", []byte("hello2"))
		deleted := r.del("tbl_test", key)
		is.False(t, deleted)

		// 删除存在的记录
		key.Get("ks2").Str = []byte("hello")
		deleted = r.del("tbl_test", key)
		is.True(t, deleted)
	}

	r.dispose()
}

// TestStringEscape 测试字符串转义和反转义功能
// 验证包含特殊字符的字符串能够被正确编码和解码
func TestStringEscape(t *testing.T) {
	in := [][]byte{
		{},
		{0},
		{1},
	}
	out := [][]byte{
		{},
		{1, 1},
		{1, 2},
	}
	for i, s := range in {
		b := escapeString(s)
		is.Equal(t, out[i], b)
		s2 := unescapeString(b)
		is.Equal(t, s, s2)
	}
}

// TestTableEncoding 测试表值的编码和解码
// 验证编码后的值保持原始排序顺序
func TestTableEncoding(t *testing.T) {
	input := []int{-1, 0, +1, math.MinInt64, math.MaxInt64}
	sort.Ints(input)

	encoded := []string{}
	for _, i := range input {
		v := Value{Type: TYPE_INT64, I64: int64(i)}
		b := encodeValues(nil, []Value{v})
		out := []Value{v}
		decodeValues(b, out)
		assert(out[0].I64 == int64(i))
		encoded = append(encoded, string(b))
	}

	// 验证编码后的值仍然保持排序
	is.True(t, sort.StringsAreSorted(encoded))
}

// TestTableScan 测试表扫描功能
// 验证在各种条件下的范围查询和索引扫描
func TestTableScan(t *testing.T) {
	r := newR()
	tdef := &TableDef{
		Name:  "tbl_test",
		Cols:  []string{"ki1", "ks2", "s1", "i2"},
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES, TYPE_INT64},
		Indexes: [][]string{
			{"ki1", "ks2"}, // 主键索引
			{"i2"},         // 二级索引
		},
	}
	r.create(tdef)

	// 添加测试数据
	size := 100
	for i := 0; i < size; i += 2 {
		rec := Record{}
		rec.AddInt64("ki1", int64(i)).AddStr("ks2", []byte("hello"))
		rec.AddStr("s1", []byte("world")).AddInt64("i2", int64(i/2))
		added := r.add("tbl_test", rec)
		assert(added)
	}

	// 测试全表扫描（不指定键范围）
	tx := r.begin()
	{
		rec := Record{} // 空记录表示全表扫描
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)

		// 收集扫描结果并与参考数据比较
		got := []Record{}
		for req.Valid() {
			rec := Record{}
			req.Deref(&rec)
			got = append(got, rec)
			req.Next()
		}
		is.Equal(t, r.ref["tbl_test"], got)
	}
	r.commit(tx)

	// 辅助函数，创建用于扫描的键
	tmpkey := func(n int) Record {
		rec := Record{}
		rec.AddInt64("ki1", int64(n)) // 部分主键
		return rec
	}
	i2key := func(n int) Record {
		rec := Record{}
		rec.AddInt64("i2", int64(n)/2) // 二级索引
		return rec
	}

	// 测试各种范围扫描场景
	tx = r.begin()
	for i := 0; i < size; i += 2 {
		ref := []int64{}
		for j := i; j < size; j += 2 {
			ref = append(ref, int64(j))

			// 创建各种扫描条件的组合
			scanners := []Scanner{
				{
					Cmp1: CMP_GE,
					Cmp2: CMP_LE,
					Key1: tmpkey(i),
					Key2: tmpkey(j),
				},
				{
					Cmp1: CMP_GE,
					Cmp2: CMP_LE,
					Key1: tmpkey(i - 1),
					Key2: tmpkey(j + 1),
				},
				{
					Cmp1: CMP_GT,
					Cmp2: CMP_LT,
					Key1: tmpkey(i - 1),
					Key2: tmpkey(j + 1),
				},
				{
					Cmp1: CMP_GT,
					Cmp2: CMP_LT,
					Key1: tmpkey(i - 2),
					Key2: tmpkey(j + 2),
				},
				{
					Cmp1: CMP_GE,
					Cmp2: CMP_LE,
					Key1: i2key(i),
					Key2: i2key(j),
				},
				{
					Cmp1: CMP_GT,
					Cmp2: CMP_LT,
					Key1: i2key(i - 2),
					Key2: i2key(j + 2),
				},
			}

			// 为每个扫描器创建反向版本（交换比较操作和键）
			for _, tmp := range scanners {
				tmp.Cmp1, tmp.Cmp2 = tmp.Cmp2, tmp.Cmp1
				tmp.Key1, tmp.Key2 = tmp.Key2, tmp.Key1
				scanners = append(scanners, tmp)
			}

			// 测试所有扫描器
			for _, sc := range scanners {
				err := tx.Scan("tbl_test", &sc)
				assert(err == nil)

				// 收集扫描结果
				keys := []int64{}
				got := Record{}
				for sc.Valid() {
					sc.Deref(&got)
					keys = append(keys, got.Get("ki1").I64)
					sc.Next()
				}

				// 如果是反向扫描，需要反转结果进行比较
				if sc.Cmp1 < sc.Cmp2 {
					// 反转数组
					for a := 0; a < len(keys)/2; a++ {
						b := len(keys) - 1 - a
						keys[a], keys[b] = keys[b], keys[a]
					}
				}

				// 验证扫描结果与预期一致
				is.Equal(t, ref, keys)
			} // scanners
		} // j
	} // i
	r.commit(tx)

	r.dispose()
}

// TestTableIndex 测试表索引功能
// 验证多索引表的查询和更新操作
func TestTableIndex(t *testing.T) {
	r := newR()
	tdef := &TableDef{
		Name:  "tbl_test",
		Cols:  []string{"ki1", "ks2", "s1", "i2"},
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES, TYPE_INT64},
		Indexes: [][]string{
			{"ki1", "ks2"}, // 主键索引
			{"ks2", "ki1"}, // 二级索引1：反向主键
			{"i2"},         // 二级索引2：单列
			{"ki1", "i2"},  // 二级索引3：复合列
		},
	}
	r.create(tdef)

	// 辅助函数，创建测试记录
	record := func(ki1 int64, ks2 string, s1 string, i2 int64) Record {
		rec := Record{}
		rec.AddInt64("ki1", ki1).AddStr("ks2", []byte(ks2))
		rec.AddStr("s1", []byte(s1)).AddInt64("i2", i2)
		return rec
	}

	// 添加测试记录
	r1 := record(1, "a1", "v1", 2)
	r2 := record(2, "a2", "v2", -2)
	r.add("tbl_test", r1)
	r.add("tbl_test", r2)

	// 测试通过二级索引查询
	tx := r.begin()
	{
		rec := Record{}
		rec.AddInt64("i2", 2)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		is.True(t, req.Valid())

		// 验证查询结果
		out := Record{}
		req.Deref(&out)
		is.Equal(t, r1, out)

		req.Next()
		is.False(t, req.Valid())
	}
	r.commit(tx)

	// 测试范围查询（无结果）
	tx = r.begin()
	{
		rec1 := Record{}
		rec1.AddInt64("i2", 2)
		rec2 := Record{}
		rec2.AddInt64("i2", 4)
		req := Scanner{
			Cmp1: CMP_GT, Cmp2: CMP_LE,
			Key1: rec1, Key2: rec2,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		is.False(t, req.Valid())
	}
	r.commit(tx)

	// 更新记录并测试索引更新
	r.add("tbl_test", record(1, "a1", "v1", 1))
	tx = r.begin()
	{
		// 原索引值不再有效
		rec := Record{}
		rec.AddInt64("i2", 2)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		is.False(t, req.Valid())
	}
	r.commit(tx)

	// 测试新索引值
	tx = r.begin()
	{
		rec := Record{}
		rec.AddInt64("i2", 1)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		is.True(t, req.Valid())
	}
	r.commit(tx)

	// 测试删除记录后索引更新
	{
		rec := Record{}
		rec.AddInt64("ki1", 1).AddStr("ks2", []byte("a1"))
		ok := r.del("tbl_test", rec)
		assert(ok)
	}

	// 验证删除后索引已更新
	tx = r.begin()
	{
		rec := Record{}
		rec.AddInt64("i2", 1)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		is.False(t, req.Valid())
	}
	r.commit(tx)

	r.dispose()
}

// TestRecord 测试Record结构的基本功能
func TestRecord(t *testing.T) {
	r := Record{}
	r.AddStr("name", []byte("jack")).AddInt64("age", 18)
	fmt.Println(r.Get("name").Type, string(r.Get("name").Str), fmt.Sprint(r.Get("name").I64))
	for i := 0; i < len(r.Cols); i++ {
		if r.Vals[i].Type == 1 {
			fmt.Println("col:", string(r.Cols[i]), "val:", string(r.Vals[i].Str))
		} else if r.Vals[i].Type == 2 {
			fmt.Println("col:", string(r.Cols[i]), "val:", r.Vals[i].I64)
		}
	}
}

// TestScanner 测试表的范围扫描功能
// 验证使用二级索引进行范围查询的正确性
func TestScanner(t *testing.T) {
	r := newR()
	defer r.dispose() // 只需要一个 dispose

	tdef := &TableDef{
		Name:  "tbl_test",
		Cols:  []string{"ki1", "ks2", "s1", "i2"},
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_BYTES, TYPE_INT64},
		Indexes: [][]string{
			{"ki1", "ks2"}, // 主键索引
			{"i2"},         // 二级索引
		},
	}
	r.create(tdef)

	// 添加测试数据
	size := 100
	for i := 0; i < size; i += 2 {
		rec := Record{}
		rec.AddInt64("ki1", int64(i)).AddStr("ks2", []byte("hello"))
		rec.AddStr("s1", []byte("world")).AddInt64("i2", int64(i/2))
		added := r.add("tbl_test", rec)
		assert(added)
	}

	// 测试范围扫描（使用i2索引，范围5-15）
	tx := r.begin()
	{
		rec1 := Record{}
		rec1.AddInt64("i2", 5)

		rec2 := Record{}
		rec2.AddInt64("i2", 15)

		req := Scanner{
			Cmp1: CMP_GE,
			Cmp2: CMP_LE,
			Key1: rec1,
			Key2: rec2,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)

		got := []Record{}
		for req.Valid() {
			rec := Record{}
			req.Deref(&rec)
			got = append(got, rec)
			req.Next()
		}

		for j := 0; j < len(got); j++ {
			r := got[j]
			for i := 0; i < len(r.Cols); i++ {
				if r.Vals[i].Type == 1 {
					fmt.Println("col:", string(r.Cols[i]), "val:", string(r.Vals[i].Str))
				} else if r.Vals[i].Type == 2 {
					fmt.Println("col:", string(r.Cols[i]), "val:", r.Vals[i].I64)
				}
			}
		}

		// 从参考数据中筛选出应该在范围内的记录
		expected := []Record{}
		for _, rec := range r.ref["tbl_test"] {
			i2val := rec.Get("i2").I64
			if i2val >= 5 && i2val <= 15 {
				expected = append(expected, rec)
			}
		}

		// 比较扫描结果和预期结果
		is.Equal(t, expected, got)
	}
	r.commit(tx)
}

type RecordTestData struct {
	ID     int64
	Name   string
	Age    int64
	Height int64
}

func GenerateTestData(filepath string, count int) error {
	generateRandomString := func() string {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		b := make([]byte, 6)
		b[0] = letters[rand.Intn(26)+26]
		for i := 1; i < 6; i++ {
			b[i] = letters[rand.Intn(26)]
		}
		return string(b)
	}

	records := make([]RecordTestData, count)
	for i := 0; i < count; i++ {
		records[i] = RecordTestData{
			ID:     int64(i + 1),
			Name:   generateRandomString(), 
			Age:    rand.Int63n(30) + 15,   
			Height: rand.Int63n(40) + 150, 
		}
	}

	// 写入文件部分保持不变...
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("could not write data to file: %v", err)
	}

	if err := os.Chmod(filepath, 0644); err != nil {
		return fmt.Errorf("could not set file permissions: %v", err)
	}

	return nil
}

func ReadTestDataFromFile(filePath string) ([]RecordTestData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	var records []RecordTestData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return nil, fmt.Errorf("could not read data from file: %v", err)
	}

	return records, nil
}

func PrintIndexQuery(req Scanner) {
	got := []Record{}
	for req.Valid() {
		rec := Record{}
		req.Deref(&rec)
		got = append(got, rec)
		req.Next()
	}

	// 打印表头
	fmt.Printf("%-10s %-20s %-10s %-10s\n", "ID", "Name", "Age", "Height")
	fmt.Println("---------------------------------------------------------")

	// 打印每一条记录
	for _, r := range got {
		var id int64
		var name string
		var age, height int64

		for i := 0; i < len(r.Cols); i++ {
			switch string(r.Cols[i]) {
			case "id":
				if r.Vals[i].Type == 2 {
					id = r.Vals[i].I64
				}
			case "name":
				if r.Vals[i].Type == 1 {
					name = string(r.Vals[i].Str)
				}
			case "age":
				if r.Vals[i].Type == 2 {
					age = r.Vals[i].I64
				}
			case "height":
				if r.Vals[i].Type == 2 {
					height = r.Vals[i].I64
				}
			}
		}

		// 输出数据
		fmt.Printf("%-10d %-20s %-10d %-10d\n", id, name, age, height)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
func TestIndexQuery(t *testing.T) {
	// 	const (
	// 	CMP_GE = +3 // >=
	// 	CMP_GT = +2 // >
	// 	CMP_LT = -2 // <
	// 	CMP_LE = -3 // <=
	// )
	r := newR()
	defer r.dispose()

	tdef := &TableDef{
		Name:  "tbl_test",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_INT64, TYPE_INT64},
		Indexes: [][]string{
			{"id"},            // 主键索引
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	r.create(tdef)
	record := func(id int64, name string, age int64, height int64) Record {
		rec := Record{}
		rec.AddInt64("id", id).AddStr("name", []byte(name))
		rec.AddInt64("age", age).AddInt64("height", height)
		return rec
	}
	fmt.Printf("Adding test records...\n")

	// 持久化测试数据,便于观察
	filePath := "test_data.json"
	if !fileExists(filePath) {
		_ = GenerateTestData(filePath, 2000)
	}

	records, err := ReadTestDataFromFile(filePath)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	for _, rec_data := range records {
		rec := record(rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height)
		added := r.add("tbl_test", rec)
		assert(added)
	}
	fmt.Printf("Added test records successfully\n")

	// test age == 18
	fmt.Println("select * from table where age == 18")
	tx := r.begin()
	{
		rec := Record{}
		rec.AddInt64("age", 18)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		PrintIndexQuery(req)
	}
	r.commit(tx)

	fmt.Println()
	fmt.Println("select * from table where age > 43")
	tx = r.begin()
	{
		rec := Record{}
		rec.AddInt64("age", 43)

		rec2 := Record{}
		rec2.AddInt64("age", math.MaxInt64/2)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec2,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		PrintIndexQuery(req)
	}
	r.commit(tx)

	fmt.Println()
	fmt.Println("select * from table where age == 18 and height == 174")
	tx = r.begin()
	{
		rec := Record{}
		rec.AddInt64("age", 18).AddInt64("height", 174)
		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		PrintIndexQuery(req)
	}
	r.commit(tx)

	fmt.Println()
	fmt.Println("select * from table where age == 18 and height between 170 and 175 ")
	tx = r.begin()
	{
		rec := Record{}
		rec.AddInt64("age", 18).AddInt64("height", 170)

		rec2 := Record{}
		rec2.AddInt64("age", 18).AddInt64("height", 175)

		req := Scanner{
			Cmp1: CMP_GE, Cmp2: CMP_LE,
			Key1: rec, Key2: rec2,
		}
		err := tx.Scan("tbl_test", &req)
		assert(err == nil)
		PrintIndexQuery(req)
	}
	r.commit(tx)
}
