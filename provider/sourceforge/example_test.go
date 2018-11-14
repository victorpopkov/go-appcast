package sourceforge_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/provider/sourceforge"
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
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/wesnoth/rss", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, err := source.NewRemote("https://sourceforge.net/projects/wesnoth/rss")
	if err != nil {
		panic(err)
	}

	a := sourceforge.New(src)

	err = a.LoadSource()
	if err != nil {
		panic(err)
	}

	a.Unmarshal()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	fmt.Print("Filtering:\n\n")
	fmt.Printf("%12s %d total\n", "Before:", a.Releases().Len())

	// apply some filters
	a.Releases().FilterByMediaType("application/x-bzip2")
	a.Releases().FilterByTitle("wesnoth-1.14", true)
	a.Releases().FilterByUrl("dmg")
	defer a.Releases().ResetFilters() // reset

	fmt.Printf("%12s %d total\n\n", "After:", a.Releases().Len())

	if a.Releases().Len() > 0 {
		r := a.Releases().First()
		fmt.Print("First release details:\n\n")
		fmt.Printf("%12s %s\n", "Version:", r.Version())
		fmt.Printf("%12s %v\n", "Pre-release:", r.IsPreRelease())
		fmt.Printf("%12s %s\n", "Title:", r.Title())
		fmt.Printf("%12s %v\n\n", "Published:", r.PublishedDateTime())

		d := r.Downloads()[0]
		fmt.Printf("%12s %d total\n\n", "Downloads:", len(r.Downloads()))
		fmt.Printf("%12s %s\n", "URL:", d.Url())
		fmt.Printf("%12s %s\n", "Type:", d.Filetype())
		fmt.Printf("%12s %d\n", "Length:", d.Length())
	}

	// Output:
	// Type:     *sourceforge.Appcast
	// Checksum: 880cf7f2f6aa0aa0d859f1fc06e4fbcbba3d4de15fa9736bf73c07accb93ce36
	// Releases: 95 total
	//
	// Filtering:
	//
	//      Before: 95 total
	//       After: 10 total
	//
	// First release details:
	//
	//     Version: 1.13.14
	// Pre-release: false
	//       Title: /wesnoth/wesnoth-1.13.14/Wesnoth_1.13.14.dmg
	//   Published: Sun, 15 Apr 2018 08:45:18 UTC
	//
	//   Downloads: 1 total
	//
	//         URL: https://sourceforge.net/projects/wesnoth/files/wesnoth/wesnoth-1.13.14/Wesnoth_1.13.14.dmg/download
	//        Type: application/x-bzip2; charset=binary
	//      Length: 439082409
}
