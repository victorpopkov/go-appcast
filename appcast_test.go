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
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/client"
	"github.com/victorpopkov/go-appcast/release"
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
// returns its pointer. By default the content is []byte("content"). However,
// own content can be provided as an argument.
func newTestAppcast(content ...interface{}) *Appcast {
	var resultContent []byte

	d := new(release.Download)
	d.SetUrl("https://example.com/app_2.0.0-beta.dmg")
	d.SetFiletype("application/octet-stream")
	d.SetLength(100000)

	r := new(release.Release)
	r.SetBuild("200")
	r.SetTitle("Release 2.0.0-beta")
	r.SetDescription("Release 2.0.0-beta Description")
	r.SetReleaseNotesLink("https://example.com/changelogs/2.0.0-beta.html")
	r.SetMinimumSystemVersion("10.10")
	r.SetDownloads([]release.Download{*d})
	r.SetIsPreRelease(true)

	// r1
	d1 := d
	r1 := *r
	r1.SetDownloads([]release.Download{*d1})
	r1.SetVersionString("2.0.0-beta")

	t, _ := time.Parse(time.RFC1123Z, "Fri, 13 May 2016 12:00:00 +0200")
	r1.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r2
	d2 := d
	d2.SetUrl("https://example.com/app_1.1.0.dmg")

	r2 := *r
	r2.SetBuild("110")
	r2.SetTitle("Release 1.1.0")
	r2.SetDescription("Release 1.1.0 Description")
	r2.SetReleaseNotesLink("https://example.com/changelogs/1.1.0.html")
	r2.SetMinimumSystemVersion("10.9")
	r2.SetDownloads([]release.Download{*d2})
	r2.SetVersionString("1.1.0")
	r2.SetIsPreRelease(false)

	t, _ = time.Parse(time.RFC1123Z, "Thu, 12 May 2016 12:00:00 +0200")
	r2.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r3
	d3 := d
	d3.SetUrl("https://example.com/app_1.0.1.dmg")

	r3 := *r
	r3.SetBuild("101")
	r3.SetTitle("Release 1.0.1")
	r3.SetDescription("Release 1.0.1 Description")
	r3.SetReleaseNotesLink("https://example.com/changelogs/1.0.1.html")
	r3.SetMinimumSystemVersion("10.9")
	r3.SetDownloads([]release.Download{*d3})
	r3.SetVersionString("1.0.1")
	r3.SetIsPreRelease(false)

	t, _ = time.Parse(time.RFC1123Z, "Wed, 11 May 2016 12:00:00 +0200")
	r3.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r4
	d4 := d
	d4.SetUrl("https://example.com/app_1.0.0.dmg")

	r4 := *r
	r4.SetBuild("100")
	r4.SetTitle("Release 1.0.0")
	r4.SetDescription("Release 1.0.0 Description")
	r4.SetReleaseNotesLink("https://example.com/changelogs/1.0.0.html")
	r4.SetMinimumSystemVersion("10.9")
	r4.SetDownloads([]release.Download{*d3})
	r4.SetVersionString("1.0.0")
	r4.SetIsPreRelease(false)

	t, _ = time.Parse(time.RFC1123Z, "Tue, 10 May 2016 12:00:00 +0200")
	r4.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("content")
	}

	url := "https://example.com/appcast.xml"
	request, _ := client.NewRequest(url)

	s := &Appcast{
		source: &RemoteSource{
			Source: &Source{
				content:  resultContent,
				provider: Unknown,
			},
			request: request,
			url:     url,
		},
		output: &LocalOutput{
			Output: &Output{
				content: resultContent,
				checksum: &Checksum{
					algorithm: SHA256,
					source:    resultContent,
					result:    []byte("test"),
				},
				provider: Unknown,
			},
			filepath:    "/tmp/test.txt",
			permissions: 0777,
		},
		releases: release.NewReleases([]release.Releaser{&r1, &r2, &r3, &r4}),
	}

	return s
}

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.source)

	// test (with source)
	a = New(NewLocalSource(getTestdataPath("sparkle/default.xml")))
	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.source)
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

	// test (successful) [url]
	a := New()
	p, err := a.LoadFromRemoteSource("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleAppcast{}, p)
	assert.NotEmpty(t, a.source.Content())
	assert.Equal(t, Sparkle, a.source.Provider())
	assert.NotNil(t, a.source.Checksum())
	assert.IsType(t, &SparkleAppcast{}, a.source.Appcast())

	// test (successful) [request]
	a = New()
	r, _ := client.NewRequest("https://example.com/appcast.xml")
	p, err = a.LoadFromRemoteSource(r)
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleAppcast{}, p)
	assert.NotEmpty(t, a.source.Content())
	assert.Equal(t, Sparkle, a.source.Provider())
	assert.NotNil(t, a.source.Checksum())
	assert.IsType(t, &SparkleAppcast{}, a.source.Appcast())

	// test (error) [invalid url]
	a = New()
	url := "http://192.168.0.%31/"
	p, err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Nil(t, a.source)

	// test (error) [invalid request]
	a = New()
	p, err = a.LoadFromRemoteSource("invalid")
	assert.Error(t, err)
	assert.EqualError(t, err, "Get invalid: no responder found")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Nil(t, a.source)

	// test (error) [unmarshalling failure]
	url = "https://example.com/appcast.xml"
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		url,
		httpmock.NewBytesResponder(200, getTestdata("sparkle/invalid_version.xml")),
	)

	a = New()
	p, err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.EqualError(t, err, "malformed version: invalid")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.IsType(t, &RemoteSource{}, a.source)
	assert.IsType(t, &SparkleAppcast{}, a.source.Appcast())
}

