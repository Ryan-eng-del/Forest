package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)


var (
	target = "localhost:8001"
	proxyAddr = "localhost:8002"
)


func WebsocketProxy() {
	targetUrl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	http.ListenAndServe(proxyAddr, proxy)
}