package tcc

import "time"

type Options struct {
	// 事务执行时间限制
	Timeout time.Duration
	// 轮询监控任务间隔时长
	MonitorTick time.Duration
}

type Option func(*Options)
