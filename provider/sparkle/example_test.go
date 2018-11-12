package sparkle_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/provider/sparkle"
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
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, err := source.NewRemote("https://www.adium.im/sparkle/appcast-release.xml")
	if err != nil {
		panic(err)
	}

	a := sparkle.New(src)

	err = a.LoadSource()
	if err != nil {
		panic(err)
	}

	a.Unmarshal()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	r := a.Releases().First()
	fmt.Print("First release details:\n\n")
	fmt.Printf("%23s %s\n", "Version:", r.Version())
	fmt.Printf("%23s %s\n", "Build:", r.Build())
	fmt.Printf("%23s %v\n", "Pre-release:", r.IsPreRelease())
	fmt.Printf("%23s %s\n", "Title:", r.Title())
	fmt.Printf("%23s %v\n", "Published:", r.PublishedDateTime())
	fmt.Printf("%23s %v\n", "Release notes:", r.ReleaseNotesLink())
	fmt.Printf("%23s %v\n\n", "Minimum system version:", r.MinimumSystemVersion())

	d := r.Downloads()[0]
	fmt.Printf("%23s %d total\n\n", "Downloads:", len(r.Downloads()))
	fmt.Printf("%23s %s\n", "URL:", d.Url())
	fmt.Printf("%23s %s\n", "Type:", d.Filetype())
	fmt.Printf("%23s %d\n", "Length:", d.Length())
	fmt.Printf("%23s %s\n", "DSA Signature:", d.DsaSignature())

	// Output:
	// Type:     *sparkle.Appcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Releases: 5 total
	//
	// First release details:
	//
	//                Version: 1.5.10.4
	//                  Build: 1.5.10.4
	//            Pre-release: false
	//                  Title: Adium 1.5.10.4
	//              Published: Sun, 14 May 2017 05:04:01 -0700
	//          Release notes: https://www.adium.im/changelogs/1.5.10.4.html
	// Minimum system version: 10.7.5
	//
	//              Downloads: 1 total
	//
	//                    URL: https://adiumx.cachefly.net/Adium_1.5.10.4.dmg
	//                   Type: application/octet-stream
	//                 Length: 21140435
	//          DSA Signature: MC4CFQCeqQ/MxlFt2H3rQfCPimChDPibCgIVAJhZmHcU8ZHylc7EjvbkVr3ardLp
}
