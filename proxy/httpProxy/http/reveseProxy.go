package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	downstreamAddr = "http://localhost:8001"
)
func main() {
	// ReverseProxy 实现
	downstreamAddr1, _ := url.Parse(downstreamAddr)
	proxy := NewSingleHostReverseProxy(downstreamAddr1)
	http.ListenAndServe(":80", proxy)

	
	// http 代理实现
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("request path is ", r.URL.Path)
	// 	downstreamAddr, _ := url.Parse(downstreamAddr)
	// 	r.URL.Scheme = downstreamAddr.Scheme
	// 	r.URL.Host = downstreamAddr.Host

	// 	transport := http.DefaultTransport
	// 	resp, err := transport.RoundTrip(r)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	defer resp.Body.Close()

	// 	for key, value := range resp.Header {
	// 		for _, v := range value {
	// 			w.Header().Add(key, v)
	// 		}
	// 	}
	
	// 	bufio.NewReader(resp.Body).WriteTo(w)
	// })

	// http.ListenAndServe(":80", nil)
}


func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	return &httputil.ReverseProxy{Director: director}
}


func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = joinURLPath(target.Path, req.URL.Path)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}


func joinURLPath(target string, path string) string {
	isSuffix := strings.HasPrefix(path, "/")
	isPrefix := strings.HasSuffix(target, "/")
	var p string

	if isPrefix && isSuffix {
		p = target + path[1:]
	} else if isPrefix || isSuffix {
		p = target + path
	} else {
		p = target + "/" + path
	}
	return p
}