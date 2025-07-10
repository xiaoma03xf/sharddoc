package raft

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	"google.golang.org/protobuf/proto"
)

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

func (s *Store) BatchPut(ctx context.Context, req *raftpb.BatchPutRequest) (*raftpb.BatchPutResponse, error) {
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal PutRequest: %v", err)
	}
	op := &raftpb.Operation{
		Type:  raftpb.OperationType_BATCHPUT,
		Data:  reqData,
		Term:  s.raft.CurrentTerm(),
		Index: s.raft.LastIndex(),
	}
	opData, err := proto.Marshal(op)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal OPData: %v", err)
	}
	f := s.raft.Apply(opData, defaultRaftTimeout)
	if err := f.Error(); err != nil {
		return nil, fmt.Errorf("raft apply error: %v", err)
	}
	resp, ok := f.Response().(*raftpb.BatchPutResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}
