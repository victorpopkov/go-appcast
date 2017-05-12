package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessProviderFromContent(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"github_default.xml": GitHubAtomFeed,

		// SourceForge RSS Feed
		"sourceforge_default.xml": SourceForgeRSSFeed,
		"sourceforge_empty.xml":   SourceForgeRSSFeed,
		"sourceforge_single.xml":  SourceForgeRSSFeed,

		// Sparkle RSS Feed
		"sparkle_attributes_as_elements.xml": SparkleRSSFeed,
		"sparkle_default_asc.xml":            SparkleRSSFeed,
		"sparkle_default.xml":                SparkleRSSFeed,
		"sparkle_incorrect_namespace.xml":    SparkleRSSFeed,
		"sparkle_multiple_enclosure.xml":     SparkleRSSFeed,
		"sparkle_no_releases.xml":            SparkleRSSFeed,
		"sparkle_single.xml":                 SparkleRSSFeed,
		"sparkle_without_comments.xml":       SparkleRSSFeed,
		"sparkle_without_namespaces.xml":     SparkleRSSFeed,

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
	assert.Regexp(t, "Sparkle RSS Feed", SparkleRSSFeed.String())
	assert.Regexp(t, "SourceForge RSS Feed", SourceForgeRSSFeed.String())
	assert.Regexp(t, "GitHub Atom Feed", GitHubAtomFeed.String())
}
