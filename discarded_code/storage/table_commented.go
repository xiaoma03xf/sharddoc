package storage

// /**
//  * table_commented.go - 基于KV存储构建关系型数据库表的核心实现（详细注释版本）
//  *
//  * 这个文件实现了一个完整的关系型数据库存储层，基于底层的Key-Value存储构建
//  * 主要功能包括：
//  * 1. 表结构定义和管理
//  * 2. 行记录的编码/解码
//  * 3. 主键和二级索引的管理
//  * 4. 事务性的增删改查操作
//  * 5. 范围查询和扫描
//  * 6. 数据库快照导入导出
//  */
// package storage

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"encoding/gob"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"math"
// 	"os"
// 	"path/filepath"
// 	"slices"
// 	"strings"
// 	"sync"
// )

// // ================================
// // 核心数据结构定义
// // ================================

// // DB 数据库实例，封装底层KV存储并提供关系型数据库接口
// type DB struct {
// 	Path string                    // 数据库文件路径
// 	// 内部状态
// 	kv     KV                     // 底层KV存储引擎（B+树实现）
// 	mu     sync.Mutex             // 保护表结构缓存的互斥锁
// 	tables map[string]*TableDef   // 表结构缓存，避免重复读取磁盘
// }

// // DBTX 数据库事务，封装KV事务并提供表级别的操作接口
// type DBTX struct {
// 	kv KVTX  // 底层KV事务
// 	db *DB   // 关联的数据库实例
// }

// // ================================
// // 事务管理方法 - 数据库事务的生命周期管理
// // ================================

// // Begin 开始一个新的数据库事务
// // 初始化事务上下文，为后续的数据操作做准备
// func (db *DB) Begin(tx *DBTX) {
// 	tx.db = db
// 	db.kv.Begin(&tx.kv)
// }

// // Commit 提交事务，使所有修改持久化到磁盘
// // 这是ACID中的持久性（Durability）保证
// func (db *DB) Commit(tx *DBTX) error {
// 	return db.kv.Commit(&tx.kv)
// }

// // Abort 回滚事务，撤销所有修改
// // 这是ACID中的原子性（Atomicity）保证
// func (db *DB) Abort(tx *DBTX) {
// 	db.kv.Abort(&tx.kv)
// }

// // ================================
// // 表结构定义 - 关系型数据库的表模式
// // ================================

// // TableDef 表结构定义，描述一个关系型表的元数据
// // 这是整个系统的核心数据结构，定义了表的逻辑结构
// type TableDef struct {
// 	// 用户定义的表信息
// 	Name    string     // 表名，如 "users", "orders"
// 	Types   []uint32   // 列类型数组，与Cols一一对应，如 [TYPE_INT64, TYPE_BYTES, TYPE_INT64]
// 	Cols    []string   // 列名数组，如 ["id", "name", "age"]
// 	Indexes [][]string // 索引定义，第一个必须是主键，后续为二级索引
// 	                   // 如 [["id"], ["name"], ["age", "id"]]

// 	// 系统自动分配的B+树键前缀
// 	// 每个表的每个索引都有唯一的前缀，确保在底层KV存储中不会冲突
// 	// 例如：表users的主键前缀=100，name索引前缀=101，age索引前缀=102
// 	Prefixes []uint32
// }

// // ================================
// // 数据类型定义 - 支持的列数据类型
// // ================================

// const (
// 	TYPE_ERROR = 0    // 未初始化类型，表示错误状态
// 	TYPE_BYTES = 1    // 字节串类型（用于存储字符串）
// 	TYPE_INT64 = 2    // 64位有符号整数类型
// 	TYPE_INF   = 0xff // 无穷大标记，用于范围查询的边界处理
// )

// // Value 表示表中的一个值，支持整数和字节串两种基本类型
// // 这是数据库系统中最小的数据单元
// type Value struct {
// 	Type uint32  // 值的类型标识符
// 	I64  int64   // 整数值（当Type为TYPE_INT64时使用）
// 	Str  []byte  // 字节串值（当Type为TYPE_BYTES时使用）
// }

// // Record 表示一行记录，是关系型数据库中的基本数据单位
// // 例如用户表中的一行：id=1, name="jack", age=18, height=180
// type Record struct {
// 	Cols []string  // 列名数组，如 ["id", "name", "age", "height"]
// 	Vals []Value   // 对应的值数组，与Cols一一对应
// }

// // gob序列化注册，用于数据导入导出功能
// // gob是Go语言的二进制编码格式，用于高效的数据序列化
// func init() {
// 	gob.Register(Record{})
// 	gob.Register(Value{})
// }

// // ================================
// // Record 操作方法 - 记录构建和访问
// // ================================

// // AddStr 向记录添加一个字符串类型的列
// // 这是构建记录的便捷方法，用于链式调用
// // 例如：rec.AddStr("name", []byte("jack")).AddStr("email", []byte("jack@example.com"))
// func (rec *Record) AddStr(col string, val []byte) *Record {
// 	rec.Cols = append(rec.Cols, col)
// 	rec.Vals = append(rec.Vals, Value{Type: TYPE_BYTES, Str: val})
// 	return rec
// }

// // AddInt64 向记录添加一个整数类型的列
// // 例如：rec.AddInt64("age", 18).AddInt64("score", 95)
// func (rec *Record) AddInt64(col string, val int64) *Record {
// 	rec.Cols = append(rec.Cols, col)
// 	rec.Vals = append(rec.Vals, Value{Type: TYPE_INT64, I64: val})
// 	return rec
// }

// // Get 根据列名获取记录中的值
// // 这是记录的核心访问方法，类似于SQL中的列引用
// func (rec *Record) Get(key string) *Value {
// 	for i, c := range rec.Cols {
// 		if c == key {
// 			return &rec.Vals[i]
// 		}
// 	}
// 	return nil // 列不存在时返回nil
// }

// // ================================
// // 数据提取和校验 - 类型安全保证
// // ================================

