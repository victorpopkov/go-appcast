package release

import (
	"regexp"
	"sort"
)

// Releaseser is the interface that wraps the Releases methods.
//
// TODO: Find a better Releaseser interface name.
type Releaseser interface {
	SortByVersions(s Sort)
	FilterByTitle(regexpStr string, inversed ...interface{})
	FilterByMediaType(regexpStr string, inversed ...interface{})
	FilterByUrl(regexpStr string, inversed ...interface{})
	FilterByPrerelease(inversed ...interface{})
	ResetFilters()
	Len() int
	First() Releaser
	Filtered() []Releaser
	SetFiltered(filtered []Releaser)
	Original() []Releaser
	SetOriginal(original []Releaser)
}

// Releases represents the appcast releases which holds both the filtered and
// the original ones.
type Releases struct {
	// filtered specifies a slice of all application releases. All filtered
	// releases are stored here.
	filtered []Releaser

	// original specifies a slice holding a copy of the Releases.filtered. It is
	// used to restore the Releases.filtered using the Releases.ResetFilters
	// method.
	original []Releaser
}

// Sort holds different supported sorting behaviours.
type Sort int

const (
	// ASC represents the ascending order.
	ASC Sort = iota

	// DESC represents the descending order.
	DESC
)

// NewReleases returns a new Releases instance pointer. Requires []Releaser
// slice to be passed as a parameter which will be set to the Releases.filtered
// and the Releases.original.
func NewReleases(releases []Releaser) *Releases {
	return &Releases{
		filtered: releases,
		original: releases,
	}
}

// SortByVersions sorts the Releases.filtered slice by versions. Can be
// useful if the versions order is inconsistent.
func (r *Releases) SortByVersions(s Sort) {
	if s == ASC {
		sort.Sort(ByVersion(r.filtered))
	} else if s == DESC {
		sort.Sort(sort.Reverse(ByVersion(r.filtered)))
	}
}

// FilterByTitle filters all Releases.filtered by matching the release title
// with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) FilterByTitle(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	r.filterBy(func(r Releaser) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(r.Title()) {
			return true
		}
		return false
	}, inverse)
}

// filterBy filters all Releases.filtered using the passed function.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) filterBy(f func(r Releaser) bool, inverse bool) {
	var result []Releaser

	for _, r := range r.filtered {
		if inverse == false && f(r) {
			result = append(result, r)
			continue
		}

		if inverse == true && !f(r) {
			result = append(result, r)
			continue
		}
	}

	r.filtered = result
}

// FilterByMediaType filters all Releases.filtered by matching the downloads
// media type with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) FilterByMediaType(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	r.filterDownloadsBy(func(d Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.Filetype()) {
			return true
		}
		return false
	}, inverse)
}

// filterDownloadsBy filters all Downloads for Releases.filtered using the
// passed function.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) filterDownloadsBy(f func(d Download) bool, inverse bool) {
	var result []Releaser

	for _, r := range r.filtered {
		for _, download := range r.Downloads() {
			if inverse == false && f(download) {
				result = append(result, r)
				continue
			}

			if inverse == true && !f(download) {
				result = append(result, r)
				continue
			}
		}
	}

	r.filtered = result
}

// FilterByUrl filters all Releases.filtered by matching the release download
// URL with the provided RegExp string.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) FilterByUrl(regexpStr string, inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	r.filterDownloadsBy(func(d Download) bool {
		re := regexp.MustCompile(regexpStr)
		if re.MatchString(d.Url()) {
			return true
		}
		return false
	}, inverse)
}

// FilterByPrerelease filters all Releases.filtered by matching only the
// pre-releases.
//
// When inversed bool is set to true, the unmatched releases will be used
// instead.
func (r *Releases) FilterByPrerelease(inversed ...interface{}) {
	inverse := false
	if len(inversed) > 0 {
		inverse = inversed[0].(bool)
	}

	r.filterBy(func(r Releaser) bool {
		if r.IsPreRelease() == true {
			return true
		}
		return false
	}, inverse)
}

// ResetFilters resets the Releases.filtered to their original state before
// applying any filters.
func (r *Releases) ResetFilters() {
	r.filtered = r.original
}

// Len is a convenience method to get the Releases.filtered slice length.
func (r *Releases) Len() int {
	return len(r.filtered)
}

// First is a convenience method to get the first release from the
// Releases.filtered slice.
func (r *Releases) First() Releaser {
	return r.filtered[0]
}

// Filtered is a Releases.filtered getter.
func (r *Releases) Filtered() []Releaser {
	return r.filtered
}

// SetFiltered is a Releases.filtered setter.
func (r *Releases) SetFiltered(filtered []Releaser) {
	r.filtered = filtered
}

// Original is a Releases.original getter.
func (r *Releases) Original() []Releaser {
	return r.original
}

// SetOriginal is a Releases.original setter.
func (r *Releases) SetOriginal(original []Releaser) {
	r.original = original
}
