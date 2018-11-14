package github

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
)

// workingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func workingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// testdata returns a file content as a byte slice from the provided testdata
// paths. If the file is not found, prints an error to os.Stdout and exits with
// exit status 1.
func testdata(paths ...string) []byte {
	path := testdataPath(paths...)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// testdataPath returns a full path for the provided testdata paths.
func testdataPath(paths ...string) string {
	return filepath.Join(workingDir(), "./testdata/", filepath.Join(paths...))
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

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.Source())

	// test (with source)
	src := new(appcaster.Source)
	src.SetContent([]byte("content"))
	src.SetProvider(appcaster.Provider(0))

	a = New(src)
	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.Source())
}

func TestAppcast_Unmarshal(t *testing.T) {
	testCases := map[string]map[string][]string{
		"default.xml": {
			"2.0.0": {"2016-05-13T12:00:00+02:00"},
			"1.1.0": {"2016-05-12T12:00:00+02:00"},
			"1.0.1": {"2016-05-11T12:00:00+02:00"},
			"1.0.0": {"2016-05-10T12:00:00+02:00"},
		},
		"empty.xml": {},
		"invalid_pubdate.xml": {
			"2.0.0": {"2016-05-13T12:00:00+02:00"},
			"1.1.0": {""},
			"1.0.1": {"2016-05-11T12:00:00+02:00"},
			"1.0.0": {"2016-05-10T12:00:00+02:00"},
		},
		"prerelease.xml": {
			"2.0.0-beta": {"2016-05-13T12:00:00+02:00"},
			"1.1.0":      {"2016-05-12T12:00:00+02:00"},
			"1.0.1":      {"2016-05-11T12:00:00+02:00"},
			"1.0.0":      {"2016-05-10T12:00:00+02:00"},
		},
	}

	errorTestCases := map[string][]string{
		"invalid_tag.xml": {
			"XML syntax error on line 18: element <thumbnail> closed by </entry>",
		},
		"invalid_version.xml": {
			"release #2 (malformed version: invalid)",
		},
	}

	// test (successful)
	for path, releases := range testCases {
		// preparations
		a := newTestAppcast(testdata("unmarshal", path))

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())
		assert.Empty(t, a.Releases())

		p, err := a.Unmarshal()

		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, p)
		assert.IsType(t, &Appcast{}, a.Source().Appcast())

		assert.Len(t, releases, a.Releases().Len())
		for _, r := range a.Releases().Filtered() {
			v := r.Version().String()
			assert.Equal(t, fmt.Sprintf("%s", v), r.Title())
			assert.NotEmpty(t, r.Description())
			assert.Equal(t, releases[v][0], r.PublishedDateTime().String())
		}
	}

	// test (error) [unmarshalling failure]
	for path, errorMsgs := range errorTestCases {
		// preparations
		a := newTestAppcast(testdata("unmarshal", path))

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())

		_, errors := a.Unmarshal()

		assert.Len(t, errors, len(errorMsgs))
		for i, errorMsg := range errorMsgs {
			err := errors[i]
			assert.Error(t, err)
			assert.EqualError(t, err, errorMsg)
			assert.IsType(t, &Appcast{}, a.Source().Appcast())
		}
	}

	// test (error) [no source]
	a := new(Appcast)

	p, errors := a.Unmarshal()

	assert.Len(t, errors, 1)
	err := errors[0]

	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.Source())
}
