// Package sparkle adds support for the Sparkle Framework releases RSS feed.
package sparkle

import (
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/appcaster"
)

// Appcaster is the interface that wraps the Appcast methods.
type Appcaster interface {
	appcaster.Appcaster
	Channel() *Channel
	SetChannel(channel *Channel)
}

// Appcast represents the appcast itself.
type Appcast struct {
	appcaster.Appcast
	channel *Channel
}

// Channel represents the appcast channel.
type Channel struct {
	Title       string
	Link        string
	Description string
	Language    string
}

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases and
// Appcast.channel.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) Unmarshal() (appcaster.Appcaster, error) {
	return unmarshal(a)
}

// UnmarshalReleases unmarshals the Appcast.source.content into the
// Appcast.releases and Appcast.channel.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
//
// Deprecated: Use Appcast.Unmarshal instead.
func (a *Appcast) UnmarshalReleases() (appcaster.Appcaster, error) {
	return a.Unmarshal()
}

// Uncomment uncomments XML tags in Appcast.source.content.
func (a *Appcast) Uncomment() error {
	if a.Source() == nil || len(a.Source().Content()) == 0 {
		return fmt.Errorf("no source")
	}

	regex := regexp.MustCompile(`(<!--([[:space:]]*)?)|(([[:space:]]*)?-->)`)
	if regex.Match(a.Source().Content()) {
		a.Source().SetContent(regex.ReplaceAll(a.Source().Content(), []byte("")))
	}

	return nil
}

// Channel is a Appcast.channel getter.
func (a *Appcast) Channel() *Channel {
	return a.channel
}

// SetChannel is a Appcast.channel setter.
func (a *Appcast) SetChannel(channel *Channel) {
	a.channel = channel
}