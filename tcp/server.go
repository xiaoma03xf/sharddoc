package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/xiaoma03xf/sharddoc/cluster"
	"github.com/xiaoma03xf/sharddoc/lib/logger"
)

const (
	TypeExec    = 0x01
	TypeJoin    = 0x02
	TypeStatus  = 0x03
	TypeOKResp  = 0x81
	TypeBadResp = 0x82
)

type Node interface {
	Exec(sql string) *cluster.ExecSQLRsp
	Join(nodeID string, addr string) error
	Status() (cluster.StoreStatus, error)
}

// Servcie provides HTTP service
type Service struct {
	Addr string
	Node Node
}

func NewService(addr string, node Node) *Service {
	return &Service{Addr: addr, Node: node}
}

func ReadRequest(conn net.Conn) (byte, map[string]interface{}, error) {
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
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return 0, nil, err
	}
	return msgType, data, nil
}

// val 可以为ExecPayload, JoinPayload, 主要用于客户端构造请求
func BuildTcpInfo(datatype byte, val map[string]interface{}) ([]byte, error) {
	// tcp 消息必须是指定类型
	switch datatype {
	case TypeExec:
		if !MapFieldCheck([]string{"sql"}, val) {
			return nil, fmt.Errorf("Exec[%v]消息类型缺少某些字段", TypeExec)
		}
	case TypeJoin:
		if !MapFieldCheck([]string{"node_id", "addr"}, val) {
			return nil, fmt.Errorf("Join[%v]消息类型缺少某些字段", TypeJoin)
		}
	case TypeStatus:
		// Valid types, do nothing
	default:
		return nil, fmt.Errorf("datatype unsupported: %v", datatype)
	}

	// Json 编码
	payload, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	// 构造 TCP 数据：type(1 byte) + length(4 bytes) + payload
	var buf bytes.Buffer
	buf.WriteByte(datatype)
	binary.Write(&buf, binary.BigEndian, uint32(len(payload)))
	buf.Write(payload)
	return buf.Bytes(), nil
}

// 主要用于服务端构造答复
func SendResponse(conn net.Conn, datatype byte, payload []byte) error {
	var buf bytes.Buffer
	buf.WriteByte(datatype)

	if err := binary.Write(&buf, binary.BigEndian, uint32(len(payload))); err != nil {
		return err
	}
	buf.Write(payload)
	_, err := conn.Write(buf.Bytes())
	return err
}
func (s *Service) Handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		msgType, data, err := ReadRequest(conn)
		if err != nil {
			if err == io.EOF {
				logger.Info("client disabled connection")
			} else {
				logger.Info("read request err", err)
			}
		}
		switch msgType {
		case TypeExec:
			sqlStr, ok := data["sql"].(string)
			if !ok {
				fmt.Println("Exec 请求缺少 sql 字段或格式错误")
				continue
			}
			result := s.Node.Exec(sqlStr)
			if result.Err != nil {
				_ = SendResponse(conn, TypeBadResp, []byte(err.Error()))
			}
			resp, _ := json.Marshal(result)
			SendResponse(conn, TypeOKResp, resp)

		case TypeJoin:
			nodeID, ok1 := data["node_id"].(string)
			addr, ok2 := data["addr"].(string)
			if !ok1 || !ok2 {
				fmt.Println("Join 请求缺少 node_id 或 addr 字段或格式错误")
				continue
			}
			err := s.Node.Join(nodeID, addr)
			resp := map[string]string{"status": "ok"}
			if err != nil {
				resp["status"] = "fail"
				resp["error"] = err.Error()
			}
			respData, _ := json.Marshal(resp)
			SendResponse(conn, respData)

		case TypeStatus:
			status, err := s.Node.Status()
			var resp []byte
			if err != nil {
				resp, _ = json.Marshal(map[string]string{
					"error": err.Error(),
				})
			} else {
				resp, _ = json.Marshal(status)
			}
			SendResponse(conn, resp)
		default:
			fmt.Printf("未知消息类型: %v\n", msgType)
		}
	}
}

func MapFieldCheck(fields []string, mp map[string]interface{}) bool {
	for _, field := range fields {
		if _, ok := mp[field]; !ok {
			return false
		}
	}
	return true
}
