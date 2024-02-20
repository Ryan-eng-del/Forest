package main

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
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


var DefaultTransport http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          129,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	ModifyResponse :=  func( r *http.Response) error {
		if r.StatusCode == 200 {
			srcBody, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			newBody := []byte(string(srcBody) + " Hello")
			r.Body = io.NopCloser(bytes.NewBuffer(newBody))
			len := int64(len(newBody))

			r.ContentLength = len
			r.Header.Set("Content-Length", strconv.FormatInt(len, 10))
		}
		return nil
	}

	ErrorHandler := func(w http.ResponseWriter, r *http.Request, e error) {
		http.Error(w, "ErrorHandler " + e.Error(), http.StatusBadGateway)
	}

	return &httputil.ReverseProxy{Director: director, ModifyResponse: ModifyResponse, ErrorHandler: ErrorHandler, Transport: DefaultTransport}
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