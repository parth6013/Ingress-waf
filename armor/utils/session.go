package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
)

// Session is a wrapper over builtin http module from go
type Session struct {
	Client  http.Client
	Headers map[string]string
}

// NewSession returs a new session instance
func NewSession() *Session {
	return &Session{
		Client:  http.Client{},
		Headers: map[string]string{},
	}
}

// DissableSSL dissables ssl verification for the client
func (s *Session) DissableSSL() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	s.Client.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}

// BasicAuth adds basic auth to the client headders
func (s *Session) BasicAuth(username string, password string) {
	auth := username + ":" + password
	s.Headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// createNewRequest Creats a new request object
func (s *Session) createNewRequest(method string, url string, body *[]byte) (*http.Request, error) {
	if body != nil {
		return http.NewRequest(method, url, bytes.NewBuffer(*body))
	}
	return http.NewRequest(method, url, nil)
}

// SetHeaders adds headers to the request
func (s *Session) SetHeaders(req *http.Request, headers *map[string]string) {
	// Adding default heaers
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	// Adding headers passed to the function
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
}

// SetQueryParams adds query parameter to the request. Query parameters need to be passed as a map of string to string
func (s *Session) SetQueryParams(req *http.Request, query *map[string]string) {
	if query != nil {
		q := req.URL.Query()
		for k, v := range *query {
			q.Add(k, v)
		}
		// Re-encode URL
		req.URL.RawQuery = q.Encode()
	}
}

// GetBody returns the string representation of the response boddy and closes the response body buffer
func (s *Session) GetBody(resp *http.Response) (string, error) {
	// Read the respons body and covert it to a string
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return string(respBody), nil
}

// Request returns a http.Response object with the proper method, url, body, headers and query params
func (s *Session) Request(method string, url string, body *[]byte, headers *map[string]string, query *map[string]string) (*http.Response, error) {
	// Prepare request body
	req, err := s.createNewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	s.SetHeaders(req, headers)

	// Set query params
	s.SetQueryParams(req, query)

	// Make HTTP call
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response
	return resp, nil
}

// Get resuest
func (s *Session) Get(url string, headers *map[string]string, query *map[string]string) (*http.Response, error) {
	// Make get request
	return s.Request("GET", url, nil, headers, query)
}

// Post resuest
func (s *Session) Post(url string, body *[]byte, headers *map[string]string, query *map[string]string) (*http.Response, error) {
	// Make post request
	return s.Request("POST", url, body, headers, query)
}

// Put resuest
func (s *Session) Put(url string, body *[]byte, headers *map[string]string, query *map[string]string) (*http.Response, error) {
	// Make put request
	return s.Request("PUT", url, body, headers, query)
}

// Delete resuest
func (s *Session) Delete(url string, body *[]byte, headers *map[string]string, query *map[string]string) (*http.Response, error) {
	// Make delete request
	return s.Request("DELETE", url, body, headers, query)
}