// // getValues 从记录中提取指定列的值，并校验类型是否匹配表定义
// // 这是一个关键的类型安全函数，确保运行时的数据类型正确性
// // tdef: 表定义，包含列类型信息
// // rec: 输入记录
// // cols: 要提取的列名列表
// func getValues(tdef *TableDef, rec Record, cols []string) ([]Value, error) {
// 	vals := make([]Value, len(cols))
// 	for i, c := range cols {
// 		v := rec.Get(c)
// 		if v == nil {
// 			return nil, fmt.Errorf("missing column: %s", tdef.Cols[i])
// 		}
// 		// 关键：校验实际类型与表定义是否匹配
// 		expectedType := tdef.Types[slices.Index(tdef.Cols, c)]
// 		if v.Type != expectedType {
// 			return nil, fmt.Errorf("bad column type: %s", c)
// 		}
// 		vals[i] = *v
// 	}
// 	return vals, nil
// }

// // ================================
// // 字符串转义处理 - 键值编码的基础
// // ================================

// // escapeString 转义字符串中的控制字节，确保编码后的键值不包含分隔符
// // 这对于键值编码至关重要，因为空字节(0x00)用作字段分隔符
// // 转义规则：
// // 0x00 -> 0x01 0x01 (空字节转为两字节序列)
// // 0x01 -> 0x01 0x02 (转义字节本身也需要转义)
// func escapeString(in []byte) []byte {
// 	toEscape := bytes.Count(in, []byte{0}) + bytes.Count(in, []byte{1})
// 	if toEscape == 0 {
// 		return in // 快速路径：无需转义时直接返回
// 	}

// 	// 分配足够的空间存储转义后的结果
// 	out := make([]byte, len(in)+toEscape)
// 	pos := 0
// 	for _, ch := range in {
// 		if ch <= 1 {
// 			// 使用0x01作为转义前缀
// 			out[pos+0] = 0x01
// 			out[pos+1] = ch + 1
// 			pos += 2
// 		} else {
// 			out[pos] = ch
// 			pos += 1
// 		}
// 	}
// 	return out
// }

// // unescapeString 反转义字符串，恢复原始内容
// // 这是escapeString的逆操作
// func unescapeString(in []byte) []byte {
// 	if bytes.Count(in, []byte{1}) == 0 {
// 		return in // 快速路径：无转义字符时直接返回
// 	}

// 	out := make([]byte, 0, len(in))
// 	for i := 0; i < len(in); i++ {
// 		if in[i] == 0x01 {
// 			// 解析转义序列：0x01 0x01 -> 0x00, 0x01 0x02 -> 0x01
// 			i++
// 			assert(in[i] == 1 || in[i] == 2)
// 			out = append(out, in[i]-1)
// 		} else {
// 			out = append(out, in[i])
// 		}
// 	}
// 	return out
// }

// // ================================
// // 键值编码/解码 - 核心算法，保序编码的实现
// // ================================

// // encodeValues 对值序列进行保序编码
// // 保序编码是关系型数据库索引的核心技术，确保：
// // 编码后的字节序列的字典序 == 原值的逻辑序
// // 这样B+树就可以正确地按逻辑序进行范围查询
// func encodeValues(out []byte, vals []Value) []byte {
// 	for _, v := range vals {
// 		// 1. 写入类型标记（1字节）
// 		out = append(out, byte(v.Type))

// 		switch v.Type {
// 		case TYPE_INT64: // 处理64位整数
// 			var buf [8]byte
// 			// 关键技术：符号位翻转 + 大端序
// 			// 目标：使 -∞ < ... < -1 < 0 < 1 < ... < +∞ 的字典序成立
// 			// 方法：加上2^63，将有符号数映射到无符号数
// 			// 原理：-2^63 变成 0，-1 变成 2^63-1，0 变成 2^63，2^63-1 变成 2^64-1
// 			u := uint64(v.I64) + (1 << 63)
// 			binary.BigEndian.PutUint64(buf[:], u) // 大端序：高位在前，保证字典序
// 			out = append(out, buf[:]...)

// 		case TYPE_BYTES: // 处理字节串（字符串）
// 			out = append(out, escapeString(v.Str)...) // 转义处理
// 			out = append(out, 0)                      // 添加终止符分隔不同字段

// 		default:
// 			panic("unsupported type in encodeValues")
// 		}
// 	}
// 	return out
// }

// // encodeKey 编码完整的索引键
// // 格式：[4字节表/索引前缀][编码后的列值序列]
// // 前缀的作用：在同一个B+树中区分不同表和不同索引的数据
// func encodeKey(out []byte, prefix uint32, vals []Value) []byte {
// 	// 1. 写入4字节前缀，用于命名空间隔离
// 	var buf [4]byte
// 	binary.BigEndian.PutUint32(buf[:], prefix)
// 	out = append(out, buf[:]...)

// 	// 2. 编码列值序列
// 	out = encodeValues(out, vals)
// 	return out
// }

// // encodeKeyPartial 编码部分键，专用于范围查询
// // 当查询键是索引键的前缀时，需要处理"缺失"列的边界情况
// // cmp参数决定缺失列的处理方式：
// // - CMP_GT/CMP_LE: 缺失列视为+∞（在范围的上边界）
// // - CMP_GE/CMP_LT: 缺失列视为-∞（在范围的下边界）
// func encodeKeyPartial(
// 	out []byte, prefix uint32, vals []Value, cmp int,
// ) []byte {
// 	out = encodeKey(out, prefix, vals)
// 	if cmp == CMP_GT || cmp == CMP_LE {
// 		out = append(out, 0xff) // 添加+∞标记（0xFF是最大的字节值）
// 	}
// 	// -∞用空字符串表示（什么都不添加）
// 	return out
// }

// // decodeValues 解码值序列，与encodeValues完全相反
// // out参数必须预先分配好正确的类型信息
// func decodeValues(in []byte, out []Value) {
// 	for i := range out {
// 		// 验证类型一致性
// 		assert(out[i].Type == uint32(in[0]))
// 		in = in[1:] // 跳过类型标记

