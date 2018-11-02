package appcast

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestOutput creates a new Output instance for testing purposes and returns
// its pointer.
func newTestOutput() *Output {
	content := []byte("test")
	o := &Output{
		content: content,
		checksum: &Checksum{
			algorithm: SHA256,
			source:    content,
			result:    []byte("test"),
		},
		provider: Unknown,
		appcast:  &SparkleAppcast{},
	}

	return o
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

	// test
	o.GenerateChecksum(SHA256)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", o.Checksum().String())
}

func TestOutput_Content(t *testing.T) {
	o := newTestOutput()
	assert.Equal(t, []byte("test"), o.Content())
}

func TestOutput_SetContent(t *testing.T) {
	o := newTestOutput()
	o.SetContent([]byte("new test"))
	assert.Equal(t, []byte("new test"), o.content)
}

func TestOutput_Provider(t *testing.T) {
	o := newTestOutput()
	assert.Equal(t, Unknown, o.Provider())
}

func TestOutput_SetProvider(t *testing.T) {
	o := newTestOutput()
	o.SetProvider(Sparkle)
	assert.Equal(t, Sparkle, o.provider)
}

func TestOutput_Appcast(t *testing.T) {
	o := newTestOutput()
	assert.IsType(t, &SparkleAppcast{}, o.Appcast())
}

func TestOutput_SetAppcast(t *testing.T) {
	o := newTestOutput()
	o.SetAppcast(&GitHubAtomFeedAppcast{})
	assert.IsType(t, &GitHubAtomFeedAppcast{}, o.appcast)
}
