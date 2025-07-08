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
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"github.com/xiaoma03xf/sharddoc/server/etcd"

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

	raftDir      string
	raftAddr     string
	raft         *raft.Raft
	grpcaddr     string
	nodeID       string
	registry     *etcd.ServiceRegistry
	isRegistered int32 // atomic

	logger *log.Logger
	// undo map
	applyWaiters map[string]chan struct{}
}

func NewStore(c *CommonConfig, n *NodeConfig) (*Store, error) {
	// 如果raft目录不存在就新建一个
	// 如果 db 文件夹不存在也新建一个
	if err := os.MkdirAll(n.RaftDir, 0755); err != nil {
		return nil, fmt.Errorf("create raft_dir failed: %w", err)
	}
	dir := filepath.Dir(n.KVPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	kv := kv.KV{Path: n.KVPath, Snapshot: n.KVLogPath}
	if err := kv.Open(); err != nil {
		return nil, err
	}

	s := &Store{}
	s.kv = &kv
	s.raftDir = n.RaftDir
	s.raftAddr = n.RaftAddr
	s.grpcaddr = n.GrpcAddr
	s.nodeID = n.NodeID
	s.logger = log.New(os.Stderr, "[store] ", log.LstdFlags)
	s.applyWaiters = make(map[string]chan struct{})

	// 开启grpc服务，并定期检测注册或退出etcd
	grpcsignal := make(chan struct{})
	go s.grpcListenAndServe(n.GrpcAddr, grpcsignal)
	<-grpcsignal

	registry, err := etcd.NewServiceRegistry(c.Etcd.Endpoints, n.ClusterID, n.GrpcAddr, c.Etcd.TTL)
	if err != nil {
		return nil, err
	}
	s.registry = registry

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
func (s *Store) IsLeader() bool {
	return s.raft.State() == raft.Leader
}

// startLeaderWatcher 监听状态，同步etcd
func (s *Store) startLeaderWatcher() {
	ch := make(chan raft.Observation, 1)
	observer := raft.NewObserver(ch, false, nil)
	s.raft.RegisterObserver(observer)

	go func() {
		for obs := range ch {
			if leaderObs, ok := obs.Data.(raft.LeaderObservation); ok {
				if string(leaderObs.LeaderAddr) == s.raftAddr {
					if atomic.LoadInt32(&s.isRegistered) == 0 {
						err := s.registry.Register(s.IsLeader, func() {
							atomic.StoreInt32(&s.isRegistered, 0)
						})
						if err == nil {
							logger.Info(s.raftAddr, "注册到etcd成功")
							atomic.StoreInt32(&s.isRegistered, 1)
						} else {
							log.Printf("etcd 注册失败: %v", err)
						}
					}
				}
			}
		}
	}()
}

func (s *Store) Open(c *CommonConfig, n *NodeConfig) error {
	// 创建 Raft 配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(n.NodeID)
	config.SnapshotInterval = c.SnapshotInterval   // 至少间隔 3 分钟检查一次是否需要快照
	config.SnapshotThreshold = c.SnapshotThreshold // 日志条目超过 100000 才创建快照
	config.Logger = hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Off, // 完全关闭日志
	})

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
	if n.Bootstrap {
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

	go s.startLeaderWatcher()
	return nil
}

func BootstrapCluster(cfgPath string, nodeID string) {
	cfg, err := LoadNodeConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	var nodeCfg *NodeConfig
	for _, node := range cfg.Nodes {
		if node.NodeID == nodeID {
			nodeCfg = node
			break
		}
	}
	if nodeCfg == nil {
		panic(fmt.Errorf("未找到指定 NodeID=%s 的节点配置", nodeID))
	}

	s, err := NewStore(cfg.Common, nodeCfg)
	if err != nil {
		panic(err)
	}
	if err := s.Open(cfg.Common, nodeCfg); err != nil {
		panic(err)
	}

	// If join was specified, make the join request.
	if nodeCfg.JoinAddr != "" {
		// 创建 发往 JoinAddr 的grpc Join请求
		client, conn, err := BuildGrpcConn(nodeCfg.JoinAddr)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		resp, err := client.Join(context.Background(), &pb.JoinRequest{
			NodeId:  nodeCfg.NodeID,
			Address: nodeCfg.RaftAddr,
		})
		if err != nil || !resp.Success {
			logger.Warn("failed to connect, err:", err)
			panic(err)
		}
	}
	// we're up and running!
	logger.Info(fmt.Sprintf("node started successfully, listening on grpc: %s", nodeCfg.GrpcAddr))
	fmt.Println(s.kv.Path, s.kv.Snapshot)
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	logger.Info("raft node exiting")
}

func Assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}

func LoadNodeConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	baseDir := filepath.Dir(absPath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	for _, node := range cfg.Nodes {
		node.RaftDir = resolvePath(baseDir, node.RaftDir)
		node.KVPath = resolvePath(baseDir, node.KVPath)
		node.KVLogPath = resolvePath(baseDir, node.KVLogPath)
	}

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
