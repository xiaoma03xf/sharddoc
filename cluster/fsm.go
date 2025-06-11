package cluster

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/lib"
	"github.com/xiaoma03xf/sharddoc/storage"
)

type snapshot struct {
	db     *storage.DB
	logger *log.Logger
}

func (f *Store) Apply(l *raft.Log) interface{} {
	var cmd SQLCommand
	if err := json.Unmarshal(l.Data, &cmd); err != nil {
		f.logger.Println("FSM Apply: Failed to unmarshal:", err)
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	// SQL query
	if cmd.IsSelect {
		var result *storage.QueryResult
		if result = f.db.Raw(cmd.SQL); result.Err != nil {
			f.logger.Println("FSM Apply: Query error:", result.Err)
			return result.Err
		}
		// 序列化查询结果
		b, err := json.Marshal(result)
		if err != nil {
			f.logger.Println("FSM Apply: Marshal result error:", err)
			return err
		}
		return &ExecSQLRsp{data: b, Err: err}
	}

	if err := f.db.Exec(cmd.SQL); err != nil {
		f.logger.Println("FSM Apply: ExecSQL error:", err)
		return err
	}
	return &ExecSQLRsp{data: []byte("OK"), Err: nil}
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
	if err = lib.UntarGz(rc, tmpDir); err != nil {
		return err
	}

	// 删除之前的数据库
	dbname := f.db.Path
	backup := dbname + ".backup"
	if err := os.Rename(dbname, backup); err != nil {
		return err
	}

	db, err := storage.ImportDB(tmpDir, dbname)
	if err != nil {
		_ = os.Rename(backup, dbname)
		return err
	}
	f.db = db
	_ = os.Remove(backup)
	return nil
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()
	s.logger.Println("Starting snapshot export...")
	snapshotDir, err := s.db.ExportDB()
	if err != nil {
		return err
	}
	defer os.RemoveAll(snapshotDir)
	// 将snapshot 压缩成 tar.gz 并写入 sink
	start := time.Now()
	if err := lib.TarGz(snapshotDir, sink); err != nil {
		_ = sink.Cancel()
		return err
	}
	s.logger.Printf("Snapshot export complete, duration=%v", time.Since(start))
	return nil
}

func (s *snapshot) Release() {}
