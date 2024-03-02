package main

import (
	loadbalance "go-gateway/loadBalance"
	"go-gateway/proxy/grpcProxy"
	"log"
	"net"

	"google.golang.org/grpc"
)


func main() {
	listen, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}
	lb := loadbalance.LoadBalanceFactory(loadbalance.LbRoundRobin)
	lb.Add("127.0.0.1:8001")
	handler := grpcProxy.NewGrpcLoadBalanceHandler(lb)
	s := grpc.NewServer(grpc.UnknownServiceHandler(handler))
	s.Serve(listen)
}