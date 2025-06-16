package utils

import (
	"io"
	"strings"
)

// raft 库 会一直胡乱输出<timestamp> Rollback failed: tx closed
// raft-boltdb/bolt_store.go 具体这行代码 log.Printf("Rollback failed: %v", err)
// 有人 提issues但是还没解决，暂时先屏蔽这个日志输出。 不过似乎对正常运转没影响？
type InterceptWriter struct {
	W     io.Writer
	Block string
}

func (iw *InterceptWriter) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), iw.Block) {
		// 拦截掉这条日志，不输出，返回长度表示写入成功（防止程序卡死）
		return len(p), nil
	}
	return iw.W.Write(p)
}
