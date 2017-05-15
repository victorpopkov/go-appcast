package appcast

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

// getTestdata returns a file content as a byte array from provided testdata
// filename. If file not found, prints an error to os.Stdout and exits with exit
// status 1.
func getTestdata(filename string) []byte {
	path := filepath.Join(getWorkingDir(), testdataPath, filename)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// ReadLine reads a provided line number from io.Reader and returns it alongside
// with an error. Error should be "nil", if the line has been retrieved
// successfully.
func readLine(r io.Reader, lineNum int) (line string, err error) {
	var lastLine int

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), nil
		}
	}

	return "", fmt.Errorf("There is no line \"%d\" in specified io.Reader", lineNum)
}

// getLineFromString returns a specified line from the passed string content and
// an error. Error should be "nil", if the line has been retrieved successfully.
func getLineFromString(lineNum int, content string) (line string, err error) {
	r := bytes.NewReader([]byte(content))

	return readLine(r, lineNum)
}

func TestNew(t *testing.T) {
	a := New()
	assert.IsType(t, BaseAppcast{}, *a)
	assert.Equal(t, Unknown, a.Provider)
}

func TestLoadFromURL(t *testing.T) {
	// mock the request
	content := string(getTestdata("sparkle_default.xml"))
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://example.com/appcast.xml", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// test (successful)
	a := New()
	err := a.LoadFromURL("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Content)
	assert.Equal(t, SparkleRSSFeed, a.Provider)
	assert.Empty(t, a.Checksum.Result)

	// test "Invalid URL" error
	a = New()
	err = a.LoadFromURL("http://192.168.0.%31/")
	assert.Error(t, err)
	assert.Equal(t, "parse http://192.168.0.%31/: invalid URL escape \"%31\"", err.Error())
	assert.Equal(t, Unknown, a.Provider)
	assert.Empty(t, a.Checksum.Result)

	// test "Invalid request" error
	a = New()
	err = a.LoadFromURL("invalid")
	assert.Error(t, err)
	assert.Equal(t, "Get invalid: no responder found", err.Error())
	assert.Equal(t, Unknown, a.Provider)
	assert.Empty(t, a.Checksum.Result)
}
