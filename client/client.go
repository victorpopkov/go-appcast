// Package client provides functionality for making remote client requests.
package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Client represents an HTTP client that makes requests.
type Client struct {
	// HTTPClient specifies the http.Client which will be used for communicating
	// with the remote server during the file transfer.
	HTTPClient *http.Client

	// UserAgent specifies the User-Agent string which will be set in the headers
	// of all requests made by this client.
	//
	// The user agent string may be overridden in the headers of each request.
	UserAgent string

	// Timeout specifies a time limit for the client request to happen.
	Timeout time.Duration
}

// New returns a new Client instance pointer with the default Client.UserAgent
// and Client.HTTPClient.
func New() *Client {
	return &Client{
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		HTTPClient: http.DefaultClient,
	}
}

// InsecureSkipVerify forces client to bypass TLS certificates verifications.
// Not recommended though for security reasons.
func (c *Client) InsecureSkipVerify() {
	c.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

// Do sends an HTTP request and returns an HTTP response using http.Client.Do. A
// successful call returns a nil error.
func (c *Client) Do(req *Request) (*http.Response, error) {
	// set UserAgent
	if c.UserAgent != "" && req.HTTPRequest.Header.Get("User-Agent") == "" {
		req.HTTPRequest.Header.Set("User-Agent", c.UserAgent)
	}

	// set Timeout
	if c.Timeout > 0 {
		c.HTTPClient.Timeout = time.Duration(10 * time.Second)
	}

	// make request and return response
	resp, err := c.HTTPClient.Do(req.HTTPRequest)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
