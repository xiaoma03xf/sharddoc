package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type DB struct {
	Path string
	// internals
	kv     KV
	mu     sync.Mutex           // for the cached table schema
	tables map[string]*TableDef // cached table schemas
}

// DB transaction
type DBTX struct {
	kv KVTX
	db *DB
}

func (db *DB) Begin(tx *DBTX) {
	tx.db = db
	db.kv.Begin(&tx.kv)
}
func (db *DB) Commit(tx *DBTX) error {
	return db.kv.Commit(&tx.kv)
}
func (db *DB) Abort(tx *DBTX) {
	db.kv.Abort(&tx.kv)
}

// table schema
type TableDef struct {
	// user defined
	Name    string
	Types   []uint32   // column types
	Cols    []string   // column names
	Indexes [][]string // the first index is the primary key
	// auto-assigned B-tree key prefixes for different tables and indexes
	Prefixes []uint32
}

const (
	TYPE_ERROR = 0 // uninitialized
	TYPE_BYTES = 1
	TYPE_INT64 = 2
	TYPE_INF   = 0xff // do not use
)

// table cell
type Value struct {
	Type uint32
	I64  int64
	Str  []byte
}

// table row
type Record struct {
	Cols []string
	Vals []Value
}

func init() {
	gob.Register(Record{})
	gob.Register(Value{})
}

func (rec *Record) AddStr(col string, val []byte) *Record {
	rec.Cols = append(rec.Cols, col)
	rec.Vals = append(rec.Vals, Value{Type: TYPE_BYTES, Str: val})
	return rec
}
func (rec *Record) AddInt64(col string, val int64) *Record {
	rec.Cols = append(rec.Cols, col)
	rec.Vals = append(rec.Vals, Value{Type: TYPE_INT64, I64: val})
	return rec
}

func (rec *Record) Get(key string) *Value {
	for i, c := range rec.Cols {
		if c == key {
			return &rec.Vals[i]
		}
	}
	return nil
}

// extract multiple column values
func getValues(tdef *TableDef, rec Record, cols []string) ([]Value, error) {
	vals := make([]Value, len(cols))
	for i, c := range cols {
		v := rec.Get(c)
		if v == nil {
			return nil, fmt.Errorf("missing column: %s", tdef.Cols[i])
		}
		if v.Type != tdef.Types[slices.Index(tdef.Cols, c)] {
			return nil, fmt.Errorf("bad column type: %s", c)
		}
		vals[i] = *v
	}
	return vals, nil
}

// escape the null byte so that the string contains no null byte.
func escapeString(in []byte) []byte {
	toEscape := bytes.Count(in, []byte{0}) + bytes.Count(in, []byte{1})
	if toEscape == 0 {
		return in // fast path: no escape
	}

	out := make([]byte, len(in)+toEscape)
	pos := 0
	for _, ch := range in {
		if ch <= 1 {
			// using 0x01 as the escaping byte:
			// 00 -> 01 01
			// 01 -> 01 02
			out[pos+0] = 0x01
			out[pos+1] = ch + 1
			pos += 2
		} else {
			out[pos] = ch
			pos += 1
		}
	}
	return out
}

func unescapeString(in []byte) []byte {
	if bytes.Count(in, []byte{1}) == 0 {
		return in // fast path: no unescape
	}

	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == 0x01 {
			// 01 01 -> 00
			// 01 02 -> 01
			i++
			assert(in[i] == 1 || in[i] == 2)
			out = append(out, in[i]-1)
		} else {
			out = append(out, in[i])
		}
	}
	return out
}

// order-preserving encoding
func encodeValues(out []byte, vals []Value) []byte {
	for _, v := range vals {
		out = append(out, byte(v.Type)) // 1. 类型标记
		switch v.Type {
		case TYPE_INT64: // 处理整数
			var buf [8]byte
			u := uint64(v.I64) + (1 << 63)        // 符号位翻转
			binary.BigEndian.PutUint64(buf[:], u) // 大端序存储
			out = append(out, buf[:]...)
		case TYPE_BYTES: // 处理字符串
			out = append(out, escapeString(v.Str)...) // 转义处理
			out = append(out, 0)                      // 添加终止符
		default:
			panic("what?")
		}
	}
	return out
}

