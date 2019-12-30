package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/golovin-mv/mvCache/cache"
)

type ReverseProxy struct {
	config *ProxyConfig
	BaseProxy
}

// Serve a reverse proxy for a given url
func (p *ReverseProxy) Serve(res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(p.config.Target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = p.makeHandler(cache.GetKey(req), p.config.CacheErrors)
	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func (p *ReverseProxy) makeHandler(key string, cacheError bool) func(r *http.Response) error {
	return func(r *http.Response) error {
		if p.BaseProxy.ResponseHandler != nil {
			p.BaseProxy.ResponseHandler(r)
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return err
		}

		// если статус 200 или не 200 но мы кэшируем ошибки
		if isOkStatus(r.StatusCode) || (!isOkStatus(r.StatusCode) && cacheError) {
			cache.CurrentCacher.Add(key, &cache.CachedResponse{headerToArray(r.Header), body})
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return nil
	}
}
