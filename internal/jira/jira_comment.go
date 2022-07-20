package jira

import (
	"encoding/json"
	"fmt"
	"strings"
)

type commentRequestBody struct {
	Body string `json:"body"`
}

const (
	DoCommentIssueRestPath = "/rest/api/2/issue/%s/comment"
)

func (api *httpJiraApi) DoComment(issueId string, commentBody string) error {
	jsonBody, err := json.Marshal(&commentRequestBody{
		Body: commentBody,
	})
	if err != nil {
		return err
	}
	_, err = api.jiraRequest("POST", fmt.Sprintf(DoCommentIssueRestPath, issueId), &nilParams{}, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
