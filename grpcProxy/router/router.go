package router

import (
	"fmt"
	"go-gateway/grpcProxy"
	"go-gateway/grpcProxy/middlewares"
	"go-gateway/grpcProxy/unaryMiddlewares"
	"go-gateway/handler"
	"go-gateway/model"
	"go-gateway/public"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var GrpcManagerHandler *GrpcManager

type GrpcManager struct {
	ServerList []*warpGrpcServer
}

type warpGrpcServer struct {
	Addr        string
	ServiceName string
	UpdateAt    time.Time
	*grpc.Server
}


func NewGrpcManager() *GrpcManager {
	return &GrpcManager{}
}


func (g *GrpcManager) GrpcServerRun() {
	serviceList := handler.ServiceManagerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go g.grpcOneServerRun(tempItem)
	}
	handler.ServiceManagerHandler.Regist(g)
}


func (g *GrpcManager) grpcOneServerRun(service *model.ServiceDetail) {
	addr := fmt.Sprintf(":%d", service.GRPCRule.Port)
	rb, err := handler.LoadBalancerHandler.GetLoadBalancer(service)

	if err != nil {
		log.Printf(" [ERROR] GetTcpLoadBalancer %v err:%v\n", addr, err)
		return
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf(" [ERROR] GrpcListen %v err:%v\n", addr, err)
		return
	}

	grpcHandler := grpcProxy.NewGrpcLoadBalanceHandler(rb)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryMiddlewares.GrpcFlowCountMiddleware(service),
			unaryMiddlewares.GrpcFlowLimitMiddleware(service),
			unaryMiddlewares.GrpcWhiteListMiddleware(service),
			unaryMiddlewares.GrpcBlackListMiddleware(service),
			unaryMiddlewares.GrpcJwtAuthTokenMiddleware(service),
			unaryMiddlewares.GrpcHeaderTransferMiddleware(service),
		),
		grpc.ChainStreamInterceptor(
			middlewares.GrpcFlowCountMiddleware(service),
			middlewares.GrpcFlowLimitMiddleware(service),
			middlewares.GrpcWhiteListMiddleware(service),
			middlewares.GrpcBlackListMiddleware(service),
			middlewares.GrpcJwtAuthTokenMiddleware(service),
			middlewares.GrpcHeaderTransferMiddleware(service),

		),
		grpc.UnknownServiceHandler(grpcHandler),
	)

	GrpcManagerHandler.ServerList = append(GrpcManagerHandler.ServerList, &warpGrpcServer{
		Addr:        addr,
		ServiceName: service.Info.ServiceName,
		UpdateAt:    time.Time(service.Info.UpdateAt),
		Server:      s,
	})

	log.Printf(" [INFO] grpc_proxy_run %v\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Printf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
	}
}

func (g *GrpcManager) Update(e *handler.ServiceEvent) {
	delList := e.DeleteService
	for _, delService := range delList {
		if delService.Info.LoadType != public.LoadType(2) {
			continue
		}
		for _, tcpServer := range GrpcManagerHandler.ServerList {
			if delService.Info.ServiceName != tcpServer.ServiceName {
				continue
			}
			tcpServer.GracefulStop()
			log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", tcpServer.Addr)
		}
	}
	addList := e.AddService
	for _, addService := range addList {
		if addService.Info.LoadType != public.LoadType(2) {
			continue
		}
		go g.grpcOneServerRun(addService)
	}
}

func (g *GrpcManager) ServerStop() {
	for _, grpcServerItem := range GrpcManagerHandler.ServerList {
		wait := sync.WaitGroup{}
		wait.Add(1)

		go func(grpcServerItem *warpGrpcServer){
			defer func() {
				wait.Done()
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			grpcServerItem.GracefulStop()
		}(grpcServerItem)
		wait.Wait()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServerItem.Addr)
	}
}