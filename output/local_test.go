package output

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/provider"
)

// newTestLocal creates a new Local instance for testing purposes and returns
// its pointer. By default the content is []byte("test"). However, own content
// can be provided as an argument.
func newTestLocal(content ...interface{}) *Local {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	o := new(appcaster.Output)
	o.SetContent(resultContent)
	o.GenerateChecksum(appcaster.SHA256)
	o.SetProvider(provider.Unknown)

	return &Local{
		Output:      o,
		filepath:    "/tmp/test.txt",
		permissions: 0777,
	}
}

func TestNewLocal(t *testing.T) {
	// preparations
	path := "/tmp/go-appcast_TestNewLocalOutput.txt"

	// test (successful)
	o := NewLocal(provider.Sparkle, path, 0777)
	assert.IsType(t, Local{}, *o)
	assert.NotNil(t, o.Output)
	assert.Equal(t, path, o.filepath)
}

func TestLocal_Save(t *testing.T) {
	// test (successful)
	LocalWriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return nil
	}

	o := newTestLocal()
	err := o.Save()
	assert.Nil(t, err)
	assert.Equal(t, provider.Unknown, o.Provider())

	// test (error)
	LocalWriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("error")
	}

	o = newTestLocal()
	err = o.Save()
	assert.NotNil(t, err)
	assert.Equal(t, provider.Unknown, o.Provider())
	assert.Equal(t, []byte("test"), o.Content())

	LocalWriteFile = ioutil.WriteFile
}

func TestLocal_Filepath(t *testing.T) {
	o := newTestLocal()
	assert.Equal(t, o.filepath, o.Filepath())
}

func TestLocal_SetFilepath(t *testing.T) {
	src := newTestLocal()
	src.SetFilepath("")
	assert.Empty(t, src.filepath)
}

func TestLocal_Permissions(t *testing.T) {
	o := newTestLocal()
	assert.Equal(t, o.permissions, o.Permissions())
}

func TestLocal_SetPermissions(t *testing.T) {
	src := newTestLocal()
	src.SetPermissions(0666)
	assert.Equal(t, os.FileMode(0x1b6), src.permissions)
}
