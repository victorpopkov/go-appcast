package appcast

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

// newTestRemoteSource creates a new RemoteSource instance for testing purposes
// and returns its pointer. By default the content is []byte("test"). However,
// own content can be provided as an argument.
func newTestRemoteSource(content ...interface{}) *RemoteSource {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://example.com/appcast.xml"
	r, _ := NewRequest(url)

	src := &RemoteSource{
		Source: &Source{
			content: resultContent,
			checksum: &Checksum{
				algorithm: SHA256,
				source:    resultContent,
				result:    []byte("test"),
			},
			provider: Unknown,
		},
		request: r,
		url:     url,
	}

	return src
}

func TestNewRemoteSource(t *testing.T) {
	// preparations
	url := "https://example.com/appcast.xml"

	// test (successful) [URL]
	src, err := NewRemoteSource(url)
	assert.Nil(t, err)
	assert.IsType(t, RemoteSource{}, *src)
	assert.NotNil(t, src.Source)
	assert.NotNil(t, src.request)
	assert.Equal(t, url, src.url)

	// test (successful) [Request]
	r, _ := NewRequest(url)
	src, err = NewRemoteSource(r)
	assert.Nil(t, err)
	assert.IsType(t, RemoteSource{}, *src)
	assert.NotNil(t, src.Source)
	assert.NotNil(t, src.request)
	assert.Equal(t, url, src.url)

	// test (error)
	url = "http://192.168.0.%31/"
	src, err = NewRemoteSource(url)
	assert.Nil(t, src)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
}

func TestRemoteSource_Load(t *testing.T) {
	// preparations
	url := "https://example.com/appcast.xml"
	content := getTestdata("sparkle/default.xml")

	// mock the request
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// test (successful)
	src, err := NewRemoteSource(url)
	assert.Nil(t, err)
	err = src.Load()
	assert.Nil(t, err)
	assert.Equal(t, SparkleRSSFeed, src.provider)
	assert.Equal(t, content, src.content)

	// test (error)
	src = newTestRemoteSource()
	src.request.HTTPRequest.URL = nil
	err = src.Load()
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, src.provider)
	assert.Equal(t, []byte("test"), src.content)
}

func TestRemoteSource_GuessProvider(t *testing.T) {
	// test (Unknown)
	src := newTestRemoteSource()
	src.GuessProvider()
	assert.Equal(t, Unknown, src.Provider())

	// test (SparkleRSSFeed)
	src = newTestRemoteSource(getTestdata("sparkle/default.xml"))
	src.GuessProvider()
	assert.Equal(t, SparkleRSSFeed, src.Provider())

	// test (SourceForgeRSSFeed)
	src = newTestRemoteSource(getTestdata("sourceforge/default.xml"))
	src.GuessProvider()
	assert.Equal(t, SourceForgeRSSFeed, src.Provider())

	// test (GitHubAtomFeed)
	src = newTestRemoteSource(getTestdata("github/default.xml"))
	src.GuessProvider()
	assert.Equal(t, GitHubAtomFeed, src.Provider())
}

func TestRemoteSource_Request(t *testing.T) {
	src := newTestRemoteSource()
	assert.Equal(t, src.request, src.Request())
}

func TestRemoteSource_Url(t *testing.T) {
	src := newTestRemoteSource()
	assert.Equal(t, src.url, src.Url())
}
