package source

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/client"
	"github.com/victorpopkov/go-appcast/provider"
)

// newTestRemote creates a new Remote instance for testing purposes and returns
// its pointer. By default the content is []byte("test"). However, own content
// can be provided as an argument.
func newTestRemote(content ...interface{}) *Remote {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://example.com/appcast.xml"
	r, _ := client.NewRequest(url)

	s := new(appcaster.Source)
	s.SetContent(resultContent)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(provider.Unknown)

	return &Remote{
		Source:  s,
		request: r,
		url:     url,
	}
}

func TestNewRemote(t *testing.T) {
	// preparations
	url := "https://example.com/appcast.xml"

	// test (successful) [URL]
	src, err := NewRemote(url)
	assert.Nil(t, err)
	assert.IsType(t, Remote{}, *src)
	assert.NotNil(t, src.Source)
	assert.NotNil(t, src.request)
	assert.Equal(t, url, src.url)

	// test (successful) [Request]
	r, _ := client.NewRequest(url)
	src, err = NewRemote(r)
	assert.Nil(t, err)
	assert.IsType(t, Remote{}, *src)
	assert.NotNil(t, src.Source)
	assert.NotNil(t, src.request)
	assert.Equal(t, url, src.url)

	// test (error)
	url = "http://192.168.0.%31/"
	src, err = NewRemote(url)
	assert.Nil(t, src)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
}

func TestRemote_Load(t *testing.T) {
	// preparations
	url := "https://example.com/appcast.xml"
	content := []byte("test")

	// mock the request
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// test (successful)
	src, err := NewRemote(url)
	assert.Nil(t, err)
	err = src.Load()
	assert.Nil(t, err)
	assert.Nil(t, src.Provider())
	assert.Equal(t, content, src.Content())

	// test (error)
	src = newTestRemote()
	src.request.HTTPRequest.URL = nil
	err = src.Load()
	assert.NotNil(t, err)
	assert.Equal(t, provider.Unknown, src.Provider())
	assert.Equal(t, []byte("test"), src.Content())
}

func TestRemote_Request(t *testing.T) {
	src := newTestRemote()
	assert.Equal(t, src.request, src.Request())
}

func TestRemote_SetRequest(t *testing.T) {
	src := newTestRemote()
	src.SetRequest(nil)
	assert.Empty(t, src.request)
}

func TestRemote_Url(t *testing.T) {
	src := newTestRemote()
	assert.Equal(t, src.url, src.Url())
}

func TestRemote_SetUrl(t *testing.T) {
	src := newTestRemote()
	src.SetUrl("")
	assert.Empty(t, src.url)
}
