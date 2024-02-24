package main

import (
	"context"
	person "go-gateway/proxy/grpcProxy/Grpc/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HelloService struct {
	person.UnimplementedHelloServiceServer

}

func (s *HelloService) Hello(context.Context, *person.Person) (*person.Person, error) {
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