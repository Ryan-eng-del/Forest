package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)


type Server struct{
	Addr string `json:"addr"`
}

func main() {
	server := &Server{Addr: "localhost:8001"}

	server.Run()

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
		log.Fatal(server.ListenAndServe())
	}()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server Address: %s", s.Addr)))
}