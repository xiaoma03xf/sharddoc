package tcp

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xiaoma03xf/sharddoc/discarded/storage"
	"github.com/xiaoma03xf/sharddoc/lib/utils"
)

func TestBuildAndReadTcpInfo(t *testing.T) {
	// 准备测试数据
	execMap := map[string]interface{}{
		"sql": "SELECT * FROM users",
	}
	joinMap := map[string]interface{}{
		"node_id": "node123",
		"addr":    "127.0.0.1:8080",
	}
	// 构造消息
	execMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeExec,
		Payload:  execMap,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo exec error: %v", err)
	}
	joinMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeJoin,
		Payload:  joinMap,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo join error: %v", err)
	}
	statusMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeStatus,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo status error: %v", err)
	}

	// 创建本地 TCP 监听
	ln, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		t.Fatalf("net.Listen error: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Logf("accept error: %v", err)
			return
		}
		defer conn.Close()

		for i := 0; i < 3; i++ {
			_, raftReq, err := ReadRequest(conn)
			if err != nil {
				t.Errorf("ReadRequest error: %v", err)
				return
			}

			fmt.Printf("Received msgType: 0x%02x, payload: %s\n", raftReq.DataType, raftReq.Payload)
		}
	}()

	// 等待服务启动
	time.Sleep(100 * time.Millisecond)

	// 客户端连接并发送消息
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	conn.Write(execMsg)
	conn.Write(joinMsg)
	conn.Write(statusMsg)
}

func TestLoadConfig(t *testing.T) {
	nodeCfg, err := LoadNodeConfig("../node1.yaml")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(nodeCfg)
}

func TestClientTryJoin(t *testing.T) {
	cfgPath := "../node2.yaml"
	BootstrapCluster(cfgPath)
}

func TestCreateTable(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	execMap := map[string]interface{}{
		"sql": `
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`,
	}
	execMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeExec,
		Payload:  execMap,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo exec error: %v", err)
	}
	conn.Write(execMsg)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("join response read error:", err)
	} else {
		fmt.Println("join response:", string(buf[:n]))
	}
}

func bootCluster(t *testing.T) {
	// 默认测试时, node1 节点为leader
	go BootstrapCluster("../node1.yaml")
	time.Sleep(3 * time.Second)
	go BootstrapCluster("../node2.yaml")
	time.Sleep(1 * time.Second)
	go BootstrapCluster("../node3.yaml")
	time.Sleep(1 * time.Second)

	// 确保集群正常启动
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		panic(err)
	}
	statusMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeStatus,
	})
	if err != nil {
		panic(err)
	}
	conn.Write(statusMsg)
	resp, err := ReadResponse(conn)
	if err != nil {
		panic(err)
	}
	var status StoreStatus
	if err = json.Unmarshal(resp.Body, &status); err != nil {
		panic(err)
	}
	// if len(status.Followers)
	conn.Close()
	assert.Equal(t, status.Leader.ID, "node1")
	assert.Equal(t, len(status.Followers), 2)
	if status.Followers[0].ID == "node2" {
		assert.Equal(t, status.Followers[1].ID, "node3")
	} else {
		assert.Equal(t, status.Followers[1].ID, "node2")
	}
}

