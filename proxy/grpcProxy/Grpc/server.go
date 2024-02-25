package main

import (
	"context"
	"fmt"
	person "go-gateway/proxy/grpcProxy/Grpc/pb"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type HelloService struct {
	person.UnimplementedHelloServiceServer

}

func (s *HelloService) 	BidirectionalStreamHello(stream person.HelloService_BidirectionalStreamHelloServer) error {
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				stream.Send(&person.Person{Name: "end"})
			}
			break
		}
		log.Println(reply)
		stream.Send(&person.Person{Name: "Server-" + reply.Name})
	}
	return nil
}

func (s *HelloService) ClientStreamHello(stream person.HelloService_ClientStreamHelloServer) error {
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				stream.SendAndClose(&person.Person{Name: "end"})
			}
			break
		}
		log.Println(reply)
	}
	return nil
}


func (s *HelloService) ServerStreamHello(p *person.Person, stream person.HelloService_ServerStreamHelloServer) error {
	for i := 0; i <= 5; i++ {
		stream.Send(&person.Person{
			Name: strconv.Itoa(i) + "-cyan",
		})
	}
	return nil
}

func (s *HelloService) Hello(ctx context.Context, p *person.Person) (*person.Person, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println(md)
	}

	return &person.Person{
		Name: "ServerName",
	}, nil
}


func Server() {
	grpcServer := grpc.NewServer()
	person.RegisterHelloServiceServer(grpcServer, &HelloService{})
	l, err := net.Listen("tcp", "localhost:8001")

	if err != nil {
		log.Println(err)
		return
	}
	grpcServer.Serve(l)
}