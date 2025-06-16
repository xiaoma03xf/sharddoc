package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
)

type RaftNode interface {
	Exec(*RaftRequest) *ExecSQLRsp
	Join(*RaftRequest) error
	Status(*RaftRequest) (StoreStatus, error)
	Tables(*RaftRequest) ([]byte, error)
}

// Servcie provides HTTP service
type HandleFunc func(context.Context, net.Conn, *RaftRequest)
type Service struct {
	Addr     string
	Node     RaftNode
	Handlers map[byte]HandleFunc
}

func NewService(addr string, node RaftNode) *Service {
	s := new(Service)
	s.Addr = addr
	s.Node = node
	s.RegisterHandlers()
	return s
}
func (s *Service) Handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		msgTyps, raftReq, err := ReadRequest(conn)
		if msgTyps == 0 || raftReq == nil {
			_ = SendBadResponse(conn, []byte("Unexpected syntax..."))
			continue
		}
		if err != nil {
			if err == io.EOF {
				logger.Info("client disabled connection")
			} else {
				logger.Info("read request err", err)
			}
			return
		}
		handler, ok := s.Handlers[raftReq.DataType]
		if !ok {
			logger.Warn(fmt.Sprintf("未知消息类型: %v\n", raftReq.DataType))
			continue
		}
		handler(ctx, conn, raftReq)
	}
}

func (s *Service) Close() error {
	return nil
}

func BootstrapCluster(cfgPath string) {
	nodeCfg, err := LoadNodeConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	s, err := NewStore(nodeCfg)
	if err != nil {
		panic(err)
	}
	err = s.Open(nodeCfg)
	if err != nil {
		panic(err)
	}
	// If join was specified, make the join request.
	if nodeCfg.JoinAddr != "" {
		logger.Info("start join in cluster...")
		if err := Join(nodeCfg.JoinAddr, nodeCfg.RaftAddr, nodeCfg.NodeID); err != nil {
			logger.Warn(fmt.Sprintf("failed to join node at %s: %s", nodeCfg.JoinAddr, err.Error()))
		}
	}
	// start tcp server
	service := NewService(nodeCfg.HttpAddr, s)
	go ListenAndServeWithSignal(service.Addr, service)

	// we're up and running!
	logger.Info(fmt.Sprintf("node started successfully, listening on: http://%s", nodeCfg.HttpAddr))

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	logger.Info("hraftd exiting")
}

func ReadRequest(conn net.Conn) (byte, *RaftRequest, error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(conn, header); err != nil {
		return 0, nil, err
	}
	msgType := header[0]
	length := binary.BigEndian.Uint32(header[1:5])
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return 0, nil, err
	}
	// payload 数据反序列化为map
	var raftRequest RaftRequest
	if err := json.Unmarshal(payload, &raftRequest); err != nil {
		return 0, nil, err
	}
	return msgType, &raftRequest, nil
}

type RaftRequest struct {
	RequestID string
	DataType  byte
	Payload   map[string]interface{}
}

// val 可以为ExecPayload, JoinPayload, 主要用于客户端构造请求
func BuildTcpInfo(req *RaftRequest) ([]byte, error) {
	if req.RequestID == "" {
		req.RequestID = uuid.NewString()
	}
	switch req.DataType {
	case TypeExec:
		if !MapFieldCheck([]string{"sql"}, req.Payload) {
			return nil, fmt.Errorf("Exec[%v]消息类型缺少某些字段", TypeExec)
		}
	case TypeJoin:
		if !MapFieldCheck([]string{"node_id", "addr"}, req.Payload) {
			return nil, fmt.Errorf("Join[%v]消息类型缺少某些字段", TypeJoin)
		}
	case TypeStatus, TypeShowTbl:
		// Valid types, do nothing
	default:
		return nil, fmt.Errorf("datatype unsupported: %v", req.DataType)
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// 构造 TCP 数据：type(1 byte) + length(4 bytes) + payload
	// TODO 外type结合内type
	var buf bytes.Buffer
	buf.WriteByte(req.DataType)
	binary.Write(&buf, binary.BigEndian, uint32(len(payload)))
	buf.Write(payload)
	return buf.Bytes(), nil
}

// 主要用于服务端构造答复
func SendResponse(conn net.Conn, datatype byte, payload []byte) error {
	var buf bytes.Buffer
	buf.WriteByte(datatype)

	if err := binary.Write(&buf, binary.BigEndian, uint32(len(payload))); err != nil {
		logger.Warn(fmt.Sprintf("sendresponse err:%v", err))
		return err
	}
	buf.Write(payload)
	_, err := conn.Write(buf.Bytes())
	return err
}
func SendBadResponse(conn net.Conn, msg []byte) error {
	return SendResponse(conn, TypeBadResp, msg)
}

func MapFieldCheck(fields []string, mp map[string]interface{}) bool {
	for _, field := range fields {
		if _, ok := mp[field]; !ok {
			return false
		}
	}
	return true
}

func Join(joinAddr, raftAddr, nodeID string) error {
	data := make(map[string]interface{})
	data["node_id"] = nodeID
	data["addr"] = raftAddr
	//TypeJoin, data
	joinMsg, err := BuildTcpInfo(&RaftRequest{
		RequestID: uuid.New().String(),
		DataType:  TypeJoin,
		Payload:   data,
	})
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", joinAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Write(joinMsg)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logger.Info(fmt.Sprintf("join response read error:%v", err))
	} else {
		logger.Info("join response:", string(buf[:n]))
	}
	return nil
}
