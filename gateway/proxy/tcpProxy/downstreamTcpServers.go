package main

import (
	"context"
	"fmt"
	tcpProxy "go-gateway/gateway/proxy/tcpProxy/tcp"

	"log"
	"net"
	"os"
	"os/signal"
)


var (
	targetAddr = "localhost:8"
	proxyAddr = "localhost:83"
)
type customHandler struct {

}

func (h *customHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("custom tcp downstream server: %s", conn.LocalAddr().String())))
}


func main() {
	go func () {
		t := &tcpProxy.TCPServer{Addr: "127.0.0.1:83", Handler: &customHandler{}}
		log.Println("server start at http://localhost:83")
		t.ListenAndServe()
	}()

	go func () {
		t := &tcpProxy.TCPServer{Addr: "127.0.0.1:84", Handler: &customHandler{}}
		log.Println("server start at http://localhost:84")
		t.ListenAndServe()
	}()

	// go func () {
	// 	tpcProxy := tcpProxy.NewTCPReverseProxy(targetAddr)
	// 	t := &tcpProxy.TCPServer{Addr: proxyAddr, Handler: tpcProxy}
	// 	log.Println("server start at http://localhost:83")
	// 	t.ListenAndServe()
	// }()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<- s
	log.Println("已经退出⏏️")
}