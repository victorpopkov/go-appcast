package appcast

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.IsType(t, Client{}, *c)
	assert.IsType(t, http.Client{}, *c.HTTPClient)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36", c.UserAgent)
	assert.Equal(t, time.Duration(0), c.Timeout)
}

func TestInsecureSkipVerify(t *testing.T) {
	c := NewClient()
	c.InsecureSkipVerify()
	assert.IsType(t, &http.Transport{}, c.HTTPClient.Transport)
}

func TestDo(t *testing.T) {
	c := NewClient()
	c.Timeout = time.Duration(time.Second)

	// mock the request
	httpmock.ActivateNonDefault(c.HTTPClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.com/", httpmock.NewStringResponder(200, `Test`))

	// test (successful)
	req, err := NewRequest("https://example.com/")
	assert.Nil(t, err)

	resp, err := c.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Test", string(body))

	// test (error)
	req, _ = NewRequest("invalid")
	resp, err = c.Do(req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}
