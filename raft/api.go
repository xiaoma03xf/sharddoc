package raft

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/protobuf/proto"
)

// Join raft cluster
func (s *Store) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	// TODO 在一个raft组中，follower向leader发送grpc消息会报错
	//[ERROR] raft-net: failed to decode incoming command: error="unknown rpc type 80"
	if s.raft.State() != raft.Leader {
		// return nil, errors.New("not leader")
		leaderAddr := string(s.raft.Leader())
		if leaderAddr == "" {
			return nil, errors.New("no leader found")
		}
		// 向leader节点发起grpc请求, 转发Join请求
		client, conn, err := BuildGrpcConn(getGrpcAddrFromRaftAddress(leaderAddr))
		defer conn.Close()
		Assert(err == nil)
		return client.Join(ctx, req)
	}
	// get cluster config
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return nil, err
	}
	// 检查 cluster 中是否已经有相同 NodeID 或 Address 节点
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(req.NodeId) || srv.Address == raft.ServerAddress(req.Address) {
			// 如果节点 ID 和 Address 完全相同，说明该节点已经是集群的一部分
			// 不然就是冲突的节点（相同地址、不同 ID），移除该节点
			if srv.Address == raft.ServerAddress(req.Address) && srv.ID == raft.ServerID(req.NodeId) {
				s.logger.Printf("node %s at %s already member of cluster, ignoring join request", req.NodeId, req.Address)
				return nil, nil // 无需加入，返回 nil
			}
			if err := s.raft.RemoveServer(srv.ID, 0, 0).Error(); err != nil {
				return nil, fmt.Errorf("error removing existing node %s at %s: %s", req.NodeId, req.Address, err)
			}
		}
	}
	// 向集群中添加新的节点作为 Voter（选举者）
	f := s.raft.AddVoter(raft.ServerID(req.NodeId), raft.ServerAddress(req.Address), 0, defaultRaftTimeout)
	if f.Error() != nil {
		return nil, f.Error()
	}
	s.logger.Printf("node %s at %s joined successfully", req.NodeId, req.Address)
	return &pb.JoinResponse{Success: true}, nil
}

func getGrpcAddrFromRaftAddress(raftaddress string) string {
	//TODO 假设我们有一个配置映射，提供 Raft 地址与 gRPC 地址的映射
	grpcAddressMap := map[string]string{
		"127.0.0.1:28001": "127.0.0.1:29001",
		"127.0.0.1:28002": "127.0.0.1:29002",
		"127.0.0.1:28003": "127.0.0.1:29003",
	}
	return grpcAddressMap[raftaddress]
}

// Status return current raft cluster Status (me, leader, follower)
// 返回的是raft内部信息，而不是暴露在外的 grpc 信息
func (s *Store) Status(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	leaderServerAddr, leaderID := s.raft.LeaderWithID()
	leader := &pb.Node{Id: string(leaderID), Address: string(leaderServerAddr)}
	servers := s.raft.GetConfiguration().Configuration().Servers
	followers := []*pb.Node{}
	me := &pb.Node{
		Address: s.raftAddr,
	}
	for _, server := range servers {
		if server.ID != leaderID {
			followers = append(followers, &pb.Node{
				Id:      string(server.ID),
				Address: string(server.Address),
			})
		}
		if string(server.Address) == s.raftAddr {
			me = &pb.Node{
				Id:      string(server.ID),
				Address: string(server.Address),
			}
		}
	}
	status := &pb.StatusResponse{
		Me:       me,
		Leader:   leader,
		Follower: followers,
	}
	return status, nil
}

func (s *Store) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal PutRequest: %v", err)
	}
	op := &pb.Operation{
		Type:  pb.OperationType_PUT,
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
	resp, ok := f.Response().(*pb.PutResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}
func (s *Store) BatchPut(ctx context.Context, req *pb.BatchPutRequest) (*pb.BatchPutResponse, error) {
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal PutRequest: %v", err)
	}
	op := &pb.Operation{
		Type:  pb.OperationType_BATCHPUT,
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
	resp, ok := f.Response().(*pb.BatchPutResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}
func (s *Store) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal PutRequest: %v", err)
	}
	op := &pb.Operation{
		Type:  pb.OperationType_BATCHPUT,
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
	resp, ok := f.Response().(*pb.DeleteResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}

// Get, Scan等操作不用等apply结束
func (s *Store) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// s.raft.CommitIndex() 集群中已经被大多数节点应用并提交的日志条目的最大索引
	// s.raft.LastIndex() 当前leader节点的最后一条日志的索引
	commitIndex, lastIndex := s.raft.CommitIndex(), s.raft.LastIndex()
	for commitIndex != lastIndex {
		select {
		case <-time.After(50 * time.Millisecond):
			commitIndex, lastIndex = s.raft.CommitIndex(), s.raft.LastIndex()
		case <-ctx.Done():
			return nil, errors.New("get operation canceled")
		}
	}
	// 执行Get操作
	tx := kv.KVTX{}
	s.kv.Begin(&tx)
	val, found := tx.Get(req.Key)
	err := s.kv.Commit(&tx)
	Assert(err == nil)
	return &pb.GetResponse{Value: val, Found: found}, nil
}

func (s *Store) Scan(context.Context, *pb.ScanRequest) (*pb.ScanResponse, error) {
	return &pb.ScanResponse{}, nil
}
