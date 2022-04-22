package jira

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	FindUser = "/rest/api/2/user/assignable/search"
)

var UserSearchDeserializeErr = errors.New("Cannot deserialize jira user search response.")

type findUserQueryParams struct {
	Project    string `url:"project"`
	MaxResults int    `url:"maxResults"`
}

func (api httpJiraApi) FindUser(project string) ([]JiraUser, error) {
	queryParams := &findUserQueryParams{
		Project:    project,
		MaxResults: 10000,
	}
	response, err := api.jiraRequest("GET", FindUser, queryParams, nil)
	if err != nil {
		return nil, err
	}
	var users []JiraUser
	if err := json.Unmarshal(response, &users); err != nil {
		log.Fatalln(err)
		return nil, UserSearchDeserializeErr
	}
	return users, nil
}
