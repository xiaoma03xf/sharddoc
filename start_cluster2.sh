#!/bin/bash
go run cmd/main.go -config ./cluster2.yaml -node node4 &
sleep 2

go run cmd/main.go -config ./cluster2.yaml -node node5 &
sleep 2

go run cmd/main.go -config ./cluster2.yaml -node node6 &
