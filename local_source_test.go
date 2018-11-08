package appcast

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
)

// newTestLocalSource creates a new LocalSource instance for testing purposes
// and returns its pointer. By default the content is []byte("test"). However,
// own content can be provided as an argument.
func newTestLocalSource(content ...interface{}) *LocalSource {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	s := new(appcaster.Source)
	s.SetContent(resultContent)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(Unknown)

	return &LocalSource{
		Source:   s,
		filepath: "/tmp/test.txt",
	}
}

func TestNewLocalSource(t *testing.T) {
	// preparations
	path := "/tmp/test.txt"

	// test (successful)
	src := NewLocalSource(path)
	assert.IsType(t, LocalSource{}, *src)
	assert.NotNil(t, src.Source)
	assert.Equal(t, path, src.filepath)
}

func TestLocalSource_Load(t *testing.T) {
	// preparations
	path := getTestdataPath("sparkle/default.xml")
	content := getTestdata("sparkle/default.xml")

	// test (successful)
	localSourceReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	src := NewLocalSource(path)
	err := src.Load()
	assert.Nil(t, err)
	assert.Equal(t, Sparkle, src.Provider())
	assert.Equal(t, string(content), string(src.Content()))

	// test (error)
	localSourceReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}

	src = newTestLocalSource()
	err = src.Load()
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, src.Provider())
	assert.Equal(t, []byte("test"), src.Content())

	localSourceReadFile = ioutil.ReadFile
}

func TestLocalSource_Filepath(t *testing.T) {
	src := newTestLocalSource()
	assert.Equal(t, src.filepath, src.Filepath())
}
