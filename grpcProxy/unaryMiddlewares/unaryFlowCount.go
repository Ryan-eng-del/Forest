package unaryMiddlewares

import (
	"context"
	"go-gateway/handler"
	constLib "go-gateway/lib/const"
	"go-gateway/model"

	"google.golang.org/grpc"
)


func GrpcFlowCountMiddleware(serviceDetail *model.ServiceDetail) grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, gHandler grpc.UnaryHandler)  (resp any, err error){
		totalCounter, err := handler.ServerCountHandler.GetCounter(constLib.FlowTotal)
		if err != nil {
			return nil, err
		}
		totalCounter.Increase()
		serviceCounter, err := handler.ServerCountHandler.GetCounter(constLib.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return nil, err
		}
		serviceCounter.Increase()
		return gHandler(ctx, req)
	}
}