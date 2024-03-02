package main

import (
	"log"
	"net/rpc"
)

func Call() {
	conn, err := rpc.Dial("tcp", "127.0.0.1:8004")

	if err != nil {
		log.Println("Dial error: ", err)
		return
	}

	defer conn.Close()
	var reply string

	conn.Call("hello.HelloWorld", "小李", &reply)
	log.Println(reply, "response")
}