package github_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/provider/github"
	"github.com/victorpopkov/go-appcast/source"
)

func testdataPath(paths ...string) string {
	testdataPath := "./testdata/"

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return filepath.Join(pwd, testdataPath, filepath.Join(paths...))
}

func testdata(paths ...string) []byte {
	content, err := ioutil.ReadFile(testdataPath(paths...))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return content
}

func Example() {
	// mock the request
	content := testdata("unmarshal/example.xml")
	httpmock.ActivateNonDefault(source.DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://github.com/atom/atom/releases.atom", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, err := source.NewRemote("https://github.com/atom/atom/releases.atom")
	if err != nil {
		panic(err)
	}

	a := github.New(src)

	err = a.LoadSource()
	if err != nil {
		panic(err)
	}

	p, errors := a.Unmarshal()
	if p == nil && len(errors) > 0 {
		panic(errors[0])
	}

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	r := a.Releases().First()
	fmt.Print("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", r.Version())
	fmt.Printf("%12s %v\n", "Pre-release:", r.IsPreRelease())
	fmt.Printf("%12s %s\n", "Title:", r.Title())
	fmt.Printf("%12s %v\n\n", "Published:", r.PublishedDateTime())

	fmt.Printf("%12s %d total\n", "Downloads:", len(r.Downloads()))

	// Output:
	// Type:     *github.Appcast
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Releases: 10 total
	//
	// First release details:
	//
	//     Version: 1.28.0-beta3
	// Pre-release: true
	//       Title: 1.28.0-beta3
	//   Published: 2018-06-06T20:09:54+03:00
	//
	//   Downloads: 0 total
}
