// Package appcast provides functionality for working with appcasts to retrieve
// valuable information about software releases.
//
// Currently supports 3 providers: Sparkle RSS Feed, SourceForge RSS Feed and
// GitHub Atom Feed.
//
// See README.md for more info.
package appcast

// A BaseAppcast represents the appcast itself and should be inherited by
// provider specific appcasts.
type BaseAppcast struct {
	// Request specifies a Request to be sent by a Client to the server. The
	// response should never be modified in the Request itself.
	Request Request

	// Content specifies the copy of the server response from the
	// Request.HTTPRequest. Unlike the response content from the Request, this can
	// be modified if needed.
	Content string

	// Provider specifies one of the supported providers or Provider.Unknown if
	// the appcast is not recognized by this library.
	Provider Provider

	// Checksum specifies the hash checksum for the original content from
	// Request.HTTPRequest. It also includes the used algorithm, source and the
	// checksum itself.
	Checksum Checksum

	// Releases specify an array of all application releases.
	Releases []Release
}
