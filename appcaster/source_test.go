package appcaster

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestSource creates a new Source instance for testing purposes and returns
// its pointer.
func newTestSource() *Source {
	content := []byte("content")

	return &Source{
		content: content,
		checksum: &Checksum{
			algorithm: SHA256,
			source:    content,
			result:    []byte("test"),
		},
		provider: Provider(0),
		appcast:  &Appcast{},
	}
}

func TestSource_Load(t *testing.T) {
	src := newTestSource()
	assert.Panics(t, func() {
		src.Load()
	})
}

func TestSource_GenerateChecksum(t *testing.T) {
	// preparations
	src := newTestSource()
	assert.Equal(t, hex.EncodeToString([]byte("test")), src.Checksum().String())
	expected := "ed7002b439e9ac845f22357d822bac1444730fbdb6016d3ec9432297b9ec9f73"

	// test
	c := src.GenerateChecksum(SHA256)
	assert.Equal(t, expected, c.String())
	assert.Equal(t, expected, src.Checksum().String())
}

func TestSource_Content(t *testing.T) {
	src := newTestSource()
	assert.Equal(t, src.content, src.Content())
}

func TestSource_SetContent(t *testing.T) {
	src := newTestSource()
	src.SetContent([]byte("new test"))
	assert.Equal(t, []byte("new test"), src.content)
}

func TestSource_Checksum(t *testing.T) {
	src := newTestSource()
	assert.Equal(t, hex.EncodeToString([]byte("test")), src.Checksum().String())
}

func TestSource_Provider(t *testing.T) {
	src := newTestSource()
	assert.Equal(t, src.provider, src.Provider())
}

func TestSource_SetProvider(t *testing.T) {
	// preparations
	src := newTestSource()
	p := Provider(1)

	// test
	src.SetProvider(p)
	assert.Equal(t, p, src.provider)
}

func TestSource_Appcast(t *testing.T) {
	src := newTestSource()
	assert.Equal(t, src.appcast, src.Appcast())
}

func TestSource_SetAppcast(t *testing.T) {
	src := newTestSource()
	src.SetAppcast(nil)
	assert.Nil(t, src.appcast)
}
