package jira

import (
	"net/http"
)

const (
	Authorization   = "Authorization"
	XAtlassianToken = "X-Atlassian-Token"
)

type authInterceptor struct {
	core  http.RoundTripper
	token string
}

func (a *authInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		if r.Body != nil {
			_ = r.Body.Close()
		}
	}()
	newRequest := a.modifyRequest(r)
	return a.core.RoundTrip(newRequest)
}

func (a *authInterceptor) modifyRequest(r *http.Request) *http.Request {
	r.Header.Set(Authorization, "Basic "+a.token)
	r.Header.Set(XAtlassianToken, "no-check")
	return r
}
