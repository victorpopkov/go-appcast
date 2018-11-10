package sourceforge

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
)

var testdataPath = "./testdata/"

// getWorkingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func getWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// getTestdata returns a file content as a byte slice from the provided testdata
// paths. If the file is not found, prints an error to os.Stdout and exits with
// exit status 1.
func getTestdata(paths ...string) []byte {
	path := getTestdataPath(paths...)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// getTestdataPath returns a full path for the provided testdata paths.
func getTestdataPath(paths ...string) string {
	return filepath.Join(getWorkingDir(), testdataPath, filepath.Join(paths...))
}

// newTestAppcast creates a new Appcast instance for testing purposes and
// returns its pointer. By default the content is []byte("test"). However, own
// content can be provided as an argument.
func newTestAppcast(content ...interface{}) *Appcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	s := new(appcaster.Source)
	s.SetContent(resultContent)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(appcaster.Provider(0))

	a := new(Appcast)
	a.SetSource(s)

	return a
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
		"prerelease.xml": {
			"2.0.0-beta": {"Fri, 13 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/2.0.0-beta/app_2.0.0-beta.dmg/download"},
			"1.1.0":      {"Thu, 12 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.1.0/app_1.1.0.dmg/download"},
			"1.0.1":      {"Wed, 11 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.1/app_1.0.1.dmg/download"},
			"1.0.0":      {"Tue, 10 May 2016 12:00:00 UTC", "https://sourceforge.net/projects/example/files/app/1.0.0/app_1.0.0.dmg/download"},
		},
	}

	errorTestCases := map[string]string{
		"invalid_tag.xml":     "XML syntax error on line 21: element <content> closed by </item>",
		"invalid_version.xml": "no version in the #2 release",
	}

	// test (successful)
	for path, releases := range testCases {
		// preparations
		a := newTestAppcast(getTestdata("unmarshal", path))

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())
		assert.Empty(t, a.Releases())

		p, err := a.Unmarshal()

		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, p)
		assert.IsType(t, &Appcast{}, a.Source().Appcast())

		assert.Len(t, releases, a.Releases().Len())
		for _, release := range a.Releases().Filtered() {
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
		a := newTestAppcast(getTestdata("unmarshal", path))

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())

		p, err := a.Unmarshal()

		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, p)
		assert.IsType(t, &Appcast{}, a.Source().Appcast())
	}

	// test (error) [no source]
	a := new(Appcast)

	p, err := a.Unmarshal()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.Source())
}