// 		switch out[i].Type {
// 		case TYPE_INT64:
// 			// 逆转编码过程：大端序解码 + 符号位还原
// 			u := binary.BigEndian.Uint64(in[:8])
// 			out[i].I64 = int64(u - (1 << 63)) // 减去偏移量还原符号
// 			in = in[8:]

// 		case TYPE_BYTES:
// 			// 查找字符串终止符
// 			idx := bytes.IndexByte(in, 0)
// 			assert(idx >= 0) // 必须能找到终止符
// 			out[i].Str = unescapeString(in[:idx]) // 反转义
// 			in = in[idx+1:] // 跳过终止符

// 		default:
// 			panic("unsupported type in decodeValues")
// 		}
// 	}
// 	assert(len(in) == 0) // 确保所有数据都被正确解码
// }

// // decodeKey 解码索引键，跳过前缀部分
// func decodeKey(in []byte, out []Value) {
// 	decodeValues(in[4:], out) // 跳过4字节前缀，解码值部分
// }

// // ================================
// // 单行查询操作 - 主键查找的实现
// // ================================

// // dbGet 通过主键获取单行记录
// // 这是最基础的数据库操作，所有其他操作都建立在此基础上
// func dbGet(tx *DBTX, tdef *TableDef, rec *Record) (bool, error) {
// 	// 1. 从记录中提取主键字段的值
// 	// tdef.Indexes[0] 始终是主键索引
// 	values, err := getValues(tdef, *rec, tdef.Indexes[0])
// 	if err != nil {
// 		return false, err // 主键字段不完整或类型错误
// 	}

// 	// 2. 构建精确查找的扫描器
// 	// 通过将起始键和结束键设为相同值，实现点查询
// 	sc := Scanner{
// 		Cmp1: CMP_GE, // 大于等于起始键
// 		Cmp2: CMP_LE, // 小于等于结束键
// 		Key1: Record{tdef.Indexes[0], values}, // 起始键
// 		Key2: Record{tdef.Indexes[0], values}, // 结束键（相同）
// 	}

// 	// 3. 执行扫描
// 	if err := dbScan(tx, tdef, &sc); err != nil || !sc.Valid() {
// 		return false, err // 查询失败或记录不存在
// 	}

// 	// 4. 提取查询结果
// 	sc.Deref(rec)
// 	return true, nil
// }

// // ================================
// // 系统表定义 - 数据库元数据管理
// // ================================

// // TDEF_META 系统元数据表，存储数据库级别的配置信息
// // 主要用途：管理表前缀的分配，确保每个表和索引都有唯一的前缀
// var TDEF_META = &TableDef{
// 	Name:     "@meta",                              // 表名以@开头表示系统表
// 	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},     // key-value结构
// 	Cols:     []string{"key", "val"},               // 列名
// 	Indexes:  [][]string{{"key"}},                  // 主键是key列
// 	Prefixes: []uint32{1},                         // 系统分配的固定前缀
// }

// // TDEF_TABLE 系统表结构表，存储所有用户表的定义
// // 这是数据库的"数据字典"，记录了所有表的元数据
// var TDEF_TABLE = &TableDef{
// 	Name:     "@table",
// 	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},     // name -> definition(JSON)
// 	Cols:     []string{"name", "def"},
// 	Indexes:  [][]string{{"name"}},                 // 主键是表名
// 	Prefixes: []uint32{2},                         // 系统分配的固定前缀
// }

// // 内部系统表注册表
// var INTERNAL_TABLES map[string]*TableDef = map[string]*TableDef{
// 	"@meta":  TDEF_META,
// 	"@table": TDEF_TABLE,
// }

// // ================================
// // 表结构管理 - 表定义的缓存和加载
// // ================================

// // getTableDef 获取表结构定义，实现三级查找策略：
// // 1. 内部系统表（直接返回）
// // 2. 内存缓存（快速访问）
// // 3. 数据库存储（磁盘读取）
// func getTableDef(tx *DBTX, name string) *TableDef {
// 	// 第1级：检查是否是内部系统表
// 	if tdef, ok := INTERNAL_TABLES[name]; ok {
// 		return tdef // 系统表定义是硬编码的，直接返回
// 	}

// 	// 第2级和第3级需要加锁保护缓存
// 	tx.db.mu.Lock()
// 	defer tx.db.mu.Unlock()

// 	// 第2级：检查内存缓存
// 	tdef := tx.db.tables[name]
// 	if tdef == nil {
// 		// 第3级：从数据库加载并更新缓存
// 		if tdef = getTableDefDB(tx, name); tdef != nil {
// 			tx.db.tables[name] = tdef // 写入缓存
// 		}
// 	}
// 	return tdef
// }

// // getTableDefDB 从数据库存储中加载表结构定义
// // 这个函数访问@table系统表，获取表的JSON定义并反序列化
// func getTableDefDB(tx *DBTX, name string) *TableDef {
// 	// 1. 构造查询记录，查找指定表名的定义
// 	rec := (&Record{}).AddStr("name", []byte(name))

// 	// 2. 从@table系统表中查询
// 	ok, err := dbGet(tx, TDEF_TABLE, rec)
// 	assert(err == nil) // 系统表查询不应该失败
// 	if !ok {
// 		return nil // 表不存在
// 	}

// 	// 3. 反序列化JSON格式的表定义
// 	tdef := &TableDef{}
// 	err = json.Unmarshal(rec.Get("def").Str, tdef)
// 	assert(err == nil) // JSON格式应该是正确的
// 	return tdef
// }

// // Get 公共API：通过主键查询记录
// // 这是用户层面的查询接口
// func (tx *DBTX) Get(table string, rec *Record) (bool, error) {
// 	// 1. 获取表结构定义（带缓存）
// 	tdef := getTableDef(tx, table)
// 	if tdef == nil {
// 		return false, fmt.Errorf("table not found: %s", table)
// 	}

// 	// 2. 委托给底层实现
// 	return dbGet(tx, tdef, rec)
// }

// // ================================
// // 表结构校验和创建 - 表的生命周期管理
// // ================================

// const TABLE_PREFIX_MIN = 100 // 用户表前缀的起始编号（0-99保留给系统）

