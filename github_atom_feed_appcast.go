package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/release"
)

// GitHubAtomFeedAppcaster is the interface that wraps the
// GitHubAtomFeedAppcaster methods.
type GitHubAtomFeedAppcaster interface {
	Appcaster
}

// GitHubAtomFeedAppcast represents appcast for "GitHub Atom Feed" that is
// created by GitHub.
type GitHubAtomFeedAppcast struct {
	Appcast
}

// unmarshalGitHubAtomFeed represents an Atom itself.
type unmarshalGitHubAtomFeed struct {
	Entries []unmarshalGitHubAtomFeedEntry `xml:"entry"`
}

// unmarshalGitHubAtomFeedEntry represents an Atom entry.
type unmarshalGitHubAtomFeedEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

// Unmarshal unmarshals the GitHubAtomFeedAppcast.source.content into the
// GitHubAtomFeedAppcast.releases.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *GitHubAtomFeedAppcast) Unmarshal() (Appcaster, error) {
	var feed unmarshalGitHubAtomFeed

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
func (a *GitHubAtomFeedAppcast) createReleases(feed unmarshalGitHubAtomFeed) ([]release.Releaser, error) {
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
