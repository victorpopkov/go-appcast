package appcast

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestOutput creates a new Output instance for testing purposes and returns
// its pointer.
func newTestOutput() *Output {
	content := []byte("content")

	return &Output{
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

func TestOutput_Save(t *testing.T) {
	o := newTestOutput()
	assert.Panics(t, func() {
		o.Save()
	})
}

func TestOutput_GenerateChecksum(t *testing.T) {
	// preparations
	o := newTestOutput()
	assert.Equal(t, hex.EncodeToString([]byte("test")), o.Checksum().String())
	expected := "ed7002b439e9ac845f22357d822bac1444730fbdb6016d3ec9432297b9ec9f73"

	// test
	c := o.GenerateChecksum(SHA256)
	assert.Equal(t, expected, c.String())
	assert.Equal(t, expected, o.Checksum().String())
}

func TestOutput_Content(t *testing.T) {
	o := newTestOutput()
	assert.Equal(t, o.content, o.Content())
}

func TestOutput_SetContent(t *testing.T) {
	// preparations
	o := newTestOutput()
	c := []byte("new test")

	// test
	o.SetContent(c)
	assert.Equal(t, c, o.content)
}

func TestOutput_Provider(t *testing.T) {
	o := newTestOutput()
	assert.Equal(t, o.provider, o.Provider())
}

func TestOutput_SetProvider(t *testing.T) {
	// preparations
	o := newTestOutput()
	p := Provider(1)

	// test
	o.SetProvider(p)
	assert.Equal(t, p, o.provider)
}

func TestOutput_Appcast(t *testing.T) {
	o := newTestOutput()
	assert.Equal(t, o.appcast, o.Appcast())
}

func TestOutput_SetAppcast(t *testing.T) {
	o := newTestOutput()
	o.SetAppcast(nil)
	assert.Nil(t, o.appcast)
}
