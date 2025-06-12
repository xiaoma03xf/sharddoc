package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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
