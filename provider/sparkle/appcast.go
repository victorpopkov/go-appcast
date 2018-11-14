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

// New returns a new Appcast instance pointer. The source can be passed as a
// parameter.
func New(src ...interface{}) *Appcast {
	a := new(Appcast)

	if len(src) > 0 {
		src := src[0].(appcaster.Sourcer)
		a.SetSource(src)
	}

	return a
}

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases and
// Appcast.channel.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an errors slice.
func (a *Appcast) Unmarshal() (appcaster.Appcaster, []error) {
	return unmarshal(a)
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
