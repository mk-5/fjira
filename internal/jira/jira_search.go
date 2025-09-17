package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/mk-5/fjira/internal/app"
)

const (
	SearchJira      = "/rest/api/3/search/jql"
	JiraIssueRegexp = "^[a-zA-Z0-9]{1,10}-[0-9]{1,20}$"
)

var SearchDeserializeErr = errors.New("Cannot deserialize jira search response.")

type searchQueryParams struct {
	Jql        string `url:"jql"`
	MaxResults int32  `url:"maxResults"`
	Fields     string `url:"fields"`
	StartAt    int32  `url:"startAt"`
}

type searchResponse struct {
	Total      int32   `json:"total"`
	MaxResults int32   `json:"maxResults"`
	Issues     []Issue `json:"issues"`
	IsLast     bool    `json:"isLast"`
}

func (api *httpApi) Search(query string) ([]Issue, int32, error) {
	isJqlAboutIssue, _ := regexp.Match(JiraIssueRegexp, []byte(query))
	jql := fmt.Sprintf("summary~\"%s*\"", query)
	if isJqlAboutIssue {
		jql = fmt.Sprintf("key=\"%s\"", query)
	}
	issues, total, _, err := api.SearchJqlPageable(jql, 0, 100)
	return issues, total, err
}

func (api *httpApi) SearchJql(jql string) ([]Issue, error) {
	issues, _, _, err := api.SearchJqlPageable(jql, 0, 100)
	return issues, err
}

func (api *httpApi) SearchJqlPageable(jql string, page int32, pageSize int32) ([]Issue, int32, int32, error) {
	queryParams := searchQueryParams{
		Jql:        jql,
		MaxResults: pageSize,
		StartAt:    page * pageSize,
		Fields:     "id,key,summary,issuetype,project,reporter,status,assignee",
	}
	body, err := api.jiraRequest("GET", SearchJira, queryParams, nil)
	if err != nil {
		return nil, -1, pageSize, err
	}
	var sResponse searchResponse
	if err := json.Unmarshal(body, &sResponse); err != nil {
		app.Error(err.Error())
		return nil, -1, pageSize, SearchDeserializeErr
	}
	return sResponse.Issues, sResponse.Total, sResponse.MaxResults, err
}
