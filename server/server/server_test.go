package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"testing"
	"time"

	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft"

	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

func TestEtcdBasic(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
	}()
	conf := "../cluster1.yaml"

	go raft.BootstrapCluster(conf, "node1")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node2")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node3")
	time.Sleep(3 * time.Second)

	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		client, _, err := db.getLeader("cluster1")
		if err != nil {
			t.Error(err)
		}
		resp, err := client.Status(context.Background(), &raftpb.StatusRequest{})
		if err != nil {
			t.Error(err)
		}
		fmt.Println(resp.Me)
		time.Sleep(1 * time.Second)
	}

	// default CachedTimeOut = 10*time.Second
	time.Sleep(10 * time.Second)
	client, _, err := db.getLeader("cluster1")
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Status(context.Background(), &raftpb.StatusRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Me)
}

func TestTableNew(t *testing.T) {
	ClearEtcd([]string{"118.89.66.104:2379"})

	defer func() {
		os.RemoveAll("../../clusterdb")
	}()
	conf := "../../cluster1.yaml"

	go raft.BootstrapCluster(conf, "node1")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node2")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node3")
	time.Sleep(3 * time.Second)

	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}

	tdef1 := &kv.TableDef{
		Name:  "TableNew",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	meta1, err := db.TablesDefDiscovery.GetMetaKey(NEXT_PREFIX)
	if err != nil {
		t.Error(err)
	}
	if meta1 != nil {
		fmt.Println("meta before", bytesToUint32LE(meta1))
	}
	err = db.TableNew(tdef1)
	if err != nil {
		t.Error(err)
	}
	meta2, err := db.TablesDefDiscovery.GetMetaKey(NEXT_PREFIX)
	if err != nil {
		t.Error(err)
	}
	if meta2 != nil {
		fmt.Println("meta after", bytesToUint32LE(meta2))
	}
}

func BootCluster(conf string, nodes []string) {
	for _, node := range nodes {
		go raft.BootstrapCluster(conf, node)
		time.Sleep(3 * time.Second)
	}
}

type RecordTestData struct {
	ID     int64
	Name   string
	Age    int64
	Height int64
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
func ClearEtcd(endpoints []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create etcd client: %w", err))
	}
	_, err = client.Delete(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
}
func DBServerPrintQuery(recs []*raftpb.Record) {
	// 打印表头
	fmt.Printf("%-10s %-20s %-10s %-10s\n", "ID", "Name", "Age", "Height")
	fmt.Println("---------------------------------------------------------")

	for _, r := range recs {
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
		fmt.Printf("%-10d %-20s %-10d %-10d\n", id, name, age, height)
	}
}
func TestDBBasic(t *testing.T) {
	// go BootCluster("../../cluster1.yaml", []string{"node1", "node2", "node3"})
	// go BootCluster("../../cluster2.yaml", []string{"node4", "node5", "node6"})
	// go BootCluster("../../cluster3.yaml", []string{"node7", "node8", "node9"})
	// time.Sleep(10 * time.Second)

	// ClearEtcd([]string{"118.89.66.104:2379"})

	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1", "cluster2", "cluster3"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)
	tdef := &kv.TableDef{
		Name:  "user",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	err = db.TableNew(tdef)
	if err != nil {
		t.Error(err)
	}

	// 插入数据
	record := func(id int64, name string, age int64, height int64) kv.Record {
		rec := kv.Record{}
		rec.AddInt64("id", id).AddStr("name", []byte(name))
		rec.AddInt64("age", age).AddInt64("height", height)
		return rec
	}
	fmt.Printf("Adding test records...\n")

	// 持久化测试数据,便于观察
	filePath := "./test_data.json"
	records, err := ReadTestDataFromFile(filePath)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	for _, rec_data := range records {
		rec := record(rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height)
		ok, err := db.Insert("user", rec)
		if !ok || err != nil {
			t.Error(err)
			return
		}
	}
	fmt.Printf("Added test records successfully\n")

	// 测试范围查询
	// test age == 18
	{
		fmt.Println("select * from table where age == 18")
		rec := kv.Record{}
		rec.AddInt64("age", 18)

		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec, Key2: rec,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}

	{
		fmt.Println("select * from table where age >= 43")
		rec1 := kv.Record{}
		rec1.AddInt64("age", 43)

		rec2 := kv.Record{}
		rec2.AddInt64("age", math.MaxInt64/2)
		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec1, Key2: rec2,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}

	{
		const MIN_NAME = ""
		const MAX_NAME = "\xff\xff\xff\xff\xff\xff\xff\xff" // 足够长的最大字节

		fmt.Println("select * from table where name > yang")
		rec := kv.Record{}
		rec.AddStr("name", []byte("Yang"))

		rec2 := kv.Record{}
		rec2.AddStr("name", []byte(MAX_NAME))
		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec, Key2: rec2,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}

	{
		fmt.Println("select * from table where age == 18 and height == 174")
		rec := kv.Record{}
		rec.AddInt64("age", 18).AddInt64("height", 174)
		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec, Key2: rec,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}

	{
		fmt.Println("select * from table where age == 18 and height between 170 and 175 ")
		rec := kv.Record{}
		rec.AddInt64("age", 18).AddInt64("height", 170)

		rec2 := kv.Record{}
		rec2.AddInt64("age", 18).AddInt64("height", 175)

		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec, Key2: rec2,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}

	{
		fmt.Println("select * from table where age == 18 and height < 175 ")
		rec := kv.Record{}
		rec.AddInt64("age", 18).AddInt64("height", 0)

		rec2 := kv.Record{}
		rec2.AddInt64("age", 18).AddInt64("height", 175)

		req := Scanner{
			Cmp1: kv.CMP_GE, Cmp2: kv.CMP_LE,
			Key1: rec, Key2: rec2,
		}
		recs, err := db.Scan("user", &req)
		assert(err == nil)
		DBServerPrintQuery(recs)
	}
}

func TestEtcdKey(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"118.89.66.104:2379"}, // 替换为你的 etcd 地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("连接 etcd 失败: %v", err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取所有键值对
	resp, err := cli.Get(ctx, "", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("查询键值对失败: %v", err)
	}

	// 打印所有键值对
	for _, kv := range resp.Kvs {
		fmt.Printf("🔑 %s = %s\n", kv.Key, kv.Value)
	}
}

func TestSqlExec(t *testing.T) {
	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1", "cluster2", "cluster3"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)

	{
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
		b, err := db.ExecSQL(creatTable)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(b))
	}

	{
		// insert data
		filePath := "./test_data.json"
		records, err := ReadTestDataFromFile(filePath)
		if err != nil {
			fmt.Println("Error reading data:", err)
			return
		}
		for _, rec_data := range records {
			cols := []string{"id", "name", "age", "height"}
			args := []interface{}{rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height}
			sql := BuildInsertSQL("users", cols, args)
			b, err := db.ExecSQL(sql)
			assert(err == nil && string(b) == INSERT_OK)
		}
	}

	// SELECT name, id FROM tbl_test WHERE age >= 25;
	{
		sql := `
	SELECT name, id FROM users WHERE age >= 25;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		// check
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 >= 25)
		}
	}

	// SELECT name, id FROM tbl_test WHERE age <= 25;
	{
		sql := `
	SELECT name, id FROM users WHERE age <= 25;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 <= 25)
		}
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age BETWEEN 18 AND 25;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 >= 18 && rec.Vals[ageindex].I64 <= 25)
		}
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age > 18 AND age < 25;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 > 18 && rec.Vals[ageindex].I64 < 25)
		}
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age = 18 AND height < 175;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex, heightindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 < 175)
		}
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age = 18 AND height >= 175;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex, heightindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 >= 175)
		}
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age = 18 AND height BETWEEN 170 AND 175;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex, heightindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 <= 175 &&
				rec.Vals[heightindex].I64 >= 170)
		}
	}

	db.Close()
	ClearEtcd([]string{"118.89.66.104:2379"})
}
func TestClearEtcd(t *testing.T) {
	ClearEtcd([]string{"118.89.66.104:2379"})
}

