package raft

import "time"

type Config struct {
	Common *CommonConfig `yaml:"common"`
	Nodes  []*NodeConfig `yaml:"nodes"`
}

type CommonConfig struct {
	Etcd struct {
		Endpoints []string `yaml:"endpoints"`
		TTL       int64    `yaml:"ttl"`
	} `yaml:"etcd"`
	RaftTimeout       time.Duration `yaml:"raft_timeout"`
	SnapshotInterval  time.Duration `yaml:"snapshot_interval"`
	SnapshotThreshold uint64        `yaml:"snapshot_threshold"`
}

type NodeConfig struct {
	NodeID    string `yaml:"node_id"`
	ClusterID string `yaml:"cluster_id"`
	Bootstrap bool   `yaml:"bootstrap"`
	JoinAddr  string `yaml:"join_addr"`
	RaftAddr  string `yaml:"raft_addr"`
	GrpcAddr  string `yaml:"grpc_addr"`
	RaftDir   string `yaml:"raft_dir"`
	KVPath    string `yaml:"kv_path"`
	KVLogPath string `yaml:"kv_logpath"`
}