// for primary keys and indexes
func encodeKey(out []byte, prefix uint32, vals []Value) []byte {
	// 1. 4字节表前缀
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], prefix)
	out = append(out, buf[:]...)

	// 2. 编码列值
	out = encodeValues(out, vals)
	return out
}

// for the input range, which can be a prefix of the index key.
func encodeKeyPartial(
	out []byte, prefix uint32, vals []Value, cmp int,
) []byte {
	out = encodeKey(out, prefix, vals)
	if cmp == CMP_GT || cmp == CMP_LE { // encode missing columns as infinity
		out = append(out, 0xff) // unreachable +infinity
	} // else: -infinity is the empty string
	return out
}

func decodeValues(in []byte, out []Value) {
	for i := range out {
		assert(out[i].Type == uint32(in[0])) // 检查出入数据类型一致
		in = in[1:]                          // 处理完一个字节，移除该字节
		switch out[i].Type {
		case TYPE_INT64:
			u := binary.BigEndian.Uint64(in[:8]) // 读取 8 个字节，转换为 uint64
			out[i].I64 = int64(u - (1 << 63))    // 通过减去最大偏移量将其转换为负数或正数
			in = in[8:]                          // 移除已经处理的字节
		case TYPE_BYTES:
			idx := bytes.IndexByte(in, 0)         // 查找字节数组中的第一个 0 字节
			assert(idx >= 0)                      // 确保找到了 0 字节
			out[i].Str = unescapeString(in[:idx]) // 解码字符串
			in = in[idx+1:]                       // 移除已处理的字节（包括 0 字节）
		default:
			panic("what?") // 如果遇到未知类型，则抛出异常
		}
	}
	assert(len(in) == 0) // 确保处理后的输入字节数组为空
}

func decodeKey(in []byte, out []Value) {
	decodeValues(in[4:], out)
}

// get a single row by the primary key
func dbGet(tx *DBTX, tdef *TableDef, rec *Record) (bool, error) {
	// 提取索引字段的值（这里只查主键）
	values, err := getValues(tdef, *rec, tdef.Indexes[0])
	if err != nil {
		return false, err // not a primary key
	}
	// just a shortcut for the scan operation
	// 	Cmp1 == GE，Cmp2 == LE，代表执行一个闭区间扫描；
	// 起止键都一样，所以就是 `查这一条` 的意思。
	sc := Scanner{
		Cmp1: CMP_GE,
		Cmp2: CMP_LE,
		Key1: Record{tdef.Indexes[0], values},
		Key2: Record{tdef.Indexes[0], values},
	}
	if err := dbScan(tx, tdef, &sc); err != nil || !sc.Valid() {
		return false, err
	}
	sc.Deref(rec)
	return true, nil
}

// internal table: metadata
var TDEF_META = &TableDef{
	Name:     "@meta",
	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},
	Cols:     []string{"key", "val"},
	Indexes:  [][]string{{"key"}},
	Prefixes: []uint32{1},
}

// internal table: table schemas
var TDEF_TABLE = &TableDef{
	Name:     "@table",
	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},
	Cols:     []string{"name", "def"},
	Indexes:  [][]string{{"name"}},
	Prefixes: []uint32{2},
}

var INTERNAL_TABLES map[string]*TableDef = map[string]*TableDef{
	"@meta":  TDEF_META,
	"@table": TDEF_TABLE,
}

// get the table schema by name, 获取表结构定义（带缓存）
func getTableDef(tx *DBTX, name string) *TableDef {
	// 1. 检查是否是内部表（如 @meta, @table）
	if tdef, ok := INTERNAL_TABLES[name]; ok {
		return tdef // expose internal tables, 直接返回内部表定义
	}
	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()

	// 3. 检查缓存中是否存在表定义
	tdef := tx.db.tables[name]
	if tdef == nil {
		// 4. 缓存未命中时，从数据库加载表定义
		if tdef = getTableDefDB(tx, name); tdef != nil {
			tx.db.tables[name] = tdef
		}
	}
	return tdef
}

