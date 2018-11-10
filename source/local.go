package source

import (
	"io/ioutil"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/provider"
)

var LocalReadFile = ioutil.ReadFile

// Localer is the interface that wraps the Local methods.
type Localer interface {
	appcaster.Sourcer
	Filepath() string
	SetFilepath(filepath string)
}

// Local represents an appcast source from the local file.
type Local struct {
	*appcaster.Source
	filepath string
}

// NewLocal returns a new Local instance pointer with the Local.filepath set.
func NewLocal(path string) *Local {
	return &Local{
		Source:   &appcaster.Source{},
		filepath: path,
	}
}

// Load loads an appcast content into the Local.Source.content from the local
// file by using the path specified in Local.filepath set earlier.
func (l *Local) Load() error {
	data, err := LocalReadFile(l.filepath)
	if err != nil {
		return err
	}

	l.SetContent(data)
	l.GuessProvider()
	l.GenerateChecksum(appcaster.SHA256)

	return nil
}

// GuessProvider attempts to guess the supported provider based on the
// Source.content. By default returns an Unknown provider.
func (l *Local) GuessProvider() {
	l.SetProvider(provider.GuessProviderByContent(l.Content()))
}

// Filepath is a Local.filepath getter.
func (l *Local) Filepath() string {
	return l.filepath
}

// SetFilepath is a Local.filepath setter.
func (l *Local) SetFilepath(filepath string) {
	l.filepath = filepath
}
