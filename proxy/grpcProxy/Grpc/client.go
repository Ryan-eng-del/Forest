package main

import (
	"context"
	person "go-gateway/proxy/grpcProxy/Grpc/pb"
	"io"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	person.HelloServiceClient
	conn *grpc.ClientConn
}

func NewClient () (*Client, error)  {
	grpcConn, err := grpc.Dial("localhost:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := person.NewHelloServiceClient(grpcConn)
	c := &Client{client, grpcConn}
	return c, nil
}

func ClientFunc() {
	c, _ := NewClient()
	defer c.conn.Close()
	c.MyHello()
	c.MyServerStreamHello()
	c.MyClientStreamHello()
	c.MyBidirectionalStreamHello()
}

func (c *Client) MyHello() {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	reply, err := c.Hello(ctx, &person.Person{Name: "Client"})
	if err != nil {
		log.Println(err, reply)
		return
	}

	log.Println(reply)
}

func (c *Client) MyServerStreamHello() {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	reply, err := c.ServerStreamHello(ctx, &person.Person{Name: ""})
	if err != nil {
		log.Println(err)
		return
	}
	
	isReadAll := false

	for {
		msg, err := reply.Recv()
		if err != nil {
			if err == io.EOF {
				isReadAll = true
			}
			break
		}
		log.Println(msg, "reply MyServerStreamHello")
	}

	if isReadAll {
		log.Println("消息全部读出")
	}
}


func (c *Client) MyClientStreamHello() {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	reply, err := c.ClientStreamHello(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	

	for i := 0; i < 5; i++{
		err := reply.Send(&person.Person{
			Name: strconv.FormatInt(int64(i), 10) + "-cyan",
		})

		if err != nil {
			log.Println(err)
			break
		}
	}

	msg, err := reply.CloseAndRecv()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(msg, "reply MyClientStreamHello")
}

func (c *Client) MyBidirectionalStreamHello() {

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	reply, err := c.BidirectionalStreamHello(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	
	sendOk, receOk := make(chan struct{}), make(chan struct{})

	go func (){
		for i := 0; i < 5; i++{
			err := reply.Send(&person.Person{
				Name: strconv.FormatInt(int64(i), 10) + "-cyan",
			})
	
			if err != nil {
				log.Println(err)
				break
			}
		}
		sendOk <- struct{}{}
	}()


	go func ()  {
		for {
			msg, err := reply.Recv()
			if err != nil {
				break
			}
			log.Println(msg)
		}
		receOk <- struct{}{}
	}()

	<- sendOk
	<- receOk

}