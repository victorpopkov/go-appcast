package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"time"
)

// A SparkleRSSFeedAppcast represents appcast for "Sparkle RSS Feed" that is
// generated by Sparkle Framework.
type SparkleRSSFeedAppcast struct {
	BaseAppcast
}

// A SparkleRSSFeedXML represents an RSS itself.
type SparkleRSSFeedXML struct {
	Items []SparkleRSSFeedXMLItem `xml:"channel>item"`
}

// A SparkleRSSFeedXMLItem represents an RSS item.
type SparkleRSSFeedXMLItem struct {
	Title                string                     `xml:"title"`
	Description          string                     `xml:"description"`
	MinimumSystemVersion string                     `xml:"minimumSystemVersion"`
	PubDate              string                     `xml:"pubDate"`
	Enclosure            SparkleRSSFeedXMLEnclosure `xml:"enclosure"`
	Version              string                     `xml:"version"`
	ShortVersionString   string                     `xml:"shortVersionString"`
}

// A SparkleRSSFeedXMLEnclosure represents an RSS enclosure in item.
type SparkleRSSFeedXMLEnclosure struct {
	Version            string `xml:"version,attr"`
	ShortVersionString string `xml:"shortVersionString,attr"`
	URL                string `xml:"url,attr"`
	Length             int    `xml:"length,attr"`
	Type               string `xml:"type,attr"`
}

// Uncomment uncomments XML tags in SparkleRSSFeedAppcast.Content.
func (a *SparkleRSSFeedAppcast) Uncomment() {
	if a.Content == "" {
		return
	}

	regex := regexp.MustCompile(`(<!--([[:space:]]*)?)|(([[:space:]]*)?-->)`)
	if regex.MatchString(a.Content) {
		a.Content = regex.ReplaceAllString(a.Content, "")
		return
	}
}

// ExtractReleases parses the Sparkle RSS Feed content from
// SparkleRSSFeedAppcast.Content and stores the extracted releases as an array
// in SparkleRSSFeedAppcast.Releases. Returns an error, if extracting was
// unsuccessful.
func (a *SparkleRSSFeedAppcast) ExtractReleases() error {
	var x SparkleRSSFeedXML
	var version, build string

	xml.Unmarshal([]byte(a.Content), &x)

	items := make([]Release, len(x.Items))
	for i, item := range x.Items {
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
			return fmt.Errorf("Version is required, but it's not specified in release #%d", i+1)
		} else if version == "" && build != "" {
			version = build
		}

		r, err := NewRelease(version, build)
		if err != nil {
			return err
		}

		r.Title = item.Title
		r.Description = item.Description

		d := NewDownload(item.Enclosure.URL, item.Enclosure.Type, item.Enclosure.Length)
		r.AddDownload(*d)

		// published date and time
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			r.PublishedDateTime = parsedTime
		}

		items[i] = *r
	}

	a.Releases = items

	return nil
}
