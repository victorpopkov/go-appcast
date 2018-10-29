// Package release provides functionality for appcast releases.
package release

import (
	"errors"
	"regexp"
	"time"
)

var PublishedDateTimeFormats = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
	"Monday, January 02, 2006 15:04:05 MST",
}

// PublishedDateTimer is the interface that wraps the PublishedDateTime methods.
type PublishedDateTimer interface {
	Parse(dateTime string) (err error)
	Time() *time.Time
	SetTime(time *time.Time)
	Format() string
	SetFormat(format string)
}

// PublishedDateTime represents a release published date and time.
type PublishedDateTime struct {
	// original specifies the original datetime string.
	original string

	// time specifies the original time.
	time *time.Time

	// format represents the format in which
	format string
}

// NewPublishedDateTime returns a new PublishedDateTime instance pointer.
// Optionally, the time can be passed as a parameter.
func NewPublishedDateTime(a ...interface{}) *PublishedDateTime {
	d := new(PublishedDateTime)

	if len(a) > 0 {
		t := a[0].(*time.Time)
		d.time = t
	}

	return d
}

// Parse parses the provided dateTime string using the predefined time formats
// in the PublishedDateTimeFormats global variable and sets the original string
// value as the PublishedDateTime.original value alongside with the
// PublishedDateTime.time and PublishedDateTime.format.
func (p *PublishedDateTime) Parse(dateTime string) (err error) {
	var re *regexp.Regexp

	// remove suffixes "st|nd|rd|th" from day digit
	re = regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	if re.MatchString(dateTime) {
		// extract last part that represents version
		versionMatches := re.FindAllStringSubmatch(dateTime, 1)
		dateTime = re.ReplaceAllString(dateTime, versionMatches[0][1])
	}

	// change "UT" to "UTC"
	re = regexp.MustCompile(`UT$`)
	if re.MatchString(dateTime) {
		dateTime = re.ReplaceAllString(dateTime, "UTC")
	}

	// parse by predefined formats
	for _, format := range PublishedDateTimeFormats {
		parsedTime, err := time.Parse(format, dateTime)
		if err == nil {
			p.time = &parsedTime
			p.format = format

			return nil
		}
	}

	return errors.New("parsing of the published datetime failed")
}

// Time is a PublishedDateTime.time getter.
func (p *PublishedDateTime) Time() *time.Time {
	return p.time
}

// SetTime is a PublishedDateTime.time setter.
func (p *PublishedDateTime) SetTime(time *time.Time) {
	p.time = time
}

// Format is a PublishedDateTime.format getter.
func (p *PublishedDateTime) Format() string {
	return p.format
}

// SetFormat is a PublishedDateTime.format setter.
func (p *PublishedDateTime) SetFormat(format string) {
	p.format = format
}

// String returns the string representation of the PublishedDateTime.
func (p *PublishedDateTime) String() string {
	if p.time == nil {
		return ""
	}

	if p.format != "" {
		return p.time.Format(p.format)
	}

	return p.time.String()
}
