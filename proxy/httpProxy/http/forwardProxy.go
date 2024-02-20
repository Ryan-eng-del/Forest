package main

import (
	"fmt"
	"net/http"
)

// 需要在部署网关的物理机器上，配置代理
// 目的是让请求，转发到网关
// windows 配置代理 -> 网络 -> 代理 -> 手动设置代理 -> 输入 IP 和 端口
// linux 配置代理 -> vim /ect/profile 或者 ~/.bashrc
// export http_proxy = "ip:port"
// source /etc/profile 或者 source ~/.bashrc

type Proxy struct {

}

// func main() {
// 	proxy :=  &Proxy{}
// 	// http.HandleFunc("/", proxy.ServeHTTP)
// 	http.Handle("/", proxy)
// 	http.ListenAndServe(":8080", nil)
// }


func (p *Proxy) ServeHTTP(rsw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s:\n", req.Method, req.Host, req.RemoteAddr)
	// 浅拷贝，代理 request
	outReq := &http.Request{}
	*outReq = *req

	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)

	if err != nil {
		rsw.WriteHeader(http.StatusBadGateway)
		return
	}

	for key, value := range res.Header {
		for _, v := range value {
			rsw.Header().Add(key, v)
		}
	}

	rsw.WriteHeader(res.StatusCode)
}