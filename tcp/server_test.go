package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/xiaoma03xf/sharddoc/cluster"
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
	statusMap := map[string]interface{}{}

	// 构造消息
	execMsg, err := BuildTcpInfo(TypeExec, execMap)
	if err != nil {
		t.Fatalf("BuildTcpInfo exec error: %v", err)
	}
	joinMsg, err := BuildTcpInfo(TypeJoin, joinMap)
	if err != nil {
		t.Fatalf("BuildTcpInfo join error: %v", err)
	}
	statusMsg, err := BuildTcpInfo(TypeStatus, statusMap)
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
			msgType, payload, err := ReadRequest(conn)
			if err != nil {
				t.Errorf("ReadRequest error: %v", err)
				return
			}

			fmt.Printf("Received msgType: 0x%02x, payload: %s\n", msgType, payload)
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
	nodeCfg, err := cluster.LoadNodeConfig("../configs/node1.yaml")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(nodeCfg)
}

func TestClientTryJoin(t *testing.T) {
	cfgPath := "../configs/node2.yaml"
	BootstrapCluster(cfgPath)
}

func TestSelectNodeStatus(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	statusMap := map[string]interface{}{}
	statusMsg, err := BuildTcpInfo(TypeStatus, statusMap)
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
	select {}
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
	execMsg, err := BuildTcpInfo(TypeExec, execMap)
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

func bootCluster() {
	go BootstrapCluster("../node1.yaml")
	time.Sleep(1 * time.Second)
	go BootstrapCluster("../node2.yaml")
	time.Sleep(1 * time.Second)
	go BootstrapCluster("../node3.yaml")
	time.Sleep(1 * time.Second)

	// 确保集群正常启动
	conn, err := net.Dial("tcp", "127.0.0.1:29001")
	if err != nil {
		panic(err)
	}
	statusMap := map[string]interface{}{}
	statusMsg, err := BuildTcpInfo(TypeStatus, statusMap)
	if err != nil {
		panic(err)
	}
	conn.Write(statusMsg)

	resp, err := ReadResponse(conn)
	if err != nil {
		panic(err)
	}

	var status cluster.StoreStatus
	if err = json.Unmarshal(resp.Body, &status); err != nil {
		panic(err)
	}
	conn.Close()
	fmt.Println(status)
}
func TestCteatables(t *testing.T) {
	defer func() {
		os.RemoveAll("./clusterdb")
	}()
	bootCluster()
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
	execMsg, err := BuildTcpInfo(TypeExec, execMap)
	if err != nil {
		t.Fatalf("BuildTcpInfo exec error: %v", err)
	}
	conn.Write(execMsg)

	resp, err := ReadResponse(conn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(resp.Body), resp.Type)

	// 检查各个节点创建表的情况
	tableMap := make(map[string]interface{})
	tablesMsg, err := BuildTcpInfo(TypeShowTbl, tableMap)
	if err != nil {
		t.Fatalf("BuildTcpInfo status error: %v", err)
	}
	conn.Write(tablesMsg)
	resp, err = ReadResponse(conn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("create table server reply", string(resp.Body))

	var tables []storage.TableDef
	if err = json.Unmarshal(resp.Body, &tableMap); err != nil {
		t.Error(err)
	}
	fmt.Println("select tables reply", tables)

	select {}
}
