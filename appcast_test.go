package appcast

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
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

// ReadLine reads a provided line number from io.Reader and returns it alongside
// with an error.
func readLine(r io.Reader, lineNum int) (line string, err error) {
	var lastLine int

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), nil
		}
	}

	return "", fmt.Errorf("there is no line \"%d\" in specified io.Reader", lineNum)
}

// getLine returns a specified line from the passed content.
func getLine(lineNum int, content []byte) (line string, err error) {
	return readLine(bytes.NewReader(content), lineNum)
}

// getLineFromString returns a specified line from the passed string content.
func getLineFromString(lineNum int, content string) (line string, err error) {
	return getLine(lineNum, []byte(content))
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

	url := "https://example.com/appcast.xml"
	r, _ := NewRequest(url)

	s := &Appcast{
		source: &RemoteSource{
			Source: &Source{
				content:  resultContent,
				provider: Unknown,
			},
			request: r,
			url:     url,
		},
	}

	return s
}

func TestNew(t *testing.T) {
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.source)
}

func TestAppcast_LoadFromRemoteSource(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("sparkle/default.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// test (successful) [URL]
	a := New()
	err := a.LoadFromRemoteSource("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test (successful) [Request]
	a = New()
	r, _ := NewRequest("https://example.com/appcast.xml")
	err = a.LoadFromRemoteSource(r)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test "Invalid URL" error
	a = New()
	url := "http://192.168.0.%31/"
	err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
	assert.Nil(t, a.Source())

	// test "Invalid request" error
	a = New()
	err = a.LoadFromRemoteSource("invalid")
	assert.Error(t, err)
	assert.EqualError(t, err, "Get invalid: no responder found")
	assert.Nil(t, a.Source())
}

func TestAppcast_LoadFromURL(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("sparkle/default.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// test (successful) [URL]
	a := New()
	err := a.LoadFromURL("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test (successful) [Request]
	a = New()
	r, _ := NewRequest("https://example.com/appcast.xml")
	err = a.LoadFromURL(r)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test "Invalid URL" error
	a = New()
	url := "http://192.168.0.%31/"
	err = a.LoadFromURL(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
	assert.Nil(t, a.Source())

	// test "Invalid request" error
	a = New()
	err = a.LoadFromURL("invalid")
	assert.Error(t, err)
	assert.EqualError(t, err, "Get invalid: no responder found")
	assert.Nil(t, a.Source())
}

func TestAppcast_LoadFromLocalSource(t *testing.T) {
	// test (successful)
	a := New()
	err := a.LoadFromLocalSource(filepath.Join(getWorkingDir(), testdataPath, "sparkle/default.xml"))
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test (error)
	a = New()
	err = a.LoadFromLocalSource("unexisting_file.xml")
	assert.Error(t, err)
	assert.EqualError(t, err, "open unexisting_file.xml: no such file or directory")
	assert.Nil(t, a.Source())
}

func TestAppcast_LoadFromFile(t *testing.T) {
	// test (successful)
	a := New()
	err := a.LoadFromFile(filepath.Join(getWorkingDir(), testdataPath, "sparkle/default.xml"))
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())

	// test (error)
	a = New()
	err = a.LoadFromFile("unexisting_file.xml")
	assert.Error(t, err)
	assert.EqualError(t, err, "open unexisting_file.xml: no such file or directory")
	assert.Nil(t, a.Source())
}

func TestAppcast_GenerateSourceChecksum(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast()
	assert.Nil(t, a.Source().Checksum())

	// test
	result := a.GenerateSourceChecksum(MD5)
	assert.Equal(t, result.String(), a.Source().Checksum().String())
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", result.String())
	assert.Equal(t, MD5, a.Source().Checksum().Algorithm())
}

func TestAppcast_GenerateChecksum(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast()
	assert.Nil(t, a.Source().Checksum())

	// test
	result := a.GenerateChecksum(MD5)
	assert.Equal(t, result.String(), a.Source().Checksum().String())
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", result.String())
	assert.Equal(t, MD5, a.Source().Checksum().Algorithm())
}

func TestAppcast_Uncomment_Unknown(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// test
	err := a.Uncomment()
	assert.EqualError(t, err, "uncommenting is not available for the \"Unknown\" provider")
	a.SetSource(nil)
	err = a.Uncomment()
	assert.EqualError(t, err, "no source")
}

func TestAppcast_Uncomment_SparkleRSSFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("sparkle/with_comments.xml"))
	a.source.SetProvider(SparkleRSSFeed)

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// test
	err := a.Uncomment()
	assert.Nil(t, err)
	for _, commentLine := range []int{13, 20} {
		line, _ := getLine(commentLine, a.Source().Content())
		check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
		assert.False(t, check)
	}
}

func TestAppcast_Uncomment_SourceForgeRSSFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("sourceforge/default.xml"))
	a.source.SetProvider(SourceForgeRSSFeed)

	// test
	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "uncommenting is not available for the \"SourceForge RSS Feed\" provider")
}

