#!/bin/bash
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go build -o bin/go_gateway
ps aux | grep go_gateway | grep -v 'grep' | awk '{print $2}' | xargs kill
nohup ./bin/gateway -config=./conf/prod/ -endpoint=server >> logs/server.log 2>&1 &
echo "nohup ./bin/gateway -config=./conf/prod/ -endpoint=server >> logs/server.log 2 > &1 &"
