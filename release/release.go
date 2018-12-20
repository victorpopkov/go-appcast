// Package release provides an appcast release(s) that suits most appcast
// providers.
//
// Officially, it supports 3 providers: "GitHub Atom Feed", "SourceForge RSS
// Feed" and "Sparkle RSS Feed". However, it can be extended to your own needs
// if necessary.
package release

import (
	"fmt"

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
	PublishedDateTime() *PublishedDateTime
	SetPublishedDateTime(publishedDateTime *PublishedDateTime)
	ReleaseNotesLink() string
	SetReleaseNotesLink(releaseNotesLink string)
	MinimumSystemVersion() string
	SetMinimumSystemVersion(minimumSystemVersion string)
	AddDownload(d Download)
	Downloads() []Download
	SetDownloads(downloads []Download)
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

	// publishedDateTime specifies the release published date and time.
	publishedDateTime *PublishedDateTime

	// releaseNotesLink specifies a link to the release notes.
	releaseNotesLink string

	// minimumSystemVersion specifies the required system version for the
	// current app release.
	minimumSystemVersion string

	// downloads specify a slice of Download structs which represents a list of
	// all current release downloads.
	downloads []Download

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
		return fmt.Errorf("malformed version: %s", value)
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

// PublishedDateTime is a Release.publishedDateTime getter.
func (r *Release) PublishedDateTime() *PublishedDateTime {
	return r.publishedDateTime
}

// SetPublishedDateTime is a Release.publishedDateTime setter.
func (r *Release) SetPublishedDateTime(publishedDateTime *PublishedDateTime) {
	r.publishedDateTime = publishedDateTime
}

// ReleaseNotesLink is a Release.releaseNotesLink getter.
func (r *Release) ReleaseNotesLink() string {
	return r.releaseNotesLink
}

// SetReleaseNotesLink is a Release.releaseNotesLink setter.
func (r *Release) SetReleaseNotesLink(releaseNotesLink string) {
	r.releaseNotesLink = releaseNotesLink
}

// MinimumSystemVersion is a Release.minimumSystemVersion getter.
func (r *Release) MinimumSystemVersion() string {
	return r.minimumSystemVersion
}

// SetMinimumSystemVersion is a Release.minimumSystemVersion setter.
func (r *Release) SetMinimumSystemVersion(minimumSystemVersion string) {
	r.minimumSystemVersion = minimumSystemVersion
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

// IsPreRelease is a Release.isPreRelease getter.
func (r *Release) IsPreRelease() bool {
	return r.isPreRelease
}

// SetIsPreRelease is a Release.isPreRelease setter.
func (r *Release) SetIsPreRelease(isPreRelease bool) {
	r.isPreRelease = isPreRelease
}
