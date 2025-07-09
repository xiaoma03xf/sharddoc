#!/bin/bash

echo "🛑 正在关闭所有节点..."

# 关闭运行中的 Go 节点进程
pkill -f './bin/node'

# 等待进程结束
sleep 1

echo "🧹 清理日志和二进制文件..."

# 删除日志和二进制文件
rm -rf logs
rm -rf bin
rm -rf cluster


# 如果你用 etcd 存储数据，也可以顺便清理（可选）
# export ETCDCTL_API=3
# etcdctl --endpoints=localhost:2379 del "" --prefix

# 删除你自己的数据库数据目录（如果有）
# rm -rf ./data

echo "✅ 所有节点已关闭并清理干净"