func TestAppcast_LoadFromLocalSource(t *testing.T) {
	// test (successful)
	path := getTestdataPath("sparkle/default.xml")
	content := getTestdata("sparkle/default.xml")

	localSourceReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	a := New()
	p, err := a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleAppcast{}, p)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.source.Content())
	assert.Equal(t, Sparkle, a.source.Provider())
	assert.NotNil(t, a.source.Checksum())
	assert.IsType(t, &SparkleAppcast{}, a.source.Appcast())

	// test (error) [reading failure]
	localSourceReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Nil(t, a.source)

	// test (error) [unmarshalling failure]
	path = getTestdataPath("sparkle/invalid_version.xml")
	content = getTestdata("sparkle/invalid_version.xml")

	localSourceReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.Error(t, err)
	assert.EqualError(t, err, "malformed version: invalid")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.IsType(t, &LocalSource{}, a.source)
	assert.IsType(t, &SparkleAppcast{}, a.source.Appcast())

	localSourceReadFile = ioutil.ReadFile
}

func TestAppcast_GenerateSourceChecksum(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.source.Checksum())

	// test
	result := a.GenerateSourceChecksum(MD5)
	assert.Equal(t, result.String(), a.source.Checksum().String())
	assert.Equal(t, "9a0364b9e99bb480dd25e1f0284c8555", result.String())
	assert.Equal(t, MD5, a.source.Checksum().Algorithm())
}

func TestAppcast_LoadSource(t *testing.T) {
	// preparations
	a := New(NewLocalSource(getTestdataPath("sparkle/default.xml")))
	assert.Nil(t, a.source.Content())

	// test
	a.LoadSource()
	assert.NotNil(t, a.source.Content())
}

func TestAppcast_Unmarshal(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sparkle/attributes_as_elements.xml": {
			"provider": Sparkle,
			"appcast":  &SparkleAppcast{},
			"checksum": "d59d258ce0b06d4c6216f6589aefb36e2bd37fbd647f175741cc248021e0e8b4",
			"releases": 4,
		},
		"sourceforge/default.xml": {
			"provider": SourceForge,
			"appcast":  &SourceForgeAppcast{},
			"checksum": "d4afcf95e193a46b7decca76786731c015ee0954b276e4c02a37fa2661a6a5d0",
			"releases": 4,
		},
		"github/default.xml": {
			"provider": GitHub,
			"appcast":  &GitHubAppcast{},
			"checksum": "c28ff87daf2c02471fd2c836b7ed3776d927a8febbb6b8961daf64ce332f6185",
			"releases": 4,
		},
		"unknown.xml": {
			"provider": Unknown,
			"checksum": "c29665078d79a8e67b37b46a51f2a34c6092719833ccddfdda6109fd8f28043c",
			"error":    "releases for the \"Unknown\" provider can't be unmarshaled",
		},
		"sparkle/invalid_version.xml": {
			"provider": Sparkle,
			"checksum": "65d754f5bd04cfad33d415a3605297069127e14705c14b8127a626935229b198",
			"error":    "malformed version: invalid",
		},
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test
	for path, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(path)),
		)

		// preparations
		a := New()

		assert.Nil(t, a.source)
		assert.Empty(t, a.releases)

		src, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(src)
		a.source.Load()

		assert.Nil(t, err)
		assert.Equal(t, src, a.source)
		assert.Empty(t, a.releases)
		assert.Equal(t, data["provider"], a.source.Provider())
		assert.NotEmpty(t, a.source.Content())
		assert.Equal(t, data["checksum"], a.source.Checksum().String())

		p, err := a.Unmarshal()
		p, err = a.UnmarshalReleases()

		if data["error"] == nil {
			// test (successful)
			assert.Nil(t, err)
			assert.IsType(t, &Appcast{}, a)
			assert.IsType(t, data["appcast"], p)
			assert.Equal(t, a.releases.Len(), data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", path))
			assert.IsType(t, data["appcast"], a.source.Appcast())
		} else {
			// test (error)
			assert.Error(t, err)
			assert.EqualError(t, err, data["error"].(string))
			assert.IsType(t, &Appcast{}, a)
			assert.Nil(t, p)
		}
	}
}

