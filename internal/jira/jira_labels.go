package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type LabelsSuggestionsResponseBody struct {
	Token       string `json:"token"`
	Suggestions []struct {
		Label string `json:"label"`
		Html  string `json:"html"`
	} `json:"suggestions"`
}

type labelRequestBody struct {
	Update struct {
		Labels []labelAdd `json:"labels"`
	} `json:"update"`
}

type labelAdd struct {
	Add string `json:"add"`
}

type findLabelsQueryParams struct {
	Query string `url:"query"`
}

const (
	LabelsForIssuePath   = "/rest/api/1.0/labels/%s/suggest"
	LabelsForProjectPath = "/rest/api/1.0/labels/suggest"
	DoLabelPath          = "/rest/api/2/issue/%s"
)

func (api *httpApi) FindLabels(issue *Issue, query string) ([]string, error) {
	path := LabelsForProjectPath
	if issue != nil {
		path = fmt.Sprintf(LabelsForIssuePath, url.QueryEscape(issue.Id))
	}
	response, err := api.jiraRequest("GET", path, &findLabelsQueryParams{Query: query}, nil)
	if err != nil {
		return nil, err
	}
	var responseBody LabelsSuggestionsResponseBody
	if err := json.Unmarshal(response, &responseBody); err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(responseBody.Suggestions))
	for _, label := range responseBody.Suggestions {
		labels = append(labels, label.Label)
	}
	return labels, nil
}

func (api *httpApi) AddLabel(issueId string, label string) error {
	request := &labelRequestBody{}
	request.Update.Labels = make([]labelAdd, 0, 1)
	request.Update.Labels = append(request.Update.Labels, labelAdd{Add: label})
	jsonBody, _ := json.Marshal(request)
	_, err := api.jiraRequest("PUT", fmt.Sprintf(DoLabelPath, url.QueryEscape(issueId)), &nilParams{}, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
