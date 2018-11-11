// Package appcast provides functionality for working with appcasts to retrieve
// valuable information about software releases.
//
// Currently supports 3 providers: "GitHub Atom Feed", "SourceForge RSS Feed"
// and "Sparkle RSS Feed". However, it can be extended to your own needs
// if necessary.
//
// See README.md for more info.
package appcast

import (
	"fmt"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/provider"
	"github.com/victorpopkov/go-appcast/provider/github"
	"github.com/victorpopkov/go-appcast/provider/sourceforge"
	"github.com/victorpopkov/go-appcast/provider/sparkle"
	"github.com/victorpopkov/go-appcast/source"
)

// Appcaster is the interface that wraps the Appcast methods.
//
// This interface should be embedded by provider-specific Appcaster interfaces.
type Appcaster interface {
	appcaster.Appcaster
	LoadFromRemoteSource(i interface{}) (appcaster.Appcaster, error)
	LoadFromLocalSource(path string) (appcaster.Appcaster, error)
}

// Appcast represents the non provider-specific appcast.
type Appcast struct {
	appcaster.Appcast
}

// New returns a new Appcast instance pointer. The Source can be passed as
// a parameter.
func New(src ...interface{}) *Appcast {
	a := new(Appcast)

	if len(src) > 0 {
		src := src[0].(appcaster.Sourcer)
		a.SetSource(src)
	}

	return a
}

// LoadFromRemoteSource creates a new RemoteSource instance and loads the data
// from the remote location by using the RemoteSource.Load method.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) LoadFromRemoteSource(i interface{}) (appcaster.Appcaster, error) {
	src, err := source.NewRemote(i)
	if err != nil {
		return nil, err
	}

	err = src.Load()
	if err != nil {
		return nil, err
	}

	a.SetSource(src)
	a.GuessSourceProvider()

	appcast, err := a.Unmarshal()
	if err != nil {
		return nil, err
	}

	return appcast, nil
}

// LoadFromLocalSource creates a new LocalSource instance and loads the data
// from the local file by using the LocalSource.Load method.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) LoadFromLocalSource(path string) (appcaster.Appcaster, error) {
	src := source.NewLocal(path)
	err := src.Load()
	if err != nil {
		return nil, err
	}

	a.SetSource(src)
	a.GuessSourceProvider()

	appcast, err := a.Unmarshal()
	if err != nil {
		return nil, err
	}

	return appcast, nil
}

// LoadSource calls the Appcast.source.Load method.
func (a *Appcast) LoadSource() error {
	err := a.Source().Load()
	if err == nil {
		a.GuessSourceProvider()
		return nil
	}

	return err
}

// GuessSourceProvider attempts to guess the supported provider based on the
// Appcast.source.content.
func (a *Appcast) GuessSourceProvider() {
	switch src := a.Source().(type) {
	case *source.Remote:
		src.SetProvider(provider.GuessProviderByUrl(src.Url()))
		if src.Provider() == provider.Unknown {
			src.SetProvider(provider.GuessProviderByContent(src.Content()))
		}
	case *source.Local:
		src.SetProvider(provider.GuessProviderByContent(src.Content()))
	default:
		src.SetProvider(provider.Unknown)
	}
}

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases by
// calling the appropriate provider-specific Unmarshal method from the supported
// providers.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) Unmarshal() (appcaster.Appcaster, error) {
	var appcast appcaster.Appcaster

	p := a.Source().Provider()

	switch p {
	case provider.Sparkle:
		appcast = &sparkle.Appcast{Appcast: a.Appcast}
		break
	case provider.SourceForge:
		appcast = &sourceforge.Appcast{Appcast: a.Appcast}
		break
	case provider.GitHub:
		appcast = &github.Appcast{Appcast: a.Appcast}
		break
	default:
		name := p.String()
		if name == "-" {
			name = "Unknown"
		}

		return nil, fmt.Errorf("releases for the \"%s\" provider can't be unmarshaled", name)
	}

	appcast, err := appcast.Unmarshal()
	if err != nil {
		return nil, err
	}

	a.Source().SetAppcast(appcast)
	a.SetReleases(appcast.Releases())

	return appcast, nil
}

// Uncomment uncomments the commented out lines by calling the appropriate
// provider-specific Uncomment method from the supported providers.
func (a *Appcast) Uncomment() error {
	if a.Source() == nil {
		return fmt.Errorf("no source")
	}

	p := a.Source().Provider()

	switch p {
	case provider.Sparkle:
		appcast := sparkle.Appcast{Appcast: a.Appcast}
		appcast.Uncomment()
		a.Source().SetContent(appcast.Appcast.Source().Content())
		return nil
	}

	name := p.String()
	if name == "-" {
		name = "Unknown"
	}

	return fmt.Errorf("uncommenting is not available for the \"%s\" provider", name)
}
