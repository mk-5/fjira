package jira

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func Test_httpApi_jiraRequest_should_return_error_when_http_client_error(t *testing.T) {
	// given
	api := &httpApi{
		client: &http.Client{
			Transport: &authInterceptor{core: http.DefaultTransport, token: "test"},
		},
		restUrl: &url.URL{},
	}

	// when
	_, err := api.jiraRequest("POST", "test", struct{}{}, nil)

	// then
	assert.NotNil(t, err)
}

func Test_httpApi_jiraRequest_should_return_error_when_invalid_params(t *testing.T) {
	// given
	api := &httpApi{
		client: &http.Client{
			Transport: &authInterceptor{core: http.DefaultTransport, token: "test"},
		},
		restUrl: &url.URL{},
	}

	// when
	_, err := api.jiraRequest("POST", "test", "invalid params", nil)

	// then
	assert.NotNil(t, err)
}
