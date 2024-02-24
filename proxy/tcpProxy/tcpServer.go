package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type contextKey struct {
	name string
}

var (
	// ServerContextKey is a context key. It can be used in HTTP
	// handlers with Context.Value to access the server that
	// started the handler. The associated value will be of
	// type *Server.
	ServerContextKey = &contextKey{"http-server"}

	// LocalAddrContextKey is a context key. It can be used in
	// HTTP handlers with Context.Value to access the local
	// address the connection arrived on.
	// The associated value will be of type net.Addr.
	LocalAddrContextKey = &contextKey{"local-addr"}
)

var ErrServerClosed = errors.New("http: Server closed")
var ErrAbortHandler = errors.New("net/http: abort Handler")


type TCPHandler interface {
	ServeTCP(context.Context, *conn)
}

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }


type TCPServer struct {
	Addr string
	Handler TCPHandler
	BaseCtx context.Context
	err error
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	KeepAliveTimeout time.Duration
	mu sync.Mutex
	doneChan chan struct{}
	inShutdown int32
	l net.Listener
}

func (srv *TCPServer) ListenAndServe() error {
	if srv.shuttingDown() {
		return http.ErrServerClosed
	}

	addr := srv.Addr

	if srv.Handler == nil {
		srv.Handler = &tcpHandler{}
	}

	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	return srv.Serve(ln)
}

func (srv *TCPServer) Serve (l net.Listener) error {
	srv.l = &onceCloseListener{Listener: l}
	defer l.Close()

	if srv.BaseCtx == nil {
		srv.BaseCtx = context.Background()
	}

	baseCtx := srv.BaseCtx

	ctx := context.WithValue(baseCtx, ServerContextKey, srv)

	for {
		rw, err := l.Accept()
		if err != nil {
			log.Println(err)
			if srv.shuttingDown() {
				return ErrServerClosed
			}	
		}

		c := srv.newConn(rw)
		go c.serve(ctx)
	}
}

func (srv *TCPServer) Close () error {
	atomic.StoreInt32(&srv.inShutdown, 1)
	close(srv.doneChan)
	srv.l.Close()
	return nil
}

func (srv *TCPServer) shuttingDown() bool{
	return atomic.LoadInt32(&srv.inShutdown) != 0
}

func (srv *TCPServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
		remoteAddr: rwc.RemoteAddr().String(),
	}

	if t:= srv.ReadTimeout; t != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(t))
	}

	if t := srv.WriteTimeout; t != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(t))
	}

	if t := srv.KeepAliveTimeout; t != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(t)
		}
	}
	return c
}


type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, conn *conn){
	// http.ListenAndServe()
	conn.rwc.Write([]byte("Pong! Tcp handler here"))
}


// type serverHandler struct {
// 	Srv *TCPServer
// }

// func (srv serverHandler) ServeTCP (ctx context.Context, conn *conn) {

// }

type conn struct {
	server *TCPServer
	cancelCtx context.CancelFunc
	rwc net.Conn
	mu sync.Mutex
	remoteAddr string
}

func (c *conn) serve(ctx context.Context) {
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())

	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}

	}()
	c.server.Handler.ServeTCP(ctx, c)
	// serverHandler{c.server}.ServeTCP(ctx, c)
}


func ListenAndServeTCP(addr string, handler TCPHandler) {
	t := &TCPServer{Addr: addr, Handler: handler, BaseCtx: context.Background()}
	t.ListenAndServe()
}




