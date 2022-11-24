package jira

import (
	"encoding/json"
	"errors"
	"github.com/mk-5/fjira/internal/app"
)

const (
	FindUser = "/rest/api/2/user/assignable/search"
)

var UserSearchDeserializeErr = errors.New("Cannot deserialize jira user search response.")

type findUserQueryParams struct {
	Project    string `url:"project"`
	MaxResults int    `url:"maxResults"`
}

func (api httpApi) FindUsers(project string) ([]User, error) {
	queryParams := &findUserQueryParams{
		Project:    project,
		MaxResults: 10000,
	}
	response, err := api.jiraRequest("GET", FindUser, queryParams, nil)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(response, &users); err != nil {
		app.Error(err.Error())
		return nil, UserSearchDeserializeErr
	}
	return users, nil
}
