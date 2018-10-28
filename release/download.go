package release

// Downloader is the interface that wraps the Download methods.
type Downloader interface {
	Url() string
	SetUrl(url string)
	Filetype() string
	SetFiletype(filetype string)
	Length() int
	SetLength(length int)
	DsaSignature() string
	SetDsaSignature(dsaSignature string)
	Md5() string
	SetMd5(dsaSignature string)
}

// Download holds a single release download data.
type Download struct {
	// url specifies a remote file URL.
	url string

	// filetype specifies a request MIME type.
	filetype string

	// length specifies a request length.
	length int

	// dsaSignature specifies a file DSA signature value.
	dsaSignature string

	// md5 specifies a file MD5 checksum.
	md5 string
}

// NewDownload returns a new Download instance pointer. Requires an url to be
// passed as a parameter. Optionally, the filetype can be passed as a second
// parameter, the length as a third one, the dsaSignature as a fourth and the
// md5 as a fifth.
func NewDownload(url string, a ...interface{}) *Download {
	d := &Download{
		url: url,
	}

	if len(a) > 0 {
		d.filetype = a[0].(string)
	}

	if len(a) > 1 {
		d.length = a[1].(int)
	}

	if len(a) > 2 {
		d.dsaSignature = a[2].(string)
	}

	if len(a) > 3 {
		d.md5 = a[3].(string)
	}

	return d
}

// Url is a Download.url getter.
func (d *Download) Url() string {
	return d.url
}

// SetUrl is a Download.url setter.
func (d *Download) SetUrl(url string) {
	d.url = url
}

// Filetype is a Download.filetype filetype.
func (d *Download) Filetype() string {
	return d.filetype
}

// SetFiletype is a Download.filetype setter.
func (d *Download) SetFiletype(filetype string) {
	d.filetype = filetype
}

// Length is a Download.length getter.
func (d *Download) Length() int {
	return d.length
}

// SetLength is a Download.length setter.
func (d *Download) SetLength(length int) {
	d.length = length
}

// DsaSignature is a Download.dsaSignature getter.
func (d *Download) DsaSignature() string {
	return d.dsaSignature
}

// SetDsaSignature is a Download.dsaSignature setter.
func (d *Download) SetDsaSignature(dsaSignature string) {
	d.dsaSignature = dsaSignature
}

// Md5 is a Download.md5 getter.
func (d *Download) Md5() string {
	return d.md5
}

// SetMd5 is a Download.md5 setter.
func (d *Download) SetMd5(md5 string) {
	d.md5 = md5
}