func TestAppcast_Uncomment_GitHubAtomFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("github/default.xml"))
	a.source.SetProvider(GitHubAtomFeed)

	// test
	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "uncommenting is not available for the \"GitHub Atom Feed\" provider")
}

func TestAppcast_UnmarshalReleases_Unknown(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// provider "Unknown"
	err := a.UnmarshalReleases()
	assert.Error(t, err)
	assert.EqualError(t, err, "releases can't be extracted from the \"Unknown\" provider")
}

func TestAppcast_UnmarshalReleases_SparkleRSSFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sparkle/attributes_as_elements.xml": {
			"checksum": "8c42d7835109ff61fe85bba66a44689773e73e0d773feba699bceecefaf09359",
			"releases": 4,
		},
		"sparkle/default_asc.xml": {
			"checksum": "9f94a728eab952284b47cc52acfbbb64de71f3d38e5b643d1f3523ef84495d9f",
			"releases": 4,
		},
		"sparkle/default.xml": {
			"checksum": "83c1fd76a250dd50334db793a0db5da7575fc83d292c7c58fd9d31d5bcef6566",
			"releases": 4,
		},
		"sparkle/incorrect_namespace.xml": {
			"checksum": "2e66ef346c49a8472bf8bf26e6e778c5b4d494723223c84c35d9f272a7792430",
			"releases": 4,
		},
		"sparkle/invalid_pubdate.xml": {
			"checksum": "e0273ccbce5a6fb6a5fe31b5edffb8173d88afa308566cf9b4373f3fed909705",
			"releases": 4,
		},
		// "sparkle/multiple_enclosure.xml": {
		// 	"checksum": "48fc8531b253c5d3ed83abfe040edeeafb327d103acbbacf12c2288769dc80b9",
		// 	"releases": 4,
		// },
		"sparkle/no_releases.xml": {
			"checksum": "befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"releases": 0,
		},
		"sparkle/only_version.xml": {
			"checksum": "5c3e7cf62383d4c0e10e5ec0f7afd1a5e328137101e8b6bade050812e4e7451f",
			"releases": 4,
		},
		"sparkle/prerelease.xml": {
			"checksum": "56f95889fe5ddabd847adfe995304fd78dbeeefe47354c2e1c8bde0f003ecf5c",
			"releases": 4,
		},
		"sparkle/single.xml": {
			"checksum": "ac649bebe55f84d85767072e3a1122778a04e03f56b78226bd57ab50ce9f9306",
			"releases": 1,
		},
		"sparkle/without_namespaces.xml": {
			"checksum": "ee2d28f74e7d557bd7259c0f24a261658a9f27a710308a5c539ab761dae487c1",
			"releases": 4,
		},
	}

	errorTestCases := map[string]string{
		"sparkle/invalid_version.xml": "Malformed version: invalid",
		"sparkle/with_comments.xml":   "version is required, but it's not specified in release #1",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		s, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(s)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		err = a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
	}
}

func TestAppcast_UnmarshalReleases_SourceForgeRSSFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sourceforge/default.xml": {
			"checksum": "c15a5e4755b424b20e3e7138c36045893aec70f9569acd5946796199c6f79596",
			"releases": 4,
		},
		"sourceforge/empty.xml": {
			"checksum": "12bbf7be638d5cf251c320aacd68c90acef450e3a9a22cc6cbfa29ffa4ee7f6a",
			"releases": 0,
		},
		"sourceforge/invalid_pubdate.xml": {
			"checksum": "de0f431e001f7aded7fe01c3aec7412e39898d3f97acf809765fc7e2752ffc2c",
			"releases": 4,
		},
		"sourceforge/single.xml": {
			"checksum": "5f3df25c0979faae5b5abef266f5929f4ac6aeb4df74e054461f93e0dbc51183",
			"releases": 1,
		},
	}

	errorTestCases := map[string]string{
		"sourceforge/invalid_version.xml": "version is required, but it's not specified in release #2",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		s, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(s)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, SourceForgeRSSFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		err = a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
	}
}

func TestAppcast_UnmarshalReleases_GitHubAtomFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"github/default.xml": {
			"checksum": "c28ff87daf2c02471fd2c836b7ed3776d927a8febbb6b8961daf64ce332f6185",
			"releases": 4,
		},
		"github/invalid_pubdate.xml": {
			"checksum": "52f87bba760a4e5f8ee418cdbc3806853d79ad10d3f961e5c54d1f5abf09b24b",
			"releases": 4,
		},
	}

	errorTestCases := map[string]string{
		"github/invalid_version.xml": "Malformed version: invalid",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		s, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(s)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, GitHubAtomFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		err = a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.EqualError(t, err, errorMsg)
	}
}

