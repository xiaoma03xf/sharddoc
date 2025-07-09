package tcp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
	"github.com/xiaoma03xf/sharddoc/storage"
)

type Response struct {
	Type byte
	Body []byte
}

func ReadResponse(conn net.Conn) (*Response, error) {
	header := make([]byte, 5) // 1 字节类型 + 4 字节长度

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("faild to read header: %w", err)
	}
	msgType := header[0]
	bodyLen := binary.BigEndian.Uint32(header[1:5])

	// 读取正文
	body := make([]byte, bodyLen)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("faild to read body: %w", err)
	}
	return &Response{
		Type: msgType,
		Body: body,
	}, nil
}
func buildRawRequest(sql string) ([]byte, error) {
	return BuildTcpInfo(&RaftRequest{
		RequestID: uuid.NewString(),
		DataType:  TypeExec,
		Payload: map[string]interface{}{
			"sql": sql,
		},
	})
}
func buildExecRequest(sql string) ([]byte, error) {
	return BuildTcpInfo(&RaftRequest{
		RequestID: uuid.NewString(),
		DataType:  TypeExec,
		Payload: map[string]interface{}{
			"sql": sql,
		},
	})
}

type Client struct {
	Conn net.Conn
}

func Open(remoteAddr string) (*Client, error) {
	// new Client
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	c := Client{}
	c.Conn = conn
	return &c, nil
}
func (c *Client) Raw(sql string) ([]storage.Record, error) {
	tcpreq, err := buildRawRequest(sql)
	if err != nil {
		return nil, err
	}
	_, err = c.Conn.Write(tcpreq)
	if err != nil {
		return nil, err
	}
	res, err := ReadResponse(c.Conn)
	if err != nil {
		return nil, err
	}
	if res.Type == TypeBadResp {
		return nil, errors.New(string(res.Body))
	}
	var recs storage.QueryResult
	_ = json.Unmarshal(res.Body, &recs)
	return recs.Recs, recs.Err
}
func (c *Client) Exec(sql string) (string, error) {
	tcpreq, err := buildExecRequest(sql)
	if err != nil {
		return "", err
	}
	_, err = c.Conn.Write(tcpreq)
	if err != nil {
		return "", err
	}
	res, err := ReadResponse(c.Conn)
	if err != nil {
		return "", err
	}
	if res.Type == TypeBadResp {
		return "", errors.New(string(res.Body))
	}
	return string(res.Body), nil
}
func (c *Client) Close() error {
	return c.Conn.Close()
}
