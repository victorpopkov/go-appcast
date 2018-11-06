package release

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestPublishedDateTime creates a new PublishedDateTime instance for testing
// purposes and returns its pointer.
func newTestPublishedDateTime() *PublishedDateTime {
	s := "Sun, 14 May 2017 12:00:00 -0200"
	f := time.RFC1123Z
	d, _ := time.Parse(f, s)

	return &PublishedDateTime{
		original: s,
		time:     &d,
		format:   f,
	}

}

func TestNewPublishedDateTime(t *testing.T) {
	// preparations
	now := time.Now()
	d := NewPublishedDateTime(&now)

	// test
	assert.IsType(t, PublishedDateTime{}, *d)
	assert.Equal(t, now.String(), d.String())
	assert.Empty(t, d.format)
}

func TestPublishedDateTime_Parse(t *testing.T) {
	testCases := map[string]string{
		"Sun, 14 May 2017 05:04:01 -0700": "2017-05-14 12:04:01 +0000 UTC", // RFC1123Z
		"Thu, 25 May 2017 19:26:48 UTC":   "2017-05-25 19:26:48 +0000 UTC", // RFC1123
		"2016-05-13T12:00:00+02:00":       "2016-05-13 10:00:00 +0000 UTC", // RFC3339

		// custom
		"Thu, 25 May 2017 19:26:48 UT":              "2017-05-25 19:26:48 +0000 UTC",
		"Monday, January 12th, 2010 23:30:00 GMT-5": "2010-01-12 23:30:00 +0000 UTC",
	}

	// test (successful)
	for dateTime, expected := range testCases {
		d := new(PublishedDateTime)
		err := d.Parse(dateTime)
		assert.Nil(t, err)
		assert.Equal(t, expected, d.time.UTC().String())
	}

	// test (error)
	d := new(PublishedDateTime)
	err := d.Parse("invalid")
	assert.Error(t, err)
	assert.EqualError(t, err, "parsing of the published datetime failed")
}

func TestPublishedDateTime_Time(t *testing.T) {
	d := newTestPublishedDateTime()
	assert.Equal(t, "2017-05-14 14:00:00 +0000 UTC", d.Time().UTC().String())
}

func TestPublishedDateTime_SetTime(t *testing.T) {
	// preparations
	now := time.Now()
	d := newTestPublishedDateTime()

	// test
	assert.Equal(t, "2017-05-14 14:00:00 +0000 UTC", d.time.UTC().String())
	d.SetTime(&now)
	assert.Equal(t, now.UTC().String(), d.time.UTC().String())
}

func TestPublishedDateTime_Format(t *testing.T) {
	d := newTestPublishedDateTime()
	assert.Equal(t, time.RFC1123Z, d.Format())
}

func TestPublishedDateTime_SetFormat(t *testing.T) {
	d := newTestPublishedDateTime()
	d.SetFormat(time.RFC1123)
	assert.Equal(t, time.RFC1123, d.format)
}

func TestPublishedDateTime_String(t *testing.T) {
	// test (successful)
	d := newTestPublishedDateTime()
	assert.Equal(t, "Sun, 14 May 2017 12:00:00 -0200", d.String())

	// test (nil)
	d.time = nil
	assert.Equal(t, "", d.String())
}
