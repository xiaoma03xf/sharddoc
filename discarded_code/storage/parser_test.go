package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/xiaoma03xf/sharddoc/parser/ast"
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
	fmt.Println(r.db.tables["tbl_test"])

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
	// creatTable := `
	// CREATE TABLE users (
	//     id INT64,
	//     name BYTES,
	//     age INT64,
	// 	height INT64,
	// 	PRIMARY KEY (id),
	//     INDEX (age, height)
	// );
	// `
	creatTable := `
	CREATE TABLE users (
	    id INT64,
	    name BYTES,
	    age INT64,
		height INT64,
		PRIMARY KEY (id),
		INDEX (name),
	    INDEX (age, height)
	);`
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
	SELECT name, id FROM tbl_test WHERE age >= 35;
	`

		// 	type SelectInfo struct {
		// 	TableName   string
		// 	SelectField []string
		// 	Scan        *Scanner
		// }

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
func TestUtils(t *testing.T) {
	cols := []string{"id", "name", "age", "height"}
	args := []interface{}{1, "Alice", 20, 168}
	fmt.Println(BuildInsertSQL("users", cols, args))
}

func TestDBRaw(t *testing.T) {
	r := NewR()
	defer r.dispose()

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
	err := r.db.Exec(creatTable)
	if err != nil {
		t.Errorf("creat table err:%v", err)
	}

	// 插入100条数据(读取测试json中数据)
	type User struct {
		ID     int    `json:"ID"`
		Name   string `json:"Name"`
		Age    int    `json:"Age"`
		Height int    `json:"Height"`
	}
	var u []User
	data, err := os.ReadFile("./test_data.json") // 你的 JSON 文件路径
	if err != nil {
		log.Fatalf("failed to read json: %v", err)
	}
	err = json.Unmarshal(data, &u)
	if err != nil {
		log.Fatalf("failed to unmarshal: %v", err)
	}
	for _, user := range u {
		cols := []string{"id", "name", "age", "height"}
		args := []interface{}{user.ID, user.Name, user.Age, user.Height}

		sql := BuildInsertSQL("users", cols, args)

		err := r.db.Exec(sql)
		if err != nil {
			log.Printf("insert error on user %v: %v", user.ID, err)
		}
	}

	// 简单查询
	type SUser struct {
		ID     int64
		Name   string
		Age    int64
		Height int64
	}
	var users []SUser
	selectsql := "SELECT id, name, age FROM users WHERE age=18 AND height > 175"
	err = r.db.Raw(selectsql).Scan(&users)
	if err != nil {
		t.Errorf("select data err:%v", err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

	// 测试简单的更新
	updatasql := "UPDATE users SET age=10086 WHERE id>1995"
	err = r.db.Exec(updatasql)
	if err != nil {
		t.Errorf("update data err:%v", err)
	}
	var users2 []SUser
	selectsql2 := "SELECT id, name, age FROM users WHERE id > 1995"
	err = r.db.Raw(selectsql2).Scan(&users2)
	if err != nil {
		t.Errorf("select data err:%v", err)
	}
	for _, user := range users2 {
		fmt.Println(user)
	}

	// 测试简单的删除
	delsql := "DELETE FROM users WHERE age=18 AND height > 180"
	err = r.db.Exec(delsql)
	if err != nil {
		t.Errorf("del data err:%v", err)
	}
	var users3 []SUser
	selectsql3 := "SELECT id, name, age FROM users WHERE age=18 AND height > 175"
	err = r.db.Raw(selectsql3).Scan(&users3)
	if err != nil {
		t.Errorf("select data err:%v", err)
	}
	for _, user := range users3 {
		fmt.Println(user)
	}
}