func TestBasicCRUD(t *testing.T) {
	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)

	{
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
		b, err := db.ExecSQL(creatTable)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(b))
	}

	{
		// insert data
		filePath := "./test_data.json"
		records, err := ReadTestDataFromFile(filePath)
		if err != nil {
			fmt.Println("Error reading data:", err)
			return
		}
		for _, rec_data := range records {
			cols := []string{"id", "name", "age", "height"}
			args := []interface{}{rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height}
			sql := BuildInsertSQL("users", cols, args)
			b, err := db.ExecSQL(sql)
			assert(err == nil && string(b) == INSERT_OK)
		}
	}

	{
		// 测试简单的更新
		updatasql := "UPDATE users SET age=10086 WHERE id>295"
		_, err := db.ExecSQL(updatasql)
		if err != nil {
			t.Error(err)
		}

		selectsql2 := "SELECT id, name, age FROM users WHERE id > 295"
		b, err := db.ExecSQL(selectsql2)
		if err != nil {
			t.Errorf("select data err:%v", err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		DBServerPrintQuery(resp.Records)
	}

	{
		sql := `
	SELECT name, id FROM users WHERE age = 18 AND height > 175;
	`
		b, err := db.ExecSQL(sql)
		if err != nil {
			t.Error(err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		var ageindex, heightindex int
		for i, v := range resp.Records[0].Cols {
			if v == "age" {
				ageindex = i
			}
			if v == "height" {
				heightindex = i
			}
		}
		fmt.Println(sql)
		DBServerPrintQuery(resp.Records)
		for _, rec := range resp.Records {
			assert(rec.Vals[ageindex].I64 == 18 && rec.Vals[heightindex].I64 > 175)
		}
	}

	{
		// 测试简单的删除
		delsql := "DELETE FROM users WHERE age=18 AND height > 180"
		b, err := db.ExecSQL(delsql)
		if err != nil {
			t.Error(err)
		}
		assert(string(b) == DELETE_OK)

		selectsql2 := `
	SELECT name, id FROM users WHERE age = 18 AND height > 175;
	`
		b, err = db.ExecSQL(selectsql2)
		if err != nil {
			t.Errorf("select data err:%v", err)
		}
		var resp raftpb.ScanResponse
		err = proto.Unmarshal(b, &resp)
		if err != nil {
			t.Errorf("failed to unmarshal ScanResponse: %v", err)
		}
		DBServerPrintQuery(resp.Records)
	}
	db.Close()
	ClearEtcd([]string{"118.89.66.104:2379"})
}

func TestServerBasic(t *testing.T) {
	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)

	{
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
		b, err := db.ExecSQL(creatTable)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(b))
	}
	record := func(id int64, name string, age int64, height int64) kv.Record {
		rec := kv.Record{}
		rec.AddInt64("id", id).AddStr("name", []byte(name))
		rec.AddInt64("age", age).AddInt64("height", height)
		return rec
	}

	{
		rec := record(1, "jack", 23, 170)
		_, err := db.Insert("users", rec)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		rec := record(1, "jack", 23, 175)
		_, err := db.Update("users", rec)
		if err != nil {
			t.Fatal(err)
		}
	}
}
