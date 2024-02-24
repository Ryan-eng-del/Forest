package main

import (
	"go-gateway/proxy/grpcProxy/jsonRpc/greeterRpc/inter"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"google.golang.org/grpc"
)
type Client struct {
	c *rpc.Client
}

func (client *Client) HelloWorld(req string, resp *string) error {
	return client.c.Call(inter.HelloServiceMethod, req, resp)
}

func NewClient () (*Client, error) {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8004")

	if err != nil {
		log.Println("Dial error: ", err)
		return nil, err
	}

	client := &Client{
		c: conn,
	}
	return client, err
}


func Call() {
	client, err := NewClient()
	if err != nil {
		log.Println("Dial error: ", err)
		return
	}
	grpc.NewServer()

	defer client.c.Close()
	var reply string
	client.HelloWorld("小李", &reply)
	log.Println(reply, "response")
}