package cluster

import (
	"fmt"
	"os"
	"time"

	"github.com/stretchr/testify/assert/yaml"
)


type Node struct {
	ID      string `json:"id"`      // 节点唯一标识
	Address string `json:"address"` // 节点的 Raft 通信地址（如 "localhost:12000"）
}
type StoreStatus struct {
	Me        Node   `json:"me"`        // 当前节点信息
	Leader    Node   `json:"leader"`    // 领导者节点信息
	Followers []Node `json:"followers"` // 跟随者节点列表
}
type NodeConfig struct {
	NodeID    string `yaml:"node_id"`
	RaftAddr  string `yaml:"raft_addr"`
	HttpAddr  string `yaml:"http_addr"`
	RaftDir   string `yaml:"raft_dir"`
	DBPath    string `yaml:"db_path"`
	Bootstrap bool   `yaml:"bootstrap"`
	JoinAddr  string `yaml:"join_addr,omitempty"`

	// common config
	RaftTimeout       time.Duration `yaml:"raft_timeout"`
	SnapshotInterval  time.Duration `yaml:"snapshot_interval"`
	SnapshotThreshold uint64        `yaml:"snapshot_threshold"`
}

func LoadNodeConfig(path string) (*NodeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg NodeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	return &cfg, nil
}
