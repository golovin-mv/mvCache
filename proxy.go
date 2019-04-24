package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Serve a reverse proxy for a given url
func ServeReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = makeHandler(GetKey(req))
	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
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

func makeHandler(key string) func(r *http.Response) error {
	return func(r *http.Response) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return err
		}

		CurrentCacher.Add(key, &CachedResponse{headerToArray(r.Header), body})

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		return nil
	}
}
