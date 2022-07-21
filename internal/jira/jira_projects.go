package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

const (
	ProjectsJira      = "/rest/api/2/project"
	ProjectsByKeyJira = "/rest/api/2/project/%s"
)

var (
	ProjectNotFoundError = errors.New("Project not found.")
)

func (api *httpApi) FindProjects() ([]Project, error) {
	response, err := api.jiraRequest("GET", ProjectsJira, nil, nil)
	if err != nil {
		return nil, err
	}
	var projects []Project
	if err := json.Unmarshal(response, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (api *httpApi) FindProject(projectKey string) (*Project, error) {
	url := fmt.Sprintf(ProjectsByKeyJira, url.QueryEscape(projectKey))
	response, err := api.jiraRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	var project *Project
	if err := json.Unmarshal(response, &project); err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ProjectNotFoundError
	}
	return project, nil
}
