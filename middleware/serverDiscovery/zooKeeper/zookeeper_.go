package zookeeper

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var NodeName = "/server"

func ZooKeeperServer() {
	zk := NewZkManager([]string{"127.0.0.1:2181"})
	zk.GetConnection()
	defer zk.Close()
	zList, err := zk.GetServerListPath(NodeName)
	if err != nil {
		log.Println(err)
	}
	log.Printf("/server children node is %v", zList)

	chanList, chanErr := zk.WatchServerListByPath(NodeName)

	go func () {
		for {
			select {
			case changeErr := <- chanErr:
				log.Printf("ChangeError: %v", changeErr)
			case changeList := <-chanList:
				log.Printf("ChangeList: %v", changeList)
			}
		} 
	} ()
		
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit
}

func ZooKeeperRegister() {
	zk := NewZkManager([]string{"127.0.0.1:2181"})
	zk.GetConnection()
	defer zk.Close()
	i := 0

	for {
		err := zk.RegisterServerPath(NodeName, strconv.Itoa(i))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Register", i)
		time.Sleep(5 * time.Second)
		i++
	}
}