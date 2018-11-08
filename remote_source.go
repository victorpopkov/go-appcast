package appcast

import (
	"io/ioutil"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/client"
)

// RemoteSourcer is the interface that wraps the RemoteSource methods.
type RemoteSourcer interface {
	appcaster.Sourcer
	Request() *client.Request
	Url() string
}

// RemoteSource represents an appcast source from the remote location.
type RemoteSource struct {
	*appcaster.Source
	request *client.Request
	url     string
}

// NewRemoteSource returns a new RemoteSource instance pointer with the prepared
// RemoteSource.request and RemoteSource.url ready to be used RemoteSource.load.
//
// Supports both the remote URL string or Request struct pointer as an argument.
func NewRemoteSource(src interface{}) (*RemoteSource, error) {
	var request *client.Request

	switch v := src.(type) {
	case *client.Request:
		request = v
	case string:
		newReq, err := client.NewRequest(v)
		if err != nil {
			return nil, err
		}
		request = newReq
	}

	s := &RemoteSource{
		Source:  &appcaster.Source{},
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
	s.SetContent(body)

	s.GuessProvider()
	s.GenerateChecksum(appcaster.SHA256)

	return nil
}

// GuessProvider attempts to guess the supported provider based on the
// RemoteSource.url and RemoteSource.Source.content. By default returns an
// Unknown provider.
func (s *RemoteSource) GuessProvider() {
	s.SetProvider(GuessProviderByUrl(s.url))
	if s.Provider() == Unknown {
		s.SetProvider(GuessProviderByContent(s.Content()))
	}
}

// Request is a RemoteSource.request getter.
func (s *RemoteSource) Request() *client.Request {
	return s.request
}

// Url is a RemoteSource.url getter.
func (s *RemoteSource) Url() string {
	return s.url
}
