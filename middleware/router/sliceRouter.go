package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type SliceRouter struct {
	Groups []*SliceRoute
}

func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

func (r *SliceRouter) Group (path string) *SliceRoute {
	route := &SliceRoute{
		SliceRouter: r,
		Path: path,
	}

	r.Groups = append(r.Groups, route)
	return route
}


type HandlerFunc func (*SliceRouteContext)

type SliceRoute struct {
	*SliceRouter
	Path string
	Handler []HandlerFunc
}

func (r *SliceRoute) Use(middlewares ...HandlerFunc) *SliceRoute{
	r.Handler = append(r.Handler, middlewares...)
	log.Println("中间件数量: ", len(r.Handler))
	return r
}


type SliceRouteContext struct {
	*SliceRoute
	Ctx context.Context
	Req *http.Request
	Rw http.ResponseWriter
	Index int
}

func (c *SliceRouteContext) Reset() {
	c.Index = -1;
}

func (c *SliceRouteContext) Next() {
	c.Index++
	for c.Index < len(c.Handler) {
		c.Handler[c.Index](c)
		c.Index++
	}
}



type Handler func (*SliceRouteContext) http.Handler

type SliceRouteHandler struct {
	h Handler
	router *SliceRouter
}

func (h *SliceRouteHandler) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	c := NewSliceRouterContext(w, req, h.router)
	if h.h != nil {
		c.Handler = append(c.Handler, func (c *SliceRouteContext) {
			h.h(c).ServeHTTP(w, req)
		})
	}
	c.Reset()
	c.Next()
}


func NewSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouteContext {
	sr := &SliceRoute{}
	matchUrlLen := 0
	// 最长前缀匹配
	for _, route := range r.Groups {
		if strings.HasPrefix(req.RequestURI, route.Path) {
			pathLen := len(route.Path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*sr = *route
			}
		}
	}

	c := &SliceRouteContext{
		Rw: rw,
		Req: req,
		Ctx: req.Context(),
		SliceRoute: sr,
	}

	return c
}


func NewSliceRouteHandler(h Handler, router *SliceRouter) *SliceRouteHandler {
	return &SliceRouteHandler{
		h: h,
		router: router,
	}
}