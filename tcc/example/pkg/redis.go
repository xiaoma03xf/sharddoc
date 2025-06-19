package pkg

import (
	"fmt"
	"sync"

	"github.com/xiaoma03xf/sharddoc/lib/redislock"
)

const (
	network  = "tcp"
	address  = ""
	password = ""
)

var (
	redisClient *redislock.Client
	once        sync.Once
)

func NewRedisClient(network, address, password string) *redislock.Client {
	return redislock.NewClient(network, address, password)
}

func GetRedisClient() *redislock.Client {
	once.Do(func() {
		redisClient = redislock.NewClient(network, address, password)
	})
	return redisClient
}

// 构造事务 id key，用于幂等去重
func BuildTXKey(componentID, txID string) string {
	return fmt.Sprintf("txKey:%s:%s", componentID, txID)
}

func BuildTXDetailKey(componentID, txID string) string {
	return fmt.Sprintf("txDetailKey:%s:%s", componentID, txID)
}

// 构造请求 id，用于记录状态机
func BuildDataKey(componentID, txID, bizID string) string {
	return fmt.Sprintf("txKey:%s:%s:%s", componentID, txID, bizID)
}

// 构造事务锁 key
func BuildTXLockKey(componentID, txID string) string {
	return fmt.Sprintf("txLockKey:%s:%s", componentID, txID)
}

func BuildTXRecordLockKey() string {
	return "gotcc:txRecord:lock"
}
