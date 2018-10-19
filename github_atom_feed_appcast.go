package appcast

import (
	"encoding/xml"
	"regexp"
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

// GitHubAtomFeedAppcastXML represents an Atom itself.
type GitHubAtomFeedAppcastXML struct {
	Entries []GitHubAtomFeedAppcastXMLEntry `xml:"entry"`
}

// GitHubAtomFeedAppcastXMLEntry represents an Atom entry.
type GitHubAtomFeedAppcastXMLEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

// UnmarshalReleases unmarshals the Appcast.source.content into the
// Appcast.releases for the "GitHub Atom Feed" provider.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *GitHubAtomFeedAppcast) UnmarshalReleases() (Appcaster, error) {
	var x GitHubAtomFeedAppcastXML

	xml.Unmarshal(a.source.Content(), &x)

	items := make([]Releaser, len(x.Entries))
	for i, entry := range x.Entries {
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
		r, err := NewRelease(version, "")
		if err != nil {
			return nil, err
		}

		r.SetTitle(entry.Title)
		r.SetDescription(entry.Content)
		r.ParsePublishedDateTime(entry.Updated)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// add release
		items[i] = r
	}

	a.releases = items

	return a, nil
}
