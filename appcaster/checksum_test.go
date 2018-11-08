package appcaster

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecksum(t *testing.T) {
	// preparations
	a := SHA256
	src := []byte("test")

	// test
	c := NewChecksum(a, src)
	assert.IsType(t, Checksum{}, *c)
	assert.Equal(t, a, c.algorithm)
	assert.Equal(t, src, c.source)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.String())
}

func TestChecksum_Algorithm(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, c.algorithm, c.Algorithm())
}

func TestChecksum_Source(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, c.source, c.Source())
}

func TestChecksum_Result(t *testing.T) {
	// preparations
	c := NewChecksum(SHA256, []byte("test"))
	expected, _ := hex.DecodeString("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")

	// test
	assert.Equal(t, expected, c.Result())
}

func TestChecksumAlgorithm_String(t *testing.T) {
	assert.Equal(t, "SHA256", SHA256.String())
	assert.Equal(t, "MD5", MD5.String())
}

func TestChecksum_String(t *testing.T) {
	c := NewChecksum(SHA256, []byte("test"))
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.String())
}
