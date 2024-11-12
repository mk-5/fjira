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

func Test_jiraRequest_combinePaths(t *testing.T) {
	tests := []struct {
		name     string
		apiUrl   string
		restPath string
		wanted   string
	}{
		{"should add label without error", "http://localhost", "/api1", "http://localhost/api1"},
		{"should add label without error", "http://localhost/", "/api1", "http://localhost/api1"},
		{"should add label without error", "http://localhost/jira-api-v2", "/api1", "http://localhost/jira-api-v2/api1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.apiUrl)
			api := &httpApi{restUrl: u}

			result, _ := api.jiraRequestUrl(tt.restPath, nil)

			assert.Equal(t, tt.wanted, result)
		})
	}
}
