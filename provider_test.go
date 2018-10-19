package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessProviderByContent(t *testing.T) {
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
		assert.Equal(t, provider, GuessProviderByContent(getTestdata(filename)), fmt.Sprintf("Provider doesn't match: %s", filename))
	}
}

func TestGuessProviderByContentString(t *testing.T) {
	assert.Equal(t, SparkleRSSFeed, GuessProviderByContentString(string(getTestdata("sparkle/default.xml"))))
}

func TestGuessProviderByUrl(t *testing.T) {
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
		assert.Equal(t, provider, GuessProviderByUrl(url), fmt.Sprintf("Provider doesn't match: %s", url))
	}
}

func TestProvider_String(t *testing.T) {
	assert.Equal(t, "-", Unknown.String())
	assert.Equal(t, "Sparkle RSS Feed", SparkleRSSFeed.String())
	assert.Equal(t, "SourceForge RSS Feed", SourceForgeRSSFeed.String())
	assert.Equal(t, "GitHub Atom Feed", GitHubAtomFeed.String())
}
