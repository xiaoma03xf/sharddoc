package raft

import (
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/protobuf/proto"
)

func (f *Store) Apply(l *raft.Log) interface{} {
	var op pb.Operation
	_ = proto.Unmarshal(l.Data, &op)
	switch op.Type {
	case pb.OperationType_PUT:
		return f.applyPut(l.Data)
	}
	return nil
}

func (f *Store) applyPut(data []byte) *pb.PutResponse {
	f.mu.Lock()
	defer f.mu.Unlock()

	var putReq pb.PutRequest
	if err := proto.Unmarshal(data, &putReq); err != nil {
		f.logger.Println("unmarshal putRequest data err:", err)
		return &pb.PutResponse{Success: false}
	}
	fmt.Println("key:", string(putReq.Key), "value:", string(putReq.Value))

	tx := kv.KVTX{}
	f.kv.Begin(&tx)
	_, err := tx.Set(putReq.Key, putReq.Value)
	if err != nil {
		f.logger.Printf("key: %v, val: %v, put err:%v", string(putReq.Key), string(putReq.Value), err)
		return &pb.PutResponse{Success: false}
	}
	err = f.kv.Commit(&tx)
	if err != nil {
		return &pb.PutResponse{Success: false}
	}
	return &pb.PutResponse{Success: true, IsUpdated: true}
}

func (f *Store) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}
func (f *Store) Restore(rc io.ReadCloser) error {
	return nil
}
