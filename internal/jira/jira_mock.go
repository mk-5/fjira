package jira

import (
	"net/http"
	"net/http/httptest"
)

func NewJiraApiMock(handler func(w http.ResponseWriter, r *http.Request)) Api {
	stubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler != nil {
			handler(w, r)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("")) //nolint:errcheck
	}))
	api, err := NewApi(stubServer.URL, "test", "test")
	if err != nil {
		panic(err)
	}
	return api
}
