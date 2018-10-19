package appcast

import (
	"encoding/xml"
	"fmt"
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

// SourceForgeRSSFeedXML represents an RSS itself.
type SourceForgeRSSFeedXML struct {
	Items []SourceForgeRSSFeedXMLItem `xml:"channel>item"`
}

// SourceForgeRSSFeedXMLItem represents an RSS item.
type SourceForgeRSSFeedXMLItem struct {
	Title       SourceForgeRSSFeedXMLTitle       `xml:"title"`
	Description SourceForgeRSSFeedXMLDescription `xml:"description"`
	Content     SourceForgeRSSFeedXMLContent     `xml:"content"`
	PubDate     string                           `xml:"pubDate"`
}

// SourceForgeRSSFeedXMLTitle represents an RSS item title.
type SourceForgeRSSFeedXMLTitle struct {
	Chardata string `xml:",chardata"`
}

// SourceForgeRSSFeedXMLDescription represents an RSS item description.
type SourceForgeRSSFeedXMLDescription struct {
	Chardata string `xml:",chardata"`
}

// SourceForgeRSSFeedXMLContent represents an RSS item content.
type SourceForgeRSSFeedXMLContent struct {
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
	var x SourceForgeRSSFeedXML

	xml.Unmarshal(a.source.Content(), &x)

	items := make([]Releaser, len(x.Items))
	for i, item := range x.Items {
		// extract version
		versions, err := ExtractSemanticVersions(item.Title.Chardata)
		if err != nil {
			return nil, fmt.Errorf("version is required, but it's not specified in release #%d", i+1)
		}

		// new release
		r, _ := NewRelease(versions[0], "")

		r.SetTitle(item.Title.Chardata)
		r.SetDescription(item.Description.Chardata)
		r.ParsePublishedDateTime(item.PubDate)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// downloads
		d := NewDownload(item.Content.URL, item.Content.Type, item.Content.Filesize)
		r.AddDownload(*d)

		// add release
		items[i] = r
	}

	a.releases = items

	return a, nil
}
