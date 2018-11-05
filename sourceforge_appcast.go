package appcast

import (
	"encoding/xml"
	"fmt"

	"github.com/victorpopkov/go-appcast/release"
)

// SourceForgeAppcaster is the interface that wraps the SourceForgeAppcast
// methods.
type SourceForgeAppcaster interface {
	Appcaster
}

// SourceForgeAppcast represents appcast for "SourceForge RSS Feed" that is
// created by the SourceForge applications and software distributor.
type SourceForgeAppcast struct {
	Appcast
}

// unmarshalSourceForge represents an RSS itself.
type unmarshalSourceForge struct {
	Items []unmarshalSourceForgeItem `xml:"channel>item"`
}

// unmarshalSourceForgeItem represents an RSS item.
type unmarshalSourceForgeItem struct {
	Title       unmarshalSourceForgeTitle       `xml:"title"`
	Description unmarshalSourceForgeDescription `xml:"description"`
	Content     unmarshalSourceForgeContent     `xml:"content"`
	PubDate     string                          `xml:"pubDate"`
}

// unmarshalSourceForgeTitle represents an RSS item title.
type unmarshalSourceForgeTitle struct {
	Chardata string `xml:",chardata"`
}

// unmarshalSourceForgeDescription represents an RSS item description.
type unmarshalSourceForgeDescription struct {
	Chardata string `xml:",chardata"`
}

// unmarshalSourceForgeContent represents an RSS item content.
type unmarshalSourceForgeContent struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Filesize int    `xml:"filesize,attr"`
}

// Unmarshal unmarshals the SourceForgeAppcast.source.content into the
// SourceForgeAppcast.releases.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *SourceForgeAppcast) Unmarshal() (Appcaster, error) {
	var feed unmarshalSourceForge

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

func (a *SourceForgeAppcast) createReleases(feed unmarshalSourceForge) ([]release.Releaser, error) {
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
