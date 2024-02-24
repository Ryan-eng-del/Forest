package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	serverStart := make(chan struct{}, 1)


	go func ()  {
		Server()
	}()

	go func () {
		<- time.After(5*time.Second)
		serverStart <- struct{}{}
		log.Println("服务启动")
	}()

	go func ()  {
		<- serverStart
		log.Println("开始调用")
		Client()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit 

	// select {}
}