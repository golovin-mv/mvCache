package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RetryProxy struct {
	config  *ProxyConfig
	target  *httputil.ReverseProxy
	reserve *httputil.ReverseProxy
	BaseProxy
}

type transport struct {
	current *http.Request
	handler func(r *http.Response)
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	res, _ := http.DefaultTransport.RoundTrip(req)
	if res == nil || !isOkStatus(res.StatusCode) {
		return nil, errors.New("Error Status")
	}

	t.handler(res)

	return res, nil
}

func (p *RetryProxy) Serve(res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(p.config.Target)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

	p.target.ErrorHandler = p.errorHandler

	p.target.Transport = &transport{handler: p.ResponseHandler}
	p.target.ServeHTTP(res, req)
}

func NewRetryProxy(c *ProxyConfig) *RetryProxy {
	tUrl, err := url.Parse(c.Target)

	if err != nil {
		panic(errors.New("Can not parse target url: " + c.Target))
	}

	rUrl, err := url.Parse(c.Reserve)

	if err != nil {
		panic(errors.New("Can not parse reserve url: " + c.Reserve))
	}

	pr := RetryProxy{config: c, target: httputil.NewSingleHostReverseProxy(tUrl), reserve: httputil.NewSingleHostReverseProxy(rUrl)}

	return &pr
}

func (p *RetryProxy) errorHandler(res http.ResponseWriter, req *http.Request, e error) {
	p.reserve.ServeHTTP(res, req)
}
