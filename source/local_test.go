package source

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/provider"
)

// newTestLocal creates a new Local instance for testing purposes and returns
// its pointer. By default the content is []byte("test"). However, own content
// can be provided as an argument.
func newTestLocal(content ...interface{}) *Local {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	s := new(appcaster.Source)
	s.SetContent(resultContent)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(provider.Unknown)

	return &Local{
		Source:   s,
		filepath: "/tmp/test.txt",
	}
}

func TestNewLocal(t *testing.T) {
	// preparations
	path := "/tmp/test.txt"

	// test (successful)
	src := NewLocal(path)
	assert.IsType(t, Local{}, *src)
	assert.NotNil(t, src.Source)
	assert.Equal(t, path, src.filepath)
}

func TestLocal_Load(t *testing.T) {
	// preparations
	path := "test.xml"
	content := []byte("test")

	// test (successful)
	LocalReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	src := NewLocal(path)
	err := src.Load()
	assert.Nil(t, err)
	assert.Equal(t, provider.Unknown, src.Provider())
	assert.Equal(t, string(content), string(src.Content()))

	// test (error)
	LocalReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}

	src = newTestLocal()
	err = src.Load()
	assert.NotNil(t, err)
	assert.Equal(t, provider.Unknown, src.Provider())
	assert.Equal(t, []byte("test"), src.Content())

	LocalReadFile = ioutil.ReadFile
}

func TestLocal_Filepath(t *testing.T) {
	src := newTestLocal()
	assert.Equal(t, src.filepath, src.Filepath())
}

func TestLocal_SetFilepath(t *testing.T) {
	src := newTestLocal()
	src.SetFilepath("")
	assert.Empty(t, src.filepath)
}
