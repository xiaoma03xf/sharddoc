#!/bin/bash
go run cmd/main.go -config ./cluster1.yaml -node node1 &
sleep 2

go run cmd/main.go -config ./cluster1.yaml -node node2 &
sleep 2

go run cmd/main.go -config ./cluster1.yaml -node node3 &