func TestAppcast_ExtractReleases(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// provider "Unknown"
	err := a.ExtractReleases()
	assert.Error(t, err)
	assert.EqualError(t, err, "releases can't be extracted from the \"Unknown\" provider")
}

func TestAppcast_SortReleasesByVersions(t *testing.T) {
	testCases := []string{
		"sparkle/attributes_as_elements.xml",
		"sparkle/default_asc.xml",
		"sparkle/default.xml",
		"sparkle/incorrect_namespace.xml",
		// "sparkle/multiple_enclosure.xml",
		"sparkle/without_namespaces.xml",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, filename := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")
		err := a.UnmarshalReleases()
		assert.Nil(t, err)

		// test (ASC)
		a.SortReleasesByVersions(ASC)
		assert.Equal(t, "1.0.0", a.releases[0].Version.String())

		// test (DESC)
		a.SortReleasesByVersions(DESC)
		assert.Equal(t, "2.0.0", a.releases[0].Version.String())
	}
}

func TestAppcast_Filters(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("sparkle/prerelease.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// preparations
	a := New()
	a.LoadFromRemoteSource("https://example.com/appcast.xml")
	a.UnmarshalReleases()

	// Appcast.FilterReleasesByTitle
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByTitle("Release 1.0")
	assert.Len(t, a.releases, 2)
	a.FilterReleasesByTitle("Release 1.0.0", true)
	assert.Len(t, a.releases, 1)
	assert.Equal(t, "Release 1.0.1", a.releases[0].Title)
	a.ResetFilters()

	// Appcast.FilterReleasesByMediaType
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByMediaType("application/octet-stream")
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByMediaType("test", true)
	assert.Len(t, a.releases, 4)
	a.ResetFilters()

	// Appcast.FilterReleasesByURL
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByURL(`app_1.*dmg$`)
	assert.Len(t, a.releases, 3)
	a.FilterReleasesByURL(`app_1.0.*dmg$`, true)
	assert.Len(t, a.releases, 1)
	a.ResetFilters()

	// Appcast.FilterReleasesByPrerelease
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByPrerelease()
	assert.Len(t, a.releases, 1)
	a.ResetFilters()

	assert.Len(t, a.releases, 4)
	a.FilterReleasesByPrerelease(true)
	assert.Len(t, a.releases, 3)
	a.ResetFilters()
}

func TestAppcast_GetReleasesLength(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("sparkle/default.xml"))
	a.UnmarshalReleases()

	// test
	assert.Len(t, a.releases, a.GetReleasesLength())
}

func TestAppcast_GetFirstRelease(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast(getTestdata("sparkle/default.xml"))
	a.UnmarshalReleases()

	// test
	assert.Equal(t, a.releases[0].GetVersionString(), a.GetFirstRelease().GetVersionString())
}

func TestExtractSemanticVersions(t *testing.T) {
	testCases := map[string][]string{
		// single
		"Version 1":           nil,
		"Version 1.0":         nil,
		"Version 1.0.2":       {"1.0.2"},
		"Version 1.0.2-alpha": {"1.0.2-alpha"},
		"Version 1.0.2-beta":  {"1.0.2-beta"},
		"Version 1.0.2-dev":   {"1.0.2-dev"},
		"Version 1.0.2-rc1":   {"1.0.2-rc1"},

		// multiples
		"First is v1.0.1, second is v1.0.2, third is v1.0.3": {"1.0.1", "1.0.2", "1.0.3"},
	}

	// test
	for data, versions := range testCases {
		actual, err := ExtractSemanticVersions(data)
		if versions == nil {
			assert.Error(t, err)
			assert.EqualError(t, err, "no semantic versions found")
		} else {
			assert.Nil(t, err)
			assert.Equal(t, versions, actual)
		}
	}
}

func TestAppcast_Source(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.source, a.Source())
}

func TestAppcast_SetSource(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.source)

	// test
	a.SetSource(nil)
	assert.Nil(t, a.source)
}

func TestAppcast_Releases(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.releases, a.Releases())
}

func TestAppcast_SetReleases(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.originalReleases)

	// test
	a.SetReleases([]Release{{}})
	assert.Len(t, a.releases, 1)
}

func TestAppcast_OriginalReleases(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.originalReleases, a.OriginalReleases())
}

func TestAppcast_SetOriginalReleases(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.originalReleases)

	// test
	a.SetOriginalReleases([]Release{{}})
	assert.Len(t, a.originalReleases, 1)
}

func TestAppcast_GetChecksum(t *testing.T) {
	a := newTestAppcast()
	a.GenerateSourceChecksum(SHA256)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", a.GetChecksum().String())
}

func TestAppcast_GetProvider(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, Unknown, a.GetProvider())
}
