package release

import (
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

// newTestRelease creates a new Release instance for testing purposes and
// returns its pointer.
func newTestRelease() *Release {
	v, _ := version.NewVersion("1.0.0")
	t, _ := time.Parse(time.RFC1123Z, "Fri, 13 May 2016 12:00:00 +0200")

	return &Release{
		version:              v,
		build:                "1000",
		title:                "Test",
		description:          "Test",
		publishedDateTime:    NewPublishedDateTime(t),
		releaseNotesLink:     "https://example.com/changelogs/1.0.0.html",
		minimumSystemVersion: "10.9",
		downloads: []Download{
			*NewDownload("https://example.com/1.0.0/one.dmg", "application/octet-stream", 100000),
			*NewDownload("https://example.com/1.0.0/two.dmg", "application/octet-stream", 100000),
		},
		isPreRelease: false,
	}
}

func TestNewRelease(t *testing.T) {
	// preparations
	v := "1.0.0"
	b := "1000"

	// test (successful)
	r, err := New(v, b)
	assert.Nil(t, err)
	assert.IsType(t, Release{}, *r)
	assert.Equal(t, v, r.version.String())
	assert.Equal(t, b, r.build)

	// test (error)
	r, err = New("invalid", b)
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
	r := newTestRelease()
	assert.Equal(t, r.version, r.Version())
}

func TestRelease_SetVersion(t *testing.T) {
	// preparations
	r := newTestRelease()
	v, _ := version.NewVersion("1.0.0")

	// test
	r.SetVersion(v)
	assert.Equal(t, v, r.version)
}

func TestRelease_SetVersionString(t *testing.T) {
	// preparations
	v := "1.0.1"

	// test (successful)
	r := newTestRelease()
	assert.Equal(t, "1.0.0", r.version.String())
	err := r.SetVersionString(v)
	assert.Nil(t, err)
	assert.Equal(t, v, r.version.String())

	// test (error)
	r = newTestRelease()
	assert.Equal(t, "1.0.0", r.version.String())
	err = r.SetVersionString("invalid")
	assert.Error(t, err)
	assert.Equal(t, "1.0.0", r.version.String())
}

func TestRelease_Build(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.build, r.Build())
}

func TestRelease_SetBuild(t *testing.T) {
	r := newTestRelease()
	r.SetBuild("1001")
	assert.Equal(t, "1001", r.build)
}

func TestRelease_Title(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.title, r.Title())
}

func TestRelease_SetTitle(t *testing.T) {
	r := newTestRelease()
	r.SetTitle("Title")
	assert.Equal(t, "Title", r.title)
}

func TestRelease_Description(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.description, r.Description())
}

func TestRelease_SetDescription(t *testing.T) {
	r := newTestRelease()
	r.SetDescription("Description")
	assert.Equal(t, "Description", r.description)
}

func TestRelease_PublishedDateTime(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.publishedDateTime, r.PublishedDateTime())
}

func TestRelease_SetPublishedDateTime(t *testing.T) {
	// preparations
	now := time.Now()
	r := newTestRelease()

	// test
	r.SetPublishedDateTime(NewPublishedDateTime(now))
	assert.Equal(t, now.UTC(), r.publishedDateTime.time.UTC())
}

func TestRelease_ReleaseNotesLink(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.releaseNotesLink, r.ReleaseNotesLink())
}

func TestRelease_SetReleaseNotesLink(t *testing.T) {
	r := newTestRelease()
	r.SetReleaseNotesLink("test")
	assert.Equal(t, "test", r.releaseNotesLink)
}

func TestRelease_MinimumSystemVersion(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.minimumSystemVersion, r.MinimumSystemVersion())
}

func TestRelease_SetMinimumSystemVersion(t *testing.T) {
	r := newTestRelease()
	r.SetMinimumSystemVersion("10.13.6")
	assert.Equal(t, "10.13.6", r.minimumSystemVersion)
}

func TestRelease_AddDownload(t *testing.T) {
	// preparations
	r := newTestRelease()

	// test
	assert.Len(t, r.downloads, 2)
	r.AddDownload(*NewDownload("https://example.com/1.0.0/three.dmg", "application/octet-stream", 100000))
	r.AddDownload(*NewDownload("https://example.com/1.0.0/four.dmg", "application/octet-stream", 100000))
	assert.Len(t, r.downloads, 4)
}

func TestRelease_Downloads(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.downloads, r.Downloads())
}

func TestRelease_SetDownloads(t *testing.T) {
	// preparations
	d := []Download{
		*NewDownload("https://example.com/1.0.0/one.dmg", "application/octet-stream", 100000),
	}

	// test
	r := newTestRelease()
	r.SetDownloads(d)
	assert.Equal(t, d, r.downloads)
}

func TestRelease_IsPreRelease(t *testing.T) {
	r := newTestRelease()
	assert.Equal(t, r.isPreRelease, r.IsPreRelease())
}

func TestRelease_SetIsPreRelease(t *testing.T) {
	r := newTestRelease()
	r.SetIsPreRelease(true)
	assert.Equal(t, true, r.isPreRelease)
}
