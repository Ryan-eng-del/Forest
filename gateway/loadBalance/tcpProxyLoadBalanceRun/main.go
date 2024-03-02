package main

import (
	"context"
	loadbalance "go-gateway/gateway/loadBalance"
	tcpProxy "go-gateway/gateway/proxy/tcpProxy/tcp"
	"log"
)

func main() {
	lb := loadbalance.LoadBalanceFactory(loadbalance.LbRoundRobin)

	lb.Add([]string{"127.0.0.1:83", "127.0.0.1:84"}...)

	proxy := tcpProxy.NewTcpLoadbalanceReverseProxy(context.Background(), lb)
	server := tcpProxy.TCPServer{Addr: "127.0.0.1:82", Handler: proxy}
	log.Fatal(server.ListenAndServe())
}