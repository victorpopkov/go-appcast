// Package appcast provides functionality for working with appcasts to retrieve
// valuable information about software releases.
//
// Currently supports 3 providers: Sparkle RSS Feed, SourceForge RSS Feed and
// GitHub Atom Feed.
//
// See README.md for more info.
package appcast

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/victorpopkov/go-appcast/client"
	"github.com/victorpopkov/go-appcast/release"
)

// DefaultClient is the default Client that is used for making requests in the
// appcast package.
var DefaultClient = client.New()

// Appcaster is the interface that wraps the Appcast methods.
//
// This interface should be embedded by provider-specific Appcaster interfaces.
type Appcaster interface {
	LoadFromRemoteSource(i interface{}) (Appcaster, error)
	LoadFromLocalSource(path string) (Appcaster, error)
	GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum
	LoadSource() error
	Unmarshal() (Appcaster, error)
	UnmarshalReleases() (Appcaster, error)
	Uncomment() error
	SortReleasesByVersions(s Sort)
	FilterReleasesByTitle(regexpStr string, inversed ...interface{})
	FilterReleasesByURL(regexpStr string, inversed ...interface{})
	FilterReleasesByPrerelease(inversed ...interface{})
	Source() Sourcer
	SetSource(src Sourcer)
	Output() Outputer
	SetOutput(src Outputer)
	Releases() []release.Releaser
	SetReleases(releases []release.Releaser)
	FirstRelease() release.Releaser
	OriginalReleases() []release.Releaser
	SetOriginalReleases(originalReleases []release.Releaser)
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

	// releases specify a slice of all application releases. All filtered
	// releases are stored here.
	releases []release.Releaser

	// originalReleases specify a slice holding a copy of the Appcast.releases.
	// It is used to restore the Appcast.releases using the Appcast.ResetFilters
	// method.
	originalReleases []release.Releaser
}

// Sort holds different supported sorting behaviours.
type Sort int

const (
	// ASC represents the ascending order.
	ASC Sort = iota

	// DESC represents the descending order.
	DESC
)

// New returns a new Appcast instance pointer. The Source can be passed as
// a parameter.
func New(src ...interface{}) *Appcast {
	a := &Appcast{}

	if len(src) > 0 {
		src := src[0].(Sourcer)
		a.SetSource(src)
	}

	return a
}

// LoadFromRemoteSource creates a new RemoteSource instance and loads the data
// from the remote location by using the RemoteSource.Load method.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) LoadFromRemoteSource(i interface{}) (Appcaster, error) {
	src, err := NewRemoteSource(i)
	if err != nil {
		return nil, err
	}

	err = src.Load()
	if err != nil {
		return nil, err
	}

	a.source = src
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
func (a *Appcast) LoadFromLocalSource(path string) (Appcaster, error) {
	src := NewLocalSource(path)
	err := src.Load()
	if err != nil {
		return nil, err
	}

	a.source = src
	appcast, err := a.Unmarshal()
	if err != nil {
		return nil, err
	}

	return appcast, nil
}

// GenerateSourceChecksum creates a new Checksum instance in the Appcast.source
// based on the provided algorithm and returns its pointer.
func (a *Appcast) GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum {
	a.source.GenerateChecksum(algorithm)
	return a.source.Checksum()
}

// LoadSource calls the Appcast.source.Load method.
func (a *Appcast) LoadSource() error {
	return a.source.Load()
}

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases by
// calling the appropriate provider-specific Unmarshal method from the supported
// providers.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) Unmarshal() (Appcaster, error) {
	var appcast Appcaster

	p := a.source.Provider()

	switch p {
	case Sparkle:
		appcast = &SparkleAppcast{Appcast: *a}
		break
	case SourceForge:
		appcast = &SourceForgeAppcast{Appcast: *a}
		break
	case GitHub:
		appcast = &GitHubAppcast{Appcast: *a}
		break
	default:
		provider := p.String()
		if provider == "-" {
			provider = "Unknown"
		}

		return nil, fmt.Errorf("releases for the \"%s\" provider can't be unmarshaled", provider)
	}

	appcast, err := appcast.Unmarshal()
	if err != nil {
		return nil, err
	}

	a.source.SetAppcast(appcast)
	a.releases = appcast.Releases()
	a.originalReleases = a.releases

	return appcast, nil
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
	if a.source == nil {
		return fmt.Errorf("no source")
	}

	p := a.source.Provider()
	provider := p.String()

	switch p {
	case Sparkle:
		appcast := SparkleAppcast{Appcast: *a}
		appcast.Uncomment()
		a.source.SetContent(appcast.Appcast.source.Content())

		return nil
	default:
		if provider == "-" {
			provider = "Unknown"
		}
		break
	}

	return fmt.Errorf("uncommenting is not available for the \"%s\" provider", provider)
}

