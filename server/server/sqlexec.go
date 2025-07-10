package server

import (
	"fmt"
	"strconv"

	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	"google.golang.org/protobuf/proto"
)

const (
	TableNewType int = iota
	InsertData
	SelectData
	UpdateData
	DeleteData
	UnKnowType
)

const (
	INSERT_OK       = "INSERT DATA OK"
	DELETE_OK       = "DELETE DATA OK"
	CREATE_TABLE_OK = "CTEATE TABLE OK"
	UPDATE_OK       = "UPDATE DATA OK"
)

func (db *DBServer) ExecSQL(sql string) ([]byte, error) {
	tree := db.SQLParser.BuildTree(sql)
	c := db.SQLParser.Visit(tree)
	switch res := c.(type) {
	case *kv.TableDef:
		return db.ExecTableNew(res)
	case *InsertRes:
		return db.ExecInsert(res.TableName, *res.Rec)
	case *DelRes:
		return db.ExecDeleteRecords(res.TableName, res.Scan)
	case *SelectInfo:
		return db.ExecSelectData(res.TableName, res.SelectField, res.Scan)
	case *UpdateRes:
		return db.ExecUpdateData(res.TableName, res.UpdateMp, res.Scan)
	default:
	}
	return []byte("UNKNOW SQL TYPE"), nil
}

func (db *DBServer) ExecUpdateData(tablename string, updateMP map[string]string, scan *Scanner) ([]byte, error) {
	if scan == nil {
		return nil, fmt.Errorf("scan is nil")
	}
	recs, err := db.Scan(tablename, scan)
	if err != nil {
		return nil, fmt.Errorf("Scan data err:%v", err)
	}
	for _, rec := range recs {
		for col, newval := range updateMP {
			for i, v := range rec.Cols {
				if v == col {
					switch rec.Vals[i].Type {
					case kv.TYPE_INT64:
						newvalInt, _ := strconv.ParseInt(newval, 10, 64)
						rec.Vals[i].I64 = newvalInt
					case kv.TYPE_BYTES:
						rec.Vals[i].Str = []byte(newval)
					default:
						return nil, fmt.Errorf("unsupported type for update")
					}
				}
			}
		}
		// Update 仅当存在时才更新
		_, err := db.Update(tablename, pbRecordToKvRecord(rec))
		if err != nil {
			return nil, err
		}
	}
	return []byte("UPDATE DATA OK"), nil
}
func (db *DBServer) ExecSelectData(tablename string, selectFields []string, scan *Scanner) ([]byte, error) {
	// TODO 先不筛选 selectFields
	pbRecords, err := db.Scan(tablename, scan)
	if err != nil {
		return nil, fmt.Errorf("DB SCAN ERROR:%v", err)
	}
	recordsList := &raftpb.ScanResponse{Records: pbRecords}
	data, err := proto.Marshal(recordsList)
	if err != nil {
		return nil, fmt.Errorf("UNEXPECTED OCCUR:%v", err)
	}
	return data, nil
}

func (db *DBServer) ExecTableNew(table *kv.TableDef) ([]byte, error) {
	if err := db.TableNew(table); err != nil {
		return nil, err
	}
	return []byte(CREATE_TABLE_OK), nil
}
func (db *DBServer) ExecInsert(tablename string, rec kv.Record) ([]byte, error) {
	_, err := db.Insert(tablename, rec)
	if err != nil {
		return nil, err
	}
	return []byte(INSERT_OK), nil
}

func (db *DBServer) ExecDeleteRecords(tablename string, scan *Scanner) ([]byte, error) {
	recs, err := db.Scan(tablename, scan)
	if err != nil {
		return nil, err
	}
	for _, rec := range recs {
		ok, err := db.Delete(tablename, pbRecordToKvRecord(rec))
		assert(ok)
		if err != nil {
			return nil, err
		}
	}
	return []byte(DELETE_OK), nil
}
func pbRecordToKvRecord(pbRec *raftpb.Record) kv.Record {
	kvRec := kv.Record{
		Cols: make([]string, len(pbRec.Cols)),
		Vals: make([]kv.Value, len(pbRec.Vals)),
	}
	copy(kvRec.Cols, pbRec.Cols)
	for i, pbVal := range pbRec.Vals {
		kvRec.Vals[i] = kv.Value{
			Type: pbVal.Type,
			I64:  pbVal.I64,
			Str:  pbVal.Str,
		}
	}

	return kvRec
}
