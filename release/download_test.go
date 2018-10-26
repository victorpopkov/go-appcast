package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDownload(t *testing.T) {
	d := NewDownload("https://example.com/", "application/octet-stream", 100000)
	assert.IsType(t, Download{}, *d)
	assert.Equal(t, "https://example.com/", d.URL)
	assert.Equal(t, "application/octet-stream", d.Type)
	assert.Equal(t, 100000, d.Length)
}
