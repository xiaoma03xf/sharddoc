#!/bin/bash

# 编译主程序
go build -o ./bin/node cmd/main.go
if [ $? -ne 0 ]; then
  echo "❌ 编译失败"
  exit 1
fi

echo "🚀 启动集群节点..."

mkdir -p logs

# 启动所有 Leader 节点（并发）
declare -a LEADERS=(
  "cluster1.yaml node1"
  "cluster2.yaml node4"
  "cluster3.yaml node7"
)

for entry in "${LEADERS[@]}"; do
  config=$(echo "$entry" | awk '{print $1}')
  node=$(echo "$entry" | awk '{print $2}')
  echo "🟢 启动 Leader $node ..."
  nohup ./bin/node -config ./$config -node $node > logs/$node.log 2>&1 &
done

# 等待所有 Leader Raft 稳定
echo "⏳ 等待 Leader 初始化 Raft..."
sleep 3

# 启动所有 Follower 节点
declare -a FOLLOWERS=(
  "cluster1.yaml node2"
  "cluster1.yaml node3"
  "cluster2.yaml node5"
  "cluster2.yaml node6"
  "cluster3.yaml node8"
  "cluster3.yaml node9"
)

for entry in "${FOLLOWERS[@]}"; do
  config=$(echo "$entry" | awk '{print $1}')
  node=$(echo "$entry" | awk '{print $2}')
  echo "🔵 启动 Follower $node ..."
  nohup ./bin/node -config ./$config -node $node > logs/$node.log 2>&1 &
  sleep 1
done

echo "✅ 所有节点已启动（日志保存在 logs/）"
