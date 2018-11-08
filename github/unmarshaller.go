package github

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/release"
)

// unmarshalFeed represents an Atom itself for the unmarshalling purposes.
type unmarshalFeed struct {
	Entries []unmarshalFeedEntry `xml:"entry"`
}

// unmarshalFeedEntry represents an Atom entry for the unmarshalling purposes.
type unmarshalFeedEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

// unmarshal unmarshals the Appcast.source.content from the provided Appcast
// pointer into its Appcast.releases and Appcast.channel fields.
func unmarshal(a *Appcast) (appcaster.Appcaster, error) {
	var feed unmarshalFeed

	if a.Source() == nil || len(a.Source().Content()) == 0 {
		return nil, fmt.Errorf("no source")
	}

	if a.Source().Appcast() == nil {
		a.Source().SetAppcast(a)
	}

	err := xml.Unmarshal(a.Source().Content(), &feed)
	if err != nil {
		return nil, err
	}

	r, err := createReleases(feed)
	if err != nil {
		return nil, err
	}

	a.SetReleases(r)

	return a, nil
}

// createReleases creates a release.Releaseser slice from the unmarshalled feed.
func createReleases(feed unmarshalFeed) (release.Releaseser, error) {
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

	return release.NewReleases(items), nil
}
