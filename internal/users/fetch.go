package users

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

type RecordsProvider interface {
	FetchUsers(projectKey string, query string) []jira.User
}

type apiRecordsProvider struct {
	api jira.Api
}

func NewApiRecordsProvider(api jira.Api) RecordsProvider {
	return &apiRecordsProvider{
		api: api,
	}
}

func (r *apiRecordsProvider) FetchUsers(projectKey string, query string) []jira.User {
	us, err := r.api.FindUsersWithQuery(projectKey, query)
	if err != nil {
		app.Error(err.Error())
	}
	return us
}
