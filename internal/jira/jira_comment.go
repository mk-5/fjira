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
	jsonBody, _ := json.Marshal(&commentRequestBody{
		Body: commentBody,
	})
	_, err := api.jiraRequest("POST", fmt.Sprintf(DoCommentIssueRestPath, issueId), &nilParams{}, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
