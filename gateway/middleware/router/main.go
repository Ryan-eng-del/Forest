package router

import (
	"log"
	"net/http"
)

func main() {
	var addr = "127.0.0.1:8006"
	log.Println("Starting HttpServer at: ", addr)

	sliceRouter := NewSliceRouter()
	sliceRoute := sliceRouter.Group("/v2")

	sliceRoute.Use(handler1, handler2, handler3)
	var routerHandler http.Handler = NewSliceRouteHandler(nil, sliceRouter)
	http.ListenAndServe(addr, routerHandler)

}

func handler1 (c *SliceRouteContext) {
	log.Println("trace in handler1")
	c.Next()
	log.Println("trace out handler1")
}

func handler2 (c *SliceRouteContext) {
	log.Println("trace in handler2")
	c.Next()
	log.Println("trace out handler2")
}

func handler3 (c *SliceRouteContext) {
	log.Println("trace in handler3")
	c.Next()
	log.Println("trace out handler3")
}