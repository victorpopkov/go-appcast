package appcast

import (
	"fmt"
	"reflect"

	"github.com/victorpopkov/go-appcast/release"
	"gopkg.in/jarcoal/httpmock.v1"
)

// Demonstrates the "Sparkle RSS Feed" appcast loading.
func Example_sparkleRSSFeedAppcast() {
	// mock the request
	content := getTestdata("sparkle/example.xml")
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")
	a.SortReleasesByVersions(release.DESC)

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
	// Type:     *appcast.SparkleAppcast
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

// Demonstrates the "SourceForge RSS Feed" appcast loading.
func Example_sourceForgeRSSFeedAppcast() {
	// mock the request
	content := getTestdata("sourceforge/example.xml")
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/filezilla/rss", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://sourceforge.net/projects/filezilla/rss")

	// apply some filters
	a.FilterReleasesByMediaType("application/x-bzip2")
	a.FilterReleasesByTitle("FileZilla_Client_Unstable", true)
	a.FilterReleasesByURL("macosx")
	defer a.ResetFilters() // reset

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

	d := r.Downloads()[0]
	fmt.Printf("%12s %d total\n\n", "Downloads:", len(r.Downloads()))
	fmt.Printf("%12s %s\n", "URL:", d.Url())
	fmt.Printf("%12s %s\n", "Type:", d.Filetype())
	fmt.Printf("%12s %d\n", "Length:", d.Length())

	// Output:
	// Type:     *appcast.SourceForgeAppcast
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Provider: SourceForge RSS Feed
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

// Demonstrates the "Github Atom Feed" appcast loading.
func Example_gitHubAtomFeedAppcast() {
	// mock the request
	content := getTestdata("github/example.xml")
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://github.com/atom/atom/releases.atom", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://github.com/atom/atom/releases.atom")

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
	// Type:     *appcast.GitHubAppcast
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

// Demonstrates the RemoteSource usage.
func ExampleRemoteSource() {
	// mock the request
	content := getTestdata("sparkle/example.xml")
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, _ := NewRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")

	a := New(src)
	a.LoadSource()
	a.Unmarshal()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	// Output:
	// Type:     *appcast.SparkleAppcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}

// Demonstrates the LocalSource usage.
func ExampleLocalSource() {
	src := NewLocalSource(getTestdataPath("sparkle/example.xml"))

	a := New(src)
	a.LoadSource()
	a.Unmarshal()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", a.Releases().Len())

	// Output:
	// Type:     *appcast.SparkleAppcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}
