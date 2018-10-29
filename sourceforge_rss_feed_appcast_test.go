package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorpopkov/go-appcast/client"
)

// newTestSourceForgeRSSFeedAppcast creates a new SourceForgeRSSFeedAppcast
// instance for testing purposes and returns its pointer. By default the content
// is []byte("test"). However, own content can be provided as an argument.
func newTestSourceForgeRSSFeedAppcast(content ...interface{}) *SourceForgeRSSFeedAppcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://sourceforge.net/projects/test/rss"
	r, _ := client.NewRequest(url)

	appcast := &SourceForgeRSSFeedAppcast{
		Appcast: Appcast{
			source: &RemoteSource{
				Source: &Source{
					content:  resultContent,
					provider: SourceForgeRSSFeed,
				},
				request: r,
				url:     url,
			},
		},
	}

	return appcast
}

func TestSourceForgeRSSFeedAppcast_UnmarshalReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"sourceforge/default.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
		"sourceforge/empty.xml": {},
		"sourceforge/invalid_pubdate.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
		"sourceforge/single.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{
		"sourceforge/invalid_version.xml": "version is required, but it's not specified in release #2",
	}

	// test (successful)
	for filename, releases := range testCases {
		// preparations
		a := newTestSourceForgeRSSFeedAppcast(getTestdata(filename))
		assert.Empty(t, a.releases)

		// test
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &SourceForgeRSSFeedAppcast{}, a)
		assert.IsType(t, &SourceForgeRSSFeedAppcast{}, p)
		assert.Len(t, a.releases, len(releases))
		for _, release := range a.releases {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("/app/%s/app_%s.dmg", v, v), release.Title())
			assert.Equal(t, fmt.Sprintf("/app/%s/app_%s.dmg", v, v), release.Description())
			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())
			assert.Equal(t, releases[v][1], release.Downloads()[0].Url())
			assert.Equal(t, "application/octet-stream", release.Downloads()[0].Filetype())
			assert.Equal(t, 100000, release.Downloads()[0].Length())
		}
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// preparations
		a := newTestSourceForgeRSSFeedAppcast(getTestdata(filename))

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.IsType(t, &SourceForgeRSSFeedAppcast{}, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
	}
}
