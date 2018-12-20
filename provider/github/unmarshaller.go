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
func unmarshal(a *Appcast) (appcaster.Appcaster, []error) {
	var feed unmarshalFeed
	var errors []error

	if a.Source() == nil || len(a.Source().Content()) == 0 {
		return nil, append(errors, fmt.Errorf("no source"))
	}

	if a.Source().Appcast() == nil {
		a.Source().SetAppcast(a)
	}

	err := xml.Unmarshal(a.Source().Content(), &feed)
	if err != nil {
		return nil, append(errors, err)
	}

	r, errors := createReleases(feed)

	a.SetReleases(r)

	return a, errors
}

// createReleases creates a release.Releaseser slice from the unmarshalled feed.
func createReleases(feed unmarshalFeed) (release.Releaseser, []error) {
	var items []release.Releaser
	var errors []error

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
			errors = append(errors, fmt.Errorf("release #%d (%s)", i+1, err.Error()))
			continue
		}

		r.SetTitle(entry.Title)
		r.SetDescription(entry.Content)

		// publishedDateTime
		p := release.NewPublishedDateTime()

		err = p.Parse(entry.Updated)
		if err != nil {
			errors = append(errors, fmt.Errorf("release #%d (%s)", i+1, err.Error()))
		}

		r.SetPublishedDateTime(p)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// add release
		if r != nil {
			items = append(items, r)
		}
	}

	return release.NewReleases(items), errors
}
