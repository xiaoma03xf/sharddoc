package engine

// func (db *DB) Begin()
// select 语句
// db.Raw("SELECT name, id FROM tbl_test WHERE age = 18").Scan(&[]User{})

func (db *DBTX) Raw(sql string) []Record {
	selectInfo := VisitTree(sql).(SelectInfo)
	err := db.Scan(selectInfo.TableName, selectInfo.Scan)
	if err != nil {
		return nil
	}
	return nil
}
