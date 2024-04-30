package grpcProxy

import (
	"context"
	load_balance "go-gateway/loadBalance"
	"log"

	"github.com/Ryan-eng-del/Streams/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)


func NewGrpcLoadBalanceHandler(lb *load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}

		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			// Setup of the proxy's Director.
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype(proxy.Name)))
			
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}

		return proxy.TransparentHandler(director)
	}()
}