func TestSelectNodeStatus(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	statusMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeStatus,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo status error: %v", err)
	}
	conn.Write(statusMsg)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[5:])
	if err != nil {
		fmt.Println("join response read error:", err)
	} else {
		fmt.Println("join response:", string(buf[:n]))
	}
}
func TestCteatables(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
	}()
	bootCluster(t)
	// 默认node1 为leader节点
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	{
		execMap := map[string]interface{}{
			"sql": `
	CREATE TABLE users (
	    id INT64,
	    name BYTES,
	    age INT64,
		height INT64,
		PRIMARY KEY (id),
		INDEX (name),
	    INDEX (age, height)
	);`,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)

		resp, err := ReadResponse(conn)
		if err != nil {
			t.Error(err)
		}
		// raft apply returned nil 130
		fmt.Println(string(resp.Body), resp.Type)
	}
	{
		execMap := map[string]interface{}{
			"sql": `
	CREATE TABLE students (
	        id INT64,
	        age INT64,
			lol INT64,
			PRIMARY KEY (id),
			INDEX (lol)
	    );
		`,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)

		resp, err := ReadResponse(conn)
		if err != nil {
			t.Error(err)
		}
		// raft apply returned nil 130
		fmt.Println(string(resp.Body), resp.Type)
	}

	// 检查各个节点创建表的情况
	time.Sleep(1 * time.Second)
	tablesMsg, err := BuildTcpInfo(&RaftRequest{
		RequestID: uuid.New().String(),
		DataType:  TypeShowTbl,
	})
	if err != nil {
		t.Fatalf("BuildTcpInfo status error: %v", err)
	}
	conn.Write(tablesMsg)
	resp, err := ReadResponse(conn)
	if err != nil {
		t.Error(err)
	}
	//create table server reply
	// [{"Name":"users","Types":[2,1,2,2],"Cols":["id","name","age","height"],"Indexes":[["id"],["age","height","id"]],"Prefixes":[100,101]}]
	// fmt.Println("create table server reply", string(resp.Body))

	var tables []storage.TableDef
	if err = json.Unmarshal(resp.Body, &tables); err != nil {
		t.Error(err)
	}
	fmt.Println("select tables reply", tables)

	time.Sleep(2 * time.Second)
	// 检查是否实现意义上的分布式
	node1db := storage.DB{Path: "../clusterdb/data/node1.db"}
	node2db := storage.DB{Path: "../clusterdb/data/node2.db"}
	node3db := storage.DB{Path: "../clusterdb/data/node3.db"}
	_ = node1db.Open()
	_ = node2db.Open()
	_ = node3db.Open()
	t1, err1 := node1db.GetAllTables()
	t2, err2 := node2db.GetAllTables()
	t3, err3 := node3db.GetAllTables()
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal(err)
	}
	assert.Equal(t, t1, t2)
	assert.Equal(t, t1, t3)
}

func PrintlnRec(recs []storage.Record) {
	// 打印表头
	fmt.Printf("%-10s %-20s %-10s %-10s\n", "ID", "Name", "Age", "Height")
	fmt.Println("---------------------------------------------------------")
	// 打印每一条记录
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
		// 输出数据
		fmt.Printf("%-10d %-20s %-10d %-10d\n", id, name, age, height)
	}
}

