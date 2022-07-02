package jira

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
)

func (api *httpJiraApi) jiraRequest(method string, restPath string, queryParams interface{}, reqBody io.Reader) ([]byte, error) {
	queryParamsValues, err := query.Values(queryParams)
	if err != nil {
		return nil, err
	}
	u := api.restUrl.ResolveReference(&url.URL{Path: restPath, RawQuery: queryParamsValues.Encode()})
	req, err := http.NewRequest(method, u.String(), reqBody)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	response, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Jira error, status: %s - request: %s", response.Status, req.RequestURI)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return body, nil
}
