// Package github adds support for the GitHub releases Atom feed.
package github

import "github.com/victorpopkov/go-appcast/appcaster"

// Appcaster is the interface that wraps the Appcast methods.
type Appcaster interface {
	appcaster.Appcaster
}

// Appcast represents the appcast itself.
type Appcast struct {
	appcaster.Appcast
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

// Unmarshal unmarshals the Appcast.source.content into the Appcast.releases.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *Appcast) Unmarshal() (appcaster.Appcaster, error) {
	return unmarshal(a)
}