func getTableDefDB(tx *DBTX, name string) *TableDef {
	// 1. 构造查询记录（主键为表名）
	rec := (&Record{}).AddStr("name", []byte(name))

	// 2. 从系统表 @table 中查询表定义
	ok, err := dbGet(tx, TDEF_TABLE, rec)
	assert(err == nil)
	if !ok {
		return nil
	}

	// 3. 反序列化 JSON 数据到 TableDef
	tdef := &TableDef{}
	err = json.Unmarshal(rec.Get("def").Str, tdef)
	assert(err == nil) // 确保反序列化成功
	return tdef
}

// get a single row by the primary key
// 获取表结构分两步, 如果要查询的是系统表则返回, 如果缓存中没有就用name主键去系统表中查询
// 然后把查询的结果放入缓存中, 然后通过主键查询记录
func (tx *DBTX) Get(table string, rec *Record) (bool, error) {
	// 1. 获取表结构定义（带缓存）
	tdef := getTableDef(tx, table)
	if tdef == nil {
		return false, fmt.Errorf("table not found: %s", table)
	}
	// 2. 通过主键查询记录
	return dbGet(tx, tdef, rec)
}

const TABLE_PREFIX_MIN = 100

func tableDefCheck(tdef *TableDef) error {
	// verify the table schema
	bad := tdef.Name == "" || len(tdef.Cols) == 0 || len(tdef.Indexes) == 0
	bad = bad || len(tdef.Cols) != len(tdef.Types)
	if bad {
		return fmt.Errorf("bad table schema: %s", tdef.Name)
	}
	// verify the indexes, 依次校验每一个索引是否有问题
	for i, index := range tdef.Indexes {
		index, err := checkIndexCols(tdef, index)
		if err != nil {
			return err
		}
		tdef.Indexes[i] = index
	}

	// for _, col := range tdef.Cols {
	// 	if strings.ToLower(col) == "id" {
	// 		return nil
	// 	}
	// }
	// return fmt.Errorf("the id field is required")
	return nil
}

func checkIndexCols(tdef *TableDef, index []string) ([]string, error) {
	if len(index) == 0 {
		return nil, fmt.Errorf("empty index")
	}
	seen := map[string]bool{}

	// 对索引的字段检查, 是否字段存在, 是否重复
	for _, c := range index {
		// check the index columns
		if slices.Index(tdef.Cols, c) < 0 {
			return nil, fmt.Errorf("unknown index column: %s", c)
		}
		if seen[c] {
			return nil, fmt.Errorf("duplicated column in index: %s", c)
		}
		seen[c] = true
	}
	// add the primary key to the index
	// 确保每个索引都包含了主键字段
	for _, c := range tdef.Indexes[0] {
		if !seen[c] {
			index = append(index, c)
		}
	}
	assert(len(index) <= len(tdef.Cols))
	return index, nil
}

func (tx *DBTX) TableNew(tdef *TableDef) error {
	// 0. sanity checks, 校验表结构,检查表定义是否合法
	if err := tableDefCheck(tdef); err != nil {
		return err
	}
	// 1. check the existing table
	table := (&Record{}).AddStr("name", []byte(tdef.Name))
	ok, err := dbGet(tx, TDEF_TABLE, table)
	assert(err == nil)
	if ok {
		return fmt.Errorf("table exists: %s", tdef.Name)
	}
	// 2. allocate new prefixes, 分配前缀
	//TDEF_META 是数据库内部的`元数据表`，用来保存系统内部的配置、状态等，比如前缀分配器
	// 去TDEF_META表中查询"key","next_prefix"字段
	prefix := uint32(TABLE_PREFIX_MIN)
	meta := (&Record{}).AddStr("key", []byte("next_prefix"))
	ok, err = dbGet(tx, TDEF_META, meta)
	assert(err == nil)
	if ok {
		prefix = binary.LittleEndian.Uint32(meta.Get("val").Str)
		assert(prefix > TABLE_PREFIX_MIN)
	} else {
		// 处理首次分配
		meta.AddStr("val", make([]byte, 4))
	}
	assert(len(tdef.Prefixes) == 0)
	for i := range tdef.Indexes {
		// 为索引分配前缀
		tdef.Prefixes = append(tdef.Prefixes, prefix+uint32(i))
	}
	// 3. update the next prefix
	// FIXME: integer overflow.跟新下一个前缀并写入
	next := prefix + uint32(len(tdef.Indexes))
	binary.LittleEndian.PutUint32(meta.Get("val").Str, next)
	_, err = dbUpdate(tx, TDEF_META, &DBUpdateReq{Record: *meta})
	if err != nil {
		return err
	}

	// 4. store the schema, 持久化表结构
	val, err := json.Marshal(tdef)
	assert(err == nil)
	table.AddStr("def", val)
	_, err = dbUpdate(tx, TDEF_TABLE, &DBUpdateReq{Record: *table})
	return err
}

