package main

import (
	lib "go-gateway/lib/mysql"
	"go-gateway/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := server.InitModule("./conf/dev/")

	db, _ := lib.GetGormPool("default")
	// db.AutoMigrate(model.Service{},model.AccessControl{}, model.App{}, model.GrpcRule{}, model.Admin{},  model.TcpRule{},  model.HttpRule{}, model.LoadBalance{})

	if err != nil {
		panic(err)
	}
	
	server.HttpServerRun()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	server.HTTPServerStop()
}
