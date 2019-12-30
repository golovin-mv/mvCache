package proxy

import (
	"errors"
	"net/http"

	"github.com/golovin-mv/mvCache/mutation"
)

type ProxyConfig struct {
	Type        string
	Target      string
	Reserve     string
	CacheErrors bool
}

type BaseProxy struct {
	ResponseHandler func(r *http.Response)
}

type Proxy interface {
	Serve(res http.ResponseWriter, req *http.Request)
}

func NewProxy(c *ProxyConfig, rMu []mutation.ResponseMutation) Proxy {
	var p Proxy
	switch c.Type {
	case "reverse":
		revP := &ReverseProxy{config: c}
		if rMu != nil && len(rMu) > 0 {
			revP.ResponseHandler = func(res *http.Response) {
				for _, r := range rMu {
					r.Change(res)
				}
			}
		}
		p = revP
	case "retry":
		retP := NewRetryProxy(c)
		if rMu != nil && len(rMu) > 0 {
			retP.ResponseHandler = func(res *http.Response) {
				for _, r := range rMu {
					r.Change(res)
				}
			}
		}

		p = retP
	default:
		panic(errors.New("Unknown Proxy Type"))
	}

	return p
}

func headerToArray(header http.Header) map[string]string {
	res := make(map[string]string)
	for name, values := range header {
		for _, value := range values {
			res[name] = value
		}
	}
	return res
}

func isOkStatus(status int) bool {
	return status >= 200 && status <= 299
}
