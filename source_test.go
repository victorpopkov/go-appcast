package appcast

import (
	"testing"
	"encoding/hex"

	"github.com/stretchr/testify/assert"
)

// newTestSource creates a new Source instance for testing purposes and returns
// its pointer.
func newTestSource() *Source {
	content := []byte("test")
	s := &Source{
		content: content,
		checksum: &Checksum{
			algorithm: SHA256,
			source:    content,
			result:    []byte("test"),
		},
		provider: Unknown,
	}

	return s
}

func TestSource_Load(t *testing.T) {
	s := newTestSource()
	assert.Panics(t, func() {
		s.Load()
	})
}

func TestSource_Content(t *testing.T) {
	s := newTestSource()
	assert.Equal(t, []byte("test"), s.Content())
}

func TestSource_SetContent(t *testing.T) {
	s := newTestSource()
	s.SetContent([]byte("new test"))
	assert.Equal(t, []byte("new test"), s.content)
}

func TestSource_Checksum(t *testing.T) {
	s := newTestSource()
	assert.Equal(t, hex.EncodeToString([]byte("test")), s.Checksum().String())
}

func TestSource_Provider(t *testing.T) {
	s := newTestSource()
	assert.Equal(t, Unknown, s.Provider())
}

func TestSource_SetProvider(t *testing.T) {
	s := newTestSource()
	s.SetProvider(SparkleRSSFeed)
	assert.Equal(t, SparkleRSSFeed, s.provider)
}

func TestSource_GenerateChecksum(t *testing.T) {
	// preparations
	s := newTestSource()
	assert.Equal(t, hex.EncodeToString([]byte("test")), s.Checksum().String())

	// test
	s.GenerateChecksum(SHA256)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", s.Checksum().String())
}

func TestSource_GuessProvider(t *testing.T) {
	// test (Unknown)
	s := newTestSource()
	assert.Equal(t, Unknown, s.Provider())

	// test (SparkleRSSFeed)
	s.SetContent(getTestdata("sparkle/default.xml"))
	s.GuessProvider()
	assert.Equal(t, SparkleRSSFeed, s.Provider())

	// test (SourceForgeRSSFeed)
	s.SetContent(getTestdata("sourceforge/default.xml"))
	s.GuessProvider()
	assert.Equal(t, SourceForgeRSSFeed, s.Provider())

	// test (GitHubAtomFeed)
	s.SetContent(getTestdata("github/default.xml"))
	s.GuessProvider()
	assert.Equal(t, GitHubAtomFeed, s.Provider())
}
