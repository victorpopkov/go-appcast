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
)

// Appcaster is the interface that wraps the Appcast methods.
//
// This interface should be embedded by provider specific Appcaster interfaces.
type Appcaster interface {
	LoadFromRemoteSource(i interface{}) error
	LoadFromLocalSource(path string) error
	GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum
	Uncomment() error
	ExtractReleases() error
	SortReleasesByVersions(s Sort)
	FilterReleasesByTitle(regexpStr string, inversed ...interface{})
	FilterReleasesByURL(regexpStr string, inversed ...interface{})
	FilterReleasesByPrerelease(inversed ...interface{})
	Source() Sourcer
	SetSource(source Sourcer)
	OriginalReleases() []Release
	SetOriginalReleases(originalReleases []Release)
}

// An Appcast represents the appcast itself and should be inherited by provider
// specific appcasts.
type Appcast struct {
	// source specifies an appcast source which holds the information about the
	// retrieved appcast. Can be any use-case specific Sourcer interface
	// implementation.
	//
	// By default, two sourcers are supported: LocalSource (for loading appcast
	// from the local file) and RemoteSource (for loading appcast from the
	// remote location by URL).
	source Sourcer

	// Releases specify a slice of all application releases. All filtered
	// releases are stored here.
	Releases []Release

	// originalReleases specify a slice holds a copy of the Appcast.Releases. It
	// is used to restore the Appcast.Releases using the Appcast.ResetFilters
	// function.
	originalReleases []Release
}

// Sort holds different supported sorting behaviors.
type Sort int

const (
	// ASC represents the ascending order.
	ASC Sort = iota

	// DESC represents the descending order.
	DESC
)

// New returns a new Appcast instance pointer.
func New() *Appcast {
	a := &Appcast{}

	return a
}

// LoadFromRemoteSource creates a new RemoteSource instance and loads the data
// from the remote location by using the RemoteSource.Load method.
func (a *Appcast) LoadFromRemoteSource(i interface{}) error {
	s, err := NewRemoteSource(i)
	if err != nil {
		return err
	}

	err = s.Load()
	if err != nil {
		return err
	}

	a.source = s

	return nil
}

// LoadFromURL creates a new RemoteSource instance and loads the data from the
// remote location by using the RemoteSource.Load method.
//
// Deprecated: Use Appcast.LoadFromRemoteSource instead.
func (a *Appcast) LoadFromURL(i interface{}) error {
	return a.LoadFromRemoteSource(i)
}

// LoadFromLocalSource creates a new LocalSource instance and loads the data
// from the local file by using the LocalSource.Load method.
func (a *Appcast) LoadFromLocalSource(path string) error {
	s := NewLocalSource(path)
	err := s.Load()
	if err != nil {
		return err
	}

	a.source = s

	return nil
}

// LoadFromFile creates a new LocalSource instance and loads the data from the
// local file by using the LocalSource.Load method.
//
// Deprecated: Use Appcast.LoadFromLocalSource instead.
func (a *Appcast) LoadFromFile(path string) error {
	return a.LoadFromLocalSource(path)
}

// GenerateSourceChecksum creates a new Checksum instance in the Appcast.source
// based on the provided algorithm and returns its pointer.
func (a *Appcast) GenerateSourceChecksum(algorithm ChecksumAlgorithm) *Checksum {
	a.source.GenerateChecksum(algorithm)
	return a.source.Checksum()
}

// GenerateChecksum creates a new Checksum instance in the Appcast.source based
// on the provided algorithm and returns its pointer.
//
// Deprecated: Use Appcast.GenerateSourceChecksum instead.
func (a *Appcast) GenerateChecksum(algorithm ChecksumAlgorithm) *Checksum {
	return a.GenerateSourceChecksum(algorithm)
}

// Uncomment uncomments the commented out lines by calling the appropriate
// provider specific Uncomment function from the supported providers.
func (a *Appcast) Uncomment() error {
	if a.source == nil {
		return fmt.Errorf("no source")
	}

	provider := a.source.Provider()
	providerString := provider.String()

	switch provider {
	case SparkleRSSFeed:
		s := SparkleRSSFeedAppcast{Appcast: *a}
		s.Uncomment()
		a.source.SetContent(s.Appcast.source.Content())

		return nil
	default:
		if providerString == "-" {
			providerString = "Unknown"
		}
		break
	}

	return fmt.Errorf("uncommenting is not available for the \"%s\" provider", providerString)
}

