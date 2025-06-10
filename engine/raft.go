package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/xiaoma03xf/sharddoc/lib/utils"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type Store struct {
	db       *DB
	RaftDir  string
	RaftAddr string
	mu       sync.RWMutex
	raft     *raft.Raft
	logger   *log.Logger
}

func NewStore(raftdir, raftAddr string) *Store {
	s := &Store{}
	s.RaftDir = raftdir
	s.RaftAddr = raftAddr
	s.logger = log.New(os.Stderr, "[store] ", log.LstdFlags)
	return s
}
func (s *Store) Open(LocalID string, db *DB) error {
	// 创建 Raft 配置
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(LocalID)
	// 控制快照频率
	config.SnapshotInterval = 3 * time.Minute // 至少间隔 3 分钟检查一次是否需要快照
	config.SnapshotThreshold = 100000         // 日志条目超过 100000 才创建快照

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

	hasState, err := raft.HasExistingState(logStore, stableStore, snapshots)
	if err != nil {
		return err
	}
	if !hasState {
		config := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		rf.BootstrapCluster(config)
	}
	return nil
}

type SQLCommand struct {
	IsSelect bool
	SQL      string `json:"sql"`
}

func (s *Store) ApplyExec(sql string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not Leader")
	}
	c := &SQLCommand{
		IsSelect: false,
		SQL:      sql,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// 提供查询
func (s *Store) ApplyQuery(sql string, dest any) error {
	if s.raft.State() != raft.Leader {
		leaderAddr := string(s.raft.Leader())
		if leaderAddr == "" {
			return fmt.Errorf("no leader available")
		}
		return fmt.Errorf("not Leader, try:%v", leaderAddr)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Raw(sql).Scan(&dest)
}

func (f *Store) Apply(l *raft.Log) interface{} {
	var cmd SQLCommand
	if err := json.Unmarshal(l.Data, &cmd); err != nil {
		f.logger.Println("FSM Apply: Failed to unmarshal:", err)
		return err
	}
	// 执行sql
	f.mu.Lock()
	defer f.mu.Unlock()

	if cmd.IsSelect {
		// 执行查询
		var result interface{}
		if err := f.db.Raw(cmd.SQL).Scan(&result); err != nil {
			f.logger.Println("FSM Apply: Query error:", err)
			return err
		}
		// 序列化查询结果
		b, err := json.Marshal(result)
		if err != nil {
			f.logger.Println("FSM Apply: Marshal result error:", err)
			return err
		}
		return b
	}
	
	if err := f.db.Exec(cmd.SQL); err != nil {
		f.logger.Println("FSM Apply: ExecSQL error:", err)
		return err
	}
	return nil
}

func (f *Store) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return &snapshot{db: f.db, logger: f.logger}, nil
}

func (f *Store) Restore(rc io.ReadCloser) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	defer rc.Close()

	tmpDir, err := os.MkdirTemp("", "raft-restore-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if err = utils.UntarGz(rc, tmpDir); err != nil {
		return err
	}

	// 删除之前的数据库
	dbname := f.db.Path
	_ = os.Remove(dbname)
	db, err := ImportDB(tmpDir, dbname)
	if err != nil {
		return err
	}
	f.db = db
	return nil
}

type snapshot struct {
	db     *DB
	logger *log.Logger
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()
	s.logger.Println("Starting snapshot export...")
	snapshotDir, err := s.db.ExportDB()
	if err != nil {
		return err
	}
	// 将snapshot 压缩成 tar.gz 并写入 sink
	if err := utils.TarGz(snapshotDir, sink); err != nil {
		_ = sink.Cancel()
		return err
	}
	s.logger.Println("Snapshot export complete")
	return nil
}

func (s *snapshot) Release() {}
