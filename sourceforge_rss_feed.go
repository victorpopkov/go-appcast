package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"time"
)

// A SourceForgeRSSFeedAppcast represents appcast for "SourceForge RSS Feed"
// that is created by SourceForge applications and software distributor.
type SourceForgeRSSFeedAppcast struct {
	BaseAppcast
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
// array in SourceForgeRSSFeedAppcast.Releases. Returns an error, if extracting
// was unsuccessful.
func (a *SourceForgeRSSFeedAppcast) ExtractReleases() error {
	var x SourceForgeRSSFeedXML

	xml.Unmarshal([]byte(a.Content), &x)

	items := make([]Release, len(x.Items))
	for i, item := range x.Items {
		// extract version
		versions, err := ExtractSemanticVersions(item.Title.Chardata)
		if err != nil {
			return fmt.Errorf("Version is required, but it's not specified in release #%d", i+1)
		}

		// new release
		r, _ := NewRelease(versions[0], "")

		r.Title = item.Title.Chardata
		r.Description = item.Description.Chardata

		// downloads
		d := NewDownload(item.Content.URL, item.Content.Type, item.Content.Filesize)
		r.AddDownload(*d)

		// published date and time
		pubData := item.PubDate
		regexVersion := regexp.MustCompile(`UT$`)
		if regexVersion.MatchString(pubData) {
			pubData = regexVersion.ReplaceAllString(pubData, "UTC")
		}

		parsedTime, err := time.Parse(time.RFC1123, pubData)
		if err == nil {
			r.PublishedDateTime = parsedTime
		}

		// add release
		items[i] = *r
	}

	a.Releases = items

	return nil
}
