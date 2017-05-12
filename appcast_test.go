package appcast

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var testdataPath = "./testdata/"

// getWorkingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func getWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// getTestdata returns a file content as a byte array from provided testdata
// filename. If file not found, prints an error to os.Stdout and exits with exit
// status 1.
func getTestdata(filename string) []byte {
	content, err := ioutil.ReadFile(filepath.Join(getWorkingDir(), testdataPath, filename))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return content
}
