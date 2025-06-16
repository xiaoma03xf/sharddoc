package tcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/xiaoma03xf/sharddoc/lib"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/storage"
)

type snapshot struct {
	db     *storage.DB
	logger *log.Logger
}

func (f *Store) Apply(l *raft.Log) interface{} {
	var req RaftRequest
	_ = json.Unmarshal(l.Data, &req)

	// do task
	rsp := f.applyBusinessLogic(req.DataType, req.Payload)

	if f.raft.State() == raft.Leader {
		f.mu.Lock()
		if ch, ok := f.applyWaiters[req.RequestID]; ok {
			close(ch)
			delete(f.applyWaiters, req.RequestID)
		} else {
			logger.Warn("unexpected err occur, no channel waited?", "requestID", req.RequestID)
		}
		f.mu.Unlock()
	}
	return rsp
}

func (f *Store) applyBusinessLogic(datatype byte, payload map[string]interface{}) *ExecSQLRsp {
	// TODO 数据库相关操作封装, 目前仅支持增删查改, 建表
	sql, _ := payload["sql"].(string)

	logger.Info(fmt.Sprintf("statue:%v handle sql:%v", f.raft.State(), sql))

	// 查询相关操作
	if strings.HasPrefix(strings.ToLower(sql), "select") {
		queryRaw := f.db.Raw(sql)
		b, err := json.Marshal(queryRaw)
		return &ExecSQLRsp{Data: b, Err: err}
	}
	if err := f.db.Exec(sql); err != nil {
		f.logger.Println("FSM Apply: ExecSQL error:", err)
		return &ExecSQLRsp{Err: err}
	}
	return &ExecSQLRsp{Data: []byte("OK")}
}

func (f *Store) WaitForApply(requestID string, timeout time.Duration, ch <-chan struct{}) error {
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		f.mu.Lock()
		delete(f.applyWaiters, requestID)
		f.mu.Unlock()
		return fmt.Errorf("timeout waiting for apply:%s", requestID)
	}
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
