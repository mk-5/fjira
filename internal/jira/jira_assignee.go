package jira

import (
	"encoding/json"
	"fmt"
	"strings"
)

type assigneeRequestBody struct {
	AccountId string `json:"accountId"`
}

const (
	DoAssigneePath = "/rest/api/2/issue/%s/assignee"
)

func (api *httpApi) DoAssignee(issueId string, accountId string) error {
	url := fmt.Sprintf(DoAssigneePath, issueId)
	body := &assigneeRequestBody{AccountId: accountId}
	jsonBody, _ := json.Marshal(body)
	_, err := api.jiraRequest("PUT", url, nil, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