// // tableDefCheck 校验表结构定义的完整性和正确性
// // 这是创建表时的重要安全检查
// func tableDefCheck(tdef *TableDef) error {
// 	// 基本完整性检查
// 	bad := tdef.Name == "" || len(tdef.Cols) == 0 || len(tdef.Indexes) == 0
// 	bad = bad || len(tdef.Cols) != len(tdef.Types) // 列数与类型数必须匹配
// 	if bad {
// 		return fmt.Errorf("bad table schema: %s", tdef.Name)
// 	}

// 	// 校验所有索引定义的正确性
// 	for i, index := range tdef.Indexes {
// 		// checkIndexCols 会修正索引定义（添加主键列）
// 		correctedIndex, err := checkIndexCols(tdef, index)
// 		if err != nil {
// 			return err
// 		}
// 		tdef.Indexes[i] = correctedIndex
// 	}
// 	return nil
// }

// // checkIndexCols 校验和修正索引列定义
// // 关键功能：确保每个二级索引都包含完整的主键
// // 这是B+树索引正确工作的必要条件
// func checkIndexCols(tdef *TableDef, index []string) ([]string, error) {
// 	if len(index) == 0 {
// 		return nil, fmt.Errorf("empty index")
// 	}

// 	seen := map[string]bool{}

// 	// 1. 验证索引列的存在性和唯一性
// 	for _, c := range index {
// 		// 检查列是否在表定义中
// 		if slices.Index(tdef.Cols, c) < 0 {
// 			return nil, fmt.Errorf("unknown index column: %s", c)
// 		}
// 		// 检查列是否重复
// 		if seen[c] {
// 			return nil, fmt.Errorf("duplicated column in index: %s", c)
// 		}
// 		seen[c] = true
// 	}

// 	// 2. 关键步骤：确保索引包含完整的主键
// 	// 原因：二级索引需要能够唯一定位到主表记录
// 	// 如果索引不包含完整主键，就可能出现重复的索引键，导致查询歧义
// 	primaryKey := tdef.Indexes[0]
// 	for _, primaryCol := range primaryKey {
// 		if !seen[primaryCol] {
// 			index = append(index, primaryCol) // 自动添加缺失的主键列
// 		}
// 	}

// 	assert(len(index) <= len(tdef.Cols)) // 索引列数不能超过表列数
// 	return index, nil
// }

// // TableNew 创建新表的完整流程
// // 这是DDL操作，需要更新系统表和分配系统资源
// func (tx *DBTX) TableNew(tdef *TableDef) error {
// 	// 0. 校验表结构的正确性
// 	if err := tableDefCheck(tdef); err != nil {
// 		return err
// 	}

// 	// 1. 检查表名冲突
// 	table := (&Record{}).AddStr("name", []byte(tdef.Name))
// 	ok, err := dbGet(tx, TDEF_TABLE, table)
// 	assert(err == nil)
// 	if ok {
// 		return fmt.Errorf("table exists: %s", tdef.Name)
// 	}

// 	// 2. 分配唯一的前缀编号
// 	// 从@meta表的next_prefix计数器获取下一个可用前缀
// 	prefix := uint32(TABLE_PREFIX_MIN)
// 	meta := (&Record{}).AddStr("key", []byte("next_prefix"))
// 	ok, err = dbGet(tx, TDEF_META, meta)
// 	assert(err == nil)

// 	if ok {
// 		// 从现有计数器获取
// 		prefix = binary.LittleEndian.Uint32(meta.Get("val").Str)
// 		assert(prefix > TABLE_PREFIX_MIN)
// 	} else {
// 		// 首次分配，初始化计数器记录
// 		meta.AddStr("val", make([]byte, 4))
// 	}

// 	// 3. 为每个索引分配连续的前缀编号
// 	assert(len(tdef.Prefixes) == 0) // 确保是新表
// 	for i := range tdef.Indexes {
// 		tdef.Prefixes = append(tdef.Prefixes, prefix+uint32(i))
// 	}

// 	// 4. 更新前缀计数器，为下一个表预留编号
// 	nextPrefix := prefix + uint32(len(tdef.Indexes))
// 	binary.LittleEndian.PutUint32(meta.Get("val").Str, nextPrefix)
// 	_, err = dbUpdate(tx, TDEF_META, &DBUpdateReq{Record: *meta})
// 	if err != nil {
// 		return err
// 	}

// 	// 5. 持久化表结构到@table系统表
// 	tableDefJSON, err := json.Marshal(tdef)
// 	assert(err == nil)
// 	table.AddStr("def", tableDefJSON)
// 	_, err = dbUpdate(tx, TDEF_TABLE, &DBUpdateReq{Record: *table})
// 	return err
// }

// // ================================
// // 数据更新操作 - 表数据的增删改
// // ================================

// // DBUpdateReq 数据更新请求的参数和结果
// type DBUpdateReq struct {
// 	// 输入参数
// 	Record Record // 要更新的记录数据
// 	Mode   int    // 更新模式（仅插入/仅更新/插入或更新）

// 	// 输出结果
// 	Updated bool  // 是否更新了现有记录
// 	Added   bool  // 是否添加了新记录
// }

// // dbUpdate 更新表中的一行数据
// // 这是数据库中最复杂的操作之一，需要同时维护：
// // 1. 主表数据
// // 2. 所有相关的二级索引
// // 3. 事务一致性
// func dbUpdate(tx *DBTX, tdef *TableDef, dbreq *DBUpdateReq) (bool, error) {
// 	// 1. 标准化列顺序：主键列在前，非主键列在后
// 	// 这样做是为了保证编码的一致性
// 	primaryKeyCols := tdef.Indexes[0]
// 	nonPrimaryCols := nonPrimaryKeyCols(tdef)
// 	cols := slices.Concat(primaryKeyCols, nonPrimaryCols)

// 	// 提取完整的列值，确保类型匹配
// 	values, err := getValues(tdef, dbreq.Record, cols)
// 	if err != nil {
// 		return false, err // 列不完整或类型不匹配
// 	}

