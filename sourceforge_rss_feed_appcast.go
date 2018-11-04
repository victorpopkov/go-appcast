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

// Unmarshal unmarshals the SourceForgeRSSFeedAppcast.source.content into the
// SourceForgeRSSFeedAppcast.releases.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *SourceForgeRSSFeedAppcast) Unmarshal() (Appcaster, error) {
	var feed unmarshalSourceForgeRSSFeed

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

func (a *SourceForgeRSSFeedAppcast) createReleases(feed unmarshalSourceForgeRSSFeed) ([]release.Releaser, error) {
	items := make([]release.Releaser, len(feed.Items))
	for i, item := range feed.Items {
		// extract version
		versions, err := ExtractSemanticVersions(item.Title.Chardata)
		if err != nil {
			return nil, fmt.Errorf("no version in the #%d release", i+1)
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

	return items, nil
}
