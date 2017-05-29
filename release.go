package appcast

import (
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-version"
)

// A Release represents an application release.
type Release struct {
	// Version specifies the release version. It should follow the SemVer
	// specification.
	Version *version.Version

	// Build specifies the release build. This could have any value.
	Build string

	// Title specifies the release title.
	Title string

	// Description specifies the release description.
	Description string

	// Downloads specifies an array of downloads.
	Downloads []Download

	// PublishedDateTime specifies the release published data and time in UTC.
	PublishedDateTime time.Time

	// IsPrerelease specifies if the current release is not stable.
	//
	// By default, each release is considered to be stable, so the default value
	// is "false". If the release version, build or any other provider specific
	// value point that the release is not stable, the value should become "true".
	IsPrerelease bool
}

// NewRelease returns a new Release instance pointer. Requires both version and
// build strings. By default, Release.IsPrerelease is set to "false", so the
// release will be considered as stable.
func NewRelease(version string, build string) (*Release, error) {
	r := &Release{
		IsPrerelease: false,
	}

	// add version
	err := r.SetVersion(version)
	if err != nil {
		return nil, err
	}

	// add build, if its not empty
	if build != "" {
		r.Build = build
	}

	return r, nil
}

// SetVersion sets the Release.Version from the provided version value string.
// Returns an error, if the provided version string value doesn't follow SemVer
// specification.
func (r *Release) SetVersion(value string) error {
	v, err := version.NewVersion(value)
	if err != nil {
		return err
	}

	r.Version = v

	return nil
}

// AddDownload adds a new Download to the Release.Downloads array.
func (r *Release) AddDownload(d Download) {
	r.Downloads = append(r.Downloads, d)
}

// ParsePublishedDateTime parses the provided dateTime string using predefined
// time formats and sets the Release.PublishedDateTime in UTC.
func (r *Release) ParsePublishedDateTime(dateTime string) (err error) {
	var re *regexp.Regexp

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"Monday, January 02, 2006 15:04:05 MST",
	}

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
	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateTime)
		if err == nil {
			r.PublishedDateTime = parsedTime.UTC()
			return nil
		}
	}

	return errors.New("Parsing of the published datetime failed")
}
