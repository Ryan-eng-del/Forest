package main

import (
	"fmt"
	zookeeper "go-gateway/middleware/serverDiscovery/zooKeeper"
	"log"
	"net/http"
	"os"
	"os/signal"
)


type Server struct{
	Addr string `json:"addr"`
}

func main() {
	server := &Server{Addr: "127.0.0.1:8001"}

	go func ()  {
		server.Run()
	}()

	server1 := &Server{Addr: "127.0.0.1:8002"}

	go func ()  {
		server1.Run()
	}()

	// server.Run()

	// signals := make(chan struct{}, 1)
	// signals <- struct{}{}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit
	log.Println("已经退出⏏️")
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/real", s.ServeHTTP)
	server := &http.Server {
		Addr: s.Addr,
		Handler: mux,
	}
	
	go func ()  {
		zk := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
		err := zk.GetConnection()
		if err != nil {
			log.Println(err)
			return
		}
		defer zk.Close()
		err = zk.RegisterServerPath(zookeeper.NodeName, server.Addr)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Register", server.Addr)
		log.Fatal(server.ListenAndServe())
	}()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server Address: %s", s.Addr)))
}