// 测试增删查改
func TestRaftDB(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
		os.Remove("./test_data.json")
		os.Remove("./testdb")
	}()
	bootCluster(t)
	// 默认node1 为leader节点
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	{
		execMap := map[string]interface{}{
			"sql": `
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)

		resp, err := ReadResponse(conn)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(resp.Body))
	}

	// 持久化测试数据,便于观察
	filePath := "test_data.json"
	datacnt := 200
	recs, err := utils.GenerateData(filePath, datacnt)
	if err != nil {
		t.Error(err)
		return
	}
	// 并发插入，测试时序竞态功能
	// 多个goroutine需要使用独立的连接,不然会产生竞态条件(数据乱序,甚至阻塞  )
	fmt.Println("测试数据长度", len(recs))
	var wg sync.WaitGroup
	for _, rec_data := range recs {
		wg.Add(1)
		go func(rec_data utils.RecordTestData) {
			defer wg.Done()

			// 每个 goroutine 单独建立 TCP 连接
			conn, err := net.Dial("tcp", "127.0.0.1:29001")
			if err != nil {
				t.Errorf("dial error: %v", err)
				return
			}
			defer conn.Close()

			// 构建插入 SQL
			cols := []string{"id", "name", "age", "height"}
			args := []interface{}{rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height}
			execMap := map[string]interface{}{
				"sql": storage.BuildInsertSQL("users", cols, args),
			}
			execMsg, err := BuildTcpInfo(&RaftRequest{
				RequestID: uuid.New().String(),
				DataType:  TypeExec,
				Payload:   execMap,
			})
			if err != nil {
				t.Errorf("BuildTcpInfo error: %v", err)
				return
			}
			if _, err := conn.Write(execMsg); err != nil {
				t.Errorf("conn.Write error: %v", err)
				return
			}
			resp, err := ReadResponse(conn)
			if err != nil {
				t.Errorf("ReadResponse error: %v", err)
				return
			}
			if resp.Type == TypeBadResp {
				t.Logf("Insert error response: %s", string(resp.Body))
				return
			}
			t.Logf("Insert OK response: %s", string(resp.Body))
			// raft apply returned nil 130
			// assert.Equal(t, string(resp.Body), "OK")
			t.Log("insert response:", string(resp.Body))
		}(rec_data)
	}
	wg.Wait()
	time.Sleep(1 * time.Second)

	// 测试非分布式数据对比
	dbpath := "./testdb"
	db := storage.DB{Path: dbpath}
	_ = db.Open()
	{
		tdef := &storage.TableDef{
			Name:  "users",
			Cols:  []string{"id", "name", "age", "height"},
			Types: []uint32{storage.TYPE_INT64, storage.TYPE_BYTES, storage.TYPE_INT64, storage.TYPE_INT64},
			Indexes: [][]string{
				{"id"},            // 主键索引
				{"age", "height"}, // 二级索引（复合索引）
			},
		}
		tx := storage.DBTX{}
		db.Begin(&tx)
		tx.TableNew(tdef)
		db.Commit(&tx)

		// 添加测试数据
		for _, rec_data := range recs {
			cols := []string{"id", "name", "age", "height"}
			args := []interface{}{rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height}
			sql := storage.BuildInsertSQL("users", cols, args)

			tx := storage.DBTX{}
			db.Begin(&tx)
			if err := db.Exec(sql); err != nil {
				t.Error(err)
				return
			}
			_ = db.Commit(&tx)
		}
		res4 := db.Raw("SELECT name, id FROM users WHERE age >= 0")
		fmt.Println(len(res4.Recs))
	}
	time.Sleep(1 * time.Second)

	// 测试简单的更新, 注意要使用索引
	var updateUserFirst storage.QueryResult
	var updateUserAfter storage.QueryResult
	{
		selectsql2 := "SELECT name, id FROM users WHERE age >= 0"
		execMap := map[string]interface{}{
			"sql": selectsql2,
		}
		execMsg, _ := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.NewString(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		conn.Write(execMsg)
		resp, _ := ReadResponse(conn)
		if resp.Type == TypeBadResp {
			t.Log("update response err:", string(resp.Body))
			return
		}
		_ = json.Unmarshal(resp.Body, &updateUserFirst)
		fmt.Println("before update...")
		PrintlnRec(updateUserFirst.Recs)

		// 执行更新
		updatasql := "UPDATE users SET age=10086 WHERE age>22"
		execMap = map[string]interface{}{
			"sql": updatasql,
		}
		execMsg, err = BuildTcpInfo(&RaftRequest{
			RequestID: uuid.NewString(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)

		resp, _ = ReadResponse(conn)
		if resp.Type == TypeBadResp {
			t.Log("update response err:", string(resp.Body))
			return
		}

		if err = db.Exec(updatasql); err != nil {
			t.Error(err)
			return
		}

		// 测试更新后的数据
		selectsql2 = "SELECT name, id FROM users WHERE age >= 0"
		execMap = map[string]interface{}{
			"sql": selectsql2,
		}
		execMsg, err = BuildTcpInfo(&RaftRequest{
			RequestID: uuid.NewString(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)
		resp, _ = ReadResponse(conn)
		if resp.Type == TypeBadResp {
			t.Log("update response err:", string(resp.Body))
			return
		}
		_ = json.Unmarshal(resp.Body, &updateUserAfter)
		fmt.Println("after update...")
		PrintlnRec(updateUserAfter.Recs)
	}
	time.Sleep(1 * time.Second)
	// 测试删除
	var deleteAfter storage.QueryResult
	{
		delsql := "DELETE FROM users WHERE age=10086 AND height > 175"
		execMap := map[string]interface{}{
			"sql": delsql,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.NewString(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
			return
		}
		conn.Write(execMsg)
		resp, _ := ReadResponse(conn)
		if resp.Type == TypeBadResp {
			t.Log("update response err:", string(resp.Body))
			return
		}

		// 测试更新后的数据
		selectsql2 := "SELECT name, id FROM users WHERE age >= 0"
		execMap = map[string]interface{}{
			"sql": selectsql2,
		}
		execMsg, err = BuildTcpInfo(&RaftRequest{
			RequestID: uuid.NewString(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			t.Fatalf("BuildTcpInfo exec error: %v", err)
		}
		conn.Write(execMsg)
		resp, _ = ReadResponse(conn)
		if resp.Type == TypeBadResp {
			t.Log("update response err:", string(resp.Body))
			return
		}
		_ = json.Unmarshal(resp.Body, &deleteAfter)
		fmt.Println("len of result", len(deleteAfter.Recs))
		fmt.Println("after delete...")
		PrintlnRec(deleteAfter.Recs)
	}
}
func basicBootCluster() error {
	go BootstrapCluster("../node1.yaml")
	time.Sleep(3 * time.Second)
	go BootstrapCluster("../node2.yaml")
	time.Sleep(1 * time.Second)
	go BootstrapCluster("../node3.yaml")
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		return err
	}
	defer conn.Close()

	statusMsg, err := BuildTcpInfo(&RaftRequest{
		DataType: TypeStatus,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(statusMsg)
	if err != nil {
		return err
	}

	resp, err := ReadResponse(conn)
	if err != nil {
		return err
	}

	var status StoreStatus
	if err := json.Unmarshal(resp.Body, &status); err != nil {
		return err
	}

	if status.Leader.ID != "node1" {
		return fmt.Errorf("leader is not node1")
	}
	if len(status.Followers) != 2 {
		return fmt.Errorf("followers count != 2")
	}

	if !(status.Followers[0].ID == "node2" || status.Followers[1].ID == "node2") {
		return fmt.Errorf("followers do not contain node2")
	}
	if !(status.Followers[0].ID == "node3" || status.Followers[1].ID == "node3") {
		return fmt.Errorf("followers do not contain node3")
	}
	return nil
}

func BenchmarkRaftDB_Insert(b *testing.B) {
	defer func() {
		os.RemoveAll("../clusterdb")
		os.Remove("./test_data.json")
		os.Remove("./testdb")
	}()
	if err := basicBootCluster(); err != nil {
		panic(err)
	}

	b.ReportAllocs()

	// 先生成测试数据
	filePath := "bench_data.json"
	datacnt := 10000 // 测试量，按需调整
	recs, err := utils.GenerateData(filePath, datacnt)
	if err != nil {
		b.Fatalf("generate data error: %v", err)
	}

	{
		conn, err := net.Dial("tcp", "127.0.0.1:29001")
		if err != nil {
			b.Errorf("dial error: %v", err)
			return
		}

		execMap := map[string]interface{}{
			"sql": `
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			panic(err)
		}
		conn.Write(execMsg)

		resp, err := ReadResponse(conn)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resp.Body))
		conn.Close()
	}

	// 多协程并发插入
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 每个 goroutine 单独连接
			conn, err := net.Dial("tcp", "127.0.0.1:29001")
			if err != nil {
				b.Errorf("dial error: %v", err)
				return
			}

			// 随机取一条数据插入
			rec := recs[rand.Intn(len(recs))]
			cols := []string{"id", "name", "age", "height"}
			args := []interface{}{rec.ID, rec.Name, rec.Age, rec.Height}
			execMap := map[string]interface{}{
				"sql": storage.BuildInsertSQL("users", cols, args),
			}
			execMsg, err := BuildTcpInfo(&RaftRequest{
				RequestID: uuid.New().String(),
				DataType:  TypeExec,
				Payload:   execMap,
			})
			if err != nil {
				b.Errorf("BuildTcpInfo error: %v", err)
				conn.Close()
				return
			}
			if _, err := conn.Write(execMsg); err != nil {
				b.Errorf("conn.Write error: %v", err)
				conn.Close()
				return
			}
			// 读取响应，忽略返回体，只要不报错即可
			_, err = ReadResponse(conn)
			if err != nil {
				b.Errorf("ReadResponse error: %v", err)
			}
			conn.Close()
		}
	})
}

