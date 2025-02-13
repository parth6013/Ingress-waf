package handlers

import (
	"armor/models"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxyHandler(config *models.Config) *ProxyHandler {
	target := &url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", config.Proxy.TargetHost, config.Proxy.TargetPort)}

	return &ProxyHandler{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: forwarding request to %s", p.target.Host)
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	// Use a waitgroup to make sure no inflight requests are cancelled during application shutdown
	p.proxy.ServeHTTP(w, r)
}
