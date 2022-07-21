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

type findLabelsQueryParams struct {
	Query string `url:"query"`
}

const (
	LabelsJira      = "/rest/api/1.0/labels/%s/suggest"
	DoLabelRestPath = "/rest/api/2/issue/%s"
)

func (api *httpApi) FindLabels(issue *Issue, query string) ([]string, error) {
	response, err := api.jiraRequest("GET", fmt.Sprintf(LabelsJira, url.QueryEscape(issue.Id)), &findLabelsQueryParams{Query: query}, nil)
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
	_, err := api.jiraRequest("PUT", fmt.Sprintf(DoLabelRestPath, url.QueryEscape(issueId)), &nilParams{}, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
