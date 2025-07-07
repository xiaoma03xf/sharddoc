package kv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/kv/pb"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"google.golang.org/protobuf/proto"
)

type OperationType byte

const (
	OpDelete   OperationType = iota // DELETE
	OpUpdate                        // UPDATE
	BATCH_SIZE = 1000               // 批量写入数量
)

var snapshotBatchPool = sync.Pool{
	New: func() interface{} {
		return &pb.SnapshotBatch{Snapshots: make([]*pb.IncrementalSnapshot, 0, BATCH_SIZE)}
	},
}

func (kv *KV) asyncSnapshot(tx *KVTX, wals *pb.SnapshotBatch) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(kv.Snapshot), 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// 打开文件：如果文件不存在则创建文件，存在则以追加模式打开
	file, err := os.OpenFile(kv.Snapshot, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Warn("open walfile err", err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if fileInfo, _ := os.Stat(kv.Snapshot); fileInfo.Size() == 0 {
		// 初始化snapshot, 直接把snapshot全写入，再写入pending tree
		// 分批次写入 BATCH_SIZE
		if err = batchInsertIterTree(&tx.snapshot, file); err != nil {
			return err
		}
	}
	if err = writeBatchToFile(file, wals); err != nil {
		return err
	}

	return nil
}

func batchInsertIterTree(b *BTree, file *os.File) error {
	btreeData := snapshotBatchPool.Get().(*pb.SnapshotBatch)
	defer snapshotBatchPool.Put(btreeData)

	for iter := b.Seek(nil, CMP_GT); iter.Valid(); iter.Next() {
		key, val := iter.Deref()
		basic := pb.IncrementalSnapshot{}
		basic.Key = key
		basic.Value = val
		basic.Operation = int32(OpUpdate)
		basic.Timestamp = time.Now().UnixNano()
		btreeData.Snapshots = append(btreeData.Snapshots, &basic)

		// count >= BatchSize
		if len(btreeData.Snapshots) >= BATCH_SIZE {
			if err := writeBatchToFile(file, btreeData); err != nil {
				return err
			}
			btreeData = snapshotBatchPool.Get().(*pb.SnapshotBatch)
		}
	}
	if len(btreeData.Snapshots) > 0 {
		if err := writeBatchToFile(file, btreeData); err != nil {
			return err
		}
	}
	return nil
}

// 批量写入增量数据
func writeBatchToFile(file *os.File, wals *pb.SnapshotBatch) error {
	writer := bufio.NewWriter(file)
	data, err := proto.Marshal(wals)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("faild to write to file:%w", err)
	}
	return writer.Flush()
}

func (kv *KV) restoreFromSnapshot(snapshotPath string) (*pb.SnapshotBatch, error) {
	file, err := os.Open(snapshotPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snapshot pb.SnapshotBatch
	err = proto.Unmarshal(data, &snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}
