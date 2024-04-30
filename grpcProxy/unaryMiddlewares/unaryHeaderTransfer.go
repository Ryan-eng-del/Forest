package unaryMiddlewares

import (
	"context"
	"errors"
	"go-gateway/model"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)


func GrpcHeaderTransferMiddleware(serviceDetail *model.ServiceDetail) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, gHandler grpc.UnaryHandler)  (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}

			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			}

			if items[0] == "del" {
				delete(md, items[1])
			}
		}
		return gHandler(metadata.NewIncomingContext(ctx, md), req)
	}
}