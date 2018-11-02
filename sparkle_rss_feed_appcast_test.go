package appcast

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestSparkleRSSFeedAppcast creates a new SparkleRSSFeedAppcast instance for
// testing purposes and returns its pointer. By default the source is
// LocalSource and points to the "Sparkle RSS Feed" default.xml testdata.
func newTestSparkleRSSFeedAppcast(paths ...string) *SparkleRSSFeedAppcast {
	var path string
	var content []byte

	if len(paths) > 0 {
		path = getTestdataPath(paths...)
		content = getTestdata(paths...)
	} else {
		path = getTestdataPath("sparkle", "default.xml")
		content = getTestdata("sparkle", "default.xml")
	}

	appcast := &SparkleRSSFeedAppcast{
		Appcast: Appcast{
			source: &LocalSource{
				Source: &Source{
					content:  content,
					provider: SparkleRSSFeed,
				},
				filepath: path,
			},
		},
	}

	return appcast
}

func TestSparkleRSSFeedAppcast_UnmarshalReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"attributes_as_elements.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"default_asc.xml": {
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"default.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"incorrect_namespace.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"invalid_pubdate.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		// "multiple_enclosure.xml": {},
		"no_releases.xml": {},
		"only_version.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "2.0.0", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "1.1.0", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "1.0.1", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "1.0.0", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"prerelease.xml": {
			"2.0.0-beta": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0_beta.dmg", "10.10"},
			"1.1.0":      {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1":      {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0":      {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"single.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"without_namespaces.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
	}

	errorTestCases := map[string]string{
		"invalid_version.xml": "malformed version: invalid",
		"with_comments.xml":   "version is required, but it's not specified in release #1",
	}

	// test (successful)
	for path, releases := range testCases {
		// preparations
		a := newTestSparkleRSSFeedAppcast("sparkle", path)

		// test
		assert.IsType(t, &SparkleRSSFeedAppcast{}, a)
		assert.Nil(t, a.source.Appcast())
		assert.Nil(t, a.channel)
		assert.Empty(t, a.releases)

		p, err := a.UnmarshalReleases()

		assert.Nil(t, err)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
		//assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())

		assert.IsType(t, &SparkleRSSFeedAppcastChannel{}, a.channel)
		assert.Equal(t, "App", a.channel.Title)
		assert.Equal(t, "https://example.com/app/", a.channel.Link)
		assert.Equal(t, "App Description", a.channel.Description)
		assert.Equal(t, "en", a.channel.Language)

		assert.Len(t, a.releases, len(releases))
		for _, release := range a.releases {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("Release %s", v), release.Title())
			assert.Equal(t, fmt.Sprintf("Release %s Description", v), release.Description())
			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())
			assert.Equal(t, releases[v][1], release.Build())
			assert.Equal(t, releases[v][3], release.MinimumSystemVersion())

			// downloads
			assert.Equal(t, releases[v][2], release.Downloads()[0].Url())
			assert.Equal(t, "application/octet-stream", release.Downloads()[0].Filetype())
			assert.Equal(t, 100000, release.Downloads()[0].Length())
		}
	}

	// test (error) [unmarshalling failure]
	for path, errorMsg := range errorTestCases {
		// preparations
		a := newTestSparkleRSSFeedAppcast("sparkle", path)

		// test
		assert.IsType(t, &SparkleRSSFeedAppcast{}, a)
		assert.Nil(t, a.source.Appcast())
		assert.Nil(t, a.channel)

		p, err := a.UnmarshalReleases()

		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, p)
		//assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())
		//assert.Nil(t, a.channel)
	}

	// test (error) [no source]
	a := new(SparkleRSSFeedAppcast)

	p, err := a.UnmarshalReleases()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.source)
	assert.Nil(t, a.channel)
}

func TestSparkleRSSFeedAppcast_Uncomment(t *testing.T) {
	testCases := map[string][]int{
		"attributes_as_elements.xml": nil,
		"default_asc.xml":            nil,
		"default.xml":                nil,
		"incorrect_namespace.xml":    nil,
		"multiple_enclosure.xml":     nil,
		"single.xml":                 nil,
		"with_comments.xml":          {13, 20},
		"without_namespaces.xml":     nil,
	}

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// test (successful)
	for filename, commentLines := range testCases {
		// preparations
		a := newTestSparkleRSSFeedAppcast("sparkle", filename)

		// before
		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.source.Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.True(t, check, fmt.Sprintf("\"%s\" doesn't have a commented out line", filename))
		}

		err := a.Uncomment()

		// after
		assert.Nil(t, err)

		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.source.Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.False(t, check, fmt.Sprintf("\"%s\" didn't uncomment a \"%d\" line", filename, commentLine))
		}
	}

	// test (error) [no source]
	a := new(SparkleRSSFeedAppcast)

	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, a.source)
	assert.Nil(t, a.channel)
}

func TestSparkleRSSFeedAppcast_Channel(t *testing.T) {
	a := newTestSparkleRSSFeedAppcast()
	assert.Equal(t, a.channel, a.Channel())
}

func TestSparkleRSSFeedAppcast_SetChannel(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast()
	assert.Nil(t, a.channel)

	// test
	a.SetChannel(&SparkleRSSFeedAppcastChannel{})
	assert.NotNil(t, a.channel)
}
