package main

import (
	zookeeper "go-gateway/middleware/serverDiscovery/zooKeeper"
	"os"
	"os/signal"
)


func main() {
	zookeeper.ZooKeeperServer()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit
}
