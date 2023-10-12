package jira

import (
	"encoding/json"
	"fmt"
)

const (
	FilterUrl   = "/rest/api/2/filter/%s"
	MyFilterUrl = "/rest/api/2/filter/my"
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
	resultBytes, err := api.jiraRequest("GET", MyFilterUrl, &nilParams{}, nil)
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
