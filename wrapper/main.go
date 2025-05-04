package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// NewProxy creates a reverse proxy to the target service
func NewProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	// Modify the Director function to ensure requests are forwarded correctly
	// url.Host = "localhost:8080"
	// url.Scheme = "http"
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = url.Path + req.URL.Path
	}

	//If the proxy fails to forward the request
	// Error handler to catch failed requests
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Request failed: %s %s - Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "Service unavailable", http.StatusBadGateway)
	}

	return proxy
}

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed request: %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	target := "http://localhost:8080"
	proxy := NewProxy(target)

	// Apply logging middleware
	// http.Handle("/", LoggingMiddleware(proxy))
	http.Handle("/", proxy)

	port := "3000"
	log.Printf("Logger proxy running on port %s, forwarding to %s", port, target)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
