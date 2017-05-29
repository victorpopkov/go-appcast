package appcast

import (
	"encoding/xml"
	"regexp"
)

// A GitHubAtomFeedAppcast represents appcast for "GitHub Atom Feed" that is
// created by GitHub.
type GitHubAtomFeedAppcast struct {
	BaseAppcast
}

// A GitHubAtomFeedAppcastXML represents an Atom itself.
type GitHubAtomFeedAppcastXML struct {
	Entries []GitHubAtomFeedAppcastXMLEntry `xml:"entry"`
}

// A GitHubAtomFeedAppcastXMLEntry represents an Atom entry.
type GitHubAtomFeedAppcastXMLEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Content string `xml:"description"`
}

// ExtractReleases parses the GitHub Atom Feed content from
// GitHubAtomFeedAppcast.Content and stores the extracted releases as an
// array in GitHubAtomFeedAppcast.Releases. Returns an error, if extracting
// was unsuccessful.
func (a *GitHubAtomFeedAppcast) ExtractReleases() error {
	var x GitHubAtomFeedAppcastXML

	xml.Unmarshal([]byte(a.Content), &x)

	items := make([]Release, len(x.Entries))
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
			return err
		}

		r.Title = entry.Title
		r.Description = entry.Content
		r.ParsePublishedDateTime(entry.Updated)

		// prerelease
		if r.Version.Prerelease() != "" {
			r.IsPrerelease = true
		}

		// add release
		items[i] = *r
	}

	a.Releases = items

	return nil
}
