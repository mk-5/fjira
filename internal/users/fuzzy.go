package users

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

const (
	typeaheadSearchThreshold = 100
)

func NewFuzzyFind(projectKey string, api jira.Api) (*app.FuzzyFind, *[]jira.User) {
	var us []jira.User
	var prevQuery string
	provider := NewApiRecordsProvider(api)
	return app.NewFuzzyFindWithProvider(ui.MessageSelectUser, func(query string) []string {
		// it searches up to {typeaheadThreshold} records using typeahead - then it do regular fuzzy-find
		if len(us) > 0 && len(us) < typeaheadSearchThreshold && len(query) > len(prevQuery) {
			return FormatJiraUsers(us)
		}
		prevQuery = query
		app.GetApp().Loading(true)
		us = provider.FetchUsers(projectKey, query)
		app.GetApp().Loading(false)
		us = append(us, jira.User{DisplayName: ui.MessageAll})
		usersStrings := FormatJiraUsers(us)
		return usersStrings
	}), &us
}
