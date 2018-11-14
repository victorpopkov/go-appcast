package sparkle

import (
	"encoding/xml"
	"fmt"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/release"
)

// unmarshalFeed represents an RSS itself for the unmarshalling purposes.
type unmarshalFeed struct {
	Channel unmarshalFeedChannel `xml:"channel"`
}

// unmarshalFeedChannel represents an RSS channel for the unmarshalling
// purposes.
type unmarshalFeedChannel struct {
	Title       string              `xml:"title"`
	Link        string              `xml:"link"`
	Description string              `xml:"description"`
	Language    string              `xml:"language"`
	Items       []unmarshalFeedItem `xml:"item"`
}

// unmarshalFeedItem represents a single RSS item for the unmarshalling
// purposes.
type unmarshalFeedItem struct {
	Title                string                 `xml:"title"`
	Description          string                 `xml:"description"`
	PubDate              string                 `xml:"pubDate"`
	ReleaseNotesLink     string                 `xml:"releaseNotesLink"`
	MinimumSystemVersion string                 `xml:"minimumSystemVersion"`
	Enclosure            unmarshalFeedEnclosure `xml:"enclosure"`
	Version              string                 `xml:"version"`
	ShortVersionString   string                 `xml:"shortVersionString"`
}

// unmarshalFeedEnclosure represents a single RSS item enclosure for the
// unmarshalling purposes.
type unmarshalFeedEnclosure struct {
	DsaSignature       string `xml:"dsaSignature,attr"`
	MD5Sum             string `xml:"md5Sum,attr"`
	Version            string `xml:"version,attr"`
	ShortVersionString string `xml:"shortVersionString,attr"`
	URL                string `xml:"url,attr"`
	Length             int    `xml:"length,attr"`
	Type               string `xml:"type,attr"`
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

	a.channel = &Channel{
		Title:       feed.Channel.Title,
		Link:        feed.Channel.Link,
		Description: feed.Channel.Description,
		Language:    feed.Channel.Language,
	}

	return a, errors
}

// createReleases creates a release.Releaseser slice from the unmarshalled feed.
func createReleases(feed unmarshalFeed) (release.Releaseser, []error) {
	var version, build string
	var items []release.Releaser
	var errors []error

	for i, item := range feed.Channel.Items {
		if item.Enclosure.ShortVersionString == "" && item.ShortVersionString != "" {
			version = item.ShortVersionString
		} else {
			version = item.Enclosure.ShortVersionString
		}

		if item.Enclosure.Version == "" && item.Version != "" {
			build = item.Version
		} else {
			build = item.Enclosure.Version
		}

		if version == "" && build == "" {
			errors = append(errors, fmt.Errorf("release #%d (no version)", i+1))
			continue
		} else if version == "" && build != "" {
			version = build
		}

		// new release
		r, err := release.New(version, build)
		if err != nil {
			errors = append(errors, fmt.Errorf("release #%d (%s)", i+1, err.Error()))
			continue
		}

		r.SetTitle(item.Title)
		r.SetDescription(item.Description)
		r.SetReleaseNotesLink(item.ReleaseNotesLink)
		r.SetMinimumSystemVersion(item.MinimumSystemVersion)

		// publishedDateTime
		p := release.NewPublishedDateTime()
		p.Parse(item.PubDate)
		r.SetPublishedDateTime(p)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// downloads
		e := item.Enclosure
		d := release.NewDownload(e.URL, e.Type, e.Length, e.DsaSignature, e.MD5Sum)

		r.AddDownload(*d)

		// add release
		if r != nil {
			items = append(items, r)
		}
	}

	return release.NewReleases(items), errors
}
