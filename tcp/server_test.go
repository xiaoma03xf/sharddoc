package tcp

import (
	"fmt"
	"net"
	"testing"
	"time"
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
