package main

import (
	"go-gateway/gateway/proxy/grpcProxy/jsonRpc/greeterRpc/inter"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloWorld struct {}

func RegisterService(handler inter.HelloService) error {
	err := rpc.RegisterName("hello", handler)
	if err != nil {
		return err
	}
	return nil
}


func (hw *HelloWorld) HelloWorld(req string, res *string) error {
	*res = req + "Hello"
	return nil
}

func Server() {
	err := RegisterService(&HelloWorld{})
	if err != nil {
		log.Println("register name error:", err)
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8004")

	if err != nil {
		log.Println("net.Listen.error", err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Println("Accept error:", err)
		return
	}

	jsonrpc.ServeConn(conn)
}