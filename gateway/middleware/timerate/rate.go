package main

import (
	"go-gateway/gateway/middleware/router"

	"golang.org/x/time/rate"
)

func RateLimiter(params ...int) func (c *router.SliceRouteContext) {
	var r rate.Limit = 1
	var b = 2

	if len(params) == 2 {
		r = rate.Limit(params[0])
		b = params[1]
	}

	l := rate.NewLimiter(r, b)

	return func (c *router.SliceRouteContext) {
		if !l.Allow() {
			c.Rw.Write([]byte("已被限流，请稍后重试"))
			c.Abort()
			return
		}
		c.Next()
	}
}