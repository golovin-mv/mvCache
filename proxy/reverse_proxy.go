package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/golovin-mv/mvCache/cache"
)

type ReverseProxy struct {
	config *ProxyConfig
}

// Serve a reverse proxy for a given url
func (p *ReverseProxy) Serve(res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(p.config.Target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = makeHandler(cache.GetKey(req), p.config.CacheErrors)
	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
