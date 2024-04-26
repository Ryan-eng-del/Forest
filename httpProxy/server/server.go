package httpProxyServer

import (
	"context"
	router "go-gateway/httpProxy/router"
	libViper "go-gateway/lib/viper"
	"go-gateway/middlewares"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


var (
	HttpSrvHandler  *http.Server
)

func ServerRun() {
	gin.SetMode(libViper.ViperInstance.GetStringConf("proxy.base.debug_mode"))

	r := router.InitRouter(middlewares.RecoveryMiddleware(), middlewares.RequestLogMiddleware())

	HttpSrvHandler = &http.Server{
		Addr:           libViper.ViperInstance.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(libViper.ViperInstance.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(libViper.ViperInstance.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(libViper.ViperInstance.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] http_proxy_run %s\n", libViper.ViperInstance.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] http_proxy_run %s err:%v\n", libViper.ViperInstance.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] HTTPServerStop failed: %v\n", err)
	}
	log.Printf("[INFO] HTTPServerStop completed")
}
