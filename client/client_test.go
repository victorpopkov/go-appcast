package client

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNew(t *testing.T) {
	c := New()
	assert.IsType(t, Client{}, *c)
	assert.IsType(t, http.Client{}, *c.HTTPClient)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36", c.UserAgent)
	assert.Equal(t, time.Duration(0), c.Timeout)
}

func TestClient_InsecureSkipVerify(t *testing.T) {
	c := New()
	c.InsecureSkipVerify()
	assert.IsType(t, &http.Transport{}, c.HTTPClient.Transport)
}

func TestClient_Do(t *testing.T) {
	// mock the request
	c := New()
	c.Timeout = time.Duration(time.Second)
	httpmock.ActivateNonDefault(c.HTTPClient)
	httpmock.RegisterResponder("GET", "https://example.com/", httpmock.NewStringResponder(200, `Test`))
	defer httpmock.DeactivateAndReset()

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