func TestAppcast_Uncomment(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sparkle/with_comments.xml": {
			"lines": []int{13, 20},
		},
		"sourceforge/default.xml": {
			"error": "uncommenting is not available for the \"SourceForge RSS Feed\" provider",
		},
		"github/default.xml": {
			"error": "uncommenting is not available for the \"GitHub Atom Feed\" provider",
		},
		"unknown.xml": {
			"error": "uncommenting is not available for the \"Unknown\" provider",
		},
	}

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// test
	for path, data := range testCases {
		// preparations
		a := newTestAppcast(getTestdata(path))
		a.source.GuessProvider()

		err := a.Uncomment()

		if data["error"] == nil {
			// test (successful)
			assert.Nil(t, err)

			for _, commentLine := range data["lines"].([]int) {
				line, _ := getLine(commentLine, a.source.Content())
				check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
				assert.False(t, check)
			}
		} else {
			// test (error)
			assert.Error(t, err)
			assert.EqualError(t, err, data["error"].(string))
		}
	}

	// test (error) [no source]
	a := new(Appcast)

	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, a.source)
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
		p, err := a.Unmarshal()
		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.IsType(t, &SparkleAppcast{}, p)

		// test (ASC)
		a.SortReleasesByVersions(release.ASC)
		assert.Equal(t, "1.0.0", a.releases.First().Version().String())

		// test (DESC)
		a.SortReleasesByVersions(release.DESC)
		assert.Equal(t, "2.0.0", a.releases.First().Version().String())
	}
}

func TestAppcast_Filters(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// test (Appcast.FilterReleasesByTitle)
	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByTitle("Release 1.0")
	assert.Equal(t, 2, a.releases.Len())
	a.FilterReleasesByTitle("Release 1.0.0", true)
	assert.Equal(t, 1, a.releases.Len())
	assert.Equal(t, "Release 1.0.1", a.releases.First().Title())
	a.ResetFilters()

	// test (Appcast.FilterReleasesByMediaType)
	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByMediaType("application/octet-stream")
	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByMediaType("test", true)
	assert.Equal(t, 4, a.releases.Len())
	a.ResetFilters()

	// test (Appcast.FilterReleasesByURL)
	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByURL(`app_1.*dmg$`)
	assert.Equal(t, 3, a.releases.Len())
	a.FilterReleasesByURL(`app_1.0.*dmg$`, true)
	assert.Equal(t, 1, a.releases.Len())
	a.ResetFilters()

	// test (Appcast.FilterReleasesByPrerelease)
	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByPrerelease()
	assert.Equal(t, 1, a.releases.Len())
	a.ResetFilters()

	assert.Equal(t, 4, a.releases.Len())
	a.FilterReleasesByPrerelease(true)
	assert.Equal(t, 3, a.releases.Len())
	a.ResetFilters()
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
	assert.Equal(t, a.source, a.source)
}

func TestAppcast_SetSource(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.source)

	// test
	a.SetSource(nil)
	assert.Nil(t, a.source)
}

func TestAppcast_Output(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.output, a.Output())
}

func TestAppcast_SetOutput(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.output)

	// test
	a.SetOutput(nil)
	assert.Nil(t, a.output)
}

func TestAppcast_Releases(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.releases, a.Releases())
}

func TestAppcast_SetReleases(t *testing.T) {
	a := newTestAppcast()
	a.SetReleases(nil)
	assert.Nil(t, a.releases)
}

func TestAppcast_FirstRelease(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.releases.First(), a.FirstRelease())
}