type DBUpdateReq struct {
	// in
	Record Record
	Mode   int
	// out
	Updated bool
	Added   bool
}

// add a row to the table, 更新一行的数据
// 1.先把主键和非主键拼接，也把它们对应的值拼接起来
// 2. 获取主键的长度， 主键编码需要和prefix配合， 然后编码值
// 3. 删除原索引, 创建新索引
func dbUpdate(tx *DBTX, tdef *TableDef, dbreq *DBUpdateReq) (bool, error) {
	// reorder the columns so that they start with the primary key
	// cols 把主键和非主键拼接起来, 然后取出当前数据, values 通过 依次遍历cols中的键位获得值
	cols := slices.Concat(tdef.Indexes[0], nonPrimaryKeyCols(tdef))
	values, err := getValues(tdef, dbreq.Record, cols)
	if err != nil {
		return false, err // expect a full row
	}

	// insert the row
	// npk 表示主键长度, 先获取主键的长度(主键可能是多个key也可能是单key), 然后把主键和prefix编码
	npk := len(tdef.Indexes[0]) // number of primary key columns
	key := encodeKey(nil, tdef.Prefixes[0], values[:npk])
	val := encodeValues(nil, values[npk:])
	req := UpdateReq{Key: key, Val: val, Mode: dbreq.Mode}
	if _, err := tx.kv.Update(&req); err != nil {
		return false, err // length limit
	}
	dbreq.Added, dbreq.Updated = req.Added, req.Updated

	// maintain secondary indexes
	// 如果这个记录是被更新了而不是添加
	if req.Updated && !req.Added {
		// construct the old record
		decodeValues(req.Old, values[npk:])
		oldRec := Record{cols, values}
		// delete the indexed keys, 删除索引字段
		err := indexOp(tx, tdef, INDEX_DEL, oldRec)
		assert(err == nil) // should not run into the length limit
	}
	if req.Updated {
		// add the new indexed keys
		if err := indexOp(tx, tdef, INDEX_ADD, dbreq.Record); err != nil {
			return false, err // length limit
		}
	}
	return req.Updated, nil
}

func nonPrimaryKeyCols(tdef *TableDef) (out []string) {
	for _, c := range tdef.Cols {
		// c 这个字段名不在主键字段列表（tdef.Indexes[0]）中，也就是`非主键字段`
		if slices.Index(tdef.Indexes[0], c) < 0 {
			out = append(out, c)
		}
	}
	return
}

const (
	INDEX_ADD = 1
	INDEX_DEL = 2
)

// add or remove secondary index keys
func indexOp(tx *DBTX, tdef *TableDef, op int, rec Record) error {
	for i := 1; i < len(tdef.Indexes); i++ {
		// the indexed key
		values, err := getValues(tdef, rec, tdef.Indexes[i])
		assert(err == nil) // full row
		key := encodeKey(nil, tdef.Prefixes[i], values)
		switch op {
		case INDEX_ADD:
			req := UpdateReq{Key: key, Val: nil}
			if _, err := tx.kv.Update(&req); err != nil {
				return err // length limit
			}
			assert(req.Added) // internal consistency
		case INDEX_DEL:
			deleted, err := tx.kv.Del(&DeleteReq{Key: key})
			assert(err == nil) // should not run into the length limit
			assert(deleted)    // internal consistency
		default:
			panic("unreachable")
		}
	}
	return nil
}

