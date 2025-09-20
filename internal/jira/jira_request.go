package jira

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/google/go-querystring/query"
)

func (api *httpApi) jiraRequest(method string, restPath string, queryParams interface{}, reqBody io.Reader) ([]byte, error) {
	u, err := api.jiraRequestUrl(restPath, queryParams)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u, reqBody)
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
		return nil, fmt.Errorf("jira error, status: %s - request: %s", response.Status, req.RequestURI)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return body, nil
}

func (api *httpApi) jiraRequestUrl(restPath string, queryParams interface{}) (string, error) {
	queryParamsValues, err := query.Values(queryParams)
	if err != nil {
		return "", err
	}
	u := api.restUrl.ResolveReference(&url.URL{Path: path.Join(api.restUrl.Path, restPath), RawQuery: queryParamsValues.Encode()})
	return u.String(), err
}
