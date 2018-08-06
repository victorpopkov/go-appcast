package appcast

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.True(t, check, fmt.Sprintf("\"%s\" doesn't have a commented out line", filename))
		}

		// tested function
		a.Uncomment()

		// after SparkleRSSFeedAppcast.Uncomment
		for _, commentLine := range commentLines {
			line, _ := getLineFromString(commentLine, a.Content)
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.False(t, check, fmt.Sprintf("\"%s\" didn't uncomment a \"%d\" line", filename, commentLine))
		}
	}
}

func TestSparkleRSSFeedAppcast_ExtractReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"sparkle/attributes_as_elements.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/default_asc.xml": {
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle/default.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/incorrect_namespace.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/invalid_pubdate.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"0001-01-01 00:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		// "sparkle/multiple_enclosure.xml": {},
		"sparkle/no_releases.xml": {},
		"sparkle/only_version.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "2.0.0", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "1.1.0", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "1.0.1", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "1.0.0", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/prerelease.xml": {
			"2.0.0-beta": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0_beta.dmg", "10.10"},
			"1.1.0":      {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1":      {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0":      {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
		"sparkle/single.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
		},
		"sparkle/without_namespaces.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
		},
	}

	errorTestCases := map[string]string{
		"sparkle/invalid_version.xml": "Malformed version: invalid",
		"sparkle/with_comments.xml":   "version is required, but it's not specified in release #1",
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
			assert.Equal(t, releases[v][2], release.Downloads[0].URL)
			assert.Equal(t, "application/octet-stream", release.Downloads[0].Type)
			assert.Equal(t, 100000, release.Downloads[0].Length)
		}
	}

	// test (error)
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