// add a record
func (tx *DBTX) Set(table string, dbreq *DBUpdateReq) (bool, error) {
	tdef := getTableDef(tx, table)
	if tdef == nil {
		return false, fmt.Errorf("table not found: %s", table)
	}
	return dbUpdate(tx, tdef, dbreq)
}
func (tx *DBTX) Insert(table string, rec Record) (bool, error) {
	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_INSERT_ONLY})
}
func (tx *DBTX) Update(table string, rec Record) (bool, error) {
	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_UPDATE_ONLY})
}
func (tx *DBTX) Upsert(table string, rec Record) (bool, error) {
	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_UPSERT})
}

// delete a record by its primary key
func dbDelete(tx *DBTX, tdef *TableDef, rec Record) (bool, error) {
	// 获取列值
	values, err := getValues(tdef, rec, tdef.Indexes[0])
	if err != nil {
		return false, err
	}
	// delete the row
	req := DeleteReq{Key: encodeKey(nil, tdef.Prefixes[0], values)}
	if deleted, _ := tx.kv.Del(&req); !deleted {
		return false, nil // `deleted` is also false if the key is too long
	}
	// maintain secondary indexes
	for _, c := range nonPrimaryKeyCols(tdef) {
		tp := tdef.Types[slices.Index(tdef.Cols, c)]
		values = append(values, Value{Type: tp})
	}
	decodeValues(req.Old, values[len(tdef.Indexes[0]):])
	err = indexOp(tx, tdef, INDEX_DEL, Record{tdef.Cols, values})
	assert(err == nil) // should not run into the length limit
	return true, nil
}

func (tx *DBTX) Delete(table string, rec Record) (bool, error) {
	tdef := getTableDef(tx, table)
	if tdef == nil {
		return false, fmt.Errorf("table not found: %s", table)
	}
	return dbDelete(tx, tdef, rec)
}

// 初始化数据库实例, 打开底层KV存储
func (db *DB) Open() error {
	db.kv.Path = db.Path
	db.tables = map[string]*TableDef{} // 初始化表缓存
	return db.kv.Open()                // 打开KV存储文件
}

func (db *DB) Close() {
	db.kv.Close()
}

// the iterator for range queries
type Scanner struct {
	// the range, from Key1 to Key2
	Cmp1 int // CMP_??
	Cmp2 int
	Key1 Record
	Key2 Record
	// internal
	tx    *DBTX
	tdef  *TableDef
	index int    // which index?
	iter  KVIter // the underlying KV iterator
}

// within the range or not?
func (sc *Scanner) Valid() bool {
	return sc.iter.Valid()
}

// move the underlying B-tree iterator
func (sc *Scanner) Next() {
	sc.iter.Next()
}

// return the current row
func (sc *Scanner) Deref(rec *Record) {
	assert(sc.Valid())
	tdef := sc.tdef // 获取当前表定义

	// prepare the output record
	// 拼接主键和非主键列，形成完整的列顺序
	// rec.Cols 把主键和非主键拼接起来, tdef.Indexes[0] 表示主键, nonPrimaryKeyCols非主键切片
	rec.Cols = slices.Concat(tdef.Indexes[0], nonPrimaryKeyCols(tdef))

	// 把rec.Vals的值清空, 准备存放新的值
	rec.Vals = rec.Vals[:0]
	for _, c := range rec.Cols {
		tp := tdef.Types[slices.Index(tdef.Cols, c)]
		rec.Vals = append(rec.Vals, Value{Type: tp})
	}

	// fetch the KV from the iterator
	// 从迭代器中获取当前键值对
	key, val := sc.iter.Deref()
	// primary key or secondary index?
	if sc.index == 0 {
		// decode the full row
		npk := len(tdef.Indexes[0])
		decodeKey(key, rec.Vals[:npk])
		decodeValues(val, rec.Vals[npk:])
	} else {
		// decode the index key
		// 解码二级索引的键
		assert(len(val) == 0)
		index := tdef.Indexes[sc.index]
		irec := Record{index, make([]Value, len(index))}
		for i, c := range index {
			irec.Vals[i].Type = tdef.Types[slices.Index(tdef.Cols, c)]
		}
		decodeKey(key, irec.Vals)

		// extract the primary key
		// 从索引键中提取主键值
		for i, c := range tdef.Indexes[0] {
			rec.Vals[i] = *irec.Get(c)
		}
		// fetch the row by the primary key
		// TODO: skip this if the index contains all the columns
		ok, err := dbGet(sc.tx, tdef, rec)
		assert(ok && err == nil) // internal consistency
	}
}

