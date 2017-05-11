package appcast

// ChecksumAlgorithm holds different available checksum algorithms.
type ChecksumAlgorithm int

const (
	// Sha256 represents a SHA256 checksum
	Sha256 ChecksumAlgorithm = iota

	// Sha256HomebrewCask represents a SHA256 checksum used in Homebrew-Cask
	Sha256HomebrewCask

	// Md5 represents an MD5 checksum
	Md5
)

var checksumAlgorithmNames = [...]string{
	"SHA256",
	"SHA256 (Homebrew-Cask checkpoint)",
	"MD5",
}

// String returns a string representation of the ChecksumAlgorithm.
func (a ChecksumAlgorithm) String() string {
	return checksumAlgorithmNames[a]
}
