package middlewares

import (
	"errors"
	"fmt"
	"go-gateway/handler"
	libConst "go-gateway/lib/const"
	"go-gateway/model"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcFlowLimitMiddleware(serviceDetail *model.ServiceDetail) grpc.StreamServerInterceptor{
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, sHandler grpc.StreamHandler) error {

		serviceFlowNum := serviceDetail.AccessControl.ServiceFlowLimit
		serviceFlowType := serviceDetail.AccessControl.ServiceFlowType

		if serviceFlowNum > 0 {		
			serviceLimiter, err := handler.FlowLimiterHandler.GetLimiter(
				libConst.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceFlowNum), serviceFlowType, true)
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return fmt.Errorf("service flow limit %v", serviceFlowNum)
			}
    }

		peerCtx, ok := peer.FromContext(ss.Context())

		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		clientFlowNum := serviceDetail.AccessControl.ClientIPFlowLimit
		clientFlowType := serviceDetail.AccessControl.ClientFlowType

		if clientFlowNum > 0 {
			clientLimiter, err := handler.FlowLimiterHandler.GetLimiter(libConst.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP, float64(clientFlowNum), clientFlowType, true)
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return fmt.Errorf("%v flow limit %v", clientIP,clientFlowNum)
			}
		}
		return sHandler(srv, ss);
	}
}