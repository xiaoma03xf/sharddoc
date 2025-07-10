package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Join raft cluster
func (s *Store) Join(ctx context.Context, req *raftpb.JoinRequest) (*raftpb.JoinResponse, error) {
	// TODO 在一个raft组中，follower向leader发送grpc消息会报错
	//[ERROR] raft-net: failed to decode incoming command: error="unknown rpc type 80"
	if s.raft.State() != raft.Leader {
		// return nil, errors.New("not leader")
		leaderAddr := string(s.raft.Leader())
		if leaderAddr == "" {
			return nil, errors.New("no leader found")
		}
		// 向leader节点发起grpc请求, 转发Join请求
		client, conn, err := BuildGrpcConn(s.getGrpcAddrFromRaftAddress(leaderAddr))
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
	return &raftpb.JoinResponse{Success: true}, nil
}

func (s *Store) getGrpcAddrFromRaftAddress(raftaddress string) string {
	//TODO 假设我们有一个配置映射，提供 Raft 地址与 gRPC 地址的映射
	var grpcAddressMap = map[string]string{
		"127.0.0.1:28001": "127.0.0.1:29001",
		"127.0.0.1:28002": "127.0.0.1:29002",
		"127.0.0.1:28003": "127.0.0.1:29003",
	}
	return grpcAddressMap[raftaddress]
}

// Status return current raft cluster Status (me, leader, follower)
// 返回的是raft内部信息，而不是暴露在外的 grpc 信息
func (s *Store) Status(ctx context.Context, req *raftpb.StatusRequest) (*raftpb.StatusResponse, error) {
	leaderServerAddr, leaderID := s.raft.LeaderWithID()
	leader := &raftpb.Node{
		Id:          string(leaderID),
		Address:     string(leaderServerAddr),
		Grpcaddress: s.getGrpcAddrFromRaftAddress(string(leaderServerAddr)),
	}

	servers := s.raft.GetConfiguration().Configuration().Servers
	followers := []*raftpb.Node{}
	me := &raftpb.Node{
		Id:          s.nodeID,
		Address:     s.raftAddr,
		Grpcaddress: s.grpcaddr,
	}
	for _, server := range servers {
		if server.ID != leaderID {
			followers = append(followers, &raftpb.Node{
				Id:          string(server.ID),
				Address:     string(server.Address),
				Grpcaddress: s.getGrpcAddrFromRaftAddress(string(server.Address)),
			})
		}
	}
	return &raftpb.StatusResponse{
		Me:       me,
		Leader:   leader,
		Follower: followers,
	}, nil
}

func deepCopyPutRequest(orig *raftpb.PutRequest) (*raftpb.PutRequest, error) {
	b, err := proto.Marshal(orig)
	if err != nil {
		return nil, err
	}
	copy := &raftpb.PutRequest{}
	if err := proto.Unmarshal(b, copy); err != nil {
		return nil, err
	}
	return copy, nil
}

func (s *Store) Put(ctx context.Context, req *raftpb.PutRequest) (*raftpb.PutResponse, error) {
	fmt.Printf("API Received Key:%v, Value:%v, Mode:%v\n", req.Key, req.Value, req.Mode)

	cleanReq, err := deepCopyPutRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to deep copy PutRequest: %v", err)
	}
	reqData, err := proto.Marshal(cleanReq)
	if err != nil {
		return nil, fmt.Errorf("faild to marshal PutRequest: %v", err)
	}
	// s.logger.Printf("SEND: len=%d, sha256=%x\n", len(reqData), sha256.Sum256(reqData))
	// s.logger.Printf("SEND: len=%d, data=%v\n", len(reqData), reqData)
	fmt.Printf("API Received Data:%v\n", reqData)

	op := &raftpb.Operation{
		Type:  raftpb.OperationType_PUT,
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
	resp, ok := f.Response().(*raftpb.PutResponse)
	if !ok {
		s.logger.Printf("unexpected response type:%v", f.Response())
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}
func (s *Store) Delete(ctx context.Context, req *raftpb.DeleteRequest) (*raftpb.DeleteResponse, error) {
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
	resp, ok := f.Response().(*raftpb.DeleteResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	return resp, nil
}

// Get, Scan等操作不用等apply结束
func (s *Store) Get(ctx context.Context, req *raftpb.GetRequest) (*raftpb.GetResponse, error) {
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
	return &raftpb.GetResponse{Value: val, Found: found}, nil
}

func nonPrimaryKeyCols(tdef *kv.TableDef) (out []string) {
	for _, c := range tdef.Cols {
		if slices.Index(tdef.Indexes[0], c) < 0 {
			out = append(out, c)
		}
	}
	return
}
func (s *Store) Scan(ctx context.Context, req *raftpb.ScanRequest) (*raftpb.ScanResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 执行Get操作
	tx := kv.KVTX{}
	s.kv.Begin(&tx)

	KVIter := tx.Seek(req.KeyStart, int(req.Cmp1), req.KeyEnd, int(req.Cmp2))
	err := s.kv.Commit(&tx)
	Assert(err == nil)

	// 获取table信息
	var tdef kv.TableDef
	err = json.Unmarshal(req.Table, &tdef)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot unmarshal kv.TableDef")
	}

	// 遍历迭代器
	got := []*raftpb.Record{}
	for KVIter.Valid() {
		rec := &kv.Record{}
		// decode the index key
		// 用的是二级索引
		// 二级索引只编码了 索引字段 + 主键字段，不存储原始值（val 是空的）。
		// 解码出主键字段后，构造主键 Record。
		// 然后用主键再去 执行一次主键查询 dbGet()，把这一行补全回来
		// fetch the KV from the iterator

		// pkCols, nonPKCols 取出主键字段和非主键字段
		// 把主键列和非主键列拼接成完整列顺序：["id", "name", "age"]
		rec.Cols = slices.Concat(tdef.Indexes[0], nonPrimaryKeyCols(&tdef))
		rec.Vals = rec.Vals[:0]
		for _, c := range rec.Cols {
			//为每个列设置对应的类型（在 Value.Type 字段里设置）
			// tdef.Cols = ["id", "name", "age"]
			// tdef.Types = [TYPE_INT, TYPE_STRING, TYPE_INT]
			// 最后 rec.Vals = [{Type: INT}, {Type: STRING}, {Type: INT}]
			tp := tdef.Types[slices.Index(tdef.Cols, c)]
			rec.Vals = append(rec.Vals, kv.Value{Type: tp})
		}

		// 主键索引：直接解码 key 和 val
		key, val := KVIter.Deref()
		// primary key or secondary index?
		if req.Index == 0 {

			// key = prefix + encode(1)（即 id）
			// key : prefix + (1字节类型标识位) + id
			// val = encode("Alice", 20)（即 name 和 age）
			// val = (1字节类型标识位) + id + (1字节类型标识位) + id...
			npk := len(tdef.Indexes[0])
			kv.DecodeKey(key, rec.Vals[:npk])
			kv.DecodeValues(val, rec.Vals[npk:])

		} else {
			// 二级索引：解 key -> 提取主键字段 -> 构造主键 key -> tx.Get(key) -> 解 val
			// prefix id 主键索引
			// prefix name age id
			// irec 的cols为索引值 如 name age
			// 然后找到name, age对应的类型 Type
			Assert(len(val) == 0) // 二级索引没存 value
			index := tdef.Indexes[req.Index]
			irec := kv.Record{
				Cols: index,
				Vals: make([]kv.Value, len(index)),
			}
			// 为索引列设置类型
			for i, c := range index {
				irec.Vals[i].Type = tdef.Types[slices.Index(tdef.Cols, c)]
			}
			kv.DecodeKey(key, irec.Vals)

			// 假如我的二级索引为name age，那么格式就为[prefix:4字节][name][age][id]
			// extract the primary key
			for i, c := range tdef.Indexes[0] {
				rec.Vals[i] = *irec.Get(c)
			}

			// fetch the row by the primary key
			// TODO: skip this if the index contains all the columns
			pkKey := kv.EncodeKey(nil, tdef.Prefixes[0], rec.Vals[:len(tdef.Indexes[0])])
			val, ok := tx.Get(pkKey)
			if !ok || val == nil {
				return nil, status.Errorf(codes.NotFound, "primary key not found for index entry: %+v", pkKey)
			}

			npk := len(tdef.Indexes[0])
			kv.DecodeValues(val, rec.Vals[npk:])
		}

		got = append(got, kvRecordToraftpbRecord(*rec))
		KVIter.Next()
	}

	return &raftpb.ScanResponse{
		Records: got,
	}, nil
}

func kvRecordToraftpbRecord(kvRec kv.Record) *raftpb.Record {
	pbRec := &raftpb.Record{
		Cols: make([]string, len(kvRec.Cols)),
		Vals: make([]*raftpb.Value, len(kvRec.Vals)),
	}
	// 复制 Cols
	copy(pbRec.Cols, kvRec.Cols)
	// 转换 Vals
	for i, kvVal := range kvRec.Vals {
		pbRec.Vals[i] = &raftpb.Value{
			Type: kvVal.Type,
			I64:  kvVal.I64,
			Str:  kvVal.Str,
		}
	}
	return pbRec
}
