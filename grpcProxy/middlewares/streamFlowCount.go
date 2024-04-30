package middlewares

import (
	"go-gateway/handler"
	constLib "go-gateway/lib/const"
	"go-gateway/model"

	"google.golang.org/grpc"
)

func GrpcFlowCountMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor{
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, sHandler grpc.StreamHandler) error{
		totalCounter, err := handler.ServerCountHandler.GetCounter(constLib.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		serviceCounter, err := handler.ServerCountHandler.GetCounter(constLib.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()
		return sHandler(srv, ss)
	}
}