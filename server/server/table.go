package server

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
)

const (
	NEXT_PREFIX = "next_prefix"
)

// get the table schema by name
func (db *DBServer) getTableDef(name string) (*kv.TableDef, error) {
	if tdef, ok := kv.INTERNAL_TABLES[name]; ok {
		return tdef, nil // expose internal tables
	}

	return db.TablesDefDiscovery.GetTable(name)
}

const (
	INDEX_ADD = 1
	INDEX_DEL = 2
)

// add or remove secondary index keys
func indexOp(client pb.KVStoreClient, tdef *kv.TableDef, op int, rec kv.Record) error {
	for i := 1; i < len(tdef.Indexes); i++ {
		// the indexed key
		values, err := kv.GetValues(tdef, rec, tdef.Indexes[i])
		assert(err == nil) // full row
		key := kv.EncodeKey(nil, tdef.Prefixes[i], values)
		switch op {
		case INDEX_ADD:
			req, err := client.Put(context.Background(), &pb.PutRequest{Key: key, Value: nil})
			if err != nil {
				return err
			}
			assert(req.Added) // internal consistency
		case INDEX_DEL:
			req, err := client.Delete(context.Background(), &pb.DeleteRequest{Key: key})
			assert(err == nil)  // should not run into the length limit
			assert(req.Success) // internal consistency
		default:
			panic("unreachable")
		}
	}
	return nil
}

// etcd 提取 TABLE_PREFIX_MIN 转换
func uint32ToBytesLE(n uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, n)
	return buf
}

func bytesToUint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}
func (db *DBServer) TableNew(tdef *kv.TableDef) error {
	// 0. sanity checks, 表结构检查
	if err := kv.TableDefCheck(tdef); err != nil {
		return err
	}
	// 1. check the existing table, 检查是否有重复
	existing, err := db.TablesDefDiscovery.GetTable(tdef.Name)
	assert(err == nil)
	if existing != nil {
		return fmt.Errorf("table exists: %s", tdef.Name)
	}
	// 2. allocate new prefixes
	prefix := uint32(kv.TABLE_PREFIX_MIN)
	meta, err := db.TablesDefDiscovery.GetMetaKey(NEXT_PREFIX)
	assert(err == nil)
	if meta != nil {
		// 接着上次的用
		prefix = binary.LittleEndian.Uint32(meta)
	}
	assert(len(tdef.Prefixes) == 0)
	for i := range tdef.Indexes {
		tdef.Prefixes = append(tdef.Prefixes, prefix+uint32(i))
	}
	// 3. update the next prefix
	// FIXME: integer overflow.
	next := prefix + uint32(len(tdef.Indexes))
	err = db.TablesDefDiscovery.PutMetaKey(NEXT_PREFIX, uint32ToBytesLE(next))
	assert(err == nil)

	// 存储tabledef
	return db.TablesDefDiscovery.RegisterTable(tdef)
}

type DBUpdateReq struct {
	// in
	Record kv.Record
	Mode   int
	// out
	Updated bool
	Added   bool
}

func nonPrimaryKeyCols(tdef *kv.TableDef) (out []string) {
	for _, c := range tdef.Cols {
		if slices.Index(tdef.Indexes[0], c) < 0 {
			out = append(out, c)
		}
	}
	return
}
func (db *DBServer) dbUpdate(tdef *kv.TableDef, dbreq *DBUpdateReq) (bool, error) {
	// reorder the columns so that they start with the primary key
	cols := slices.Concat(tdef.Indexes[0], nonPrimaryKeyCols(tdef))
	values, err := kv.GetValues(tdef, dbreq.Record, cols)
	if err != nil {
		return false, err // expect a full row
	}
	// insert the row
	npk := len(tdef.Indexes[0]) // number of primary key columns
	key := kv.EncodeKey(nil, tdef.Prefixes[0], values[:npk])
	val := kv.EncodeValues(nil, values[npk:])

	// get grpc serve
	client, _, err := db.getGrpcClientForPrimaryKey(tdef.Indexes[0])
	if err != nil {
		return false, err
	}
	req, err := client.Put(context.Background(), &pb.PutRequest{Key: key, Value: val, Mode: int32(dbreq.Mode)})
	if err != nil {
		return false, err
	}
	dbreq.Updated, dbreq.Added = req.Updated, req.Added

	// maintain secondary indexes
	if req.Updated && !req.Added {
		// construct the old record
		kv.DecodeValues(req.Old, values[npk:])
		oldRec := kv.Record{cols, values}
		// delete the indexed keys
		err := indexOp(client, tdef, INDEX_DEL, oldRec)
		assert(err == nil) // should not run into the length limit
	}
	if req.Updated {
		// add the new indexed keys
		if err := indexOp(client, tdef, INDEX_ADD, dbreq.Record); err != nil {
			return false, err // length limit
		}
	}
	return req.Updated, nil
}

