package appcast

import "regexp"

// Provider holds different supported providers.
type Provider int

const (
	// Unknown represents an unknown appcast provider.
	Unknown Provider = iota

	// SparkleRSSFeed represents an RSS feed that is generated by Sparkle
	// framework.
	SparkleRSSFeed

	// SourceForgeRSSFeed represents an RSS feed of the releases generated by
	// SourceForge.
	SourceForgeRSSFeed

	// GitHubAtomFeed represents an Atom feed of the releases generated by
	// GitHub.
	GitHubAtomFeed
)

var providerNames = [...]string{
	"-",
	"Sparkle RSS Feed",
	"SourceForge RSS Feed",
	"GitHub Atom Feed",
}

// GuessFromContent attempts to guess the supported provider from the passed
// content. By default returns Provider.Unknown.
func GuessProviderFromContent(content string) Provider {
	regexSparkleRSSFeed := regexp.MustCompile(`(?s)(<rss.*xmlns:sparkle)|(?s)(<rss.*<enclosure)`)
	regexSourceForgeRSSFeed := regexp.MustCompile(`(?s)(<rss.*xmlns:sf)|(?s)(<channel.*xmlns:sf)`)
	regexGitHubAtomFeed := regexp.MustCompile(`(?s)<feed.*<id>tag:github.com`)

	if regexSparkleRSSFeed.MatchString(content) {
		return SparkleRSSFeed
	}

	if regexSourceForgeRSSFeed.MatchString(content) {
		return SourceForgeRSSFeed
	}

	if regexGitHubAtomFeed.MatchString(content) {
		return GitHubAtomFeed
	}

	return Unknown
}

// GuessFromUrl attempts to guess the supported provider from the passed url.
// Only appcasts that are web-service specific can be guessed. By default
// returns Provider.Unknown.
func GuessProviderFromUrl(url string) Provider {
	regexSourceForgeRSSFeed := regexp.MustCompile(`.*sourceforge.net\/projects\/.*\/rss`)
	regexGitHubAtomFeed := regexp.MustCompile(`.*github\.com\/(?P<user>.*?)\/(?P<repo>.*?)\/releases\.atom`)

	if regexSourceForgeRSSFeed.MatchString(url) {
		return SourceForgeRSSFeed
	}

	if regexGitHubAtomFeed.MatchString(url) {
		return GitHubAtomFeed
	}

	return Unknown
}

// String returns the string representation of the Provider.
func (p Provider) String() string {
	return providerNames[p]
}
