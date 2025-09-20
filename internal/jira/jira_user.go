package jira

import (
	"encoding/json"
	"errors"

	"github.com/mk-5/fjira/internal/app"
)

type User struct {
	AccountId    string            `json:"accountId"`
	Active       bool              `json:"active"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	DisplayName  string            `json:"displayName"`
	EmailAddress string            `json:"emailAddress"`
	Locale       string            `json:"locale"`
	Self         string            `json:"self"`
	TimeZone     string            `json:"timeZone"`
	Key          string            `json:"key"`  // field used by on-premise installation
	Name         string            `json:"name"` // field used by on-premise installation
}

const (
	FindUser = "/rest/api/2/user/assignable/search"
)

var ErrUserSearchDeserialize = errors.New("cannot deserialize jira user search response")

type findUserQueryParams struct {
	Project    string  `url:"project"`
	MaxResults int     `url:"maxResults"`
	Query      *string `url:"query"`
	Username   *string `url:"username"`
}

func (api *httpApi) FindUsers(project string) ([]User, error) {
	return api.FindUsersWithQuery(project, "")
}

func (api *httpApi) FindUsersWithQuery(project string, query string) ([]User, error) {
	queryParams := &findUserQueryParams{
		Project:    project,
		MaxResults: 10000,
	}
	if query != "" && !api.IsJiraServer() {
		queryParams.Query = &query
	}
	if query != "" && api.IsJiraServer() {
		queryParams.Username = &query
	}
	response, err := api.jiraRequest("GET", FindUser, queryParams, nil)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(response, &users); err != nil {
		app.Error(err.Error())
		return nil, ErrUserSearchDeserialize
	}
	return users, nil
}
