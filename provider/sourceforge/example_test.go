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
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/filezilla/rss", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, err := source.NewRemote("https://sourceforge.net/projects/filezilla/rss")
	if err != nil {
		panic(err)
	}

	a := sourceforge.New(src)

	err = a.LoadSource()
	if err != nil {
		panic(err)
	}

	a.Unmarshal()

	// apply some filters
	a.Releases().FilterByMediaType("application/x-bzip2")
	a.Releases().FilterByTitle("FileZilla_Client_Unstable", true)
	a.Releases().FilterByUrl("macosx")
	defer a.Releases().ResetFilters() // reset

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

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

	// Output:
	// Type:     *sourceforge.Appcast
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Releases: 5 total
	//
	// First release details:
	//
	//     Version: 3.25.2
	// Pre-release: false
	//       Title: /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2
	//   Published: Sun, 30 Apr 2017 12:07:25 UTC
	//
	//   Downloads: 1 total
	//
	//         URL: https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2/download
	//        Type: application/x-bzip2; charset=binary
	//      Length: 8453714
}
