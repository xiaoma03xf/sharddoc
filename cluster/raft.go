package cluster

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/xiaoma03xf/sharddoc/storage"
)

const (
	retainSnapshotCount = 2
	defaultRaftTimeout  = 10 * time.Second
)

type Store struct {
	db       *storage.DB
	RaftDir  string
	RaftAddr string
	mu       sync.RWMutex
	raft     *raft.Raft
	logger   *log.Logger
}

func NewStore(cfg *NodeConfig) (*Store, error) {
	// 如果raft目录不存在就新建一个
	// 如果 db 文件夹不存在也新建一个
	if err := os.MkdirAll(cfg.RaftDir, 0755); err != nil {
		return nil, fmt.Errorf("create raft_dir failed: %w", err)
	}
	dir := filepath.Dir(cfg.DBPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	db := &storage.DB{Path: cfg.DBPath}
	if err := db.Open(); err != nil {
		return nil, err
	}
	s := &Store{}
	s.db = db
	s.RaftDir = cfg.RaftDir
	s.RaftAddr = cfg.RaftAddr
	s.logger = log.New(os.Stderr, "[store] ", log.LstdFlags)
	return s, nil
}
func (s *Store) Open(cfg *NodeConfig) error {
	// 创建 Raft 配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.NodeID)
	config.SnapshotInterval = cfg.SnapshotInterval   // 至少间隔 3 分钟检查一次是否需要快照
	config.SnapshotThreshold = cfg.SnapshotThreshold // 日志条目超过 100000 才创建快照

	// 设置网络传输, raftAddr 转成TCP地址结构
	// 使用 raft.NewTCPTransport 创建一个 TCP 传输通道，让 Raft 节点可以收发消息
	addr, err := net.ResolveTCPAddr("tcp", s.RaftAddr)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	// 创建快照存储, Raft 会定期对当前 FSM 的状态进行快照，减少日志体积
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return err
	}
	//logStore 用于持久化日志条目（比如命令日志）
	//stableStore 用于存储持久化元数据，比如当前的 Term、投票记录等
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft-log.db"))
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft-stable.db"))
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

type SQLCommand struct {
	IsSelect bool
	SQL      string `json:"sql"`
}
type ExecSQLRsp struct {
	Data []byte
	Err  error
}

func (s *Store) Exec(sql string) *ExecSQLRsp {
	// 测试阶段允许各节点执行
	// if s.raft.State() != raft.Leader {
	// 	return &ExecSQLRsp{Err: fmt.Errorf("not Leader")}
	// }
	c := &SQLCommand{
		IsSelect: false,
		SQL:      sql,
	}
	// 判断是否查询语句, 如果是修改cmd的IsSelect为true
	sqlcopy := sql
	if strings.HasPrefix(strings.ToLower(sqlcopy), "select") {
		c.IsSelect = true
	}
	b, err := json.Marshal(c)
	if err != nil {
		return &ExecSQLRsp{Err: err}
	}
	f := s.raft.Apply(b, defaultRaftTimeout)
	resp := f.Response()
	if resp == nil {
		return &ExecSQLRsp{Err: fmt.Errorf("raft apply returned nil")}
	}
	return resp.(*ExecSQLRsp)
}

func (s *Store) Join(nodeID, addr string) error {
	s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	// 获取当前集群配置
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}
	// 检查是否已存在相同ID 或地址的节点
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// 如果 ID 和 地址完全相同, 无需加入
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}
			// 移除冲突的节点
			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}
	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, defaultRaftTimeout)
	if f.Error() != nil {
		return f.Error()
	}
	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)

	return nil
}

func (s *Store) Status() (StoreStatus, error) {
	// 获取领导者信息
	leaderServerAddr, leaderId := s.raft.LeaderWithID()
	leader := Node{
		ID:      string(leaderId),
		Address: string(leaderServerAddr),
	}
	servers := s.raft.GetConfiguration().Configuration().Servers
	followers := []Node{}
	me := Node{
		Address: s.RaftAddr,
	}
	for _, server := range servers {
		if server.ID != leaderId {
			followers = append(followers, Node{
				ID:      string(server.ID),
				Address: string(server.Address),
			})
		}
		if string(server.Address) == s.RaftAddr {
			me = Node{
				ID:      string(server.ID),
				Address: string(server.Address),
			}
		}
	}
	status := StoreStatus{
		Me:        me,
		Leader:    leader,
		Followers: followers,
	}
	return status, nil
}
func (s *Store) Tables() ([]byte, error) {
	alltables, err := s.db.GetAllTables()
	jsonbytes, _ := json.Marshal(alltables)
	return jsonbytes, err
}