// 	// 2. 更新主表数据
// 	npk := len(tdef.Indexes[0]) // 主键列数量

// 	// 构造主表的键值对
// 	primaryKey := encodeKey(nil, tdef.Prefixes[0], values[:npk])   // 主键编码
// 	nonPrimaryData := encodeValues(nil, values[npk:])             // 非主键数据编码

// 	// 执行底层KV更新
// 	req := UpdateReq{Key: primaryKey, Val: nonPrimaryData, Mode: dbreq.Mode}
// 	if _, err := tx.kv.Update(&req); err != nil {
// 		return false, err // 可能是键太长或其他底层错误
// 	}
// 	dbreq.Added, dbreq.Updated = req.Added, req.Updated

// 	// 3. 维护二级索引的一致性
// 	if req.Updated && !req.Added {
// 		// 如果是更新操作，需要先删除旧记录的索引键
// 		// 重建旧记录以便删除索引
// 		decodeValues(req.Old, values[npk:]) // 恢复旧的非主键数据
// 		oldRecord := Record{cols, values}

// 		// 删除旧记录的所有二级索引键
// 		err := indexOp(tx, tdef, INDEX_DEL, oldRecord)
// 		assert(err == nil) // 删除操作不应该失败
// 	}

// 	if req.Updated {
// 		// 为新记录添加二级索引键
// 		if err := indexOp(tx, tdef, INDEX_ADD, dbreq.Record); err != nil {
// 			return false, err // 可能是键太长
// 		}
// 	}

// 	return req.Updated, nil
// }

// // nonPrimaryKeyCols 提取非主键列名
// // 辅助函数，用于分离主键和非主键列
// func nonPrimaryKeyCols(tdef *TableDef) (out []string) {
// 	primaryKeySet := make(map[string]bool)
// 	for _, c := range tdef.Indexes[0] {
// 		primaryKeySet[c] = true
// 	}

// 	for _, c := range tdef.Cols {
// 		if !primaryKeySet[c] {
// 			out = append(out, c) // 只返回非主键列
// 		}
// 	}
// 	return
// }

// // ================================
// // 索引管理 - 二级索引的维护
// // ================================

// const (
// 	INDEX_ADD = 1  // 添加索引键操作
// 	INDEX_DEL = 2  // 删除索引键操作
// )

// // indexOp 批量添加或删除一条记录的所有二级索引键
// // 这是保持索引一致性的核心函数
// func indexOp(tx *DBTX, tdef *TableDef, op int, rec Record) error {
// 	// 遍历所有二级索引（跳过主键索引tdef.Indexes[0]）
// 	for i := 1; i < len(tdef.Indexes); i++ {
// 		indexColumns := tdef.Indexes[i]
// 		indexPrefix := tdef.Prefixes[i]

// 		// 1. 提取当前索引涉及的列值
// 		values, err := getValues(tdef, rec, indexColumns)
// 		assert(err == nil) // 完整记录应该包含所有必要的列

// 		// 2. 构造索引键（注意：索引值为空，因为二级索引只需要键）
// 		indexKey := encodeKey(nil, indexPrefix, values)

// 		// 3. 根据操作类型执行相应的KV操作
// 		switch op {
// 		case INDEX_ADD:
// 			// 添加索引键，值设为nil（二级索引不存储数据）
// 			req := UpdateReq{Key: indexKey, Val: nil}
// 			if _, err := tx.kv.Update(&req); err != nil {
// 				return err // 可能是键长度限制
// 			}
// 			assert(req.Added) // 索引键应该是新添加的

// 		case INDEX_DEL:
// 			// 删除索引键
// 			deleted, err := tx.kv.Del(&DeleteReq{Key: indexKey})
// 			assert(err == nil) // 删除操作不应该失败
// 			assert(deleted)    // 索引键应该存在并被删除

// 		default:
// 			panic("invalid index operation")
// 		}
// 	}
// 	return nil
// }

// // ================================
// // 公共CRUD接口 - 用户层面的数据操作API
// // ================================

// // Set 通用设置接口，支持不同的更新模式
// func (tx *DBTX) Set(table string, dbreq *DBUpdateReq) (bool, error) {
// 	tdef := getTableDef(tx, table)
// 	if tdef == nil {
// 		return false, fmt.Errorf("table not found: %s", table)
// 	}
// 	return dbUpdate(tx, tdef, dbreq)
// }

// // Insert 插入新记录，仅当记录不存在时成功
// // 对应SQL的INSERT语句
// func (tx *DBTX) Insert(table string, rec Record) (bool, error) {
// 	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_INSERT_ONLY})
// }

// // Update 更新现有记录，仅当记录已存在时成功
// // 对应SQL的UPDATE语句
// func (tx *DBTX) Update(table string, rec Record) (bool, error) {
// 	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_UPDATE_ONLY})
// }

// // Upsert 插入或更新记录，无论记录是否存在都会成功
// // 对应MySQL的INSERT ... ON DUPLICATE KEY UPDATE
// func (tx *DBTX) Upsert(table string, rec Record) (bool, error) {
// 	return tx.Set(table, &DBUpdateReq{Record: rec, Mode: MODE_UPSERT})
// }

// // ================================
// // 删除操作 - 记录的完整删除
// // ================================

// // dbDelete 通过主键删除记录，同时清理所有相关索引
// func dbDelete(tx *DBTX, tdef *TableDef, rec Record) (bool, error) {
// 	// 1. 从输入记录中提取主键值
// 	primaryKeyValues, err := getValues(tdef, rec, tdef.Indexes[0])
// 	if err != nil {
// 		return false, err // 主键不完整
// 	}

// 	// 2. 删除主表记录
// 	primaryKey := encodeKey(nil, tdef.Prefixes[0], primaryKeyValues)
// 	deleteReq := DeleteReq{Key: primaryKey}
// 	if deleted, _ := tx.kv.Del(&deleteReq); !deleted {
// 		return false, nil // 记录不存在或键太长
// 	}

// 	// 3. 重建完整记录以删除相关索引
// 	// 需要恢复被删除记录的所有列值
// 	allValues := primaryKeyValues // 开始时只有主键值