func BenchmarkRaftDB_Query(b *testing.B) {
	b.ReportAllocs()

	// 连接复用测试，单连接多次查询
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		b.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		querySQL := "SELECT name, id FROM users WHERE age >= 0 LIMIT 100"
		execMap := map[string]interface{}{
			"sql": querySQL,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			b.Fatalf("BuildTcpInfo error: %v", err)
		}
		if _, err := conn.Write(execMsg); err != nil {
			b.Fatalf("conn.Write error: %v", err)
		}
		_, err = ReadResponse(conn)
		if err != nil {
			b.Fatalf("ReadResponse error: %v", err)
		}
	}
}

func BenchmarkRaftDB_Update(b *testing.B) {
	b.ReportAllocs()

	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		b.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updateSQL := "UPDATE users SET age=100 WHERE age > 20 LIMIT 10"
		execMap := map[string]interface{}{
			"sql": updateSQL,
		}
		execMsg, err := BuildTcpInfo(&RaftRequest{
			RequestID: uuid.New().String(),
			DataType:  TypeExec,
			Payload:   execMap,
		})
		if err != nil {
			b.Fatalf("BuildTcpInfo error: %v", err)
		}
		if _, err := conn.Write(execMsg); err != nil {
			b.Fatalf("conn.Write error: %v", err)
		}
		_, err = ReadResponse(conn)
		if err != nil {
			b.Fatalf("ReadResponse error: %v", err)
		}
	}
}
