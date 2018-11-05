package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorpopkov/go-appcast/client"
)

// newTestGitHubAtomFeedAppcast creates a new GitHubAppcast instance for
// testing purposes and returns its pointer. By default the content is
// []byte("test"). However, own content can be provided as an argument.
func newTestGitHubAtomFeedAppcast(content ...interface{}) *GitHubAppcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://github.com/user/repo/releases.atom"
	r, _ := client.NewRequest(url)

	appcast := &GitHubAppcast{
		Appcast: Appcast{
			source: &RemoteSource{
				Source: &Source{
					content:  resultContent,
					provider: GitHub,
				},
				request: r,
				url:     url,
			},
		},
	}

	return appcast
}

func TestGitHubAppcast_Unmarshal(t *testing.T) {
	testCases := map[string]map[string][]string{
		"default.xml": {
			"2.0.0": {"2016-05-13T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"2016-05-12T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"2016-05-11T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"2016-05-10T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{
		"invalid_tag.xml":     "XML syntax error on line 18: element <thumbnail> closed by </entry>",
		"invalid_version.xml": "malformed version: invalid",
	}

	// test (successful)
	for path, releases := range testCases {
		// preparations
		a := newTestGitHubAtomFeedAppcast(getTestdata("github", path))

		// test
		assert.IsType(t, &GitHubAppcast{}, a)
		assert.Nil(t, a.source.Appcast())
		assert.Empty(t, a.releases)

		p, err := a.Unmarshal()
		p, err = a.UnmarshalReleases()

		assert.Nil(t, err)
		assert.IsType(t, &GitHubAppcast{}, p)
		assert.IsType(t, &GitHubAppcast{}, a.source.Appcast())

		assert.Len(t, a.releases, len(releases))
		for _, release := range a.releases {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("%s", v), release.Title())
			assert.NotEmpty(t, release.Description())
			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())
		}
	}

	// test (error) [unmarshalling failure]
	for path, errorMsg := range errorTestCases {
		// preparations
		a := newTestGitHubAtomFeedAppcast(getTestdata("github", path))

		// test
		assert.IsType(t, &GitHubAppcast{}, a)
		assert.Nil(t, a.source.Appcast())

		p, err := a.Unmarshal()

		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, p)
		assert.IsType(t, &GitHubAppcast{}, a.source.Appcast())
	}

	// test (error) [no source]
	a := new(GitHubAppcast)

	p, err := a.Unmarshal()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.source)
}

func TestGitHubAppcast_UnmarshalReleases(t *testing.T) {
	TestGitHubAppcast_Unmarshal(t)
}
