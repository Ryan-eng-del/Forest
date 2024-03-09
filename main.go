package main

import (
	"go-gateway/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := server.InitModule("./conf/dev/")

	if err != nil {
		panic(err)
	}
	
	server.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	server.HTTPServerStop()
}
