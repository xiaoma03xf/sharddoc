package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/xiaoma03xf/sharddoc/lib/logger"
)

const (
	TypeExec    = 0x01 // 执行sql语句标识
	TypeJoin    = 0x02
	TypeStatus  = 0x03
	TypeShowTbl = 0x04
	TypeOKResp  = 0x81
	TypeBadResp = 0x82
)

func (s *Service) RegisterHandlers() {
	s.Handlers = map[byte]HandleFunc{
		TypeExec:    s.HandleExec,
		TypeJoin:    s.HandleJoin,
		TypeStatus:  s.HandleStatus,
		TypeShowTbl: s.HandleShowTables,
	}
}
func (s *Service) HandleExec(ctx context.Context, conn net.Conn, raftReq *RaftRequest) {
	sql, ok := raftReq.Payload["sql"].(string)
	if !ok {
		_ = SendBadResponse(conn, []byte("Exec request is missing the 'sql' field or the format is correct"))
		return
	}
	result := s.Node.Exec(raftReq)
	if result == nil {
		_ = SendBadResponse(conn, []byte("Exec returned nil result"))
		return
	}
	if result.Err != nil {
		_ = SendBadResponse(conn, []byte(result.Err.Error()))
		return
	}
	// record this sql
	logger.Info("SQL", sql)

	resp, _ := json.Marshal(result.Data)
	_ = SendResponse(conn, TypeOKResp, resp)
}
func (s *Service) HandleJoin(ctx context.Context, conn net.Conn, raftReq *RaftRequest) {
	_, ok1 := raftReq.Payload["node_id"].(string)
	_, ok2 := raftReq.Payload["addr"].(string)
	if !ok1 || !ok2 {
		_ = SendBadResponse(conn, []byte("Join request is missing the 'node_id' field or 'addr' field"))
		return
	}
	err := s.Node.Join(raftReq)
	if err != nil {
		_ = SendBadResponse(conn, []byte(err.Error()))
		return
	}
	_ = SendResponse(conn, TypeOKResp, []byte("OK"))
}

func (s *Service) HandleStatus(ctx context.Context, conn net.Conn, raftReq *RaftRequest) {
	status, err := s.Node.Status(raftReq)
	if err != nil {
		logger.Warn(fmt.Sprintf("get node :%v err:%v", s.Addr, err))
	} else {
		resp, _ := json.Marshal(status)
		_ = SendResponse(conn, TypeOKResp, resp)
	}
}

func (s *Service) HandleShowTables(ctx context.Context, conn net.Conn, raftReq *RaftRequest) {
	tblInfo, err := s.Node.Tables(raftReq)
	if err != nil {
		logger.Warn("get all tables err", err)
		_ = SendBadResponse(conn, []byte(err.Error()))
		return
	}
	_ = SendResponse(conn, TypeOKResp, tblInfo)
}
