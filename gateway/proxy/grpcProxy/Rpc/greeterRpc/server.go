package main

import (
	"log"
	"net"
	"net/rpc"
)

type HelloWorld struct {}


func (hw *HelloWorld) HelloWorld(req string, res *string) error {
	*res = req + "Hello"
	return nil
}

func Server() {
	err := rpc.RegisterName("hello", &HelloWorld{})

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

	rpc.ServeConn(conn)
}