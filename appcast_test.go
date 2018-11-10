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

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/client"
	"github.com/victorpopkov/go-appcast/provider"
	"github.com/victorpopkov/go-appcast/provider/github"
	"github.com/victorpopkov/go-appcast/provider/sourceforge"
	"github.com/victorpopkov/go-appcast/provider/sparkle"
	"github.com/victorpopkov/go-appcast/release"
	"github.com/victorpopkov/go-appcast/source"
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

	s := new(appcaster.Source)
	s.SetContent(resultContent)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(provider.Unknown)

	src := &source.Remote{
		Source: s,
	}

	src.SetRequest(request)
	src.SetUrl(url)

	o := new(appcaster.Output)
	o.SetContent(resultContent)
	o.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(provider.Unknown)

	output := &LocalOutput{
		Output:      o,
		filepath:    "/tmp/test.txt",
		permissions: 0777,
	}

	a := new(Appcast)
	a.SetSource(src)
	a.SetOutput(output)
	a.SetReleases(release.NewReleases([]release.Releaser{&r1, &r2, &r3, &r4}))

	return a
}

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.Source())

	// test (with source)
	a = New(source.NewLocal(getTestdataPath("../provider/sparkle/testdata/unmarshal/default.xml")))
	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.Source())
}

func TestAppcast_LoadFromRemoteSource(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("../provider/sparkle/testdata/unmarshal/default.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// test (successful) [url]
	a := New()
	p, err := a.LoadFromRemoteSource("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &sparkle.Appcast{}, p)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, provider.Sparkle, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &sparkle.Appcast{}, a.Source().Appcast())

	// test (successful) [request]
	a = New()
	r, _ := client.NewRequest("https://example.com/appcast.xml")
	p, err = a.LoadFromRemoteSource(r)
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &sparkle.Appcast{}, p)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, provider.Sparkle, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &sparkle.Appcast{}, a.Source().Appcast())

	// test (error) [invalid url]
	a = New()
	url := "http://192.168.0.%31/"
	p, err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Nil(t, a.Source())

	// test (error) [invalid request]
	a = New()
	p, err = a.LoadFromRemoteSource("invalid")
	assert.Error(t, err)
	assert.EqualError(t, err, "Get invalid: no responder found")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Nil(t, a.Source())

	// test (error) [unmarshalling failure]
	url = "https://example.com/appcast.xml"
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		url,
		httpmock.NewBytesResponder(200, getTestdata("../provider/sparkle/testdata/unmarshal/invalid_version.xml")),
	)

	a = New()
	p, err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.EqualError(t, err, "malformed version: invalid")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.IsType(t, &source.Remote{}, a.Source())
	assert.IsType(t, &sparkle.Appcast{}, a.Source().Appcast())
}

func TestAppcast_LoadFromLocalSource(t *testing.T) {
	// test (successful)
	path := getTestdataPath("../provider/sparkle/testdata/unmarshal/default.xml")
	content := getTestdata("../provider/sparkle/testdata/unmarshal/default.xml")

	source.LocalReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	a := New()
	p, err := a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &sparkle.Appcast{}, p)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, provider.Sparkle, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &sparkle.Appcast{}, a.Source().Appcast())

	// test (error) [reading failure]
	source.LocalReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Nil(t, a.Source())

	// test (error) [unmarshalling failure]
	path = getTestdataPath("../provider/sparkle/testdata/unmarshal/invalid_version.xml")
	content = getTestdata("../provider/sparkle/testdata/unmarshal/invalid_version.xml")

	source.LocalReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.Error(t, err)
	assert.EqualError(t, err, "malformed version: invalid")
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.IsType(t, &source.Local{}, a.Source())
	assert.IsType(t, &sparkle.Appcast{}, a.Source().Appcast())

	source.LocalReadFile = ioutil.ReadFile
}

func TestAppcast_LoadSource(t *testing.T) {
	// preparations
	a := New(source.NewLocal(getTestdataPath("../provider/sparkle/testdata/unmarshal/default.xml")))
	assert.Nil(t, a.Source().Content())

	// test
	a.LoadSource()
	assert.NotNil(t, a.Source().Content())
}

func TestAppcast_Unmarshal(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"../provider/github/testdata/unmarshal/default.xml": {
			"provider": provider.GitHub,
			"appcast":  &github.Appcast{},
			"checksum": "c28ff87daf2c02471fd2c836b7ed3776d927a8febbb6b8961daf64ce332f6185",
			"releases": 4,
		},
		"../provider/sourceforge/testdata/unmarshal/default.xml": {
			"provider": provider.SourceForge,
			"appcast":  &sourceforge.Appcast{},
			"checksum": "d4afcf95e193a46b7decca76786731c015ee0954b276e4c02a37fa2661a6a5d0",
			"releases": 4,
		},
		"../provider/sparkle/testdata/unmarshal/default.xml": {
			"provider": provider.Sparkle,
			"appcast":  &sparkle.Appcast{},
			"checksum": "0cb017e2dfd65e07b54580ca8d4eedbfcf6cef5824bcd9539a64afb72fa9ce8c",
			"releases": 4,
		},
		"unknown.xml": {
			"provider": provider.Unknown,
			"checksum": "c29665078d79a8e67b37b46a51f2a34c6092719833ccddfdda6109fd8f28043c",
			"error":    "releases for the \"Unknown\" provider can't be unmarshaled",
		},
		"../provider/sparkle/testdata/unmarshal/invalid_version.xml": {
			"provider": provider.Sparkle,
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

		assert.Nil(t, a.Source())
		assert.Empty(t, a.Releases())

		src, err := source.NewRemote("https://example.com/appcast.xml")
		a.SetSource(src)
		a.Source().Load()

		assert.Nil(t, err)
		assert.Equal(t, src, a.Source())
		assert.Empty(t, a.Releases())
		assert.Equal(t, data["provider"], a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.Equal(t, data["checksum"], a.Source().Checksum().String())

		p, err := a.Unmarshal()
		p, err = a.UnmarshalReleases()

		if data["error"] == nil {
			// test (successful)
			assert.Nil(t, err)
			assert.IsType(t, &Appcast{}, a)
			assert.IsType(t, data["appcast"], p)
			assert.Equal(t, a.Releases().Len(), data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", path))
			assert.IsType(t, data["appcast"], a.Source().Appcast())
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
		"../provider/github/testdata/unmarshal/default.xml": {
			"error": "uncommenting is not available for the \"GitHub Atom Feed\" provider",
		},
		"../provider/sourceforge/testdata/unmarshal/default.xml": {
			"error": "uncommenting is not available for the \"SourceForge RSS Feed\" provider",
		},
		"../provider/sparkle/testdata/unmarshal/with_comments.xml": {
			"lines": []int{13, 20},
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
		a.Source().GuessProvider()

		err := a.Uncomment()

		if data["error"] == nil {
			// test (successful)
			assert.Nil(t, err)

			for _, commentLine := range data["lines"].([]int) {
				line, _ := getLine(commentLine, a.Source().Content())
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
	assert.Nil(t, a.Source())
}
