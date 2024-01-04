package jira

import (
	"net/http"
	"net/http/httptest"
)

func NewJiraApiMock(handler func(w http.ResponseWriter, r *http.Request)) Api {
	return NewJiraApiMockWithTokenType(handler, ApiToken)
}

func NewJiraApiMockWithTokenType(handler func(w http.ResponseWriter, r *http.Request), tokenType JiraTokenType) Api {
	stubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler != nil {
			handler(w, r)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("")) //nolint:errcheck
	}))
	api, err := NewApi(stubServer.URL, "test", "test", tokenType)
	if err != nil {
		panic(err)
	}
	return api
}
