package appcast

import "io/ioutil"

// LocalSourcer is the interface that wraps the LocalSource methods.
type LocalSourcer interface {
	Sourcer
	Filepath() string
}

// LocalSource represents an appcast source from the local file.
type LocalSource struct {
	*Source
	filepath string
}

// NewLocalSource returns a new LocalSource instance pointer with the
// LocalSource.filepath set.
func NewLocalSource(path string) *LocalSource {
	src := &LocalSource{
		Source:   &Source{},
		filepath: path,
	}

	return src
}

// Load loads an appcast content into the LocalSource.Source.content from the
// local file by using the path specified in LocalSource.filepath set earlier.
func (s *LocalSource) Load() error {
	data, err := ioutil.ReadFile(s.filepath)
	if err != nil {
		return err
	}

	s.content = data
	s.GuessProvider()
	s.checksum = NewChecksum(SHA256, s.content)

	return nil
}

// Filepath is a LocalSource.filepath getter.
func (s *LocalSource) Filepath() string {
	return s.filepath
}
