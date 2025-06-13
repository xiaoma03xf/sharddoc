package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xiaoma03xf/sharddoc/storage"
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
	fmt.Println("create table server reply", string(resp.Body))

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
	fmt.Println(t1, t2, t3)
	assert.Equal(t, t1, t2)
	assert.Equal(t, t1, t3)
}

func TestInsertData(t *testing.T) {
	// defer func() {
	// 	os.RemoveAll("../clusterdb")
	// }()
	// bootCluster(t)
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

}
