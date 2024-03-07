package main

import (
	"go-gateway/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	server.HTTPServerStop()
}
