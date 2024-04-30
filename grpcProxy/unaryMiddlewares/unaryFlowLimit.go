package unaryMiddlewares

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/handler"
	libConst "go-gateway/lib/const"
	"go-gateway/model"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcFlowLimitMiddleware(serviceDetail *model.ServiceDetail) grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, gHandler grpc.UnaryHandler)  (resp any, err error) {

		serviceFlowNum := serviceDetail.AccessControl.ServiceFlowLimit
		serviceFlowType := serviceDetail.AccessControl.ServiceFlowType

		if serviceFlowNum > 0 {		
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				libConst.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceFlowNum), serviceFlowType, true)
			if err != nil {
				return nil, err
			}
			if !serviceLimiter.Allow() {
				return nil, fmt.Errorf("service flow limit %v", serviceFlowNum)
			}
    }
		
		peerCtx, ok := peer.FromContext(ctx)

		if !ok {
			return nil, errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		clientFlowNum := serviceDetail.AccessControl.ClientIPFlowLimit
		clientFlowType := serviceDetail.AccessControl.ClientFlowType

		if clientFlowNum > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(libConst.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP, float64(clientFlowNum), clientFlowType, true)
			if err != nil {
				return nil, err
			}
			if !clientLimiter.Allow() {
				return nil, fmt.Errorf("%v flow limit %v", clientIP,clientFlowNum)
			}
		}
		return gHandler(ctx, req);
	}
}