// 	// 为所有非主键列创建占位符
// 	for _, c := range nonPrimaryKeyCols(tdef) {
// 		colType := tdef.Types[slices.Index(tdef.Cols, c)]
// 		allValues = append(allValues, Value{Type: colType})
// 	}

// 	// 从删除操作返回的旧值中恢复非主键数据
// 	npk := len(tdef.Indexes[0])
// 	decodeValues(deleteReq.Old, allValues[npk:])

// 	// 4. 删除该记录的所有二级索引键
// 	completeRecord := Record{tdef.Cols, allValues}
// 	err = indexOp(tx, tdef, INDEX_DEL, completeRecord)
// 	assert(err == nil) // 索引删除不应该失败

// 	return true, nil
// }

// // Delete 公共删除接口
// func (tx *DBTX) Delete(table string, rec Record) (bool, error) {
// 	tdef := getTableDef(tx, table)
// 	if tdef == nil {
// 		return false, fmt.Errorf("table not found: %s", table)
// 	}
// 	return dbDelete(tx, tdef, rec)
// }

// // ================================
// // 数据库生命周期管理 - 数据库的启动和关闭
// // ================================

// // Open 初始化并打开数据库
// func (db *DB) Open() error {
// 	db.kv.Path = db.Path                // 设置底层存储路径
// 	db.tables = map[string]*TableDef{}  // 初始化表结构缓存
// 	return db.kv.Open()                 // 打开底层KV存储
// }

// // Close 关闭数据库，释放资源
// func (db *DB) Close() {
// 	db.kv.Close()
// }

// // ================================
// // 范围查询和扫描 - 复杂查询的实现
// // ================================

// // Scanner 范围查询迭代器
// // 提供对B+树范围扫描的高级封装
// type Scanner struct {
// 	// 查询范围定义：[Key1, Key2]
// 	Cmp1 int    // 起始边界类型（CMP_GE: >=, CMP_GT: >）
// 	Cmp2 int    // 结束边界类型（CMP_LE: <=, CMP_LT: <）
// 	Key1 Record // 范围起始键
// 	Key2 Record // 范围结束键

// 	// 内部执行状态
// 	tx    *DBTX    // 关联的数据库事务
// 	tdef  *TableDef // 目标表的结构定义
// 	index int       // 选中的索引编号（0=主键，1+=二级索引）
// 	iter  KVIter    // 底层B+树迭代器
// }

// // Valid 检查迭代器当前位置是否有效（是否还有数据）
// func (sc *Scanner) Valid() bool {
// 	return sc.iter.Valid()
// }

// // Next 移动迭代器到下一条记录
// func (sc *Scanner) Next() {
// 	sc.iter.Next()
// }

// // Deref 提取迭代器当前位置的完整记录
// // 这是范围查询的核心方法，需要处理主键扫描和二级索引扫描两种情况
// func (sc *Scanner) Deref(rec *Record) {
// 	assert(sc.Valid())
// 	tdef := sc.tdef

// 	// 1. 准备输出记录的结构（标准化列顺序）
// 	rec.Cols = slices.Concat(tdef.Indexes[0], nonPrimaryKeyCols(tdef))
// 	rec.Vals = rec.Vals[:0] // 清空之前的值

// 	// 为每一列创建对应类型的Value占位符
// 	for _, colName := range rec.Cols {
// 		colIndex := slices.Index(tdef.Cols, colName)
// 		colType := tdef.Types[colIndex]
// 		rec.Vals = append(rec.Vals, Value{Type: colType})
// 	}

// 	// 2. 从底层KV迭代器获取当前键值对
// 	currentKey, currentVal := sc.iter.Deref()

// 	// 3. 根据扫描类型进行不同的处理
// 	if sc.index == 0 {
// 		// 情况A：主键扫描（直接从主表读取）
// 		npk := len(tdef.Indexes[0])

// 		// 解码主键部分（存储在键中）
// 		decodeKey(currentKey, rec.Vals[:npk])

// 		// 解码非主键部分（存储在值中）
// 		decodeValues(currentVal, rec.Vals[npk:])

// 	} else {
// 		// 情况B：二级索引扫描（需要回表查询）
// 		assert(len(currentVal) == 0) // 二级索引的值部分为空

// 		// 解码索引键，提取索引列的值
// 		indexColumns := tdef.Indexes[sc.index]
// 		indexRecord := Record{indexColumns, make([]Value, len(indexColumns))}

// 		// 为索引记录的每一列设置正确的类型
// 		for i, colName := range indexColumns {
// 			colIndex := slices.Index(tdef.Cols, colName)
// 			indexRecord.Vals[i].Type = tdef.Types[colIndex]
// 		}

// 		// 解码索引键
// 		decodeKey(currentKey, indexRecord.Vals)

// 		// 从索引记录中提取主键值
// 		primaryKeyCols := tdef.Indexes[0]
// 		for i, primaryCol := range primaryKeyCols {
// 			rec.Vals[i] = *indexRecord.Get(primaryCol)
// 		}

// 		// 通过主键回表查询完整记录
// 		// 注意：这里可能有优化空间，如果索引包含了所有需要的列，就不需要回表
// 		ok, err := dbGet(sc.tx, tdef, rec)
// 		assert(ok && err == nil) // 索引指向的记录必须存在
// 	}
// }

// // checkTypes 验证记录的列类型是否与表定义匹配
// // 查询前的安全检查
// func checkTypes(tdef *TableDef, rec Record) error {
// 	if len(rec.Cols) != len(rec.Vals) {
// 		return fmt.Errorf("bad record: column count mismatch")
// 	}

// 	for i, colName := range rec.Cols {
// 		// 查找列在表定义中的位置
// 		colIndex := slices.Index(tdef.Cols, colName)
// 		if colIndex < 0 {
// 			return fmt.Errorf("unknown column: %s", colName)
// 		}

// 		// 检查类型匹配
// 		expectedType := tdef.Types[colIndex]
// 		actualType := rec.Vals[i].Type
// 		if expectedType != actualType {
// 			return fmt.Errorf("bad column type: %s", colName)
// 		}
// 	}
// 	return nil
// }

