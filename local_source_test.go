package appcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	s := &LocalSource{
		Source: &Source{
			content: resultContent,
			checksum: &Checksum{
				algorithm: SHA256,
				source:    resultContent,
				result:    []byte("test"),
			},
			provider: Unknown,
		},
		filepath: "/tmp/test.txt",
	}

	return s
}

func TestNewLocalSource(t *testing.T) {
	// preparations
	filepath := "/tmp/test.txt"

	// test (successful)
	s := NewLocalSource(filepath)
	assert.IsType(t, LocalSource{}, *s)
	assert.NotNil(t, s.Source)
	assert.Equal(t, filepath, s.filepath)
}

func TestLocalSource_Load(t *testing.T) {
	// preparations
	filepath := getTestdataPath("sparkle/default.xml")
	content := getTestdata("sparkle/default.xml")

	// test (successful)
	s := NewLocalSource(filepath)
	err := s.Load()
	assert.Nil(t, err)
	assert.Equal(t, SparkleRSSFeed, s.provider)
	assert.Equal(t, string(content), string(s.content))

	// test (error)
	s = newTestLocalSource()
	err = s.Load()
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, s.provider)
	assert.Equal(t, []byte("test"), s.content)
}

func TestLocalSource_Filepath(t *testing.T) {
	s := newTestLocalSource()
	assert.Equal(t, s.filepath, s.Filepath())
}
