package appcast

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestSparkleRSSFeedAppcast creates a new SparkleRSSFeedAppcast instance for
// testing purposes and returns its pointer. By default the content is
// []byte("test"). However, own content can be provided as an argument.
func newTestSparkleRSSFeedAppcast(content ...interface{}) *SparkleRSSFeedAppcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://example.com/appcast.xml"
	r, _ := NewRequest(url)

	appcast := &SparkleRSSFeedAppcast{
		Appcast: Appcast{
			source: &RemoteSource{
				Source: &Source{
					content:  resultContent,
					provider: SparkleRSSFeed,
				},
				request: r,
				url:     url,
			},
		},
	}

	return appcast
}

func TestSparkleRSSFeedAppcast_UnmarshalReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"sparkle/attributes_as_elements.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/default_asc.xml": {
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle/default.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/incorrect_namespace.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/invalid_pubdate.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"0001-01-01 00:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		// "sparkle/multiple_enclosure.xml": {},
		"sparkle/no_releases.xml": {},
		"sparkle/only_version.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "2.0.0", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "1.1.0", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "1.0.1", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "1.0.0", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/prerelease.xml": {
			"2.0.0-beta": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0_beta.dmg", "10.10"},
			"1.1.0":      {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1":      {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0":      {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/single.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle/without_namespaces.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
	}

	errorTestCases := map[string]string{
		"sparkle/invalid_version.xml": "Malformed version: invalid",
		"sparkle/with_comments.xml":   "version is required, but it's not specified in release #1",
	}

	// test (successful)
	for filename, releases := range testCases {
		// preparations
		a := newTestSparkleRSSFeedAppcast(getTestdata(filename))
		assert.Empty(t, a.releases)

		// test
		a.Uncomment()
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, a)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
		assert.Len(t, a.releases, len(releases))
		for _, release := range a.releases {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("Release %s", v), release.Title())
			assert.Equal(t, fmt.Sprintf("Release %s Description", v), release.Description())

			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())
			assert.Equal(t, releases[v][1], release.Build())
			assert.Equal(t, releases[v][3], release.MinimumSystemVersion())

			// downloads
			assert.Equal(t, releases[v][2], release.Downloads()[0].URL)
			assert.Equal(t, "application/octet-stream", release.Downloads()[0].Type)
			assert.Equal(t, 100000, release.Downloads()[0].Length)
		}
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// preparations
		a := newTestSparkleRSSFeedAppcast(getTestdata(filename))

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
	}
}

func TestSparkleRSSFeedAppcast_Uncomment(t *testing.T) {
	testCases := map[string][]int{
		"sparkle/attributes_as_elements.xml": nil,
		"sparkle/default_asc.xml":            nil,
		"sparkle/default.xml":                nil,
		"sparkle/incorrect_namespace.xml":    nil,
		"sparkle/multiple_enclosure.xml":     nil,
		"sparkle/single.xml":                 nil,
		"sparkle/with_comments.xml":          {13, 20},
		"sparkle/without_namespaces.xml":     nil,
	}

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// preparations
	a := new(SparkleRSSFeedAppcast)

	// test (when no content)
	assert.Nil(t, a.Source())
	err := a.Uncomment()
	assert.NotNil(t, err)

	// test (uncommenting)
	for filename, commentLines := range testCases {
		// preparations
		a = newTestSparkleRSSFeedAppcast(getTestdata(filename))

		// before SparkleRSSFeedAppcast.Uncomment
		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.Source().Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.True(t, check, fmt.Sprintf("\"%s\" doesn't have a commented out line", filename))
		}

		// tested method
		a.Uncomment()

		// after SparkleRSSFeedAppcast.Uncomment
		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.Source().Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.False(t, check, fmt.Sprintf("\"%s\" didn't uncomment a \"%d\" line", filename, commentLine))
		}
	}
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
