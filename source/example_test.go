package source_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/jarcoal/httpmock.v1"

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

func Example_local() {
	src := source.NewLocal(testdataPath("example.xml"))

	err := src.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("Type:", reflect.TypeOf(src))
	fmt.Println("Content:", fmt.Sprintf("%d (%s)", len(src.Content()), "length"))
	fmt.Println("Checksum:", fmt.Sprintf("%s (%s)", src.Checksum().String(), src.Checksum().Algorithm()))

	// Output:
	// Type: *source.Local
	// Content: 2048 (length)
	// Checksum: 0cb017e2dfd65e07b54580ca8d4eedbfcf6cef5824bcd9539a64afb72fa9ce8c (SHA256)
}

func Example_remote() {
	// mock the request
	content := testdata("example.xml")
	httpmock.ActivateNonDefault(source.DefaultClient.HTTPClient)
	httpmock.RegisterResponder("GET", "https://www.adium.im/sparkle/appcast-release.xml", httpmock.NewBytesResponder(200, content))
	defer httpmock.DeactivateAndReset()

	// example
	src, err := source.NewRemote("https://www.adium.im/sparkle/appcast-release.xml")
	if err != nil {
		panic(err)
	}

	err = src.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("Type:", reflect.TypeOf(src))
	fmt.Println("Content:", fmt.Sprintf("%d (%s)", len(src.Content()), "length"))
	fmt.Println("Checksum:", fmt.Sprintf("%s (%s)", src.Checksum().String(), src.Checksum().Algorithm()))

	// Output:
	// Type: *source.Remote
	// Content: 2048 (length)
	// Checksum: 0cb017e2dfd65e07b54580ca8d4eedbfcf6cef5824bcd9539a64afb72fa9ce8c (SHA256)
}
