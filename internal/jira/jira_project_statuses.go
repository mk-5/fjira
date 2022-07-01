package jira

import (
	"encoding/json"
	"github.com/mk5/fjira/internal/app"
	"strings"
)

//
// https://docs.atlassian.com/software/jira/docs/api/REST/8.5.1/#api/2/project-getAllStatuses
//

type statusesResponse struct {
	Statuses []JiraIssueStatus `json:"statuses"`
}

const (
	GetProjectStatuses = "/rest/api/2/project/{project}/statuses"
)

func (a *httpJiraApi) FindProjectStatuses(projectId string) ([]JiraIssueStatus, error) {
	responseBody, _ := a.jiraRequest("GET", strings.Replace(GetProjectStatuses, "{project}", projectId, 1), &nilParams{}, nil)
	var sResponse []statusesResponse
	distinct := make(map[string]bool)
	if err := json.Unmarshal(responseBody, &sResponse); err != nil {
		app.Error(err.Error())
		return nil, SearchDeserializeErr
	}
	var statuses = make([]JiraIssueStatus, 0, 100)
	for _, row := range sResponse {
		for _, status := range row.Statuses {
			if distinct[status.Name] == true {
				continue
			}
			statuses = append(statuses, status)
			distinct[status.Name] = true
		}
	}
	return statuses, nil
}
