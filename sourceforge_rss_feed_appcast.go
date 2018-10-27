package appcast

import (
	"encoding/xml"
	"fmt"

	"github.com/victorpopkov/go-appcast/release"
)

// SourceForgeRSSFeedAppcaster is the interface that wraps the
// SourceForgeRSSFeedAppcast methods.
type SourceForgeRSSFeedAppcaster interface {
	Appcaster
}

// SourceForgeRSSFeedAppcast represents appcast for "SourceForge RSS Feed"
// that is created by SourceForge applications and software distributor.
type SourceForgeRSSFeedAppcast struct {
	Appcast
}

// unmarshalSourceForgeRSSFeed represents an RSS itself.
type unmarshalSourceForgeRSSFeed struct {
	Items []unmarshalSourceForgeRSSFeedItem `xml:"channel>item"`
}

// unmarshalSourceForgeRSSFeedItem represents an RSS item.
type unmarshalSourceForgeRSSFeedItem struct {
	Title       unmarshalSourceForgeRSSFeedTitle       `xml:"title"`
	Description unmarshalSourceForgeRSSFeedDescription `xml:"description"`
	Content     unmarshalSourceForgeRSSFeedContent     `xml:"content"`
	PubDate     string                                 `xml:"pubDate"`
}

// unmarshalSourceForgeRSSFeedTitle represents an RSS item title.
type unmarshalSourceForgeRSSFeedTitle struct {
	Chardata string `xml:",chardata"`
}

// unmarshalSourceForgeRSSFeedDescription represents an RSS item description.
type unmarshalSourceForgeRSSFeedDescription struct {
	Chardata string `xml:",chardata"`
}

// unmarshalSourceForgeRSSFeedContent represents an RSS item content.
type unmarshalSourceForgeRSSFeedContent struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Filesize int    `xml:"filesize,attr"`
}

// UnmarshalReleases unmarshals the Appcast.source.content into the
// Appcast.releases for the "SourceForge RSS Feed" provider.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *SourceForgeRSSFeedAppcast) UnmarshalReleases() (Appcaster, error) {
	var x unmarshalSourceForgeRSSFeed

	xml.Unmarshal(a.source.Content(), &x)

	items := make([]release.Releaser, len(x.Items))
	for i, item := range x.Items {
		// extract version
		versions, err := ExtractSemanticVersions(item.Title.Chardata)
		if err != nil {
			return nil, fmt.Errorf("version is required, but it's not specified in release #%d", i+1)
		}

		// new release
		r, _ := release.New(versions[0], "")

		r.SetTitle(item.Title.Chardata)
		r.SetDescription(item.Description.Chardata)

		// publishedDateTime
		p := release.NewPublishedDateTime()
		p.Parse(item.PubDate)
		r.SetPublishedDateTime(p)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// downloads
		d := release.NewDownload(item.Content.URL, item.Content.Type, item.Content.Filesize)
		r.AddDownload(*d)

		// add release
		items[i] = r
	}

	a.releases = items

	return a, nil
}
