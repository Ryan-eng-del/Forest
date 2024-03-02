package tcpProxy

import (
	"context"
	loadbalance "go-gateway/gateway/loadBalance"
	"io"
	"log"
	"net"
	"time"
)

type TCPReverseProxy struct {
	Addr string
	DialTimeout time.Duration
	Deadline time.Duration
	KeepAlivePeriod time.Duration
	ModifyResponse func (net.Conn) error
	ErrorHandler func (net.Conn, error)
	DialContext func (ctx context.Context, network, address string) (net.Conn, error)
	Ctx  context.Context
	Director func (string) (string)
}

func NewTCPReverseProxy (addr string) *TCPReverseProxy {
	return &TCPReverseProxy{Addr: addr, DialTimeout: 10 * time.Second, Deadline: time.Hour, KeepAlivePeriod: time.Hour}
}


func (pxy *TCPReverseProxy) ServeTCP(ctx context.Context,src net.Conn) {
	var baseCtx context.Context
	 if pxy.DialTimeout  >= 0 {
		c, cancelTimeout := context.WithTimeout(ctx, pxy.DialTimeout)
		baseCtx = c
		defer cancelTimeout()
	 }

	 if pxy.Deadline >= 0 {
		c, cancelDeadline := context.WithDeadline(ctx, time.Now().Add(pxy.DialTimeout))
		baseCtx = c
		defer cancelDeadline()
	 }

	 if pxy.DialContext == nil {
		pxy.DialContext = (&net.Dialer{
			Timeout: pxy.DialTimeout,
			Deadline: time.Now().Add(pxy.DialTimeout),
			KeepAlive: pxy.KeepAlivePeriod,

		}).DialContext
	 }
	 log.Println(pxy.Addr, "addr")
	 pxy.Director(src.RemoteAddr().String())

	 conn , err := pxy.DialContext(baseCtx, "tcp", pxy.Addr)
	 if err != nil {
		pxy.getErrorHandler()(conn, err)
		return	
	 }

	 defer conn.Close()
	 if !pxy.modifyResponse(conn) {
		return 
	 }

	 _, err = io.Copy(src, conn)

	 if err != nil {
		pxy.getErrorHandler()(conn, err)
		return
	 }
}

func (pxy *TCPReverseProxy) defaultErrorHandler(conn net.Conn, err error) {
	log.Printf("TCP Error for remote %s is %s", conn.RemoteAddr().String(), err)
}

func (pxy *TCPReverseProxy) getErrorHandler() func(conn net.Conn, err error) {
	if pxy.ErrorHandler == nil {
		pxy.ErrorHandler = pxy.defaultErrorHandler
	}
	return pxy.ErrorHandler
}

func (pxy *TCPReverseProxy) modifyResponse (conn net.Conn) bool {
	if pxy.ModifyResponse == nil {
		return true
	}

	if err := pxy.ModifyResponse(conn); err != nil {
		conn.Close()
		pxy.getErrorHandler()(conn, err)
		return false
	}

	return true
}


func NewTcpLoadbalanceReverseProxy(c context.Context,lb loadbalance.LoadBalance) *TCPReverseProxy     {
	pxy := &TCPReverseProxy{
		Ctx: c,
		Deadline: time.Minute,
		DialTimeout: 10 * time.Second,
		KeepAlivePeriod: time.Hour,
	}

	pxy.Director = func(s string) string {
		nextAddr := lb.Get(s)
		pxy.Addr = nextAddr
		return nextAddr
	}

	return pxy
}