// check column types
func checkTypes(tdef *TableDef, rec Record) error {
	if len(rec.Cols) != len(rec.Vals) {
		return fmt.Errorf("bad record")
	}
	for i, c := range rec.Cols {
		j := slices.Index(tdef.Cols, c)
		if j < 0 || tdef.Types[j] != rec.Vals[i].Type {
			return fmt.Errorf("bad column: %s", c)
		}
	}
	return nil
}

func dbScan(tx *DBTX, tdef *TableDef, req *Scanner) error {
	// 0. sanity checks
	switch {
	case req.Cmp1 > 0 && req.Cmp2 < 0:
	case req.Cmp2 > 0 && req.Cmp1 < 0:
	default:
		return fmt.Errorf("bad range")
	}
	if !slices.Equal(req.Key1.Cols, req.Key2.Cols) {
		return fmt.Errorf("bad range key")
	}
	if err := checkTypes(tdef, req.Key1); err != nil {
		return err
	}
	if err := checkTypes(tdef, req.Key2); err != nil {
		return err
	}
	req.tx = tx
	req.tdef = tdef
	// 1. select the index
	isCovered := func(index []string) bool {
		key := req.Key1.Cols
		return len(index) >= len(key) && slices.Equal(index[:len(key)], key)
	}
	req.index = slices.IndexFunc(tdef.Indexes, isCovered)
	if req.index < 0 {
		return fmt.Errorf("no index")
	}
	// 2. encode the start key and the end key
	prefix := tdef.Prefixes[req.index]
	keyStart := encodeKeyPartial(nil, prefix, req.Key1.Vals, req.Cmp1)
	keyEnd := encodeKeyPartial(nil, prefix, req.Key2.Vals, req.Cmp2)
	// 3. seek to the start key
	req.iter = tx.kv.Seek(keyStart, req.Cmp1, keyEnd, req.Cmp2)
	return nil
}

func (tx *DBTX) Scan(table string, req *Scanner) error {
	tdef := getTableDef(tx, table)
	if tdef == nil {
		return fmt.Errorf("table not found: %s", table)
	}
	return dbScan(tx, tdef, req)
}

func (db *DB) GetAllTables() (t []TableDef, e error) {
	tx := DBTX{}
	db.Begin(&tx)
	sc := &Scanner{
		Cmp1: CMP_GE,
		Cmp2: CMP_LE,
		Key1: *(&Record{}).AddStr("name", []byte(MIN_NAME)),
		Key2: *(&Record{}).AddStr("name", []byte(MAX_NAME)),
	}
	if err := tx.Scan("@table", sc); err != nil {
		return nil, fmt.Errorf("query table info error")
	}
	err := db.Commit(&tx)
	if err != nil {
		return nil, err
	}
	// json unmarshal
	rec := reduceSelectData(sc)
	for _, r := range rec {
		tdef := &TableDef{}
		err = json.Unmarshal(r.Get("def").Str, tdef)
		assert(err == nil)
		t = append(t, *tdef)
	}
	return
}

