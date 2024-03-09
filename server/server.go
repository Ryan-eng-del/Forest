package server

import (
	"context"
	lib "go-gateway/lib/viper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	HttpServerHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.ViperInstance.GetStringConf("base.base.debug_mode"))
	
	router := InitRouter()

	HttpServerHandler = &http.Server{
		Addr: lib.ViperInstance.GetStringConf("base.http.addr"),
		Handler: router,
		ReadTimeout: time.Duration(lib.ViperInstance.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout: time.Duration(lib.ViperInstance.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.ViperInstance.GetIntConf("base.http.max_header_bytes")),
	}

	go func(){
		log.Printf("[INFO] HttpServerRun:%s\n", lib.ViperInstance.GetStringConf("base.http.addr"))

		if err := HttpServerHandler.ListenAndServe(); err != nil {
			log.Printf("[ERROR] HttpServerRun:%s err:%v\n", lib.ViperInstance.GetStringConf("base.http.addr"), err)
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