// SortReleasesByVersions sorts Appcast.releases slice by versions. Can be
// useful if the versions order in the content is inconsistent.
func (a *Appcast) SortReleasesByVersions(s Sort) {
	if s == ASC {
		sort.Sort(release.ByVersion(a.releases))
	} else if s == DESC {
		sort.Sort(sort.Reverse(release.ByVersion(a.releases)))
	}
}

// filterReleasesBy filters all Appcast.releases using the passed function. If
// inverse is set to true, the unmatched releases will be used instead.
func (a *Appcast) filterReleasesBy(f func(r release.Releaser) bool, inverse bool) {
	var result []release.Releaser

	for _, r := range a.releases {
		if inverse == false && f(r) {
			result = append(result, r)
			continue
		}

		if inverse == true && !f(r) {
			result = append(result, r)
			continue
		}
	}

	a.releases = result
}

// filterReleasesDownloadsBy filters all Downloads for Appcast.releases using
// the passed function. If inverse is set to true, the unmatched releases will
// be used instead.
func (a *Appcast) filterReleasesDownloadsBy(f func(d release.Download) bool, inverse bool) {
	var result []release.Releaser

	for _, r := range a.releases {
		for _, download := range r.Downloads() {
			if inverse == false && f(download) {
				result = append(result, r)
				continue
			}

			if inverse == true && !f(download) {
				result = append(result, r)
				continue
			}
		}
	}

	a.releases = result
}

// FilterReleasesByTitle filters all Appcast.releases by matching the release
// title with the provided RegExp string. If inversed bool is set to true, the
// unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByTitle(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesBy(func(r release.Releaser) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(r.Title()) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByMediaType filters all releases by matching the downloads
// media type with the provided RegExp string. If inversed bool is set to true,
// the unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByMediaType(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesDownloadsBy(func(d release.Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.Filetype()) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByURL filters all Appcast.releases by matching the release
// download URL with the provided RegExp string. If inversed bool is set to
// true, the unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByURL(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesDownloadsBy(func(d release.Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.Url()) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByPrerelease filters all Appcast.releases by matching only
// pre-releases. If inversed bool is set to true, the stable releases will be
// matched instead.
func (a *Appcast) FilterReleasesByPrerelease(inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesBy(func(r release.Releaser) bool {
		if r.IsPreRelease() == true {
			return true
		}
		return false
	}, inverse)
}

// ResetFilters resets the Appcast.releases to their original state before
// applying any filters.
func (a *Appcast) ResetFilters() {
	a.releases = a.originalReleases
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
func (a *Appcast) Releases() []release.Releaser {
	return a.releases
}

// SetReleases is an Appcast.releases setter.
func (a *Appcast) SetReleases(releases []release.Releaser) {
	a.releases = releases
}

// FirstRelease is a convenience method to get the first release pointer from
// the Appcast.releases slice.
func (a *Appcast) FirstRelease() release.Releaser {
	return a.releases[0]
}

// OriginalReleases is an Appcast.originalReleases getter.
func (a *Appcast) OriginalReleases() []release.Releaser {
	return a.originalReleases
}

// SetOriginalReleases is an Appcast.originalReleases setter.
func (a *Appcast) SetOriginalReleases(originalReleases []release.Releaser) {
	a.originalReleases = originalReleases
}