// add a record
func (db *DBServer) Set(table string, dbreq *DBUpdateReq) (bool, error) {
	tdef, err := db.getTableDef(table)
	if err != nil {
		return false, err
	}
	if tdef == nil {
		return false, fmt.Errorf("table not found: %s", table)
	}
	return db.dbUpdate(tdef, dbreq)
}

// use this
func (db *DBServer) Insert(table string, rec kv.Record) (bool, error) {
	return db.Set(table, &DBUpdateReq{Record: rec, Mode: kv.MODE_INSERT_ONLY})
}
func (db *DBServer) Update(table string, rec kv.Record) (bool, error) {
	return db.Set(table, &DBUpdateReq{Record: rec, Mode: kv.MODE_UPDATE_ONLY})
}
func (db *DBServer) Upsert(table string, rec kv.Record) (bool, error) {
	return db.Set(table, &DBUpdateReq{Record: rec, Mode: kv.MODE_UPSERT})
}

// delete a record by its primary key
func (db *DBServer) dbDelete(client pb.KVStoreClient, tdef *kv.TableDef, rec kv.Record) (bool, error) {
	values, err := kv.GetValues(tdef, rec, tdef.Indexes[0])
	if err != nil {
		return false, err
	}
	// delete the row
	req, err := client.Delete(context.Background(), &pb.DeleteRequest{
		Key: kv.EncodeKey(nil, tdef.Prefixes[0], values),
	})
	if !req.Success || err != nil {
		return false, nil // `deleted` is also false if the key is too long
	}

	// maintain secondary indexes
	for _, c := range nonPrimaryKeyCols(tdef) {
		tp := tdef.Types[slices.Index(tdef.Cols, c)]
		values = append(values, kv.Value{Type: tp})
	}
	kv.DecodeValues(req.Old, values[len(tdef.Indexes[0]):])
	err = indexOp(client, tdef, INDEX_DEL, kv.Record{tdef.Cols, values})
	assert(err == nil) // should not run into the length limit
	return true, nil
}
func (db *DBServer) Delete(table string, rec kv.Record) (bool, error) {
	tdef, err := db.getTableDef(table)
	if err != nil {
		return false, err
	}
	// get grpc serve
	client, _, err := db.getGrpcClientForPrimaryKey(tdef.Indexes[0])
	if err != nil {
		return false, err
	}
	return db.dbDelete(client, tdef, rec)
}

type Scanner struct {
	// the range, from Key1 to Key2
	Cmp1 int // CMP_??
	Cmp2 int
	Key1 kv.Record
	Key2 kv.Record

	index int // which index?
}

func (db *DBServer) dbScan(tdef *kv.TableDef, req *Scanner) ([]*pb.Record, error) {
	// 0. sanity checks
	switch {
	case req.Cmp1 > 0 && req.Cmp2 < 0:
	case req.Cmp2 > 0 && req.Cmp1 < 0:
	default:
		return nil, fmt.Errorf("bad range")
	}
	if err := kv.CheckTypes(tdef, req.Key1); err != nil {
		return nil, err
	}
	if err := kv.CheckTypes(tdef, req.Key2); err != nil {
		return nil, err
	}

	// 1. select the index, 选择合适的索引
	covered := func(key []string, index []string) bool {
		return len(index) >= len(key) && slices.Equal(index[:len(key)], key)
	}
	req.index = slices.IndexFunc(tdef.Indexes, func(index []string) bool {
		return covered(req.Key1.Cols, index) && covered(req.Key2.Cols, index)
	})
	if req.index < 0 {
		return nil, fmt.Errorf("no index")
	}
	// 2. encode the start key and the end key
	prefix := tdef.Prefixes[req.index]
	keyStart := kv.EncodeKeyPartial(nil, prefix, req.Key1.Vals, req.Cmp1)
	keyEnd := kv.EncodeKeyPartial(nil, prefix, req.Key2.Vals, req.Cmp2)

	// 3. seek to the start key
	// get grpc serve
	tdefBytes, err := json.Marshal(tdef)
	if err != nil {
		return nil, fmt.Errorf("tdef json marshal err:%v", err)
	}
	client, _, err := db.getGrpcClientForPrimaryKey(tdef.Indexes[0])
	if err != nil {
		return nil, err
	}
	clientReq, err := client.Scan(context.Background(), &pb.ScanRequest{
		KeyStart: keyStart, KeyEnd: keyEnd, Cmp1: int64(req.Cmp1), Cmp2: int64(req.Cmp2), Table: tdefBytes, Index: int64(req.index),
	})
	return clientReq.Records, err
}
func (db *DBServer) Scan(tablename string, req *Scanner) ([]*pb.Record, error) {
	tdef, err := db.getTableDef(tablename)
	if err != nil {
		return nil, fmt.Errorf("occur error:%v", err)
	}
	if tdef == nil {
		return nil, fmt.Errorf("table is nil")
	}
	return db.dbScan(tdef, req)
}
