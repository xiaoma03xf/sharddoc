package raft

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
)

const (
	retainSnapshotCount = 2
	defaultRaftTimeout  = 10 * time.Second
)

var _ pb.KVStoreServer = &Store{}

type Store struct {
	// grpc api
	pb.UnimplementedKVStoreServer
	mu sync.Mutex
	kv *kv.KV

	raftDir  string
	raftAddr string
	raft     *raft.Raft
	logger   *log.Logger

	// undo map
	applyWaiters map[string]chan struct{}
}
type NodeConfig struct {
	NodeID   string `yaml:"node_id"`
	RaftAddr string `yaml:"raft_addr"`
	RaftDir  string `yaml:"raft_dir"`

	KVPath    string `yaml:"kv_path"`
	Bootstrap bool   `yaml:"bootstrap"`
	JoinAddr  string `yaml:"join_addr,omitempty"`

	// grpc 暴露地址
	GrpcAddr string `yaml:"grpc_addr"`
	// common config
	RaftTimeout       time.Duration `yaml:"raft_timeout"`
	SnapshotInterval  time.Duration `yaml:"snapshot_interval"`
	SnapshotThreshold uint64        `yaml:"snapshot_threshold"`
}

func NewStore(cfg *NodeConfig) (*Store, error) {
	// 如果raft目录不存在就新建一个
	// 如果 db 文件夹不存在也新建一个
	if err := os.MkdirAll(cfg.RaftDir, 0755); err != nil {
		return nil, fmt.Errorf("create raft_dir failed: %w", err)
	}
	dir := filepath.Dir(cfg.KVPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	kv := kv.KV{Path: cfg.KVPath}
	if err := kv.Open(); err != nil {
		return nil, err
	}

	s := &Store{}
	s.kv = &kv
	s.raftDir = cfg.RaftDir
	s.raftAddr = cfg.RaftAddr
	s.logger = log.New(os.Stderr, "[store] ", log.LstdFlags)
	s.applyWaiters = make(map[string]chan struct{})

	grpcsignal := make(chan struct{})
	go s.grpcListenAndServe(cfg.GrpcAddr, grpcsignal)
	<-grpcsignal

	return s, nil
}

// TODO 为了强一致性, 暂时把所有的业务问题交给leader处理
// 后续可以把 Get, Scan 等业务交给 follower 节点处理
func raftLeaderInterceptor(s *Store) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		exemptMethods := map[string]struct{}{
			"/KVStore/Join":   {},
			"/KVStore/Status": {},
		}
		// 放行操作 Join, Status ...
		if _, exempt := exemptMethods[info.FullMethod]; exempt {
			return handler(ctx, req)
		}
		if s.raft.State() != raft.Leader {
			leaderAddr := s.raft.Leader()
			return nil, status.Errorf(codes.Unavailable, "not leader, please redirect to leader at %s", leaderAddr)
		}
		return handler(ctx, req)
	}
}

func (s *Store) grpcListenAndServe(grpcServeaddr string, grpcsignal chan<- struct{}) {
	listen, err := net.Listen("tcp", grpcServeaddr)
	if err != nil {
		panic(err)
	}
	s.logger.Printf("grpc server listening on %s", grpcServeaddr)
	close(grpcsignal)

	// 注册拦截器, 只让leader处理信息即可
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(raftLeaderInterceptor(s)),
	)
	pb.RegisterKVStoreServer(grpcServer, s)
	if err := grpcServer.Serve(listen); err != nil {
		panic(err)
	}
}

func (s *Store) Open(cfg *NodeConfig) error {
	// 创建 Raft 配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.NodeID)
	config.SnapshotInterval = cfg.SnapshotInterval   // 至少间隔 3 分钟检查一次是否需要快照
	config.SnapshotThreshold = cfg.SnapshotThreshold // 日志条目超过 100000 才创建快照

	// 设置网络传输, raftAddr 转成TCP地址结构
	// 使用 raft.NewTCPTransport 创建一个 TCP 传输通道，让 Raft 节点可以收发消息
	addr, err := net.ResolveTCPAddr("tcp", s.raftAddr)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.raftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	// 创建快照存储, Raft 会定期对当前 FSM 的状态进行快照，减少日志体积
	snapshots, err := raft.NewFileSnapshotStore(s.raftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return err
	}
	//logStore 用于持久化日志条目（比如命令日志）
	//stableStore 用于存储持久化元数据，比如当前的 Term、投票记录等
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(s.raftDir, "raft-log.db"))
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(s.raftDir, "raft-stable.db"))
	if err != nil {
		return err
	}

	rf, err := raft.NewRaft(config, s, logStore, stableStore, snapshots, transport)
	if err != nil {
		return err
	}
	s.raft = rf

	// 首节点直接启动集群
	if cfg.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		rf.BootstrapCluster(configuration)
	}
	return nil
}

func BootstrapCluster(cfgPath string) {
	nodeCfg, err := LoadNodeConfig(cfgPath)
	assert(err == nil)
	s, err := NewStore(nodeCfg)
	assert(err == nil)
	if err := s.Open(nodeCfg); err != nil {
		panic(err)
	}

	// If join was specified, make the join request.
	if nodeCfg.JoinAddr != "" {
		// 创建 发往 JoinAddr 的grpc Join请求
		client, conn, err := BuildGrpcConn(nodeCfg.JoinAddr)
		defer conn.Close()
		assert(err == nil)
		resp, err := client.Join(context.Background(), &pb.JoinRequest{
			NodeId:  nodeCfg.NodeID,
			Address: nodeCfg.RaftAddr,
		})
		if err != nil {
			fmt.Println(resp)
			fmt.Println(nodeCfg.JoinAddr)
			logger.Warn("failed to connect, err:", err)
			panic(err)
		}
	}
	// we're up and running!
	logger.Info(fmt.Sprintf("node started successfully, listening on grpc: %s", nodeCfg.GrpcAddr))
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	logger.Info("raft node exiting")
}

func assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}

func LoadNodeConfig(path string) (*NodeConfig, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	baseDir := filepath.Dir(absPath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg NodeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	cfg.RaftDir = resolvePath(baseDir, cfg.RaftDir)
	cfg.KVPath = resolvePath(baseDir, cfg.KVPath)

	return &cfg, nil
}

func resolvePath(baseDir, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(baseDir, p)
}

func BuildGrpcConn(addr string) (pb.KVStoreClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("faild to connect to grpc server:%v", err)
	}
	return pb.NewKVStoreClient(conn), conn, nil
}