// // dbScan 执行范围扫描查询
// // 这是范围查询的核心实现，包括索引选择、键编码、迭代器创建
// func dbScan(tx *DBTX, tdef *TableDef, req *Scanner) error {
// 	// 0. 基本有效性检查

// 	// 检查比较操作符的组合是否合理
// 	switch {
// 	case req.Cmp1 > 0 && req.Cmp2 < 0: // >= ... <=  或  > ... <
// 	case req.Cmp2 > 0 && req.Cmp1 < 0: // <= ... >=  （反向范围，也是合理的）
// 	default:
// 		return fmt.Errorf("bad range: invalid comparison operators")
// 	}

// 	// 检查起始键和结束键的结构一致性
// 	if !slices.Equal(req.Key1.Cols, req.Key2.Cols) {
// 		return fmt.Errorf("bad range: key structure mismatch")
// 	}

// 	// 检查键值的类型正确性
// 	if err := checkTypes(tdef, req.Key1); err != nil {
// 		return fmt.Errorf("bad range key1: %w", err)
// 	}
// 	if err := checkTypes(tdef, req.Key2); err != nil {
// 		return fmt.Errorf("bad range key2: %w", err)
// 	}

// 	// 设置扫描器的内部状态
// 	req.tx = tx
// 	req.tdef = tdef

// 	// 1. 索引选择算法
// 	// 选择能够覆盖查询键前缀的索引
// 	// 例如：查询键为["name", "age"]，索引为["name", "age", "id"]时可以使用
// 	isCovered := func(indexColumns []string) bool {
// 		queryColumns := req.Key1.Cols
// 		return len(indexColumns) >= len(queryColumns) &&
// 		       slices.Equal(indexColumns[:len(queryColumns)], queryColumns)
// 	}

// 	req.index = slices.IndexFunc(tdef.Indexes, isCovered)
// 	if req.index < 0 {
// 		return fmt.Errorf("no suitable index found for columns: %v", req.Key1.Cols)
// 	}

// 	// 2. 编码查询范围的边界键
// 	indexPrefix := tdef.Prefixes[req.index]
// 	keyStart := encodeKeyPartial(nil, indexPrefix, req.Key1.Vals, req.Cmp1)
// 	keyEnd := encodeKeyPartial(nil, indexPrefix, req.Key2.Vals, req.Cmp2)

// 	// 3. 创建底层KV迭代器并定位到起始位置
// 	req.iter = tx.kv.Seek(keyStart, req.Cmp1, keyEnd, req.Cmp2)
// 	return nil
// }

// // Scan 公共范围查询接口
// func (tx *DBTX) Scan(table string, req *Scanner) error {
// 	tdef := getTableDef(tx, table)
// 	if tdef == nil {
// 		return fmt.Errorf("table not found: %s", table)
// 	}
// 	return dbScan(tx, tdef, req)
// }

// // ================================
// // 数据导入导出功能 - 数据库备份和恢复
// // ================================

// const (
// 	MIN_NAME = "\x00"  // 最小可能的表名，用于全表扫描的下界
// 	MAX_NAME = "\xFF"  // 最大可能的表名，用于全表扫描的上界
// )

// // GetAllTables 获取数据库中所有用户表的定义
// // 用于数据库元数据的批量导出
// func (db *DB) GetAllTables() ([]TableDef, error) {
// 	tx := DBTX{}
// 	db.Begin(&tx)
// 	defer func() {
// 		// 确保事务被正确提交，即使发生错误
// 		if err := db.Commit(&tx); err != nil {
// 			db.Abort(&tx)
// 		}
// 	}()

// 	// 构造全表扫描：查询所有可能的表名
// 	sc := &Scanner{
// 		Cmp1: CMP_GE,
// 		Cmp2: CMP_LE,
// 		Key1: *(&Record{}).AddStr("name", []byte(MIN_NAME)),
// 		Key2: *(&Record{}).AddStr("name", []byte(MAX_NAME)),
// 	}

// 	// 从@table系统表中扫描所有表定义
// 	if err := tx.Scan("@table", sc); err != nil {
// 		return nil, fmt.Errorf("failed to scan table definitions: %w", err)
// 	}

// 	// 收集扫描结果
// 	tableRecords := reduceSelectData(sc)

// 	// 反序列化每个表定义
// 	var tableDefs []TableDef
// 	for _, record := range tableRecords {
// 		tdef := &TableDef{}
// 		definitionJSON := record.Get("def").Str
// 		if err := json.Unmarshal(definitionJSON, tdef); err != nil {
// 			return nil, fmt.Errorf("failed to parse table definition: %w", err)
// 		}
// 		tableDefs = append(tableDefs, *tdef)
// 	}

// 	return tableDefs, nil
// }

// // ExportDB 导出整个数据库到文件系统
// // 创建包含表结构和所有数据的快照目录
// // 返回快照目录路径，可用于后续的备份或迁移
// func (db *DB) ExportDB() (snapshotDir string, err error) {
// 	// 1. 确定导出目录名（在原文件名基础上加_export后缀）
// 	if strings.HasSuffix(db.Path, ".db") {
// 		snapshotDir = db.Path[:len(db.Path)-3] + "_export"
// 	} else {
// 		snapshotDir = db.Path + "_export"
// 	}

// 	// 2. 创建快照目录结构
// 	if err = os.MkdirAll(snapshotDir, 0755); err != nil {
// 		return "", fmt.Errorf("failed to create snapshot directory: %w", err)
// 	}

// 	// 3. 导出所有表的结构定义
// 	allTables, err := db.GetAllTables()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get table definitions: %w", err)
// 	}

// 	// 将表结构写入schema.json
// 	schemaFile, err := os.Create(filepath.Join(snapshotDir, "schema.json"))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create schema file: %w", err)
// 	}
// 	defer schemaFile.Close()

// 	encoder := json.NewEncoder(schemaFile)
// 	encoder.SetIndent("", "  ") // 美化JSON格式
// 	if err = encoder.Encode(allTables); err != nil {
// 		return "", fmt.Errorf("failed to encode schema: %w", err)
// 	}

