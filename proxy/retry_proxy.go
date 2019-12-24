package proxy

import "net/http"

type RetryProxy struct {
	config *ProxyConfig
}

func (p *RetryProxy) Serve(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("ho"))
}
