// Package appcaster provides the base for creating an appcast type.
//
// This package shouldn't be used on its own. It only provides core structs,
// interfaces and functions/methods that need to be extended by your own appcast
// type(s) and implementations.
package appcaster

import (
	"errors"
	"regexp"

	"github.com/victorpopkov/go-appcast/release"
)

// Appcaster is the interface that wraps the Appcast methods.
//
// This interface should be embedded by your own Appcaster interface.
type Appcaster interface {
	GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum
	LoadSource() error
	GuessSourceProvider()
	Unmarshal() (Appcaster, error)
	Uncomment() error
	Source() Sourcer
	SetSource(src Sourcer)
	Output() Outputer
	SetOutput(src Outputer)
	Releases() release.Releaseser
	SetReleases(releases release.Releaseser)
	FirstRelease() release.Releaser
}

// Appcast represents the appcast itself.
//
// This struct should be embedded by your own Appcast struct.
type Appcast struct {
	// source specifies an appcast source which holds the information about the
	// retrieved appcast. Can be any use-case specific Sourcer interface
	// implementation.
	//
	// By default, two sourcers are supported: LocalSource (for loading appcast
	// from the local file) and RemoteSource (for loading appcast from the
	// remote location by URL).
	source Sourcer

	// output specifies an appcast output which holds the information about the
	// marshaled releases. Can be any use-case specific Outputer interface
	// implementation.
	//
	// By default, only one outputer is supported: LocalOutput (for saving
	// appcast to the local file).
	output Outputer

	// releases specifies an appcast releases.
	releases release.Releaseser
}

// New returns a new Appcast instance pointer. The source can be passed as a
// parameter.
func New(src ...interface{}) *Appcast {
	a := new(Appcast)

	if len(src) > 0 {
		src := src[0].(Sourcer)
		a.SetSource(src)
	}

	return a
}

// ExtractSemanticVersions extracts semantic versions from the provided data
// string.
func ExtractSemanticVersions(data string) ([]string, error) {
	var versions []string

	re := regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)(?:(\-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-\-\.]+)?`)
	if re.MatchString(data) {
		versionMatches := re.FindAllStringSubmatch(data, -1)
		for _, match := range versionMatches {
			versions = append(versions, match[0])
		}

		return versions, nil
	}

	return nil, errors.New("no semantic versions found")
}

// GenerateSourceChecksum creates a new Checksum instance pointer based on the
// provided algorithm and sets it as an Appcast.source.checksum.
func (a *Appcast) GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum {
	return a.source.GenerateChecksum(algorithm)
}

// LoadSource sets the Appcast.source.content field value depending on the
// source type. It should call the appropriate Appcast.Source.Load methods
// chain.
func (a *Appcast) LoadSource() error {
	err := a.Source().Load()
	if err != nil {
		return err
	}

	return nil
}

// GuessSourceProvider attempts to guess the supported provider based on the
// Appcast.source.content.
func (a *Appcast) GuessSourceProvider() {
	panic("implement me")
}

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases by
// calling the appropriate provider-specific Unmarshal method from the supported
// providers.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) Unmarshal() (Appcaster, error) {
	panic("implement me")
}

// Uncomment uncomments the commented out lines by calling the appropriate
// provider-specific Uncomment method from the supported providers.
func (a *Appcast) Uncomment() error {
	panic("implement me")
}

// Source is an Appcast.source getter.
func (a *Appcast) Source() Sourcer {
	return a.source
}

// SetSource is an Appcast.source setter.
func (a *Appcast) SetSource(src Sourcer) {
	a.source = src
}

// Output is an Appcast.output getter.
func (a *Appcast) Output() Outputer {
	return a.output
}

// SetOutput is an Appcast.output setter.
func (a *Appcast) SetOutput(output Outputer) {
	a.output = output
}

// Releases is an Appcast.releases getter.
func (a *Appcast) Releases() release.Releaseser {
	return a.releases
}

// SetReleases is an Appcast.releases setter.
func (a *Appcast) SetReleases(releases release.Releaseser) {
	a.releases = releases
}

// FirstRelease is a convenience method to get the first filtered release from
// the Appcast.releases.
func (a *Appcast) FirstRelease() release.Releaser {
	return a.releases.First()
}
