package unaryMiddlewares

import (
	"context"
	"errors"
	"fmt"
	lib "go-gateway/lib/func"
	"go-gateway/model"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcBlackListMiddleware(serviceDetail *model.ServiceDetail) grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, gHandler grpc.UnaryHandler)  (resp any, err error) {
		peerCtx, ok := peer.FromContext(ctx)
		if !ok {
			return nil, errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		ipList := []string{}
		blackList := []string{}


		if serviceDetail.AccessControl.WhiteList!=""{
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.BlackList != "" {
			blackList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if len(blackList) > 0 && len(ipList) == 0 {
			if !lib.InIPSliceStr(clientIP, ipList) {
				return nil, fmt.Errorf("%s not in white ip list", clientIP)
			}
		}
		return gHandler(ctx, req)
	}
}