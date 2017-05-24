package appcast

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparkleRSSFeedAppcastUncomment(t *testing.T) {
	testCases := map[string][]int{
		"sparkle_attributes_as_elements.xml": nil,
		"sparkle_default_asc.xml":            nil,
		"sparkle_default.xml":                nil,
		"sparkle_incorrect_namespace.xml":    nil,
		"sparkle_multiple_enclosure.xml":     nil,
		"sparkle_single.xml":                 nil,
		"sparkle_with_comments.xml":          {13, 20},
		"sparkle_without_namespaces.xml":     nil,
	}

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// preparations
	a := new(SparkleRSSFeedAppcast)

	// test (when no content)
	assert.Empty(t, a.Content)
	a.Uncomment()
	assert.Empty(t, a.Content)

	// test (uncommenting)
	for filename, commentLines := range testCases {
		a = new(SparkleRSSFeedAppcast)
		a.Content = string(getTestdata(filename))

		// before SparkleRSSFeedAppcast.Uncomment
		for _, commentLine := range commentLines {
			line, _ := getLineFromString(commentLine, a.Content)
			check := (regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line))
			assert.True(t, check, fmt.Sprintf("\"%s\" doesn't have a commented out line", filename))
		}

		// tested function
		a.Uncomment()

		// after SparkleRSSFeedAppcast.Uncomment
		for _, commentLine := range commentLines {
			line, _ := getLineFromString(commentLine, a.Content)
			check := (regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line))
			assert.False(t, check, fmt.Sprintf("\"%s\" didn't uncomment a \"%d\" line", filename, commentLine))
		}
	}
}

func TestSparkleRSSFeedAppcastExtractReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"sparkle_attributes_as_elements.xml": {
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 12:00:00 +0200 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 12:00:00 +0200 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 12:00:00 +0200 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle_default_asc.xml": {
			"1.1.0": {"2016-05-12 12:00:00 +0200 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 12:00:00 +0200 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 12:00:00 +0200 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle_default.xml": {
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 12:00:00 +0200 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 12:00:00 +0200 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 12:00:00 +0200 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle_incorrect_namespace.xml": {
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 12:00:00 +0200 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 12:00:00 +0200 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 12:00:00 +0200 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		// "sparkle_multiple_enclosure.xml": {},
		"sparkle_no_releases.xml": {},
		"sparkle_single.xml": {
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle_without_namespaces.xml": {
			"2.0.0": {"2016-05-13 12:00:00 +0200 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 12:00:00 +0200 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 12:00:00 +0200 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 12:00:00 +0200 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
	}

	errorTestCases := map[string]string{
		"sparkle_invalid_pubdate.xml": "parsing time \"invalid\" as \"Mon, 02 Jan 2006 15:04:05 -0700\": cannot parse \"invalid\" as \"Mon\"",
		"sparkle_invalid_version.xml": "Malformed version: invalid",
		"sparkle_with_comments.xml":   "Version is required, but it's not specified in release #1",
	}

	// test (successful)
	for filename, releases := range testCases {
		// preparations
		a := new(SparkleRSSFeedAppcast)
		a.Content = string(getTestdata(filename))

		// test
		assert.Empty(t, a.Releases)
		a.Uncomment()
		err := a.ExtractReleases()
		assert.Nil(t, err)
		assert.Len(t, a.Releases, len(releases))
		for _, release := range a.Releases {
			v := release.Version.String()
			assert.Equal(t, fmt.Sprintf("Release %s", v), release.Title)
			assert.Equal(t, fmt.Sprintf("Release %s Description", v), release.Description)
			assert.Equal(t, releases[v][0], release.PublishedDateTime.String())
			assert.Equal(t, releases[v][1], release.Build)
			assert.Equal(t, releases[v][2], release.DownloadURLs[0])
		}
	}

	// test error
	for filename, errorMsg := range errorTestCases {
		// preparations
		a := new(SparkleRSSFeedAppcast)
		a.Content = string(getTestdata(filename))

		// test
		err := a.ExtractReleases()
		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())
	}
}
