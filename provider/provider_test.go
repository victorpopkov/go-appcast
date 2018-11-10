package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdataPath = "./"

// getWorkingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func getWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// getTestdata returns a file content as a byte slice from the provided testdata
// paths. If the file is not found, prints an error to os.Stdout and exits with
// exit status 1.
func getTestdata(paths ...string) []byte {
	path := getTestdataPath(paths...)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// getTestdataPath returns a full path for the provided testdata paths.
func getTestdataPath(paths ...string) string {
	return filepath.Join(getWorkingDir(), testdataPath, filepath.Join(paths...))
}

func TestGuessProviderByContent(t *testing.T) {
	testCases := map[string]Provider{
		// GitHub Atom Feed
		"github/testdata/unmarshal/default.xml":         GitHub,
		"github/testdata/unmarshal/empty.xml":           GitHub,
		"github/testdata/unmarshal/invalid_pubdate.xml": GitHub,
		"github/testdata/unmarshal/invalid_tag.xml":     GitHub,
		"github/testdata/unmarshal/invalid_version.xml": GitHub,
		"github/testdata/unmarshal/prerelease.xml":      GitHub,

		// SourceForge RSS Feed
		"sourceforge/testdata/unmarshal/default.xml":         SourceForge,
		"sourceforge/testdata/unmarshal/empty.xml":           SourceForge,
		"sourceforge/testdata/unmarshal/invalid_pubdate.xml": SourceForge,
		"sourceforge/testdata/unmarshal/invalid_tag.xml":     SourceForge,
		"sourceforge/testdata/unmarshal/invalid_version.xml": SourceForge,
		"sourceforge/testdata/unmarshal/prerelease.xml":      SourceForge,

		// Sparkle RSS Feed
		"sparkle/testdata/unmarshal/attributes_as_elements.xml": Sparkle,
		"sparkle/testdata/unmarshal/default_asc.xml":            Sparkle,
		"sparkle/testdata/unmarshal/default.xml":                Sparkle,
		"sparkle/testdata/unmarshal/incorrect_namespace.xml":    Sparkle,
		"sparkle/testdata/unmarshal/multiple_enclosure.xml":     Sparkle,
		"sparkle/testdata/unmarshal/no_releases.xml":            Sparkle,
		"sparkle/testdata/unmarshal/single.xml":                 Sparkle,
		"sparkle/testdata/unmarshal/with_comments.xml":          Sparkle,
		"sparkle/testdata/unmarshal/without_namespaces.xml":     Sparkle,

		// Unknown
		"../testdata/unknown.xml": Unknown,
	}

	for filename, provider := range testCases {
		assert.Equal(t, provider, GuessProviderByContent(getTestdata(filename)), fmt.Sprintf("Provider doesn't match: %s", filename))
	}
}

func TestGuessProviderByContentString(t *testing.T) {
	assert.Equal(t, Sparkle, GuessProviderByContentString(string(getTestdata("sparkle/testdata/unmarshal/default.xml"))))
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
