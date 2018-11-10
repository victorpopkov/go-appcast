// Package output holds the supported outputs.
package output

import (
	"io/ioutil"
	"os"

	"github.com/victorpopkov/go-appcast/appcaster"
)

var LocalWriteFile = ioutil.WriteFile

// Localer is the interface that wraps the Local methods.
type Localer interface {
	appcaster.Outputer
	Filepath() string
	SetFilepath(filepath string)
	Permissions() os.FileMode
	SetPermissions(permissions os.FileMode)
}

// Local represents an appcast output to the local file.
type Local struct {
	*appcaster.Output
	filepath    string
	permissions os.FileMode
}

// NewLocal returns a new Local instance pointer with the Local.Output.provider,
// Local.filepath and Local.permissions set.
func NewLocal(provider appcaster.Providerer, path string, perm os.FileMode) *Local {
	o := new(appcaster.Output)
	o.SetProvider(provider)

	return &Local{
		Output:      o,
		filepath:    path,
		permissions: perm,
	}
}

// Save saves an appcast content from the Local.Output.content to the local file
// by using the Local.filepath set earlier.
func (l *Local) Save() error {
	err := LocalWriteFile(l.filepath, l.Content(), l.permissions)
	if err != nil {
		return err
	}

	return nil
}

// Filepath is a Local.filepath getter.
func (l *Local) Filepath() string {
	return l.filepath
}

// SetFilepath is a Local.filepath setter.
func (l *Local) SetFilepath(filepath string) {
	l.filepath = filepath
}

// Permissions is a Local.permissions getter.
func (l *Local) Permissions() os.FileMode {
	return l.permissions
}

// SetPermissions is a Local.permissions setter.
func (l *Local) SetPermissions(permissions os.FileMode) {
	l.permissions = permissions
}
