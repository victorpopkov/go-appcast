package appcast

import (
	"fmt"

	"gopkg.in/jarcoal/httpmock.v1"
)

// ExampleSparkleRSSFeed demonstrates the loading and parsing of the "Sparkle
// RSS Feed" appcast.
func Example_sparkleRSSFeed() {
	// mock the request
	content := string(getTestdata("example_sparkle.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromURL("https://www.adium.im/sparkle/appcast-release.xml")
	a.GenerateChecksum(Sha256)
	a.ExtractReleases()
	a.SortReleasesByVersions(DESC)

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.Provider)

	for i, release := range a.Releases {
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
	content := string(getTestdata("example_sourceforge.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://sourceforge.net/projects/filezilla/rss", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromURL("https://sourceforge.net/projects/filezilla/rss")
	a.GenerateChecksum(Sha256)
	a.ExtractReleases()

	// apply some filters
	a.FilterReleasesByMediaType("application/x-bzip2")
	a.FilterReleasesByTitle("FileZilla_Client_Unstable", true)
	a.FilterReleasesByURL("macosx")
	defer a.ResetFilters() // reset

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.Provider)

	for i, release := range a.Releases {
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
	content := string(getTestdata("example_github.xml"))
	httpmock.ActivateNonDefault(DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://github.com/atom/atom/releases.atom", httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	a := New()
	a.LoadFromURL("https://github.com/atom/atom/releases.atom")
	a.GenerateChecksum(Sha256)
	a.ExtractReleases()

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.Provider)

	for i, release := range a.Releases {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release.Version, release.Title, release.PublishedDateTime, release.IsPrerelease)
	}

	fmt.Println("Release #1 description:", a.Releases[0].Description)

	// Output:
	// Checksum: 14dd5fa8a4f880ae7c441e2fc940516e9d50b23fa110277d7696a35380cdb102
	// Provider: GitHub Atom Feed
	// Release #1: 1.18.0-beta2 1.18.0-beta2 2017-05-25 20:39:10 +0000 UTC true
	// Release #2: 1.17.2 1.17.2 2017-05-25 20:38:59 +0000 UTC false
	// Release #3: 1.18.0-beta1 1.18.0-beta1 2017-05-24 22:32:31 +0000 UTC true
	// Release #4: 1.17.1 1.17.1 2017-05-24 22:40:15 +0000 UTC false
	// Release #5: 1.18.0-beta0 1.18.0-beta0 2017-05-16 16:41:06 +0000 UTC true
	// Release #6: 1.17.0 1.17.0 2017-05-16 16:41:24 +0000 UTC false
	// Release #7: 1.17.0-beta5 1.17.0-beta5 2017-05-05 15:53:16 +0000 UTC true
	// Release #8: 1.17.0-beta4 1.17.0-beta4 2017-04-26 20:51:47 +0000 UTC true
	// Release #9: 1.17.0-beta3 1.17.0-beta3 2017-04-18 20:20:32 +0000 UTC true
	// Release #10: 1.17.0-beta2 1.17.0-beta2 2017-04-14 16:28:23 +0000 UTC true
	// Release #1 description: <p>Fixes <a href="https://github.com/atom/find-and-replace/issues/901">Find in current buffer' returns no results - atom/find-and-replace#901</a></p>
}
