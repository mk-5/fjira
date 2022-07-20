package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type labelRequestBody struct {
	Update struct {
		Labels []labelAdd `json:"labels"`
	} `json:"update"`
}

type labelAdd struct {
	Add string `json:"add"`
}

const (
	LabelsJira      = "/rest/api/2/label"
	DoLabelRestPath = "/rest/api/2/issue/%s"
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

func (api *httpJiraApi) AddLabel(issueId string, label string) error {
	request := &labelRequestBody{}
	request.Update.Labels = make([]labelAdd, 0, 1)
	request.Update.Labels = append(request.Update.Labels, labelAdd{Add: label})
	jsonBody, _ := json.Marshal(request)
	_, err := api.jiraRequest("PUT", fmt.Sprintf(DoLabelRestPath, url.QueryEscape(issueId)), &nilParams{}, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
