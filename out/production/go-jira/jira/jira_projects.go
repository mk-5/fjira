package jira

import (
	"encoding/json"
)

const (
	ProjectsJira = "/rest/api/2/project"
)

func (api *httpJiraApi) FindProjects() ([]JiraProject, error) {
	response, err := api.jiraRequest("GET", ProjectsJira, nil, nil)
	if err != nil {
		return nil, err
	}
	var projects []JiraProject
	if err := json.Unmarshal(response, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}
