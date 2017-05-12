package appcast

import "net/http"

// A Request represents an HTTP request to be sent by a Client to the server.
type Request struct {
	// HTTPRequest specifies the http.Request to be sent to the remote server. It
	// includes all request configuration such as URL, protocol version, HTTP
	// method, request headers and authentication.
	HTTPRequest *http.Request
}

// NewRequest returns a new Request instance pointer and an error. The GET
// method is used.
func NewRequest(url string) (*Request, error) {
	// create http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return &Request{
		HTTPRequest: req,
	}, nil
}

// AddHeader adds a new header with specified key and value. The headers will
// be used while making request in Client.Do.
func (r *Request) AddHeader(key string, value string) {
	r.HTTPRequest.Header.Add(key, value)
}
