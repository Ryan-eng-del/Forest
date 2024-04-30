package unaryMiddlewares

import (
	"context"
	"errors"
	"go-gateway/handler"
	libJwt "go-gateway/lib/jwt"
	"go-gateway/model"
	"go-gateway/public"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcJwtAuthTokenMiddleware(serviceDetail *model.ServiceDetail) grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, gHandler grpc.UnaryHandler)  (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, errors.New("miss metadata from context")
		}
		appMatched := false
		authToken:=""
		auths:=md.Get("authorization")

		if len(auths)>0{
			authToken = auths[0]
		}

		token:=strings.ReplaceAll(authToken,"Bearer ","")
		if token != "" {
			jwtInstance := libJwt.NewJWT()
			appClaims, err := jwtInstance.ParseAppJWT(token)
			if err != nil {
				return nil, errors.New("fail to decode token: " + err.Error())
			}

			appList := handler.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == appClaims.AppId {
					md.Set("app", public.Obj2Json(appInfo))
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 &&  !appMatched {
			return nil, errors.New("not match valid app")
		}
		return gHandler(ctx, req)
	}
}