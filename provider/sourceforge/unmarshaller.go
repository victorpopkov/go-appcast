package sourceforge

import (
	"encoding/xml"
	"fmt"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/release"
)

// unmarshalFeed represents an RSS itself for the unmarshalling purposes.
type unmarshalFeed struct {
	Items []unmarshalFeedItem `xml:"channel>item"`
}

// unmarshalFeedItem represents an RSS item for the unmarshalling purposes.
type unmarshalFeedItem struct {
	Title       unmarshalFeedItemTitle       `xml:"title"`
	Description unmarshalFeedItemDescription `xml:"description"`
	Content     unmarshalFeedItemContent     `xml:"content"`
	PubDate     string                       `xml:"pubDate"`
}

// unmarshalFeedItemTitle represents an RSS item title for the unmarshalling
// purposes.
type unmarshalFeedItemTitle struct {
	Chardata string `xml:",chardata"`
}

// unmarshalFeedItemDescription represents an RSS item description for the
// unmarshalling purposes.
type unmarshalFeedItemDescription struct {
	Chardata string `xml:",chardata"`
}

// unmarshalFeedItemContent represents an RSS item content for the unmarshalling
// purposes.
type unmarshalFeedItemContent struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Filesize int    `xml:"filesize,attr"`
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
	items := make([]release.Releaser, len(feed.Items))
	for i, item := range feed.Items {
		// extract version
		versions, err := appcaster.ExtractSemanticVersions(item.Title.Chardata)
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

	return release.NewReleases(items), nil
}
