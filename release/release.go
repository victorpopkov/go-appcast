// Package release provides functionality for appcast releases.
package release

import (
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-version"
)

// Releaser is the interface that wraps the Release methods.
type Releaser interface {
	VersionOrBuildString() string
	Version() *version.Version
	SetVersion(version *version.Version)
	SetVersionString(value string) error
	Build() string
	SetBuild(build string)
	Title() string
	SetTitle(title string)
	Description() string
	SetDescription(description string)
	AddDownload(d Download)
	Downloads() []Download
	SetDownloads(downloads []Download)
	ParsePublishedDateTime(dateTime string) (err error)
	PublishedDateTime() time.Time
	SetPublishedDateTime(publishedDateTime time.Time)
	IsPreRelease() bool
	SetIsPreRelease(isPreRelease bool)
}

// Release represents a single application release.
type Release struct {
	// version specifies a release version. It should follow the SemVer
	// specification.
	version *version.Version

	// build specifies a release build. This could have any value.
	build string

	// title specifies a release title.
	title string

	// description specifies a release description.
	description string

	// downloads specifies a slice of Download structs which represents a list
	// of all current release downloads.
	downloads []Download

	// publishedDateTime specifies the release published data and time in UTC.
	publishedDateTime time.Time

	// isPreRelease specifies whether a release is not stable.
	//
	// By default, each release is considered to be stable, so the default value
	// is false. If the release version, build or any other provider-specific
	// value points that a release is unstable, the value should become true.
	isPreRelease bool
}

// New returns a new Release instance pointer. Requires both version and build
// strings. By default, Release.IsPrerelease is set to false, so the release
// will be considered as stable.
func New(version string, build string) (*Release, error) {
	r := &Release{
		isPreRelease: false,
	}

	// add version
	err := r.SetVersionString(version)
	if err != nil {
		return nil, err
	}

	// add build, if its not empty
	if build != "" {
		r.build = build
	}

	return r, nil
}

// VersionOrBuildString retrieves the release version string if it's available.
// Otherwise, returns the release build string.
func (r *Release) VersionOrBuildString() string {
	if r.version == nil || r.version.String() == "" {
		return r.build
	}

	return r.version.String()
}

// Version is a Release.version getter.
func (r *Release) Version() *version.Version {
	return r.version
}

// SetVersion is a Release.version setter.
func (r *Release) SetVersion(version *version.Version) {
	r.version = version
}

// SetVersionString sets the Release.version from the provided version value
// string. Returns an error, if the provided version string value doesn't follow
// the SemVer specification.
func (r *Release) SetVersionString(value string) error {
	v, err := version.NewVersion(value)
	if err != nil {
		return err
	}

	r.version = v

	return nil
}

// Build is a Release.build getter.
func (r *Release) Build() string {
	return r.build
}

// SetBuild is a Release.build setter.
func (r *Release) SetBuild(build string) {
	r.build = build
}

// Title is a Release.title getter.
func (r *Release) Title() string {
	return r.title
}

// SetTitle is a Release.title setter.
func (r *Release) SetTitle(title string) {
	r.title = title
}

// Description is a Release.description getter.
func (r *Release) Description() string {
	return r.description
}

// SetDescription is a Release.description setter.
func (r *Release) SetDescription(description string) {
	r.description = description
}

// AddDownload appends the provided Download to the Release.downloads slice.
func (r *Release) AddDownload(d Download) {
	r.downloads = append(r.downloads, d)
}

// Downloads is a Release.downloads getter.
func (r *Release) Downloads() []Download {
	return r.downloads
}

// SetDownloads is a Release.downloads setter.
func (r *Release) SetDownloads(downloads []Download) {
	r.downloads = downloads
}

// ParsePublishedDateTime parses the provided dateTime string using the
// predefined time formats and sets the Release.publishedDateTime in UTC.
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
			r.publishedDateTime = parsedTime.UTC()
			return nil
		}
	}

	return errors.New("parsing of the published datetime failed")
}

// PublishedDateTime is a Release.publishedDateTime getter.
func (r *Release) PublishedDateTime() time.Time {
	return r.publishedDateTime
}

// SetPublishedDateTime is a Release.publishedDateTime setter.
func (r *Release) SetPublishedDateTime(publishedDateTime time.Time) {
	r.publishedDateTime = publishedDateTime
}

// IsPreRelease is a Release.isPreRelease getter.
func (r *Release) IsPreRelease() bool {
	return r.isPreRelease
}

// SetIsPreRelease is a Release.isPreRelease setter.
func (r *Release) SetIsPreRelease(isPreRelease bool) {
	r.isPreRelease = isPreRelease
}
