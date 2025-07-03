package raft

import (
	"io"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/protobuf/proto"
)

func (f *Store) Apply(l *raft.Log) interface{} {
	var op pb.Operation
	err := proto.Unmarshal(l.Data, &op)
	if err != nil {
		panic(err)
	}
	switch op.Type {
	case pb.OperationType_PUT:
		return f.applyPut(op.Data)
	case pb.OperationType_BATCHPUT:
		return f.applyBatchPut(op.Data)
	case pb.OperationType_DELETE:
		return f.applyDelete(op.Data)
	}
	return nil
}

func (f *Store) applyPut(data []byte) *pb.PutResponse {
	var putreq pb.PutRequest
	if err := proto.Unmarshal(data, &putreq); err != nil {
		f.logger.Println("unmarshal putrequest data err:", err)
		return &pb.PutResponse{Success: false}
	}

	f.mu.Lock()
	tx := &kv.KVTX{}
	f.kv.Begin(tx)
	_, err := tx.Set(putreq.Key, putreq.Value)
	Assert(err == nil)
	err = f.kv.Commit(tx)
	Assert(err == nil)
	f.mu.Unlock()

	return &pb.PutResponse{Success: true, IsUpdated: true}
}

func (f *Store) applyBatchPut(data []byte) *pb.BatchPutResponse {
	var putreq pb.BatchPutRequest
	if err := proto.Unmarshal(data, &putreq); err != nil {
		f.logger.Println("unmarshal putrequest data err:", err)
		return &pb.BatchPutResponse{Success: false}
	}

	f.mu.Lock()
	tx := &kv.KVTX{}
	f.kv.Begin(tx)
	for _, pair := range putreq.Pairs {
		_, err := tx.Set(pair.Key, pair.Value)
		Assert(err == nil)
	}
	err := f.kv.Commit(tx)
	Assert(err == nil)
	f.mu.Unlock()

	return &pb.BatchPutResponse{Success: true}
}
func (f *Store) applyDelete(data []byte) *pb.DeleteResponse {
	var delreq pb.DeleteRequest
	if err := proto.Unmarshal(data, &delreq); err != nil {
		f.logger.Println("unmarshal deleteRequest data err:", err)
		return &pb.DeleteResponse{Success: false}
	}
	
	f.mu.Lock()
	tx := &kv.KVTX{}
	f.kv.Begin(tx)
	_, err := tx.Del(&kv.DeleteReq{Key: []byte(delreq.Key)})
	Assert(err == nil)
	err = f.kv.Commit(tx)
	Assert(err == nil)
	f.mu.Unlock()

	return &pb.DeleteResponse{Success: true}
}

func (f *Store) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}
func (f *Store) Restore(rc io.ReadCloser) error {
	return nil
}
