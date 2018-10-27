package appcast

import "io/ioutil"

// RemoteSourcer is the interface that wraps the RemoteSource methods.
type RemoteSourcer interface {
	Sourcer
	Request() *Request
	Url() string
}

// RemoteSource represents an appcast source from the remote location.
type RemoteSource struct {
	*Source
	request *Request
	url     string
}

// NewRemoteSource returns a new RemoteSource instance pointer with the prepared
// RemoteSource.request and RemoteSource.url ready to be used RemoteSource.load.
//
// Supports both the remote URL string or Request struct pointer as an argument.
func NewRemoteSource(src interface{}) (*RemoteSource, error) {
	var request *Request

	switch v := src.(type) {
	case *Request:
		request = v
	case string:
		newReq, err := NewRequest(v)
		if err != nil {
			return nil, err
		}
		request = newReq
	}

	s := &RemoteSource{
		Source:  &Source{},
		request: request,
		url:     request.HTTPRequest.URL.String(),
	}

	return s, nil
}

// Load loads an appcast content into the RemoteSource.Source.content from the
// remote source by using the RemoteSource.request set earlier.
func (s *RemoteSource) Load() error {
	resp, err := DefaultClient.Do(s.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	s.content = body

	s.GuessProvider()
	s.checksum = NewChecksum(SHA256, s.content)

	return nil
}

// GuessProvider attempts to guess the supported provider based on the
// RemoteSource.url and RemoteSource.Source.content. By default returns an
// Unknown provider.
func (s *RemoteSource) GuessProvider() {
	s.provider = GuessProviderByUrl(s.url)
	if s.provider == Unknown {
		s.provider = GuessProviderByContent(s.content)
	}
}

// Request is a RemoteSource.request getter.
func (s *RemoteSource) Request() *Request {
	return s.request
}

// Url is a RemoteSource.url getter.
func (s *RemoteSource) Url() string {
	return s.url
}
