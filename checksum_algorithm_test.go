package appcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumAlgorithmString(t *testing.T) {
	assert.Equal(t, "SHA256", SHA256.String())
	assert.Equal(t, "MD5", MD5.String())
}
