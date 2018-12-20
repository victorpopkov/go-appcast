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
// returns its pointer. By default the source is LocalSource and points to the
//// "GitHub Atom Feed" default.xml testdata.
func newTestAppcast(paths ...string) *Appcast {
	var content []byte

	if len(paths) > 0 {
		content = testdata(paths...)
	} else {
		content = testdata("unmarshal", "default.xml")
	}

	s := new(appcaster.Source)
	s.SetContent(content)
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
	type testCase struct {
		path     string
		appcast  appcaster.Appcaster
		releases map[string][]string
		errors   []string
	}

	testCases := []testCase{
		{
			path:    "default.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"2016-05-13T12:00:00+02:00"},
				"1.1.0": {"2016-05-12T12:00:00+02:00"},
				"1.0.1": {"2016-05-11T12:00:00+02:00"},
				"1.0.0": {"2016-05-10T12:00:00+02:00"},
			},
		},
		{
			path:    "empty.xml",
			appcast: &Appcast{},
		},
		{
			path:    "invalid_pubdate.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"2016-05-13T12:00:00+02:00"},
				"1.1.0": {""},
				"1.0.1": {"2016-05-11T12:00:00+02:00"},
				"1.0.0": {"2016-05-10T12:00:00+02:00"},
			},
			errors: []string{
				"release #2 (parsing of the published datetime failed)",
			},
		},
		{
			path: "invalid_tag.xml",
			errors: []string{
				"XML syntax error on line 18: element <thumbnail> closed by </entry>",
			},
		},
		{
			path:    "invalid_version.xml",
			appcast: &Appcast{},
			errors: []string{
				"release #2 (malformed version: invalid)",
			},
		},
		{
			path:    "prerelease.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0-beta": {"2016-05-13T12:00:00+02:00"},
				"1.1.0":      {"2016-05-12T12:00:00+02:00"},
				"1.0.1":      {"2016-05-11T12:00:00+02:00"},
				"1.0.0":      {"2016-05-10T12:00:00+02:00"},
			},
		},
	}

	// test
	for _, testCase := range testCases {
		// preparations
		a := newTestAppcast("unmarshal", testCase.path)

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())
		assert.Empty(t, a.Releases())

		appcast, errors := a.Unmarshal()

		if testCase.appcast != nil {
			assert.IsType(t, testCase.appcast, appcast, fmt.Sprintf("%s: appcast type mismatch", testCase.path))
			assert.IsType(t, testCase.appcast, a.Source().Appcast())
		} else {
			assert.Equal(t, testCase.appcast, appcast, fmt.Sprintf("%s: appcast type mismatch", testCase.path))
		}

		if len(testCase.errors) == 0 {
			// successful
			assert.Nil(t, errors, fmt.Sprintf("%s: errors not nil", testCase.path))

			releases := testCase.releases
			assert.Len(t, releases, a.Releases().Len())

			for _, r := range a.Releases().Filtered() {
				v := r.Version().String()
				assert.Equal(t, fmt.Sprintf("%s", v), r.Title())
				assert.NotEmpty(t, r.Description())
				assert.Equal(t, releases[v][0], r.PublishedDateTime().String())
			}
		} else {
			// error (unmarshalling failure)
			assert.Len(t, errors, len(testCase.errors), fmt.Sprintf("%s: errors length mismatch", testCase.path))

			for i, errorMsg := range testCase.errors {
				err := errors[i]
				assert.EqualError(t, err, errorMsg)
			}
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
