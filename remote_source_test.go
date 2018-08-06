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

	s := &RemoteSource{
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

	return s
}

func TestNewRemoteSource(t *testing.T) {
	// preparations
	url := "https://example.com/appcast.xml"

	// test (successful) [URL]
	s, err := NewRemoteSource(url)
	assert.Nil(t, err)
	assert.IsType(t, RemoteSource{}, *s)
	assert.NotNil(t, s.Source)
	assert.NotNil(t, s.request)
	assert.Equal(t, url, s.url)

	// test (successful) [Request]
	r, _ := NewRequest(url)
	s, err = NewRemoteSource(r)
	assert.Nil(t, err)
	assert.IsType(t, RemoteSource{}, *s)
	assert.NotNil(t, s.Source)
	assert.NotNil(t, s.request)
	assert.Equal(t, url, s.url)

	// test (error)
	url = "http://192.168.0.%31/"
	s, err = NewRemoteSource(url)
	assert.Nil(t, s)
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
	s, err := NewRemoteSource(url)
	assert.Nil(t, err)
	err = s.Load()
	assert.Nil(t, err)
	assert.Equal(t, SparkleRSSFeed, s.provider)
	assert.Equal(t, content, s.content)

	// test (error)
	s = newTestRemoteSource()
	s.request.HTTPRequest.URL = nil
	err = s.Load()
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, s.provider)
	assert.Equal(t, []byte("test"), s.content)
}

func TestRemoteSource_GuessProvider(t *testing.T) {
	// test (Unknown)
	s := newTestRemoteSource()
	s.GuessProvider()
	assert.Equal(t, Unknown, s.Provider())

	// test (SparkleRSSFeed)
	s = newTestRemoteSource(getTestdata("sparkle/default.xml"))
	s.GuessProvider()
	assert.Equal(t, SparkleRSSFeed, s.Provider())

	// test (SourceForgeRSSFeed)
	s = newTestRemoteSource(getTestdata("sourceforge/default.xml"))
	s.GuessProvider()
	assert.Equal(t, SourceForgeRSSFeed, s.Provider())

	// test (GitHubAtomFeed)
	s = newTestRemoteSource(getTestdata("github/default.xml"))
	s.GuessProvider()
	assert.Equal(t, GitHubAtomFeed, s.Provider())
}

func TestRemoteSource_Request(t *testing.T) {
	s := newTestRemoteSource()
	assert.Equal(t, s.request, s.Request())
}

func TestRemoteSource_Url(t *testing.T) {
	s := newTestRemoteSource()
	assert.Equal(t, s.url, s.Url())
}
