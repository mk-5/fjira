package jira

import (
	"net/http"
	"net/http/httptest"
)

func NewJiraApiMock(handler func(w http.ResponseWriter, r *http.Request)) JiraApi {
	stubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	api, err := NewJiraApi(stubServer.URL, "test", "test")
	if err != nil {
		panic(err)
	}
	return api
}
