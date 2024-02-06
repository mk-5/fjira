package jira

import (
	"encoding/json"
	"fmt"
)

type Filter struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	JQL       string `json:"jql"`
	Favourite bool   `json:"favourite"`
}

const (
	FilterUrl             = "/rest/api/2/filter/%s"
	MyFilterUrl           = "/rest/api/2/filter/my"
	MyFilterUrlJiraServer = "/rest/api/2/filter/favourite"
)

func (api *httpApi) GetFilter(filterId string) (*Filter, error) {
	resultBytes, err := api.jiraRequest("GET", fmt.Sprintf(FilterUrl, filterId), &nilParams{}, nil)
	if err != nil {
		return nil, err
	}
	var result Filter
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (api *httpApi) GetMyFilters() ([]Filter, error) {
	url := MyFilterUrl
	if api.IsJiraServer() {
		url = MyFilterUrlJiraServer
	}
	resultBytes, err := api.jiraRequest("GET", url, &nilParams{}, nil)
	if err != nil {
		return nil, err
	}
	var result []Filter
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
