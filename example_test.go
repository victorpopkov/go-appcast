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
	// Release #1: {1.5.10.4 1.5.10.4 Adium 1.5.10.4  [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}] 2017-05-14 05:04:01 -0700 -0700 false}
	// Release #2: {1.5.10 1.5.10 Adium 1.5.10  [{https://adiumx.cachefly.net/Adium_1.5.10.dmg application/octet-stream 24595712}] 0001-01-01 00:00:00 +0000 UTC false}
	// Release #3: {1.4.5 1.4.5 Adium 1.4.5  [{https://adiumx.cachefly.net/Adium_1.4.5.dmg application/octet-stream 23065688}] 0001-01-01 00:00:00 +0000 UTC false}
	// Release #4: {1.3.10 1.3.10 Adium 1.3.10  [{https://adiumx.cachefly.net/Adium_1.3.10.dmg application/octet-stream 22369877}] 0001-01-01 00:00:00 +0000 UTC false}
	// Release #5: {1.0.6 1.0.6 Adium 1.0.6  [{https://adiumx.cachefly.net/Adium_1.0.6.dmg application/octet-stream 13795246}] 0001-01-01 00:00:00 +0000 UTC false}
}
