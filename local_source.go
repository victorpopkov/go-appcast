package appcast

import (
	"io/ioutil"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/provider"
)

var localSourceReadFile = ioutil.ReadFile

// LocalSourcer is the interface that wraps the LocalSource methods.
type LocalSourcer interface {
	appcaster.Sourcer
	Filepath() string
}

// LocalSource represents an appcast source from the local file.
type LocalSource struct {
	*appcaster.Source
	filepath string
}

// NewLocalSource returns a new LocalSource instance pointer with the
// LocalSource.filepath set.
func NewLocalSource(path string) *LocalSource {
	src := &LocalSource{
		Source:   &appcaster.Source{},
		filepath: path,
	}

	return src
}

// Load loads an appcast content into the LocalSource.Source.content from the
// local file by using the path specified in LocalSource.filepath set earlier.
func (s *LocalSource) Load() error {
	data, err := localSourceReadFile(s.filepath)
	if err != nil {
		return err
	}

	s.SetContent(data)
	s.GuessProvider()
	s.GenerateChecksum(appcaster.SHA256)

	return nil
}

// GuessProvider attempts to guess the supported provider based on the
// Source.content. By default returns an Unknown provider.
func (s *LocalSource) GuessProvider() {
	s.SetProvider(provider.GuessProviderByContent(s.Content()))
}

// Filepath is a LocalSource.filepath getter.
func (s *LocalSource) Filepath() string {
	return s.filepath
}
