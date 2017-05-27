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
	"io/ioutil"
	"regexp"
	"sort"
)

// A BaseAppcast represents the appcast itself and should be inherited by
// provider specific appcasts.
type BaseAppcast struct {
	// Request specifies a Request to be sent by a Client to the server. The
	// response should never be modified in the Request itself.
	Request Request

	// Content specifies the copy of the server response from the
	// Request.HTTPRequest. Unlike the response content from the Request, this can
	// be modified if needed.
	Content string

	// Provider specifies one of the supported providers or Provider.Unknown if
	// the appcast is not recognized by this library.
	Provider Provider

	// Checksum specifies the hash checksum for the original content from
	// Request.HTTPRequest. It also includes the used algorithm, source and the
	// checksum itself.
	Checksum Checksum

	// Releases specify an array of all application releases.
	Releases []Release
}

// Sort holds different supported sorting behaviors.
type Sort int

const (
	// ASC represents the ascending order.
	ASC Sort = iota

	// DESC represents the descending order.
	DESC
)

// New returns a new BaseAppcast instance pointer.
func New() *BaseAppcast {
	a := &BaseAppcast{}

	return a
}

// LoadFromURL loads the appcast content from remote URL and attempts to guess
// the provider.
func (a *BaseAppcast) LoadFromURL(url string) error {
	req, err := NewRequest(url)
	if err != nil {
		return err
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// content
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	a.Content = string(body)
	a.Checksum.Source = a.Content

	// provider
	a.Provider = GuessProviderFromURL(url)
	if a.Provider == Unknown {
		a.Provider = GuessProviderFromContent(a.Content)
	}

	return nil
}

// GenerateChecksum generates and returns the checksum based on provided
// algorithm from BaseAppcast.Checksum.Source. The checksum is also stored as a
// BaseAppcast.Checksum.Result value.
func (a *BaseAppcast) GenerateChecksum(algorithm ChecksumAlgorithm) string {
	a.Checksum.Algorithm = algorithm
	a.Checksum.Source = a.Content
	a.Checksum.Generate()

	return a.Checksum.Result
}

// GetChecksum is a convenience function to retrieve the checksum value stored
// as BaseAppcast.Checksum.Result.
func (a *BaseAppcast) GetChecksum() string {
	return a.Checksum.Result
}

// Uncomment uncomments the commented out lines by calling the appropriate
// provider specific Uncomment function from the supported providers. A
// successful call returns a "nil" error.
func (a *BaseAppcast) Uncomment() error {
	switch a.Provider {
	case SparkleRSSFeed:
		s := SparkleRSSFeedAppcast{BaseAppcast: *a}
		s.Uncomment()
		a.Content = s.BaseAppcast.Content
		break
	default:
		p := a.Provider.String()
		if p == "-" {
			p = "Unknown"
		}
		return fmt.Errorf("Uncommenting is not available for \"%s\" provider", p)
	}

	return nil
}

// ExtractReleases parses the BaseAppcast.Content by calling the appropriate
// provider specific ExtractReleases function. A successful call returns a "nil"
// error.
func (a *BaseAppcast) ExtractReleases() error {
	switch a.Provider {
	case SparkleRSSFeed:
		s := SparkleRSSFeedAppcast{BaseAppcast: *a}
		err := s.ExtractReleases()
		if err != nil {
			return err
		}
		a.Releases = s.BaseAppcast.Releases
		break
	default:
		p := a.Provider.String()
		if p == "-" {
			p = "Unknown"
		}
		return fmt.Errorf("Releases can't be extracted from \"%s\" provider", p)
	}

	return nil
}

// SortReleasesByVersions sorts BaseAppcast.Releases array by versions. Can be
// useful if the versions order in the content is inconsistent.
func (a *BaseAppcast) SortReleasesByVersions(s Sort) {
	if s == ASC {
		sort.Sort(ByVersion(a.Releases))
	} else if s == DESC {
		sort.Sort(sort.Reverse(ByVersion(a.Releases)))
	}
}

// ExtractSemanticVersions extracts semantic versions from the provided data
// string.
func ExtractSemanticVersions(data string) ([]string, error) {
	var versions []string

	regexVersion := regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)(?:(\-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-\-\.]+)?`)
	if regexVersion.MatchString(data) {
		versionMatches := regexVersion.FindAllStringSubmatch(data, -1)
		for _, match := range versionMatches {
			versions = append(versions, match[0])
		}

		return versions, nil
	}

	return nil, errors.New("No semantic versions found")
}
