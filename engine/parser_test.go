package engine

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itxiaoma0610/sharddoc/engine/parser/ast"
)

func GenerateTestDB() *R {
	r := NewR()

	tdef := &TableDef{
		Name:  "tbl_test",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{TYPE_INT64, TYPE_BYTES, TYPE_INT64, TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
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
		return nil
	}

	for _, rec_data := range records {
		rec := record(rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height)
		added := r.add("tbl_test", rec)
		assert(added)
	}
	fmt.Printf("Added test records successfully\n")
	return r
}

func TestCreateTable(t *testing.T) {
	creatTable := `
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`
	tableDef := VisitTree(creatTable).(*TableDef)
	fmt.Println(tableDef)
}
func TestInsertTable(t *testing.T) {
	InsertSql := `
	INSERT INTO users (id, name, age, height) 
	VALUES (1, 'Alice', 30, 170);
	`
	input := antlr.NewInputStream(InsertSql)
	lexer := ast.NewSQLLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := ast.NewSQLParser(tokenStream)
	p.AddErrorListener(antlr.NewDefaultErrorListener())

	tree := p.Sql()
	if tree == nil {
		t.Errorf("parse tree is nil, check input or parser configuration")
	}
	v := new(SQLParser)
	tableDef := v.Visit(tree).(*Record)
	for i := 0; i < len(tableDef.Vals); i++ {
		if tableDef.Vals[i].Type == TYPE_INT64 {
			fmt.Println("col:", tableDef.Cols[i], "val: ", tableDef.Vals[i].I64)
		} else {
			fmt.Println("col:", tableDef.Cols[i], "val: ", string(tableDef.Vals[i].Str))
		}
	}
}

func TestSelectTable(t *testing.T) {
	r := GenerateTestDB()
	defer r.dispose()

	// SELECT name, id FROM tbl_test WHERE age > 25;
	tx := r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age >= 25;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)

		var ageindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 >= 25)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	// SELECT name, id FROM tbl_test WHERE age <= 25;
	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age <= 25;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)

		var ageindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 <= 25)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age BETWEEN 18 AND 25;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)
		var ageindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 >= 18 && rec.Vals[ageindex].I64 <= 25)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age > 18 AND age < 25;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)
		var ageindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 > 18 && rec.Vals[ageindex].I64 < 25)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age = 18 AND height < 175;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)
		var ageindex int
		var heightindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 < 175)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age = 18 AND height >= 175;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)
		var ageindex int
		var heightindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 >= 175)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)

	tx = r.begin()
	{
		timeDu := time.Now()
		sql := `
	SELECT name, id FROM tbl_test WHERE age = 18 AND height BETWEEN 170 AND 175;
	`
		selectInfo := VisitTree(sql).(*SelectInfo)
		err := tx.Scan(selectInfo.TableName, selectInfo.Scan)
		assert(err == nil)
		gotRec := reduceSelectData(selectInfo.Scan)
		var ageindex int
		var heightindex int
		for i, v := range gotRec[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		for _, rec := range gotRec {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 <= 175 &&
				rec.Vals[heightindex].I64 >= 170)
		}
		fmt.Println(sql, time.Since(timeDu))
	}
	r.commit(tx)
}

type User struct {
	ID  int
	Age int
}

func TestRef(t *testing.T) {
	records := []Record{
		{
			Cols: []string{"id", "age"},
			Vals: []Value{
				{Type: TYPE_INT64, I64: 1},
				{Type: TYPE_INT64, I64: 20},
			},
		},
		{
			Cols: []string{"id", "age"},
			Vals: []Value{
				{Type: TYPE_INT64, I64: 2},
				{Type: TYPE_INT64, I64: 25},
			},
		},
	}

	var ptrs []any
	for i := 0; i < len(records); i++ {
		ptrs = append(ptrs, &User{})
	}

	err := scanRecordsToStructs(records, ptrs)
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}

	for _, p := range ptrs {
		u := p.(*User)
		fmt.Printf("User{ID: %d, Age: %d}\n", u.ID, u.Age)
	}
}
