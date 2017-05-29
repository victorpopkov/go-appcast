package appcast

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestNewRelease(t *testing.T) {
	// test (successful)
	r, err := NewRelease("1.0.0", "1000")
	assert.Nil(t, err)
	assert.IsType(t, Release{}, *r)
	assert.Equal(t, "1.0.0", r.Version.String())
	assert.Equal(t, "1000", r.Build)

	// test (error)
	r, err = NewRelease("invalid", "1000")
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestSetVersion(t *testing.T) {
	// test (successful)
	r := new(Release)
	assert.Nil(t, r.Version)
	err := r.SetVersion("1.0.0")
	assert.Nil(t, err)
	assert.IsType(t, version.Version{}, *r.Version)
	assert.Equal(t, "1.0.0", r.Version.String())

	// test (error)
	r = new(Release)
	assert.Nil(t, r.Version)
	err = r.SetVersion("invalid")
	assert.Error(t, err)
	assert.Nil(t, r.Version)
}

func TestParsePublishedDateTime(t *testing.T) {
	testCases := map[string]string{
		"Sun, 14 May 2017 05:04:01 -0700": "2017-05-14 12:04:01 +0000 UTC", // RFC1123Z
		"Thu, 25 May 2017 19:26:48 UTC":   "2017-05-25 19:26:48 +0000 UTC", // RFC1123

		// custom
		"Thu, 25 May 2017 19:26:48 UT":              "2017-05-25 19:26:48 +0000 UTC",
		"Monday, January 12th, 2010 23:30:00 GMT-5": "2010-01-12 23:30:00 +0000 UTC",
	}

	for dateTime, expected := range testCases {
		r := new(Release)
		r.ParsePublishedDateTime(dateTime)
		assert.Equal(t, expected, r.PublishedDateTime.String())
	}
}
