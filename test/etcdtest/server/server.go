package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/xiaoma03xf/sharddoc/server/etcd"
	"github.com/xiaoma03xf/sharddoc/test/etcdtest/pb"
	"google.golang.org/grpc"
)

type HelloService struct {
	pb.UnimplementedHelloServeServer
	server   *grpc.Server
	registry *etcd.ServiceRegistry
	addr     string
}

func NewHelloService(etcdEndpoints []string, clusterID, addr string, ttl int64) (*HelloService, error) {
	// 创建 etcd 注册器
	registry, err := etcd.NewServiceRegistry(etcdEndpoints, clusterID, addr, ttl)
	if err != nil {
		return nil, err
	}

	return &HelloService{
		addr:     addr,
		registry: registry,
	}, nil
}

// 启动 gRPC 服务并注册到 etcd
func (s *HelloService) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	pb.RegisterHelloServeServer(s.server, s)

	if err := s.registry.Register(); err != nil {
		return err
	}

	// 捕获退出信号，优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("收到退出信号")
		s.Close()
	}()

	log.Printf("服务启动: %s", s.addr)
	return s.server.Serve(lis)
}

func (s *HelloService) Close() {
	log.Println("正在关闭服务...")
	if s.registry != nil {
		s.registry.Deregister()
		s.registry.Close()
	}
	if s.server != nil {
		s.server.GracefulStop()
	}
	log.Println("服务关闭完成")
}

func (s *HelloService) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("收到请求: %s", req.Name)
	return &pb.HelloResponse{Request: "你好, " + req.Name}, nil
}

func main() {
	// 初始化服务
	service, err := NewHelloService(
		[]string{"118.89.66.104:2379"}, // etcd 节点地址
		"demo-cluster",                 // etcd clusterID
		"127.0.0.1:50051",              // gRPC 服务监听地址
		5,                              // TTL 秒数
	)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	// 启动服务
	if err := service.Start(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
