#!/bin/sh

export GO111MODULE=auto && export GOPROXY=https://goproxy.cn && go mod tidy
GOOS=linux GOARCH=amd64 go build -o ./bin/forest-gateway

docker build -f DockerfileDashboard -t forest-gateway-dashboard:latest .
docker build -f DockerfileServer -t forest-gateway:latest .