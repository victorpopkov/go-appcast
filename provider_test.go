package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessProviderFromContent(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"github/default.xml": GitHubAtomFeed,

		// SourceForge RSS Feed
		"sourceforge/default.xml": SourceForgeRSSFeed,
		"sourceforge/empty.xml":   SourceForgeRSSFeed,
		"sourceforge/single.xml":  SourceForgeRSSFeed,

		// Sparkle RSS Feed
		"sparkle/attributes_as_elements.xml": SparkleRSSFeed,
		"sparkle/default_asc.xml":            SparkleRSSFeed,
		"sparkle/default.xml":                SparkleRSSFeed,
		"sparkle/incorrect_namespace.xml":    SparkleRSSFeed,
		"sparkle/multiple_enclosure.xml":     SparkleRSSFeed,
		"sparkle/no_releases.xml":            SparkleRSSFeed,
		"sparkle/single.xml":                 SparkleRSSFeed,
		"sparkle/with_comments.xml":          SparkleRSSFeed,
		"sparkle/without_namespaces.xml":     SparkleRSSFeed,

		// Unknown
		"unknown.xml": Unknown,
	}

	for filename, provider := range testCases {
		content := string(getTestdata(filename))
		assert.Equal(t, provider, GuessProviderFromContent(content), fmt.Sprintf("Provider doesn't match: %s", filename))
	}
}

func TestGuessProviderFromURL(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"http://github.com/user/repo/releases.atom":  GitHubAtomFeed,
		"https://github.com/user/repo/releases.atom": GitHubAtomFeed,

		// SourceForge RSS Feed
		"http://sourceforge.net/projects/name/rss":             SourceForgeRSSFeed,
		"https://sourceforge.net/projects/name/rss":            SourceForgeRSSFeed,
		"https://sourceforge.net/projects/name/rss?path=/name": SourceForgeRSSFeed,

		// Unknown
		"https://example.com/user/repo/releases.atom": Unknown,
		"https://github.com/user/repo/releases":       Unknown,
		"https://github.com/invalid/releases.atom":    Unknown,

		"https://example.com/projects/name/rss": Unknown,
		"https://example.com/projects/name":     Unknown,
		"https://sourceforge.net/invalid/rss":   Unknown,
	}

	for url, provider := range testCases {
		assert.Equal(t, provider, GuessProviderFromURL(url), fmt.Sprintf("Provider doesn't match: %s", url))
	}
}

func TestString(t *testing.T) {
	assert.Equal(t, "-", Unknown.String())
	assert.Equal(t, "Sparkle RSS Feed", SparkleRSSFeed.String())
	assert.Equal(t, "SourceForge RSS Feed", SourceForgeRSSFeed.String())
	assert.Equal(t, "GitHub Atom Feed", GitHubAtomFeed.String())
}
