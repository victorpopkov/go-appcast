package appcast

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"regexp"
)

// A Checksum holds everything needed to create a hash checksum.
type Checksum struct {
	// Algorithm specifies which one of the supported algorithms is used to
	// generate the resulting checksum.
	Algorithm ChecksumAlgorithm

	// Source specifies the string from which the checksum will be generated.
	Source string

	// Result represents the resulting checksum itself.
	Result string
}

// NewChecksum returns a new Checksum instance pointer. Requires an algorithm
// that will be used to generate the checksum and a source from which it will
// be generated.
func NewChecksum(algorithm ChecksumAlgorithm, source string) *Checksum {
	c := &Checksum{
		Algorithm: algorithm,
		Source:    source,
	}
	c.Generate()

	return c
}

// Generate creates and returns a checksum represented as a string from the
// Source field using the chosen algorithm. The checksum is also stored in the
// Result field of the Checksum.
func (c *Checksum) Generate() string {
	var hasher hash.Hash

	switch c.Algorithm {
	case Sha256:
		hasher = sha256.New()
		hasher.Write([]byte(c.Source))
	case Sha256HomebrewCask:
		re := regexp.MustCompile(`<pubDate>[^<]*<\/pubDate>`)
		sourceMod := re.ReplaceAllString(c.Source, "")
		hasher = sha256.New()
		hasher.Write([]byte(sourceMod))
	case Md5:
		hasher = md5.New()
		hasher.Write([]byte(c.Source))
	}

	c.Result = hex.EncodeToString(hasher.Sum(nil))

	return c.Result
}
