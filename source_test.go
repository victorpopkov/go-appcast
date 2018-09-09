package appcast

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestSource creates a new Source instance for testing purposes and returns
// its pointer.
func newTestSource() *Source {
	content := []byte("test")
	src := &Source{
		content: content,
		checksum: &Checksum{
			algorithm: SHA256,
			source:    content,
			result:    []byte("test"),
		},
		provider: Unknown,
	}

	return src
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

	// test
	src.GenerateChecksum(SHA256)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", src.Checksum().String())
}

func TestSource_GuessProvider(t *testing.T) {
	// test (Unknown)
	src := newTestSource()
	assert.Equal(t, Unknown, src.Provider())

	// test (SparkleRSSFeed)
	src.SetContent(getTestdata("sparkle/default.xml"))
	src.GuessProvider()
	assert.Equal(t, SparkleRSSFeed, src.Provider())

	// test (SourceForgeRSSFeed)
	src.SetContent(getTestdata("sourceforge/default.xml"))
	src.GuessProvider()
	assert.Equal(t, SourceForgeRSSFeed, src.Provider())

	// test (GitHubAtomFeed)
	src.SetContent(getTestdata("github/default.xml"))
	src.GuessProvider()
	assert.Equal(t, GitHubAtomFeed, src.Provider())
}

func TestSource_Content(t *testing.T) {
	src := newTestSource()
	assert.Equal(t, []byte("test"), src.Content())
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
	assert.Equal(t, Unknown, src.Provider())
}

func TestSource_SetProvider(t *testing.T) {
	src := newTestSource()
	src.SetProvider(SparkleRSSFeed)
	assert.Equal(t, SparkleRSSFeed, src.provider)
}
