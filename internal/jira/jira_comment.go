package jira

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Comment struct {
	Author  User   `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type commentRequestBody struct {
	Body string `json:"body"`
}

const (
	DoCommentIssueRestPath = "/rest/api/2/issue/%s/comment"
)

func (api *httpApi) DoComment(issueId string, commentBody string) error {
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
