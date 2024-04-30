package grpcProxy

import (
	"context"
	loadbalance "go-gateway/gateway/loadBalance"
	person "go-gateway/gateway/proxy/grpcProxy/Grpc/pb"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// func main() {
// 	listener, err := net.Listen("tcp", "localhost:8085")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	option := grpc.UnknownServiceHandler(handler)
// 	s := grpc.NewServer(option)
// 	s.Serve(listener)
// }

func (h *HandlerCompose) handler (srv any, pxyServerStream grpc.ServerStream) error {
	methodName, _ := grpc.MethodFromServerStream(pxyServerStream)

	ctx := pxyServerStream.Context()


	ctx, pxyClientConn, err := h.director(ctx, methodName)
	// pxyClientConn, err := grpc.DialContext(ctx, "localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return err
	}
	defer pxyClientConn.Close()

	md, _ := metadata.FromIncomingContext(ctx)

	outCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	outCtx = metadata.NewOutgoingContext(outCtx, md)

	pxyStreamDesc := &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}

	clientStream, err := grpc.NewClientStream(outCtx, pxyStreamDesc, pxyClientConn, methodName)

	if err != nil {
		log.Println(err)
		return err
	}

	toRealServerErr := ProxyServerToRealServer(clientStream, pxyServerStream)
	returnClientErr := ProxyServerToRealClient(pxyServerStream, clientStream)

	for i := 0; i < 2; i++ {
		select {
		case err := <- toRealServerErr:
				if err == io.EOF {
					clientStream.CloseSend()
				} else {
					if clientCancel != nil  {
						clientCancel()
					}
					return status.Errorf(codes.Internal, "failed proxy to real server:%v", err)
				}
		case err := <- returnClientErr:
			pxyServerStream.SetTrailer(clientStream.Trailer())
			if err != io.EOF {
				return err
			}
			return nil
		}
	}

	return nil
}

func ProxyServerToRealServer(dst grpc.ClientStream, src grpc.ServerStream) <- chan error {
	res := make(chan error, 1)
	go func ()  {
		msg:= &person.Person{}
		for {
			if err := src.RecvMsg(msg); err != nil {
				res <- err
				break
			}

			if err := dst.SendMsg(msg); err != nil {
				res <- err
				break
			}
		}
	}()
	return res
}

func ProxyServerToRealClient(dst grpc.ServerStream, src grpc.ClientStream) <- chan error{
	res := make(chan error, 1)
	go func ()  {
		msg := &person.Person{}

		for {
			md, _ := src.Header()
			_ = dst.SendHeader(md)

			if err := src.RecvMsg(msg); err != nil {
				res <- err
				break
			}
			if err := dst.SendMsg(msg); err != nil {
				res <- err
				break
			}
		}
	}()


	return res
}

type GrpcDirector func (context.Context, string) (context.Context, *grpc.ClientConn, error)

type HandlerCompose struct {
	director GrpcDirector
}

func TransportHandler(director GrpcDirector) grpc.StreamHandler {
	return (&HandlerCompose{director}).handler
}

func NewGrpcLoadBalanceHandler(lb loadbalance.LoadBalance) grpc.StreamHandler    {
	return func ()grpc.StreamHandler {
		director := func (ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error){
			nextAddr := lb.Get(fullMethodName)
			pxyClientConn, err := grpc.DialContext(ctx, nextAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			return ctx, pxyClientConn, err
		}
		return TransportHandler(director)
	}()
}