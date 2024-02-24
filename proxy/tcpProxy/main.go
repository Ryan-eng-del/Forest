package main

import (
	"context"
	TcpProxy "go-gateway/proxy/tcpProxy/tcp"
	"log"
	"net"
	"os"
	"os/signal"
)


var (
	targetAddr = "localhost:81"
	proxyAddr = "localhost:83"
)
type customHandler struct {

}

func (h *customHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	conn.Write([]byte("custom tcp downstream server"))
}


func main() {
	go func () {
		t := &TCPServer{Addr:targetAddr, Handler: &customHandler{}}
		log.Println("server start at http://localhost:81")
		t.ListenAndServe()
	}()


	go func () {
		tpcProxy := TcpProxy.NewTCPReverseProxy(targetAddr)
		t := &TCPServer{Addr: proxyAddr, Handler: tpcProxy}
		log.Println("server start at http://localhost:83")
		t.ListenAndServe()
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<- s
	log.Println("已经退出⏏️")
}