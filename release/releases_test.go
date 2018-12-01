package release

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestReleases creates a new Releases instance for testing purposes and
// returns its pointer.
func newTestReleases() *Releases {
	d := Download{
		url:      "https://example.com/app_2.0.0.dmg",
		filetype: "application/octet-stream",
		length:   100000,
	}

	r := Release{
		build:                "200",
		title:                "Release 2.0.0-beta",
		description:          "Release 2.0.0 Description",
		releaseNotesLink:     "https://example.com/changelogs/2.0.0.html",
		minimumSystemVersion: "10.10",
		downloads:            []Download{d},
		isPreRelease:         true,
	}

	// r1
	d1 := d
	r1 := r
	r1.downloads = []Download{d1}

	err := r1.SetVersionString("2.0.0-beta")
	if err != nil {
		panic(err)
	}

	t, _ := time.Parse(time.RFC1123Z, "Fri, 13 May 2016 12:00:00 +0200")
	r1.publishedDateTime = NewPublishedDateTime(&t)

	// r2
	d2 := d
	d2.url = "https://example.com/app_1.1.0.dmg"

	r2 := r
	r2.build = "110"
	r2.title = "Release 1.1.0"
	r2.description = "Release 1.1.0 Description"
	r2.releaseNotesLink = "https://example.com/changelogs/1.1.0.html"
	r2.minimumSystemVersion = "10.9"
	r2.downloads = []Download{d2}
	r2.isPreRelease = false

	err = r2.SetVersionString("1.1.0")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Thu, 12 May 2016 12:00:00 +0200")
	r2.publishedDateTime = NewPublishedDateTime(&t)

	// r3
	d3 := d
	d3.url = "https://example.com/app_1.0.1.dmg"

	r3 := r
	r3.build = "101"
	r3.title = "Release 1.0.1"
	r3.description = "Release 1.0.1 Description"
	r3.releaseNotesLink = "https://example.com/changelogs/1.0.1.html"
	r3.minimumSystemVersion = "10.9"
	r3.downloads = []Download{d3}
	r3.isPreRelease = false

	err = r3.SetVersionString("1.0.1")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Wed, 11 May 2016 12:00:00 +0200")
	r3.publishedDateTime = NewPublishedDateTime(&t)

	// r4
	d4 := d
	d4.url = "https://example.com/app_1.0.0.dmg"

	r4 := r
	r4.build = "100"
	r4.title = "Release 1.0.0"
	r4.description = "Release 1.0.0 Description"
	r4.releaseNotesLink = "https://example.com/changelogs/1.0.0.html"
	r4.minimumSystemVersion = "10.9"
	r4.downloads = []Download{d4}
	r4.isPreRelease = false

	err = r4.SetVersionString("1.0.0")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Tue, 10 May 2016 12:00:00 +0200")
	r4.publishedDateTime = NewPublishedDateTime(&t)

	return &Releases{
		filtered: []Releaser{&r1, &r2, &r3, &r4},
		original: []Releaser{&r1, &r2, &r3, &r4},
	}
}

func TestNewReleases(t *testing.T) {
	r := NewReleases(nil)
	assert.IsType(t, Releases{}, *r)
	assert.Nil(t, r.filtered)
	assert.Nil(t, r.original)
}

func TestReleases_SortByVersions(t *testing.T) {
	// preparations
	r := newTestReleases()

	// test (ASC)
	r.SortByVersions(ASC)
	assert.Equal(t, "1.0.0", r.filtered[0].Version().String())

	// test (DESC)
	r.SortByVersions(DESC)
	assert.Equal(t, "2.0.0-beta", r.filtered[0].Version().String())
}

func TestReleases_FilterByTitle(t *testing.T) {
	// preparations
	r := newTestReleases()

	// test
	assert.Len(t, r.filtered, 4)
	r.FilterByTitle("Release 1.0")
	assert.Len(t, r.filtered, 2)
	r.FilterByTitle("Release 1.0.0", true)
	assert.Len(t, r.filtered, 1)
	assert.Equal(t, "Release 1.0.1", r.filtered[0].Title())
}

func TestReleases_FilterByMediaType(t *testing.T) {
	// preparations
	r := newTestReleases()

	// test
	assert.Len(t, r.filtered, 4)
	r.FilterByMediaType("application/octet-stream")
	assert.Len(t, r.filtered, 4)
	r.FilterByMediaType("test", true)
	assert.Len(t, r.filtered, 4)
}

func TestReleases_FilterByUrl(t *testing.T) {
	// preparations
	r := newTestReleases()

	// test
	assert.Len(t, r.filtered, 4)
	r.FilterByUrl(`app_1.*dmg$`)
	assert.Len(t, r.filtered, 3)
	r.FilterByUrl(`app_1.0.*dmg$`, true)
	assert.Len(t, r.filtered, 1)
}

func TestReleases_FilterByPrerelease(t *testing.T) {
	// preparations
	r := newTestReleases()

	// test
	assert.Len(t, r.filtered, 4)
	r.FilterByPrerelease()
	assert.Len(t, r.filtered, 1)
	r.ResetFilters()

	assert.Len(t, r.filtered, 4)
	r.FilterByPrerelease(true)
	assert.Len(t, r.filtered, 3)
}

func TestReleases_Len(t *testing.T) {
	r := newTestReleases()
	assert.Equal(t, len(r.filtered), r.Len())
}

func TestReleases_First(t *testing.T) {
	r := newTestReleases()
	assert.Equal(t, r.filtered[0], r.First())
}

func TestReleases_Filtered(t *testing.T) {
	r := newTestReleases()
	assert.Equal(t, r.filtered, r.Filtered())
}

func TestReleases_SetFiltered(t *testing.T) {
	r := newTestReleases()
	r.SetFiltered(nil)
	assert.Nil(t, r.filtered)
}

func TestReleases_Original(t *testing.T) {
	r := newTestReleases()
	assert.Equal(t, r.original, r.Original())
}

func TestReleases_SetOriginal(t *testing.T) {
	r := newTestReleases()
	r.SetOriginal(nil)
	assert.Nil(t, r.original)
}
