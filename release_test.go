package appcast

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewRelease(t *testing.T) {
	// preparations
	v := "1.0.0"
	b := "1000"

	// test (successful)
	r, err := NewRelease(v, b)
	assert.Nil(t, err)
	assert.IsType(t, Release{}, *r)
	assert.Equal(t, v, r.version.String())
	assert.Equal(t, b, r.build)

	// test (error)
	r, err = NewRelease("invalid", b)
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestRelease_VersionOrBuildString(t *testing.T) {
	// preparations
	v := "1.0.0"
	b := "1000"
	r := new(Release)
	r.build = b

	// test (only build is set)
	assert.Equal(t, b, r.VersionOrBuildString())

	// test (both build and version are set)
	r.SetVersionString(v)
	assert.Equal(t, v, r.VersionOrBuildString())
}

func TestRelease_Version(t *testing.T) {
	// preparations
	r := new(Release)
	r.SetVersionString("1.0.0")

	// test
	assert.Equal(t, r.version, r.Version())
}

func TestRelease_SetVersion(t *testing.T) {
	// preparations
	r := new(Release)
	v, err := version.NewVersion("1.0.0")
	assert.Nil(t, err)
	assert.Nil(t, r.version)

	// test
	r.SetVersion(v)
	assert.Equal(t, v, r.version)
}

func TestRelease_SetVersionString(t *testing.T) {
	// preparations
	v := "1.0.0"

	// test (successful)
	r := new(Release)
	assert.Nil(t, r.version)
	err := r.SetVersionString(v)
	assert.Nil(t, err)
	assert.Equal(t, v, r.version.String())

	// test (error)
	r = new(Release)
	assert.Nil(t, r.version)
	err = r.SetVersionString("invalid")
	assert.Error(t, err)
	assert.Nil(t, r.version)
}

func TestRelease_Build(t *testing.T) {
	// preparations
	r := new(Release)
	r.build = "1000"

	// test
	assert.Equal(t, r.build, r.Build())
}

func TestRelease_SetBuild(t *testing.T) {
	// preparations
	b := "1000"
	r := new(Release)

	// test
	r.SetBuild(b)
	assert.Equal(t, b, r.build)
}

func TestRelease_Title(t *testing.T) {
	// preparations
	r := new(Release)
	r.title = "title"

	// test
	assert.Equal(t, r.title, r.Title())
}

func TestRelease_SetTitle(t *testing.T) {
	// preparations
	title := "title"
	r := new(Release)

	// test
	r.SetTitle(title)
	assert.Equal(t, title, r.title)
}

func TestRelease_Description(t *testing.T) {
	// preparations
	r := new(Release)
	r.description = "description"

	// test
	assert.Equal(t, r.description, r.Description())
}

func TestRelease_SetDescription(t *testing.T) {
	// preparations
	d := "description"
	r := new(Release)

	// test
	r.SetDescription(d)
	assert.Equal(t, d, r.description)
}

func TestRelease_AddDownload(t *testing.T) {
	// preparations
	r := new(Release)
	assert.Len(t, r.downloads, 0)

	// test
	r.AddDownload(*NewDownload("https://example.com/one.dmg", "application/octet-stream", 100000))
	r.AddDownload(*NewDownload("https://example.com/two.dmg", "application/octet-stream", 100000))
	assert.Len(t, r.downloads, 2)
}

func TestRelease_Downloads(t *testing.T) {
	// preparations
	r := new(Release)
	r.AddDownload(*NewDownload("https://example.com/one.dmg", "application/octet-stream", 100000))
	r.AddDownload(*NewDownload("https://example.com/two.dmg", "application/octet-stream", 100000))

	// test
	assert.Len(t, r.Downloads(), 2)
}

func TestRelease_SetDownloads(t *testing.T) {
	// preparations
	r := new(Release)
	assert.Len(t, r.downloads, 0)

	// test
	r.SetDownloads([]Download{
		*NewDownload("https://example.com/one.dmg", "application/octet-stream", 100000),
		*NewDownload("https://example.com/two.dmg", "application/octet-stream", 100000),
	})
	assert.Len(t, r.downloads, 2)
}

func TestRelease_ParsePublishedDateTime(t *testing.T) {
	testCases := map[string]string{
		"Sun, 14 May 2017 05:04:01 -0700": "2017-05-14 12:04:01 +0000 UTC", // RFC1123Z
		"Thu, 25 May 2017 19:26:48 UTC":   "2017-05-25 19:26:48 +0000 UTC", // RFC1123
		"2016-05-13T12:00:00+02:00":       "2016-05-13 10:00:00 +0000 UTC", // RFC3339

		// custom
		"Thu, 25 May 2017 19:26:48 UT":              "2017-05-25 19:26:48 +0000 UTC",
		"Monday, January 12th, 2010 23:30:00 GMT-5": "2010-01-12 23:30:00 +0000 UTC",
	}

	for dateTime, expected := range testCases {
		r := new(Release)
		r.ParsePublishedDateTime(dateTime)
		assert.Equal(t, expected, r.publishedDateTime.String())
	}
}

func TestRelease_PublishedDateTime(t *testing.T) {
	// preparations
	testTime := time.Now()
	r := new(Release)
	r.publishedDateTime = testTime

	// test
	assert.Equal(t, r.publishedDateTime, r.PublishedDateTime())
}

func TestRelease_SetPublishedDateTime(t *testing.T) {
	// preparations
	testTime := time.Now()
	r := new(Release)

	// test
	r.SetPublishedDateTime(testTime)
	assert.Equal(t, testTime, r.publishedDateTime)
}

func TestRelease_IsPreRelease(t *testing.T) {
	// preparations
	r := new(Release)
	r.isPreRelease = true

	// test
	assert.Equal(t, r.isPreRelease, r.IsPreRelease())
}

func TestRelease_SetIsPreRelease(t *testing.T) {
	// preparations
	r := new(Release)

	// test
	r.SetIsPreRelease(true)
	assert.Equal(t, true, r.isPreRelease)
}
