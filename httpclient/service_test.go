package httpclient_test

import (
	"adjust/httpclient"
	"adjust/mocks"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSendRequestSuccess(t *testing.T) {
	// build response JSON
	payload := `sample text`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(payload)))

	mokClient := &mocks.MockClient{}
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	server := httpclient.NewServer(10, []string{}, mokClient)
	resp, err := server.SendHttpRequest("https://google.com")
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, 200, resp.StatusCode)
}

func TestNormalizeUrl(t *testing.T) {
	server := httpclient.NewServer(10, []string{}, &mocks.MockClient{})
	url := "adjust.com"
	res, err := server.NormalizeUrl(url)
	assert.Nil(t, err)
	assert.EqualValues(t, res, "https://adjust.com")
}

func TestMd5(t *testing.T) {
	server := httpclient.NewServer(10, []string{}, &mocks.MockClient{})
	payload := `sample text`
	r := ioutil.NopCloser(bytes.NewReader([]byte(payload)))
	res, err := server.Md5(r)
	assert.Nil(t, err)
	assert.EqualValues(t, res, "70ee1738b6b21e2c8a43f3a5ab0eee71")
}
