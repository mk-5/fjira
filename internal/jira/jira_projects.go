package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

const (
	ProjectsJira      = "/rest/api/3/project/search"
	ProjectsByKeyJira = "/rest/api/3/project/%s"
)

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type searchProjectsQueryParams struct {
	MaxResults int32 `url:"maxResults"`
	StartAt    int32 `url:"startAt"`
}

type searchProjectsResponse struct {
	MaxResults int32     `json:"maxResults"`
	StartAt    int32     `json:"startAt"`
	Values     []Project `json:"values"`
}

var (
	ErrProjectNotFound = errors.New("project not found")
)

func (api *httpApi) FindProjects() ([]Project, error) {
	params := &searchProjectsQueryParams{}
	params.MaxResults = 100
	response, err := api.jiraRequest("GET", ProjectsJira, params, nil)
	if err != nil {
		return nil, err
	}
	var projects searchProjectsResponse
	if err := json.Unmarshal(response, &projects); err != nil {
		return nil, err
	}
	return projects.Values, nil
}

func (api *httpApi) FindProject(projectKey string) (*Project, error) {
	u := fmt.Sprintf(ProjectsByKeyJira, url.QueryEscape(projectKey))
	response, err := api.jiraRequest("GET", u, nil, nil)
	if err != nil {
		return nil, err
	}
	var project *Project
	if err := json.Unmarshal(response, &project); err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}
