package appcast

import (
	"io/ioutil"
	"os"
)

var localOutputWriteFile = ioutil.WriteFile

// LocalOutputer is the interface that wraps the LocalOutput methods.
type LocalOutputer interface {
	Outputer
	Filepath() string
	Permissions() os.FileMode
}

// LocalOutput represents an appcast output to the local file.
type LocalOutput struct {
	*Output
	filepath    string
	permissions os.FileMode
}

// NewLocalOutput returns a new LocalOutput instance pointer with the
// LocalOutput.Output.provider, LocalOutput.filepath and LocalOutput.permissions
// set.
func NewLocalOutput(provider Provider, path string, perm os.FileMode) *LocalOutput {
	o := &LocalOutput{
		Output: &Output{
			provider: provider,
		},
		filepath:    path,
		permissions: perm,
	}

	return o
}

// Save saves an appcast content from the LocalOutput.Output.content to the
// local file by using the LocalOutput.filepath set earlier.
func (o *LocalOutput) Save() error {
	err := localOutputWriteFile(o.filepath, o.content, o.permissions)
	if err != nil {
		return err
	}

	return nil
}

// Filepath is a LocalOutput.filepath getter.
func (o *LocalOutput) Filepath() string {
	return o.filepath
}

// Permissions is a LocalOutput.permissions getter.
func (o *LocalOutput) Permissions() os.FileMode {
	return o.permissions
}
