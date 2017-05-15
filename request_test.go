package appcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	r, err := NewRequest("http://example.com/")
	assert.Nil(t, err)
	assert.IsType(t, Request{}, *r)
	assert.Equal(t, "GET", r.HTTPRequest.Method)
	assert.Equal(t, "http://example.com/", r.HTTPRequest.URL.String())

	// test "Invalid URL" error
	r, err = NewRequest("http://192.168.0.%31/")
	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, "parse http://192.168.0.%31/: invalid URL escape \"%31\"", err.Error())
}

func TestAddHeader(t *testing.T) {
	r, _ := NewRequest("http://example.com/")

	// before
	assert.Len(t, r.HTTPRequest.Header, 0)

	// add header
	r.AddHeader("User-Agent", "Example")

	// after
	headers := r.HTTPRequest.Header
	assert.Len(t, headers, 1)
	assert.Equal(t, "Example", headers.Get("User-Agent"))
}