// ExtractReleases parses the Appcast.source.content by calling the appropriate
// provider specific ExtractReleases function.
func (a *Appcast) ExtractReleases() error {
	provider := a.source.Provider()

	switch provider {
	case SparkleRSSFeed:
		s := SparkleRSSFeedAppcast{Appcast: *a}
		err := s.ExtractReleases()
		if err != nil {
			return err
		}

		a.Releases = s.Appcast.Releases
		a.originalReleases = a.Releases
		break
	case SourceForgeRSSFeed:
		s := SourceForgeRSSFeedAppcast{Appcast: *a}
		err := s.ExtractReleases()
		if err != nil {
			return err
		}

		a.Releases = s.Appcast.Releases
		a.originalReleases = a.Releases
		break
	case GitHubAtomFeed:
		s := GitHubAtomFeedAppcast{Appcast: *a}
		err := s.ExtractReleases()
		if err != nil {
			return err
		}

		a.Releases = s.Appcast.Releases
		a.originalReleases = a.Releases
		break
	default:
		p := provider.String()
		if p == "-" {
			p = "Unknown"
		}

		return fmt.Errorf("releases can't be extracted from the \"%s\" provider", p)
	}

	return nil
}

// SortReleasesByVersions sorts Appcast.Releases slice by versions. Can be
// useful if the versions order in the content is inconsistent.
func (a *Appcast) SortReleasesByVersions(s Sort) {
	if s == ASC {
		sort.Sort(ByVersion(a.Releases))
	} else if s == DESC {
		sort.Sort(sort.Reverse(ByVersion(a.Releases)))
	}
}

// filterReleasesBy filters all Appcast.Releases using the passed function.
// If inverse is set to "true", the unmatched releases will be used instead.
func (a *Appcast) filterReleasesBy(f func(r Release) bool, inverse bool) {
	var result []Release

	for _, release := range a.Releases {
		if inverse == false && f(release) {
			result = append(result, release)
			continue
		}

		if inverse == true && !f(release) {
			result = append(result, release)
			continue
		}
	}

	a.Releases = result
}

// filterReleasesDownloadsBy filters all Downloads for Appcast.Releases
// using the passed function. If inverse is set to "true", the unmatched
// releases will be used instead.
func (a *Appcast) filterReleasesDownloadsBy(f func(d Download) bool, inverse bool) {
	var result []Release

	for _, release := range a.Releases {
		for _, download := range release.Downloads {
			if inverse == false && f(download) {
				result = append(result, release)
				continue
			}

			if inverse == true && !f(download) {
				result = append(result, release)
				continue
			}
		}
	}

	a.Releases = result
}

// FilterReleasesByTitle filters all Appcast.Releases by matching the
// release title with the provided RegExp string. If inversed bool is set to
// "true", the unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByTitle(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesBy(func(r Release) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(r.Title) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByMediaType filters all releases by matching the downloads
// media type with the provided RegExp string. If inversed bool is set to
// "true", the unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByMediaType(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesDownloadsBy(func(d Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.Type) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByURL filters all Appcast.Releases by matching the release
// download URL with the provided RegExp string. If inversed bool is set to
// "true", the unmatched releases will be used instead.
func (a *Appcast) FilterReleasesByURL(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesDownloadsBy(func(d Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.URL) {
			return true
		}
		return false
	}, inverse)
}

// FilterReleasesByPrerelease filters all Appcast.Releases by matching only
// prereleases. If inversed bool is set to "true", the stable releases will be
// matched instead.
func (a *Appcast) FilterReleasesByPrerelease(inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	a.filterReleasesBy(func(r Release) bool {
		if r.IsPrerelease == true {
			return true
		}
		return false
	}, inverse)
}

// ResetFilters resets the Appcast.Releases to their original state before
// applying any filters.
func (a *Appcast) ResetFilters() {
	a.Releases = a.originalReleases
}

// GetReleasesLength is a convenience function to retrieve the total number of
// releases in Appcast.Releases slice.
func (a *Appcast) GetReleasesLength() int {
	return len(a.Releases)
}

// GetFirstRelease is a convenience function to retrieve the first release
// pointer from Appcast.Releases slice.
func (a *Appcast) GetFirstRelease() *Release {
	return &a.Releases[0]
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
func (a *Appcast) SetSource(source Sourcer) {
	a.source = source
}

// OriginalReleases is an Appcast.originalReleases getter.
func (a *Appcast) OriginalReleases() []Release {
	return a.originalReleases
}

// SetOriginalReleases is an Appcast.originalReleases setter.
func (a *Appcast) SetOriginalReleases(originalReleases []Release) {
	a.originalReleases = originalReleases
}

// GetChecksum is an Appcast.source.checksum getter.
//
// Deprecated: Use Appcast.Source.Checksum instead.
func (a *Appcast) GetChecksum() *Checksum {
	return a.Source().Checksum()
}

// Provider is an Appcast.source.provider getter.
//
// Deprecated: Use Appcast.Source.Provider instead.
func (a *Appcast) GetProvider() Provider {
	return a.Source().Provider()
}
