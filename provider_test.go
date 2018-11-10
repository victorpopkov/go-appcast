package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessProviderByContent(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"../provider/github/testdata/unmarshal/default.xml":         GitHub,
		"../provider/github/testdata/unmarshal/empty.xml":           GitHub,
		"../provider/github/testdata/unmarshal/invalid_pubdate.xml": GitHub,
		"../provider/github/testdata/unmarshal/invalid_tag.xml":     GitHub,
		"../provider/github/testdata/unmarshal/invalid_version.xml": GitHub,
		"../provider/github/testdata/unmarshal/prerelease.xml":      GitHub,

		// SourceForge RSS Feed
		"../provider/sourceforge/testdata/unmarshal/default.xml":         SourceForge,
		"../provider/sourceforge/testdata/unmarshal/empty.xml":           SourceForge,
		"../provider/sourceforge/testdata/unmarshal/invalid_pubdate.xml": SourceForge,
		"../provider/sourceforge/testdata/unmarshal/invalid_tag.xml":     SourceForge,
		"../provider/sourceforge/testdata/unmarshal/invalid_version.xml": SourceForge,
		"../provider/sourceforge/testdata/unmarshal/prerelease.xml":      SourceForge,

		// Sparkle RSS Feed
		"../provider/sparkle/testdata/unmarshal/attributes_as_elements.xml": Sparkle,
		"../provider/sparkle/testdata/unmarshal/default_asc.xml":            Sparkle,
		"../provider/sparkle/testdata/unmarshal/default.xml":                Sparkle,
		"../provider/sparkle/testdata/unmarshal/incorrect_namespace.xml":    Sparkle,
		"../provider/sparkle/testdata/unmarshal/multiple_enclosure.xml":     Sparkle,
		"../provider/sparkle/testdata/unmarshal/no_releases.xml":            Sparkle,
		"../provider/sparkle/testdata/unmarshal/single.xml":                 Sparkle,
		"../provider/sparkle/testdata/unmarshal/with_comments.xml":          Sparkle,
		"../provider/sparkle/testdata/unmarshal/without_namespaces.xml":     Sparkle,

		// Unknown
		"unknown.xml": Unknown,
	}

	for filename, provider := range testCases {
		assert.Equal(t, provider, GuessProviderByContent(getTestdata(filename)), fmt.Sprintf("Provider doesn't match: %s", filename))
	}
}

func TestGuessProviderByContentString(t *testing.T) {
	assert.Equal(t, Sparkle, GuessProviderByContentString(string(getTestdata("../provider/sparkle/testdata/unmarshal/default.xml"))))
}

func TestGuessProviderByUrl(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"http://github.com/user/repo/releases.atom":  GitHub,
		"https://github.com/user/repo/releases.atom": GitHub,

		// SourceForge RSS Feed
		"http://sourceforge.net/projects/name/rss":             SourceForge,
		"https://sourceforge.net/projects/name/rss":            SourceForge,
		"https://sourceforge.net/projects/name/rss?path=/name": SourceForge,

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
	assert.Equal(t, "Sparkle RSS Feed", Sparkle.String())
	assert.Equal(t, "SourceForge RSS Feed", SourceForge.String())
	assert.Equal(t, "GitHub Atom Feed", GitHub.String())
}
