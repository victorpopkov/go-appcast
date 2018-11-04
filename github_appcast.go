package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/release"
)

// GitHubAppcaster is the interface that wraps the GitHubAppcaster methods.
type GitHubAppcaster interface {
	Appcaster
}

// GitHubAppcast represents appcast for "GitHub Atom Feed" that is created by
// GitHub.
type GitHubAppcast struct {
	Appcast
}

// unmarshalGitHub represents an Atom itself.
type unmarshalGitHub struct {
	Entries []unmarshalGitHubEntry `xml:"entry"`
}

// unmarshalGitHubEntry represents an Atom entry.
type unmarshalGitHubEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

// Unmarshal unmarshals the GitHubAppcast.source.content into the
// GitHubAppcast.releases.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *GitHubAppcast) Unmarshal() (Appcaster, error) {
	var feed unmarshalGitHub

	if a.source == nil || len(a.source.Content()) == 0 {
		return nil, fmt.Errorf("no source")
	}

	if a.source.Appcast() == nil {
		a.source.SetAppcast(a)
	}

	err := xml.Unmarshal(a.source.Content(), &feed)
	if err != nil {
		return nil, err
	}

	items, err := a.createReleases(feed)
	if err != nil {
		return nil, err
	}

	a.releases = items

	return a, nil
}

// createReleases creates a release.Releaser array from the unmarshalled feed.
func (a *GitHubAppcast) createReleases(feed unmarshalGitHub) ([]release.Releaser, error) {
	items := make([]release.Releaser, len(feed.Entries))
	for i, entry := range feed.Entries {
		version := ""

		re := regexp.MustCompile(`\/.*\/(.*$)`)
		if re.MatchString(entry.ID) {
			// extract last part that represents version
			versionMatches := re.FindAllStringSubmatch(entry.ID, 1)
			version = versionMatches[0][1]

			// remove the first "v"
			re := regexp.MustCompile(`^v`)
			version = re.ReplaceAllString(version, "")
		}

		// new release
		r, err := release.New(version, "")
		if err != nil {
			return nil, err
		}

		r.SetTitle(entry.Title)
		r.SetDescription(entry.Content)

		// publishedDateTime
		p := release.NewPublishedDateTime()
		p.Parse(entry.Updated)
		r.SetPublishedDateTime(p)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// add release
		items[i] = r
	}

	return items, nil
}
