package appcast

import (
	"fmt"

	"gopkg.in/jarcoal/httpmock.v1"
)

// ExampleSparkleRSSFeed demonstrates the loading and parsing of the "Sparkle
// RSS Feed" appcast.
func Example_sparkleRSSFeed() {
	// mock the request
	content := string(getTestdata("sparkle/example.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")
	a.UnmarshalReleases()
	a.SortReleasesByVersions(DESC)

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())

	for i, release := range a.Releases() {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release)
	}

	// Output:
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Release #1: {1.5.10.4 1.5.10.4 Adium 1.5.10.4  [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}] 2017-05-14 12:04:01 +0000 UTC false}
	// Release #2: {1.5.10 1.5.10 Adium 1.5.10  [{https://adiumx.cachefly.net/Adium_1.5.10.dmg application/octet-stream 24595712}] 2014-05-19 21:25:14 +0000 UTC false}
	// Release #3: {1.4.5 1.4.5 Adium 1.4.5  [{https://adiumx.cachefly.net/Adium_1.4.5.dmg application/octet-stream 23065688}] 2012-03-20 20:30:00 +0000 UTC false}
	// Release #4: {1.3.10 1.3.10 Adium 1.3.10  [{https://adiumx.cachefly.net/Adium_1.3.10.dmg application/octet-stream 22369877}] 2010-01-12 23:30:00 +0000 UTC false}
	// Release #5: {1.0.6 1.0.6 Adium 1.0.6  [{https://adiumx.cachefly.net/Adium_1.0.6.dmg application/octet-stream 13795246}] 2007-08-13 22:12:45 +0000 UTC false}
}

// ExampleSourceForgeRSSFeed demonstrates the loading and parsing of the
// "SourceForge RSS Feed" appcast.
func Example_sourceForgeRSSFeed() {
	// mock the request
	content := string(getTestdata("sourceforge/example.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/filezilla/rss", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://sourceforge.net/projects/filezilla/rss")
	a.UnmarshalReleases()

	// apply some filters
	a.FilterReleasesByMediaType("application/x-bzip2")
	a.FilterReleasesByTitle("FileZilla_Client_Unstable", true)
	a.FilterReleasesByURL("macosx")
	defer a.ResetFilters() // reset

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())

	for i, release := range a.Releases() {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release)
	}

	// Output:
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Provider: SourceForge RSS Feed
	// Release #1: {3.25.2  /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8453714}] 2017-04-30 12:07:25 +0000 UTC false}
	// Release #2: {3.25.1  /FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8460741}] 2017-03-20 17:11:09 +0000 UTC false}
	// Release #3: {3.25.0  /FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8461936}] 2017-03-13 14:36:41 +0000 UTC false}
	// Release #4: {3.24.1  /FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2 /FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8764178}] 2017-02-21 22:00:38 +0000 UTC false}
	// Release #5: {3.24.0  /FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2 /FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8765941}] 2017-01-13 20:20:31 +0000 UTC false}
}

// ExampleGitHubAtomFeed demonstrates the loading and parsing of the
// "Github Atom Feed" appcast.
func Example_gitHubAtomFeed() {
	// mock the request
	content := string(getTestdata("github/example.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://github.com/atom/atom/releases.atom", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromRemoteSource("https://github.com/atom/atom/releases.atom")
	a.UnmarshalReleases()

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())

	for i, release := range a.Releases() {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release.Version, release.Title, release.PublishedDateTime, release.IsPrerelease)
	}

	// Output:
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Provider: GitHub Atom Feed
	// Release #1: 1.28.0-beta3 1.28.0-beta3 2018-06-06 17:09:54 +0000 UTC true
	// Release #2: 1.28.0-beta2 1.28.0-beta2 2018-05-31 13:55:54 +0000 UTC true
	// Release #3: 1.27.2 1.27.2 2018-05-31 13:55:49 +0000 UTC false
	// Release #4: 1.28.0-beta1 1.28.0-beta1 2018-05-21 14:46:23 +0000 UTC true
	// Release #5: 1.27.1 1.27.1 2018-05-21 14:46:10 +0000 UTC false
	// Release #6: 1.28.0-beta0 1.28.0-beta0 2018-05-15 17:46:08 +0000 UTC true
	// Release #7: 1.27.0 1.27.0 2018-05-15 17:46:03 +0000 UTC false
	// Release #8: 1.27.0-beta1 1.27.0-beta1 2018-04-26 19:40:51 +0000 UTC true
	// Release #9: 1.26.1 1.26.1 2018-04-26 19:40:40 +0000 UTC false
	// Release #10: 1.26.0 1.26.0 2018-04-18 23:00:10 +0000 UTC false
}
