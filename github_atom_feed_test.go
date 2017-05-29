package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubAtomFeedAppcastExtractReleases(t *testing.T) {
	testCases := map[string]map[string][]string{
		"github_default.xml": {
			"2.0.0": {"2016-05-13 10:00:00 +0000 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0/app_2.0.0.dmg/download"},
			"1.1.0": {"2016-05-12 10:00:00 +0000 UTC", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1": {"2016-05-11 10:00:00 +0000 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0": {"2016-05-10 10:00:00 +0000 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{}

	// test (successful)
	for filename, releases := range testCases {
		// preparations
		a := new(GitHubAtomFeedAppcast)
		a.Content = string(getTestdata(filename))

		// test
		assert.Empty(t, a.Releases)
		err := a.ExtractReleases()
		assert.Nil(t, err)
		assert.Len(t, a.Releases, len(releases))
		for _, release := range a.Releases {
			v := release.Version.String()
			assert.Equal(t, fmt.Sprintf("%s", v), release.Title)
			assert.Equal(t, releases[v][0], release.PublishedDateTime.String())
		}
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// preparations
		a := new(GitHubAtomFeedAppcast)
		a.Content = string(getTestdata(filename))

		// test
		err := a.ExtractReleases()
		assert.Error(t, err)
		assert.Equal(t, errorMsg, err.Error())
	}
}
