package main

import (
	"context"
	"log"
	"time"

	"github.com/xiaoma03xf/sharddoc/server/etcd"
	"github.com/xiaoma03xf/sharddoc/test/etcdtest/pb"
	"google.golang.org/grpc"
)

func main() {
	endpoints := []string{"118.89.66.104:2379"}
	clusterID := "demo-cluster"

	sd, err := etcd.NewServiceDiscovery(endpoints, clusterID)
	if err != nil {
		log.Fatalf("创建服务发现失败: %v", err)
	}
	defer sd.Close()

	if err := sd.Start(); err != nil {
		log.Fatalf("启动服务发现失败: %v", err)
	}

	time.Sleep(2 * time.Second) // 等待发现服务

	service := sd.GetService()
	if service == nil {
		log.Fatal("未发现任何可用服务")
	}

	log.Printf("连接到服务: %s", service.Addr)

	conn, err := grpc.Dial(service.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServeClient(conn)

	resp, err := client.Hello(context.Background(), &pb.HelloRequest{Name: "DengJun"})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	log.Printf("响应: %s", resp.Request)
}
