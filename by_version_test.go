package appcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestByVersionReleases() (result []Release) {
	testReleases := [][]string{
		{"1.0.0", "100"},
		{"1.1.0", "110"},
		{"2.0.0", "200"},
	}

	for _, release := range testReleases {
		r, _ := NewRelease(release[0], release[1])
		result = append(result, *r)
	}

	return result
}

func TestByVersion_Len(t *testing.T) {
	// preparations
	testReleases := createTestByVersionReleases()

	// test
	assert.Equal(t, 3, ByVersion(testReleases).Len())
}

func TestByVersion_Swap(t *testing.T) {
	// preparations
	testReleases := createTestByVersionReleases()

	// test
	assert.Equal(t, "1.1.0", testReleases[1].Version.String())
	assert.Equal(t, "2.0.0", testReleases[2].Version.String())
	ByVersion(testReleases).Swap(1, 2)
	assert.Equal(t, "2.0.0", testReleases[1].Version.String())
	assert.Equal(t, "1.1.0", testReleases[2].Version.String())
}

func TestByVersion_Less(t *testing.T) {
	// preparations
	testReleases := createTestByVersionReleases()

	// test
	assert.True(t, ByVersion(testReleases).Less(0, 1))
	assert.True(t, ByVersion(testReleases).Less(1, 2))
}
