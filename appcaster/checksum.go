package appcaster

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// ChecksumAlgorithm holds different available checksum algorithms.
type ChecksumAlgorithm int

const (
	// SHA256 represents a SHA256 checksum
	SHA256 ChecksumAlgorithm = iota

	// MD5 represents an MD5 checksum
	MD5
)

var checksumAlgorithmNames = [...]string{
	"SHA256",
	"MD5",
}

// Checksum holds everything needed to create a hash checksum.
type Checksum struct {
	// algorithm specifies which one of the supported algorithms is used to
	// generate the resulting checksum.
	algorithm ChecksumAlgorithm

	// source specifies the data from which the checksum will be generated.
	source []byte

	// result represents checksum itself.
	result []byte
}

// NewChecksum returns a new Checksum instance pointer. Requires an algorithm
// that will be used to generate the checksum and a source from which it will
// be generated.
func NewChecksum(algorithm ChecksumAlgorithm, src []byte) *Checksum {
	c := &Checksum{
		algorithm: algorithm,
		source:    src,
	}

	c.generate()

	return c
}

// generate creates a checksum from the Checksum.source field using the
// specified algorithm as Checksum.algorithm and stores the generated checksum
// into the Checksum.result.
//
// This method is called in Checksum.NewChecksum.
func (c *Checksum) generate() {
	var hasher hash.Hash

	switch c.algorithm {
	case SHA256:
		hasher = sha256.New()
		hasher.Write(c.source)
	case MD5:
		hasher = md5.New()
		hasher.Write(c.source)
	}

	c.result = hasher.Sum(nil)
}

// Algorithm is a Checksum.algorithm getter.
func (c *Checksum) Algorithm() ChecksumAlgorithm {
	return c.algorithm
}

// Source is a Checksum.source getter.
func (c *Checksum) Source() []byte {
	return c.source
}

// Result is a Checksum.result getter.
func (c *Checksum) Result() []byte {
	return c.result
}

// String returns a string representation of the ChecksumAlgorithm.
func (a ChecksumAlgorithm) String() string {
	return checksumAlgorithmNames[a]
}

// String returns a string representation of the Checksum.result.
func (c *Checksum) String() string {
	return hex.EncodeToString(c.result)
}
