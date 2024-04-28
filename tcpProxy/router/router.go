package router

import (
	"context"
	"fmt"
	"go-gateway/handler"
	"go-gateway/model"
	"go-gateway/public"
	"go-gateway/tcpProxy"
	tcpMiddlewares "go-gateway/tcpProxy/middlewares"
	"go-gateway/tcpProxy/server"
	"log"
	"sync"
	"time"
)

type TcpManager struct {
	ServerList []*server.TcpServer
}

var TcpManagerHandler *TcpManager

func init() {
	TcpManagerHandler = NewTcpManager()
}


func NewTcpManager() *TcpManager {
	return &TcpManager{}
}

func (t *TcpManager) tcpOneServerRun(service *model.ServiceDetail) {
	addr := fmt.Sprintf(":%d", service.TCPRule.Port)
	rb, err := handler.LoadBalancerHandler.GetLoadBalancer(service)
	if err != nil {
		log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
		return
	}

	router := tcpMiddlewares.NewTcpSliceRouter()
	router.Group("/").Use(
		tcpMiddlewares.TCPFlowCountMiddleware(),
		tcpMiddlewares.TCPFlowLimitMiddleware(),
		tcpMiddlewares.TCPWhiteListMiddleware(),
		tcpMiddlewares.TCPBlackListMiddleware(),
	)

	routerHandler := tcpMiddlewares.NewTcpSliceRouterHandler(
		func(c *tcpMiddlewares.TcpSliceRouterContext) server.TCPHandler {
			return tcpProxy.NewTcpLoadBalanceReverseProxy(c, rb)
		}, router)

	baseCtx := context.WithValue(context.Background(), server.ServiceContextKey, service)
	tcpServer := &server.TcpServer{
		Addr:     addr,
		Handler:  routerHandler,
		BaseCtx:  baseCtx,
		UpdateAt: time.Time(service.Info.UpdateAt),
	}

	t.ServerList = append(t.ServerList, tcpServer)
	log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
	if err := tcpServer.ListenAndServe(); err != nil && err != server.ErrServerClosed {
		log.Printf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
	}
}

func (t *TcpManager) TcpServerRun() {
	serviceList := handler.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(service *model.ServiceDetail) {
			t.tcpOneServerRun(service)
		}(tempItem)
	}
	handler.ServiceManagerHandler.Regist(t)
}

func (t *TcpManager) TcpServerStop() {
	for _, tcpServer := range t.ServerList {
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func (tcpServer *server.TcpServer)  {
			defer func () {
				wait.Done()
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			tcpServer.Close()
		}(tcpServer)
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}

func (t *TcpManager) Update(e *handler.ServiceEvent) {
	delList := e.DeleteService
	for _, delService := range delList {
		if delService.Info.LoadType != public.LoadType(1) {
			continue
		}
		for _, tcpServer := range TcpManagerHandler.ServerList {
			if delService.Info.ServiceName != tcpServer.ServiceName {
				continue
			}
			tcpServer.Close()
			log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
		}
	}

	addList := e.AddService
	for _, addService := range addList {
		if addService.Info.LoadType != public.LoadType(1) {
			continue
		}
		go t.tcpOneServerRun(addService)
	}
}

