package appcast_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast"
	"github.com/victorpopkov/go-appcast/release"
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

// Demonstrates the "Github Atom Feed" appcast loading.
func Example_gitHub() {
	// mock the request
	content := testdata("github.xml")
	httpmock.ActivateNonDefault(source.DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://github.com/atom/atom/releases.atom", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := appcast.New()

	p, errors := a.LoadFromRemoteSource("https://github.com/atom/atom/releases.atom")
	if p != nil && len(errors) > 0 {
		panic(errors[0])
	}

	fmt.Printf("%-9s %d error(s)\n", "Result:", len(errors))
	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	r := a.Releases().First()
	fmt.Print("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", r.Version())
	fmt.Printf("%12s %v\n", "Pre-release:", r.IsPreRelease())
	fmt.Printf("%12s %s\n", "Title:", r.Title())
	fmt.Printf("%12s %v\n\n", "Published:", r.PublishedDateTime())

	fmt.Printf("%12s %d total\n", "Downloads:", len(r.Downloads()))

	// Output:
	// Result:   0 error(s)
	// Type:     *github.Appcast
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Provider: GitHub Atom Feed
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

// Demonstrates the "SourceForge RSS Feed" appcast loading.
func Example_sourceForge() {
	// mock the request
	content := testdata("sourceforge.xml")
	httpmock.ActivateNonDefault(source.DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/wesnoth/rss", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := appcast.New()

	p, errors := a.LoadFromRemoteSource("https://sourceforge.net/projects/wesnoth/rss")
	if p == nil && len(errors) > 0 {
		panic(errors[0])
	}

	fmt.Print("Errors:\n\n")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Printf("%-10s %s\n", "\nType:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
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
	// Errors:
	//
	// release #51 (no version)
	// release #71 (no version)
	// release #72 (no version)
	// release #73 (no version)
	// release #92 (no version)
	//
	// Type:     *sourceforge.Appcast
	// Checksum: 880cf7f2f6aa0aa0d859f1fc06e4fbcbba3d4de15fa9736bf73c07accb93ce36
	// Provider: SourceForge RSS Feed
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

// Demonstrates the "Sparkle RSS Feed" appcast loading.
func Example_sparkle() {
	// mock the request
	content := testdata("sparkle.xml")
	httpmock.ActivateNonDefault(source.DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := appcast.New()

	p, errors := a.LoadFromRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")
	if p != nil && len(errors) > 0 {
		panic(errors[0])
	}

	a.Releases().SortByVersions(release.DESC)

	fmt.Printf("%-9s %d error(s)\n", "Result:", len(errors))
	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
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
	// Result:   0 error(s)
	// Type:     *sparkle.Appcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
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
