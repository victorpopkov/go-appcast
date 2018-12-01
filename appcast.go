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
type Appcaster interface {
	appcaster.Appcaster
	LoadFromRemoteSource(i interface{}) (appcaster.Appcaster, []error)
	LoadFromLocalSource(path string) (appcaster.Appcaster, []error)
}

// Appcast represents the appcast itself.
type Appcast struct {
	appcaster.Appcast
}

// New returns a new Appcast instance pointer. The source can be passed as a
// parameter.
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
// Appcaster interface and an errors slice.
func (a *Appcast) LoadFromRemoteSource(i interface{}) (appcaster.Appcaster, []error) {
	var errors []error

	src, err := source.NewRemote(i)
	if err != nil {
		return nil, append(errors, err)
	}

	err = src.Load()
	if err != nil {
		return nil, append(errors, err)
	}

	a.SetSource(src)
	a.GuessSourceProvider()

	return a.Unmarshal()
}

// LoadFromLocalSource creates a new LocalSource instance and loads the data
// from the local file by using the LocalSource.Load method.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an errors slice.
func (a *Appcast) LoadFromLocalSource(path string) (appcaster.Appcaster, []error) {
	var errors []error

	src := source.NewLocal(path)
	err := src.Load()
	if err != nil {
		return nil, append(errors, err)
	}

	a.SetSource(src)
	a.GuessSourceProvider()

	return a.Unmarshal()
}

// LoadSource sets the Appcast.source.content field value depending on the
// source type. It should call the appropriate Appcast.Source.Load methods
// chain.
func (a *Appcast) LoadSource() error {
	err := a.Appcast.LoadSource()
	if err != nil {
		return err
	}

	a.GuessSourceProvider()

	return nil
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
// Appcaster interface and an errors slice.
func (a *Appcast) Unmarshal() (appcaster.Appcaster, []error) {
	var appcast appcaster.Appcaster
	var errors []error

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

		errors = append(errors, fmt.Errorf("releases for the \"%s\" provider can't be unmarshaled", name))

		return nil, errors
	}

	appcast, errors = appcast.Unmarshal()

	a.Source().SetAppcast(appcast)
	a.SetReleases(appcast.Releases())

	return appcast, errors
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

		err := appcast.Uncomment()
		if err != nil {
			return err
		}

		a.Source().SetContent(appcast.Appcast.Source().Content())

		return nil
	}

	name := p.String()
	if name == "-" {
		name = "Unknown"
	}

	return fmt.Errorf("uncommenting is not available for the \"%s\" provider", name)
}
