package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorpopkov/go-appcast/client"
)

// newTestGitHubAtomFeedAppcast creates a new GitHubAtomFeedAppcast instance for
// testing purposes and returns its pointer. By default the content is
// []byte("test"). However, own content can be provided as an argument.
func newTestGitHubAtomFeedAppcast(content ...interface{}) *GitHubAtomFeedAppcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://github.com/user/repo/releases.atom"
	r, _ := client.NewRequest(url)

	appcast := &GitHubAtomFeedAppcast{
		Appcast: Appcast{
			source: &RemoteSource{
				Source: &Source{
					content:  resultContent,
					provider: GitHubAtomFeed,
				},
				request: r,
				url:     url,
			},
		},
	}

	return appcast
}

func TestGitHubAtomFeedAppcast_UnmarshalReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"github/default.xml": {
			"2.0.0": {"2016-05-13T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"2016-05-12T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"2016-05-11T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"2016-05-10T12:00:00+02:00", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{}

	// test (successful)
	for filename, releases := range testCases {
		// preparations
		a := newTestGitHubAtomFeedAppcast(getTestdata(filename))
		assert.Empty(t, a.releases)

		// test
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &GitHubAtomFeedAppcast{}, a)
		assert.IsType(t, &GitHubAtomFeedAppcast{}, p)
		assert.Len(t, a.releases, len(releases))
		for _, release := range a.releases {
			v := release.Version().String()
			assert.Equal(t, fmt.Sprintf("%s", v), release.Title())
			assert.NotEmpty(t, release.Description())
			assert.Equal(t, releases[v][0], release.PublishedDateTime().String())
		}
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// preparations
		a := newTestGitHubAtomFeedAppcast(getTestdata(filename))

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.Nil(t, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
	}
}
