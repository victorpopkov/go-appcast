package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestDownload creates a new Download instance for testing purposes and
// returns its pointer.
func newTestDownload() *Download {
	return &Download{
		url:      "https://example.com/1.0.0/app.dmg",
		filetype: "application/octet-stream",
		length:   100000,
	}
}

func TestNewDownload(t *testing.T) {
	d := NewDownload("https://example.com/1.0.0/app.dmg", "application/octet-stream", 100000)
	assert.IsType(t, Download{}, *d)
	assert.Equal(t, "https://example.com/1.0.0/app.dmg", d.url)
	assert.Equal(t, "application/octet-stream", d.filetype)
	assert.Equal(t, 100000, d.length)
}

func TestDownload_Url(t *testing.T) {
	d := newTestDownload()
	assert.Equal(t, "https://example.com/1.0.0/app.dmg", d.Url())
}

func TestDownload_SetUrl(t *testing.T) {
	d := newTestDownload()
	d.SetUrl("https://example.com/2.0.0/app.dmg")
	assert.Equal(t, "https://example.com/2.0.0/app.dmg", d.url)
}

func TestDownload_Filetype(t *testing.T) {
	d := newTestDownload()
	assert.Equal(t, "application/octet-stream", d.Filetype())
}

func TestDownload_SetFiletype(t *testing.T) {
	d := newTestDownload()
	d.SetFiletype("application/x-bzip2; charset=binary")
	assert.Equal(t, "application/x-bzip2; charset=binary", d.filetype)
}

func TestDownload_Length(t *testing.T) {
	d := newTestDownload()
	assert.Equal(t, 100000, d.Length())
}

func TestDownload_SetLength(t *testing.T) {
	d := newTestDownload()
	d.SetLength(200000)
	assert.Equal(t, 200000, d.length)
}
