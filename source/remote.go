package source

import (
	"io/ioutil"

	"github.com/victorpopkov/go-appcast/appcaster"
	"github.com/victorpopkov/go-appcast/client"
	"github.com/victorpopkov/go-appcast/provider"
)

// DefaultClient is the default Client that is used for making requests in the
// appcast package.
var DefaultClient = client.New()

// Remoter is the interface that wraps the Remote methods.
type Remoter interface {
	appcaster.Sourcer
	Request() *client.Request
	SetRequest(request *client.Request)
	Url() string
	SetUrl(url string)
}

// Remote represents an appcast source from the remote location.
type Remote struct {
	*appcaster.Source
	request *client.Request
	url     string
}

// NewRemote returns a new Remote instance pointer with the prepared
// Remote.request and Remote.url ready to be used Remote.load.
//
// Supports both the remote URL string or Request struct pointer as an argument.
func NewRemote(src interface{}) (*Remote, error) {
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

	s := &Remote{
		Source:  &appcaster.Source{},
		request: request,
		url:     request.HTTPRequest.URL.String(),
	}

	return s, nil
}

// Load loads an appcast content into the Remote.Source.content from the remote
// source by using the Remote.request set earlier.
func (r *Remote) Load() error {
	resp, err := DefaultClient.Do(r.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r.SetContent(body)

	r.GuessProvider()
	r.GenerateChecksum(appcaster.SHA256)

	return nil
}

// GuessProvider attempts to guess the supported provider based on the
// Remote.url and Remote.Source.content. By default returns an
// Unknown provider.
func (r *Remote) GuessProvider() {
	r.SetProvider(provider.GuessProviderByUrl(r.url))
	if r.Provider() == provider.Unknown {
		r.SetProvider(provider.GuessProviderByContent(r.Content()))
	}
}

// Request is a Remote.request getter.
func (r *Remote) Request() *client.Request {
	return r.request
}

// SetRequest is a Remote.request setter.
func (r *Remote) SetRequest(request *client.Request) {
	r.request = request
}

// Url is a Remote.url getter.
func (r *Remote) Url() string {
	return r.url
}

// SetUrl is a Remote.url setter.
func (r *Remote) SetUrl(url string) {
	r.url = url
}
