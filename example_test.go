package appcast

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"reflect"
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
	a.SortReleasesByVersions(DESC)

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Build:", release.Build())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Type:     *appcast.SparkleRSSFeedAppcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
	//
	// First release details:
	//
	//     Version: 1.5.10.4
	//       Build: 1.5.10.4
	//       Title: Adium 1.5.10.4
	//   Downloads: [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}]
	//   Published: Sun, 14 May 2017 05:04:01 -0700
	// Pre-release: false
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
	fmt.Printf("%-9s %d total\n\n", "Releases:", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Type:     *appcast.SourceForgeRSSFeedAppcast
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Provider: SourceForge RSS Feed
	// Releases: 5 total
	//
	// First release details:
	//
	//     Version: 3.25.2
	//       Title: /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2
	//   Downloads: [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8453714}]
	//   Published: Sun, 30 Apr 2017 12:07:25 UTC
	// Pre-release: false
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
	fmt.Printf("%-9s %d total\n\n", "Releases:", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Type:     *appcast.GitHubAtomFeedAppcast
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Provider: GitHub Atom Feed
	// Releases: 10 total
	//
	// First release details:
	//
	//     Version: 1.28.0-beta3
	//       Title: 1.28.0-beta3
	//   Downloads: []
	//   Published: 2018-06-06T20:09:54+03:00
	// Pre-release: true
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
	a.UnmarshalReleases()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", len(a.Releases()))

	// Output:
	// Type:     *appcast.SparkleRSSFeedAppcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}

// Demonstrates the RemoteSource usage.
func ExampleLocalSource() {
	src := NewLocalSource(getTestdataPath("sparkle/example.xml"))

	a := New(src)
	a.LoadSource()
	a.UnmarshalReleases()

	fmt.Printf("%-9s %s\n", "Type:", reflect.TypeOf(a.Source().Appcast()))
	fmt.Printf("%-9s %s\n", "Checksum:", a.Source().Checksum())
	fmt.Printf("%-9s %s\n", "Provider:", a.Source().Provider())
	fmt.Printf("%-9s %d total\n\n", "Releases:", len(a.Releases()))

	// Output:
	// Type:     *appcast.SparkleRSSFeedAppcast
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}
