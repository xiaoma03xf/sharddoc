package raft

import (
	"fmt"
	"io"
	"os"

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
	f.mu.Lock()
	defer f.mu.Unlock()

	// kv.Snapshot 指的是 wal.log 的完整路径
	return &fileSnapshot{snapshotPath: f.kv.Snapshot}, nil
}
func (f *Store) Restore(rc io.ReadCloser) error {
	// rc io.ReadCloser由Raft提供的 snapshot 数据流，即之前 Persist() 写入的内容
	f.mu.Lock()
	defer f.mu.Unlock()
	defer rc.Close()

	tmpPath := f.kv.Snapshot + ".restore"
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create restore temp file failed: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, rc); err != nil {
		return fmt.Errorf("copy snapshot restore data failed: %w", err)
	}

	// 替换老快照文件
	if err := os.Rename(tmpPath, f.kv.Snapshot); err != nil {
		return fmt.Errorf("replace snapshot file failed: %w", err)
	}

	return nil
}

type fileSnapshot struct {
	snapshotPath string
}

func (f *fileSnapshot) Persist(sink raft.SnapshotSink) error {
	// 打开本地 wal.log 快照文件
	file, err := os.Open(f.snapshotPath)
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("open snapshot file failed: %w", err)
	}
	defer file.Close()
	if _, err := io.Copy(sink, file); err != nil {
		sink.Cancel()
		return fmt.Errorf("copy snapshot to sink failed: %w", err)
	}
	if err := sink.Close(); err != nil {
		return fmt.Errorf("sink close failed: %w", err)
	}
	return nil
}
func (f *fileSnapshot) Release() {}
