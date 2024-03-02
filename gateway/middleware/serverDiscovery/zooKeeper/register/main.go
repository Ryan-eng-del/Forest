package main

import (
	zookeeper "go-gateway/gateway/middleware/serverDiscovery/zooKeeper"
	"os"
	"os/signal"
)


func main() {
	zookeeper.ZooKeeperRegister()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit
}
