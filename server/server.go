package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	HttpServerHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode("debug")
	router := InitRouter()

	HttpServerHandler = &http.Server{
		Addr: "127.0.0.1:8087",
		Handler: router,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxHeaderBytes: 1 << uint(20),
	}

	go func(){
		log.Printf("[INFO] HttpServerRun:%s\n", "http://localhost:8081")

		if err := HttpServerHandler.ListenAndServe(); err != nil {
			log.Printf("[ERROR] HttpServerRun:%s err:%v\n", "http://localhost:8081", err)
		}
	}()
}

func HTTPServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := HttpServerHandler.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] HTTPServerStop failed: %v\n", err)
	}

	log.Printf("[INFO] HTTPServerStop completed")
}