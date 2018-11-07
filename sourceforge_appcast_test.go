package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/client"
)

// newTestSourceForgeRSSFeedAppcast creates a new SourceForgeAppcast
// instance for testing purposes and returns its pointer. By default the content
// is []byte("test"). However, own content can be provided as an argument.
func newTestSourceForgeRSSFeedAppcast(content ...interface{}) *SourceForgeAppcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://sourceforge.net/projects/test/rss"
	r, _ := client.NewRequest(url)

	appcast := &SourceForgeAppcast{
		Appcast: Appcast{
			source: &RemoteSource{
				Source: &Source{
					content:  resultContent,
					provider: SourceForge,
				},
				request: r,
				url:     url,
			},
		},
	}

	return appcast
}

func TestSourceForgeAppcast_Unmarshal(t *testing.T) {
	testCases := map[string]map[string][]string{
		"default.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"Thu, 12 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
		"empty.xml": {},
		"invalid_pubdate.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"Wed, 11 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"Tue, 10 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
		"single.xml": {
			"2.0.0": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{
		"invalid_tag.xml":     "XML syntax error on line 21: element <content> closed by </item>",
		"invalid_version.xml": "no version in the #2 release",
	}

	// test (successful)
	for path, releases := range testCases {
		// preparations
		a := newTestSourceForgeRSSFeedAppcast(getTestdata("sourceforge", path))

		// test
		assert.IsType(t, &SourceForgeAppcast{}, a)
		assert.Nil(t, a.source.Appcast())
		assert.Empty(t, a.releases)

		p, err := a.Unmarshal()
		p, err = a.UnmarshalReleases()

		assert.Nil(t, err)
		assert.IsType(t, &SourceForgeAppcast{}, p)
		assert.IsType(t, &SourceForgeAppcast{}, a.source.Appcast())

		assert.Len(t, releases, a.releases.Len())
		for _, release := range a.releases.Filtered() {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("/app/%s/app_%s.dmg", v, v), release.Title())
			assert.Equal(t, fmt.Sprintf("/app/%s/app_%s.dmg", v, v), release.Description())
			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())

			// downloads
			assert.Equal(t, releases[v][1], release.Downloads()[0].Url())
			assert.Equal(t, "application/octet-stream", release.Downloads()[0].Filetype())
			assert.Equal(t, 100000, release.Downloads()[0].Length())
		}
	}

	// test (error) [unmarshalling failure]
	for path, errorMsg := range errorTestCases {
		// preparations
		a := newTestSourceForgeRSSFeedAppcast(getTestdata("sourceforge", path))

		// test
		assert.IsType(t, &SourceForgeAppcast{}, a)
		assert.Nil(t, a.source.Appcast())

		p, err := a.Unmarshal()

		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, p)
		assert.IsType(t, &SourceForgeAppcast{}, a.source.Appcast())
	}

	// test (error) [no source]
	a := new(SourceForgeAppcast)

	p, err := a.Unmarshal()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.source)
}
