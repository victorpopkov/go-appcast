package appcast

// A Download holds everything describing the release download.
type Download struct {
	// URL specifies the download URL.
	URL string

	// Type specifies the download request MIME type.
	Type string

	// Length specifies the download request length.
	Length int
}

// NewDownload returns a new Download instance pointer. Requires an URL.
// Optionally, the type can be specified as a second parameter and length as a
// third one.
func NewDownload(URL string, a ...interface{}) *Download {
	d := &Download{
		URL: URL,
	}

	if len(a) > 0 {
		d.Type = a[0].(string)
	}

	if len(a) > 1 {
		d.Length = a[1].(int)
	}

	return d
}
