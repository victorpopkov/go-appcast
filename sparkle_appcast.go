package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/release"
)

// SparkleAppcaster is the interface that wraps the SparkleAppcast
// methods.
type SparkleAppcaster interface {
	appcaster.Appcaster
	Channel() *SparkleAppcastChannel
	SetChannel(channel *SparkleAppcastChannel)
}

// SparkleAppcast represents the "Sparkle RSS Feed" appcast which is generated
// by the Sparkle Framework.
type SparkleAppcast struct {
	appcaster.Appcast
	channel *SparkleAppcastChannel
}

// SparkleAppcastChannel represents the "Sparkle RSS Feed" appcast
// channel data.
type SparkleAppcastChannel struct {
	Title       string
	Link        string
	Description string
	Language    string
}

// unmarshalSparkle represents an RSS itself for unmarshal purposes.
type unmarshalSparkle struct {
	Channel unmarshalSparkleChannel `xml:"channel"`
}

// unmarshalSparkleChannel represents the "Sparkle RSS Feed" channel for
// unmarshal purposes.
type unmarshalSparkleChannel struct {
	Title       string                 `xml:"title"`
	Link        string                 `xml:"link"`
	Description string                 `xml:"description"`
	Language    string                 `xml:"language"`
	Items       []unmarshalSparkleItem `xml:"item"`
}

// unmarshalSparkleItem represents an RSS item.
type unmarshalSparkleItem struct {
	Title                string                    `xml:"title"`
	Description          string                    `xml:"description"`
	PubDate              string                    `xml:"pubDate"`
	ReleaseNotesLink     string                    `xml:"releaseNotesLink"`
	MinimumSystemVersion string                    `xml:"minimumSystemVersion"`
	Enclosure            unmarshalSparkleEnclosure `xml:"enclosure"`
	Version              string                    `xml:"version"`
	ShortVersionString   string                    `xml:"shortVersionString"`
}

// unmarshalSparkleEnclosure represents the "Sparkle RSS Feed" item enclosure
// for unmarshal purposes.
type unmarshalSparkleEnclosure struct {
	DsaSignature       string `xml:"dsaSignature,attr"`
	MD5Sum             string `xml:"md5Sum,attr"`
	Version            string `xml:"version,attr"`
	ShortVersionString string `xml:"shortVersionString,attr"`
	URL                string `xml:"url,attr"`
	Length             int    `xml:"length,attr"`
	Type               string `xml:"type,attr"`
}

// Unmarshal unmarshals the SparkleAppcast.source.content into the
// SparkleAppcast.releases and SparkleAppcast.channel.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *SparkleAppcast) Unmarshal() (appcaster.Appcaster, error) {
	var feed unmarshalSparkle

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

	r, err := a.createReleases(feed)
	if err != nil {
		return nil, err
	}

	a.SetReleases(r)

	a.channel = &SparkleAppcastChannel{
		Title:       feed.Channel.Title,
		Link:        feed.Channel.Link,
		Description: feed.Channel.Description,
		Language:    feed.Channel.Language,
	}

	return a, nil
}

// Unmarshal unmarshals the SparkleAppcast.source.content into the
// SparkleAppcast.releases and SparkleAppcast.channel.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
//
// Deprecated: Use SparkleAppcast.Unmarshal instead.
func (a *SparkleAppcast) UnmarshalReleases() (appcaster.Appcaster, error) {
	return a.Unmarshal()
}

// createReleases creates a release.Releaseser slice from the unmarshalled feed.
func (a *SparkleAppcast) createReleases(feed unmarshalSparkle) (release.Releaseser, error) {
	var version, build string

	items := make([]release.Releaser, len(feed.Channel.Items))
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
			return nil, fmt.Errorf("no version in the #%d release", i+1)
		} else if version == "" && build != "" {
			version = build
		}

		// new release
		r, err := release.New(version, build)
		if err != nil {
			return nil, err
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

		items[i] = r
	}

	return release.NewReleases(items), nil
}

// Uncomment uncomments XML tags in SparkleAppcast.source.content.
func (a SparkleAppcast) Uncomment() error {
	if a.Source() == nil || len(a.Source().Content()) == 0 {
		return fmt.Errorf("no source")
	}

	regex := regexp.MustCompile(`(<!--([[:space:]]*)?)|(([[:space:]]*)?-->)`)
	if regex.Match(a.Source().Content()) {
		a.Source().SetContent(regex.ReplaceAll(a.Source().Content(), []byte("")))
	}

	return nil
}

// Channel is a SparkleAppcast.channel getter.
func (a *SparkleAppcast) Channel() *SparkleAppcastChannel {
	return a.channel
}

// SetChannel is a SparkleAppcast.channel setter.
func (a *SparkleAppcast) SetChannel(channel *SparkleAppcastChannel) {
	a.channel = channel
}
