package appcast

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestNewRelease(t *testing.T) {
	// test (successful)
	r, err := NewRelease("1.0.0", "1000")
	assert.Nil(t, err)
	assert.IsType(t, Release{}, *r)
	assert.Equal(t, "1.0.0", r.Version.String())
	assert.Equal(t, "1000", r.Build)

	// test (error)
	r, err = NewRelease("invalid", "1000")
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestSetVersion(t *testing.T) {
	// test (successful)
	r := new(Release)
	assert.Nil(t, r.Version)
	err := r.SetVersion("1.0.0")
	assert.Nil(t, err)
	assert.IsType(t, version.Version{}, *r.Version)
	assert.Equal(t, "1.0.0", r.Version.String())

	// test (error)
	r = new(Release)
	assert.Nil(t, r.Version)
	err = r.SetVersion("invalid")
	assert.Error(t, err)
	assert.Nil(t, r.Version)
}
