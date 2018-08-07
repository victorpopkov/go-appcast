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

// A SourceForgeRSSFeedAppcast represents appcast for "SourceForge RSS Feed"
// that is created by SourceForge applications and software distributor.
type SourceForgeRSSFeedAppcast struct {
	Appcast
}

// A SourceForgeRSSFeedXML represents an RSS itself.
type SourceForgeRSSFeedXML struct {
	Items []SourceForgeRSSFeedXMLItem `xml:"channel>item"`
}

// A SourceForgeRSSFeedXMLItem represents an RSS item.
type SourceForgeRSSFeedXMLItem struct {
	Title       SourceForgeRSSFeedXMLTitle       `xml:"title"`
	Description SourceForgeRSSFeedXMLDescription `xml:"description"`
	Content     SourceForgeRSSFeedXMLContent     `xml:"content"`
	PubDate     string                           `xml:"pubDate"`
}

// A SourceForgeRSSFeedXMLTitle represents an RSS item title.
type SourceForgeRSSFeedXMLTitle struct {
	Chardata string `xml:",chardata"`
}

// A SourceForgeRSSFeedXMLDescription represents an RSS item description.
type SourceForgeRSSFeedXMLDescription struct {
	Chardata string `xml:",chardata"`
}

// A SourceForgeRSSFeedXMLContent represents an RSS item content.
type SourceForgeRSSFeedXMLContent struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Filesize int    `xml:"filesize,attr"`
}

// ExtractReleases parses the SourceForge RSS Feed content from
// SourceForgeRSSFeedAppcast.Content and stores the extracted releases as an
// array in SourceForgeRSSFeedAppcast.Releases.
func (a *SourceForgeRSSFeedAppcast) ExtractReleases() error {
	var x SourceForgeRSSFeedXML

	xml.Unmarshal(a.source.Content(), &x)

	items := make([]Release, len(x.Items))
	for i, item := range x.Items {
		// extract version
		versions, err := ExtractSemanticVersions(item.Title.Chardata)
		if err != nil {
			return fmt.Errorf("version is required, but it's not specified in release #%d", i+1)
		}

		// new release
		r, _ := NewRelease(versions[0], "")

		r.Title = item.Title.Chardata
		r.Description = item.Description.Chardata
		r.ParsePublishedDateTime(item.PubDate)

		// prerelease
		if r.Version.Prerelease() != "" {
			r.IsPrerelease = true
		}

		// downloads
		d := NewDownload(item.Content.URL, item.Content.Type, item.Content.Filesize)
		r.AddDownload(*d)

		// add release
		items[i] = *r
	}

	a.releases = items

	return nil
}
