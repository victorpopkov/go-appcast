// Package appcaster provides the base for creating an appcast type.
//
// On its own this package provides only core structs, interfaces and
// functions/methods that suit most appcast providers. It should be extended by
// the provider-specific types and shouldn't be used as is.
package appcaster

import (
	"errors"
	"regexp"

	"github.com/victorpopkov/go-appcast/release"
)

// Appcaster is the interface that wraps the Appcast methods.
//
// This interface should be embedded by provider-specific Appcaster interfaces.
type Appcaster interface {
	GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum
	LoadSource() error
	Unmarshal() (Appcaster, error)
	UnmarshalReleases() (Appcaster, error)
	Uncomment() error
	SortReleasesByVersions(s release.Sort)
	FilterReleasesByTitle(regexpStr string, inversed ...interface{})
	FilterReleasesByURL(regexpStr string, inversed ...interface{})
	FilterReleasesByPrerelease(inversed ...interface{})
	Source() Sourcer
	SetSource(src Sourcer)
	Output() Outputer
	SetOutput(src Outputer)
	Releases() release.Releaseser
	SetReleases(releases release.Releaseser)
	FirstRelease() release.Releaser
}

// Appcast represents the appcast itself and should be inherited by
// provider-specific appcasts.
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

// New returns a new Appcast instance pointer. The Source can be passed as
// a parameter.
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

// LoadSource calls the Appcast.source.Load method.
func (a *Appcast) LoadSource() error {
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

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases by
// calling the appropriate provider-specific Unmarshal method from the supported
// providers.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
//
// Deprecated: Use Appcast.Unmarshal instead.
func (a *Appcast) UnmarshalReleases() (Appcaster, error) {
	return a.Unmarshal()
}

// Uncomment uncomments the commented out lines by calling the appropriate
// provider-specific Uncomment method from the supported providers.
func (a *Appcast) Uncomment() error {
	panic("implement me")
}

// SortReleasesByVersions sorts the Appcast.releases slice by versions. Can be
// useful if the versions order is inconsistent.
//
// Deprecated: Use Appcast.Releases.SortByVersions methods chain instead.
func (a *Appcast) SortReleasesByVersions(s release.Sort) {
	a.releases.SortByVersions(s)
}

// FilterReleasesByTitle filters all Appcast.releases by matching the release
// title with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
//
// Deprecated: Use Appcast.Releases.FilterByTitle methods chain instead.
func (a *Appcast) FilterReleasesByTitle(regexpStr string, inversed ...interface{}) {
	a.releases.FilterByTitle(regexpStr, inversed...)
}

// FilterReleasesByMediaType filters all Appcast.releases by matching the
// downloads media type with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
//
// Deprecated: Use Appcast.Releases.FilterByMediaType methods chain instead.
func (a *Appcast) FilterReleasesByMediaType(regexpStr string, inversed ...interface{}) {
	a.releases.FilterByMediaType(regexpStr, inversed...)
}

// FilterReleasesByURL filters all Appcast.releases by matching the release
// download URL with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
//
// Deprecated: Use Appcast.Releases.FilterByUrl methods chain instead.
func (a *Appcast) FilterReleasesByURL(regexpStr string, inversed ...interface{}) {
	a.releases.FilterByUrl(regexpStr, inversed...)
}

// FilterReleasesByPrerelease filters all Appcast.releases by matching only the
// pre-releases.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
//
// Deprecated: Use Appcast.Releases.FilterByPrerelease methods chain instead.
func (a *Appcast) FilterReleasesByPrerelease(inversed ...interface{}) {
	a.releases.FilterByPrerelease(inversed...)
}

// ResetFilters resets the Appcast.releases to their original state before
// applying any filters.
//
// Deprecated: Use Appcast.Releases.ResetFilters methods chain instead.
func (a *Appcast) ResetFilters() {
	a.releases.ResetFilters()
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
