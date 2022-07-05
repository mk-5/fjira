package jira

import (
	"encoding/json"
)

const (
	LabelsJira = "/rest/api/2/label"
)

func (api *httpJiraApi) FindLabels() ([]string, error) {
	response, err := api.jiraRequest("GET", LabelsJira, nil, nil)
	if err != nil {
		return nil, err
	}
	var labels JiraLabelsResponse
	if err := json.Unmarshal(response, &labels); err != nil {
		return nil, err
	}
	return labels.Values, nil
}
