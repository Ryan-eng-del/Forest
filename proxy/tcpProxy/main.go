package main

import (
	"context"
	"log"
)

type customHandler struct {

}

func (h *customHandler) ServeTCP(ctx context.Context, conn *conn) {
	conn.rwc.Write([]byte("custom tcp"))
}


func main() {
	t := &TCPServer{Addr: "localhost:81", Handler: &customHandler{}}
	log.Println("server start at http://localhost:81")
	t.ListenAndServe()

}