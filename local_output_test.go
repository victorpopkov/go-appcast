package appcast

import (
	"fmt"
	"github.com/victorpopkov/go-appcast/appcaster"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestLocalOutput creates a new LocalOutput instance for testing purposes
// and returns its pointer. By default the content is []byte("test"). However,
// own content can be provided as an argument.
func newTestLocalOutput(content ...interface{}) *LocalOutput {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	o := new(appcaster.Output)
	o.SetContent(resultContent)
	o.GenerateChecksum(appcaster.SHA256)
	o.SetProvider(Unknown)

	return &LocalOutput{
		Output:      o,
		filepath:    "/tmp/test.txt",
		permissions: 0777,
	}
}

func TestNewLocalOutput(t *testing.T) {
	// preparations
	path := "/tmp/go-appcast_TestNewLocalOutput.txt"

	// test (successful)
	o := NewLocalOutput(Sparkle, path, 0777)
	assert.IsType(t, LocalOutput{}, *o)
	assert.NotNil(t, o.Output)
	assert.Equal(t, path, o.filepath)
}

func TestLocalOutput_Save(t *testing.T) {
	// preparations
	content := getTestdata("sparkle/default.xml")

	// test (successful)
	localOutputWriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return nil
	}

	o := newTestLocalOutput(content)
	err := o.Save()
	assert.Nil(t, err)
	assert.Equal(t, Unknown, o.Provider())

	// test (error)
	localOutputWriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("error")
	}

	o = newTestLocalOutput()
	err = o.Save()
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, o.Provider())
	assert.Equal(t, []byte("test"), o.Content())

	localOutputWriteFile = ioutil.WriteFile
}

func TestLocalOutput_Filepath(t *testing.T) {
	o := newTestLocalOutput()
	assert.Equal(t, o.filepath, o.Filepath())
}

func TestLocalOutput_Permissions(t *testing.T) {
	o := newTestLocalOutput()
	assert.Equal(t, o.permissions, o.Permissions())
}
