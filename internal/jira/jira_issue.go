package jira

import (
	"encoding/json"
	"fmt"
	"github.com/mk5/fjira/internal/app"
)

const (
	GetJiraIssuePath = "/rest/api/2/issue/%s"
)

func (api *httpJiraApi) GetIssueDetailed(id string) (*JiraIssue, error) {
	body, err := api.jiraRequest("GET", fmt.Sprintf(GetJiraIssuePath, id), &nilParams{}, nil)
	if err != nil {
		return nil, err
	}
	var jiraIssue JiraIssue
	if err := json.Unmarshal(body, &jiraIssue); err != nil {
		app.Error(err.Error())
		return nil, SearchDeserializeErr
	}
	return &jiraIssue, nil
}
