package appcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumAlgorithmString(t *testing.T) {
	assert.Equal(t, "SHA256", Sha256.String())
	assert.Equal(t, "SHA256 (Homebrew-Cask checkpoint)", Sha256HomebrewCask.String())
	assert.Equal(t, "MD5", Md5.String())
}