// 	// 4. 导出每张表的数据
// 	for _, tableDef := range allTables {
// 		if err := db.exportTableData(tableDef, snapshotDir); err != nil {
// 			return "", fmt.Errorf("failed to export table %s: %w", tableDef.Name, err)
// 		}
// 	}

// 	return snapshotDir, nil
// }

// // exportTableData 导出单个表的所有数据
// func (db *DB) exportTableData(tableDef TableDef, snapshotDir string) error {
// 	// 创建表数据文件
// 	dataPath := filepath.Join(snapshotDir, tableDef.Name+".data")
// 	dataFile, err := os.Create(dataPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer dataFile.Close()

// 	// 全表扫描：获取表中的所有记录
// 	tx := DBTX{}
// 	db.Begin(&tx)
// 	defer db.Commit(&tx)

// 	// 构造覆盖整个主键范围的扫描
// 	sc := &Scanner{
// 		Cmp1: CMP_GE,
// 		Cmp2: CMP_LE,
// 		Key1: *(&Record{}).AddInt64("id", math.MinInt64/2), // 避免溢出
// 		Key2: *(&Record{}).AddInt64("id", math.MaxInt64/2),
// 	}

// 	if err := tx.Scan(tableDef.Name, sc); err != nil {
// 		return err
// 	}

// 	// 使用gob格式序列化记录（高效的二进制格式）
// 	encoder := gob.NewEncoder(dataFile)
// 	records := reduceSelectData(sc)

// 	for _, record := range records {
// 		if err := encoder.Encode(record); err != nil {
// 			return fmt.Errorf("failed to encode record: %w", err)
// 		}
// 	}

// 	return nil
// }

// // LoadRecordsFromDataFile 从gob格式的数据文件中加载记录
// func LoadRecordsFromDataFile(dataPath string) ([]Record, error) {
// 	dataFile, err := os.Open(dataPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer dataFile.Close()

// 	decoder := gob.NewDecoder(dataFile)
// 	var records []Record

// 	// 逐个解码记录，直到文件结束
// 	for {
// 		var record Record
// 		err := decoder.Decode(&record)
// 		if err == io.EOF {
// 			break // 正常结束
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to decode record: %w", err)
// 		}
// 		records = append(records, record)
// 	}

// 	return records, nil
// }

// // ImportDB 从快照目录恢复完整的数据库
// // 这是ExportDB的逆操作，用于数据库恢复或迁移
// func ImportDB(snapshotDir, newDBPath string) (*DB, error) {
// 	// 1. 创建新的数据库实例
// 	db := &DB{Path: newDBPath}
// 	if err := db.Open(); err != nil {
// 		return nil, fmt.Errorf("failed to open new database: %w", err)
// 	}

// 	// 2. 验证快照目录的存在性
// 	if _, err := os.Stat(snapshotDir); os.IsNotExist(err) {
// 		db.Close()
// 		return nil, fmt.Errorf("snapshot directory %s does not exist", snapshotDir)
// 	}

// 	// 3. 读取表结构定义
// 	schemaPath := filepath.Join(snapshotDir, "schema.json")
// 	schemaData, err := os.ReadFile(schemaPath)
// 	if err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("failed to read schema file: %w", err)
// 	}

// 	var tableDefs []TableDef
// 	if err := json.Unmarshal(schemaData, &tableDefs); err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("failed to parse schema file: %w", err)
// 	}

// 	// 4. 在单个事务中导入所有表和数据
// 	tx := &DBTX{}
// 	db.Begin(tx)

// 	// 使用defer确保事务会被正确处理
// 	committed := false
// 	defer func() {
// 		if !committed {
// 			db.Abort(tx)
// 			db.Close()
// 		}
// 	}()

// 	// 并发导入所有表（在同一事务内）
// 	for _, tableDef := range tableDefs {
// 		if err := db.importSingleTable(tx, tableDef, snapshotDir); err != nil {
// 			return nil, fmt.Errorf("failed to import table %s: %w", tableDef.Name, err)
// 		}
// 	}

// 	// 5. 提交事务
// 	if err := db.Commit(tx); err != nil {
// 		return nil, fmt.Errorf("failed to commit import transaction: %w", err)
// 	}

// 	committed = true
// 	return db, nil
// }

// // importSingleTable 导入单个表及其数据
// func (db *DB) importSingleTable(tx *DBTX, tableDef TableDef, snapshotDir string) error {
// 	// 1. 创建表结构（清空前缀，让系统重新分配）
// 	tableDef.Prefixes = []uint32{} // 重置前缀，避免冲突
// 	if err := tx.TableNew(&tableDef); err != nil {
// 		return fmt.Errorf("failed to create table structure: %w", err)
// 	}

// 	// 2. 加载表数据
// 	dataPath := filepath.Join(snapshotDir, tableDef.Name+".data")
// 	records, err := LoadRecordsFromDataFile(dataPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to load table data: %w", err)
// 	}

// 	// 3. 插入所有记录
// 	for i, record := range records {
// 		if _, err := tx.Insert(tableDef.Name, record); err != nil {
// 			return fmt.Errorf("failed to insert record %d: %w", i, err)
// 		}
// 	}

// 	return nil
// }

// // ================================
// // 辅助函数
// // ================================

// // reduceSelectData 将扫描结果收集到切片中
// // 这是一个通用的辅助函数，将迭代器结果转换为内存中的记录数组
// func reduceSelectData(sc *Scanner) []Record {
// 	var records []Record

// 	// 遍历所有有效的记录
// 	for sc.Valid() {
// 		var record Record
// 		sc.Deref(&record)           // 提取当前记录
// 		records = append(records, record)
// 		sc.Next()                   // 移动到下一条记录
// 	}

// 	return records
// }

// // assert 简单的断言函数，用于内部一致性检查
// // 在生产环境中，这些断言表示"不应该发生"的情况
// func assert(condition bool) {
// 	if !condition {
// 		panic("assertion failed: internal consistency error")
// 	}
// }
