package main

import (
	"context"
	loadbalance "go-gateway/gateway/loadBalance"
	zookeeper "go-gateway/gateway/middleware/serverDiscovery/zooKeeper"
	httpProxy "go-gateway/gateway/proxy/httpProxy/http"
	"log"
	"net/http"
)


func main() {
	conf, err := loadbalance.NewLoadBalanceZkConf("http://%s/", zookeeper.NodeName, []string{"localhost:2181", "localhost:2182", "localhost:2183"}, map[string]string{
		"localhost:8001": "10",
		"localhost:8002": "20",
	})

	if err != nil {
		log.Println(err)
		return
	}

	rb := loadbalance.LoadBalanceWithConfFactory(loadbalance.LbRoundRobin, conf)
	
	reverseProxy  := httpProxy.NewLoadBalanceReverseProxy(context.Background(),  rb)
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", reverseProxy))
}