package test

import (
	server "go-gateway/server"
	"log"
	"sync"
)


var (
	addr       string    = "127.0.0.1:6111"
	initOnce   sync.Once = sync.Once{}
	serverOnce sync.Once = sync.Once{}
)

//初始化测试用例
func SetUp() {
	initOnce.Do(func() {
		if err := server.InitModule("../../conf/dev/"); err != nil {
			log.Fatal(err)
		}
	})
}

//销毁测试用例
func TearDown() {
	// server.DestroyModule()
}