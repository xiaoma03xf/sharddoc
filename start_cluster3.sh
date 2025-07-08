#!/bin/bash
go run cmd/main.go -config ./cluster3.yaml -node node7 &
sleep 2

go run cmd/main.go -config ./cluster3.yaml -node node8 &
sleep 2

go run cmd/main.go -config ./cluster3.yaml -node node9 &
