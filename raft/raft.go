package raft

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/grpc"
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
	NodeID    string `yaml:"node_id"`
	RaftAddr  string `yaml:"raft_addr"`
	RaftDir   string `yaml:"raft_dir"`
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
	err := kv.Open()
	assert(err == nil)

	s := &Store{}
	s.kv = &kv
	s.raftDir = cfg.RaftDir
	s.raftAddr = cfg.RaftAddr
	s.logger = log.New(os.Stderr, "[store] ", log.LstdFlags)
	s.applyWaiters = make(map[string]chan struct{})

	go s.grpcListenAndServe(cfg.GrpcAddr)
	return s, nil
}

func (s *Store) grpcListenAndServe(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
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

func assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}