// 获取当前db 所有表结构, 主要用于创建db逻辑快照
// snapshotDir, 假如数据库名为 r.db 其中两张表 tbl_test1,tbl_test2
// 会生成 r_export 文件目录, 下面包含有schema.json存有数据库表信息
// tbl_test1.data, tbl_test2.data 分别把每个表的数据 Record{}二进制编码
func (db *DB) ExportDB() (snapshotDir string, err error) {
	if strings.HasSuffix(db.Path, ".db") {
		snapshotDir = db.Path[:len(db.Path)-3] + "_export"
	} else {
		snapshotDir = db.Path + "_export"
	}
	// 1. 创建快照目录
	if err = os.MkdirAll(snapshotDir, 0755); err != nil {
		return
	}
	// 2. 获取所有的表结构,表结构存入json
	allTables, err := db.GetAllTables()
	if err != nil {
		return "", err
	}
	// 写入表结构 snapshotDir目录下的schema.json
	file, _ := os.Create(filepath.Join(snapshotDir, "schema.json"))
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err = encoder.Encode(allTables); err != nil {
		return "", err
	}
	// 把每张表的数据写入对应文件
	for _, table := range allTables {
		datapath := filepath.Join(snapshotDir, table.Name+".data")
		dataFile, err := os.Create(datapath)
		if err != nil {
			return "", err
		}
		// 查询当前表的所有数据
		tx := DBTX{}
		db.Begin(&tx)
		sc := &Scanner{
			Cmp1: CMP_GE,
			Cmp2: CMP_LE,
			Key1: *(&Record{}).AddInt64("id", math.MinInt64/2),
			Key2: *(&Record{}).AddInt64("id", math.MaxInt64/2),
		}
		if err := tx.Scan(table.Name, sc); err != nil {
			return snapshotDir, err
		}
		if err = db.Commit(&tx); err != nil {
			return "", err
		}
		recs := reduceSelectData(sc)
		enc := gob.NewEncoder(dataFile)
		for _, rec := range recs {
			if err := enc.Encode(rec); err != nil {
				_ = dataFile.Close()
				return "", err
			}
		}
		_ = dataFile.Close()
	}
	return
}

// LoadRecordsFromDataFile 从gob编码中读取数据库数据
func LoadRecordsFromDataFile(dataPath string) ([]Record, error) {
	dataFile, _ := os.Open(dataPath)
	defer dataFile.Close()
	dec := gob.NewDecoder(dataFile)

	var records []Record
	for {
		var rec Record
		err := dec.Decode(&rec)
		if err == io.EOF {
			break // 读取完毕
		}
		if err != nil {
			return nil, err // 中途出错
		}
		records = append(records, rec)
	}
	return records, nil
}

// ImportDB. dbDir 表示要恢复的文件夹
func ImportDB(dbDir, dbname string) (*DB, error) {
	r := &DB{Path: dbname}
	if err := r.Open(); err != nil {
		return nil, err
	}
	// 读取表结构,并创建表
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("snapshot directory %s does not exist", dbDir)
	}

	schemaPath := filepath.Join(dbDir, "schema.json")
	jsonBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("faild to read schema.json: %v", err)
	}
	// 读取表结构, 创建表
	var defs []TableDef
	if err := json.Unmarshal(jsonBytes, &defs); err != nil {
		return nil, fmt.Errorf("faild to unmarshal schema.json:%v", err)
	}

	var wg sync.WaitGroup
	tx := &DBTX{}
	r.Begin(tx)
	var _err error
	for _, def := range defs {
		wg.Add(1)
		go func(def TableDef) {
			defer wg.Done()
			// 将def的prefix清空
			def.Prefixes = []uint32{}
			if err := tx.TableNew(&def); err != nil {
				_err = fmt.Errorf("failed to create table: %v", err)
			}
			datapath := filepath.Join(dbDir, def.Name+".data")
			recs, err := LoadRecordsFromDataFile(datapath)
			if err != nil {
				_err = fmt.Errorf("load recs err:%v", err)
			}
			for _, rec := range recs {
				if _, err := tx.Insert(def.Name, rec); err != nil {
					_err = err
				}
			}
		}(def)
	}
	wg.Wait()
	if _err != nil {
		r.Abort(tx)
		return nil, _err
	}
	// 提交事务
	if err := r.Commit(tx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	return r, nil
}
