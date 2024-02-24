package main

import (
	"context"
	person "go-gateway/proxy/grpcProxy/Grpc/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func Client() {
	grpcConn, err := grpc.Dial("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return
	}

	client := person.NewHelloServiceClient(grpcConn)

	reply, err := client.Hello(context.Background(), &person.Person{Name: ""})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(reply)
}