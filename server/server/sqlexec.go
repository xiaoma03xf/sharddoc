package server

import (
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
)

const (
	TableNewType int = iota
	InsertData
	SelectData
	UpdateData
	DeleteData
	UnKnowType
)

func (db *DBServer) ExecSQLType(sql string) int {
	return db.SQLParser.CheckSQLType(sql)
}

func (db *DBServer) ExecInsert(tablename string, rec kv.Record) (bool, error) {
	return db.Insert(tablename, rec)
}
func (db *DBServer) ExecSelect(tablename string, rec kv.Record) (bool, error) {
	return db.Insert(tablename, rec)
}

func (db *DBServer) ExecDeleteRecords(tablename string, scan *Scanner) error {
	recs, err := db.Scan(tablename, scan)
	if err != nil {
		return err
	}
	for _, rec := range recs {
		ok, err := db.Delete(tablename, pbRecordToKvRecord(rec))
		assert(ok)
		if err != nil {
			return err
		}
	}
	return nil
}
func pbRecordToKvRecord(pbRec *pb.Record) kv.Record {
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
