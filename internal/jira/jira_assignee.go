package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type cloudAssigneeRequestBody struct {
	AccountId string `json:"accountId"`
}

type onPremiseAssigneeRequestBody struct {
	Fields struct {
		Assignee struct {
			Name string `json:"name"`
		} `json:"assignee"`
	} `json:"fields"`
}

const (
	DoAssigneePathCloud     = "/rest/api/2/issue/%s/assignee"
	DoAssigneePathOnPremise = "/rest/api/2/issue/%s"
)

var (
	CannotPerformAssignmentErr = errors.New("invalid assignee data. Cannot perform do-assignment request")
)

func (api *httpApi) DoAssignee(issueId string, user *User) error {
	var url string
	var body interface{}
	if user.AccountId != "" {
		url = fmt.Sprintf(DoAssigneePathCloud, issueId)
		body = &cloudAssigneeRequestBody{AccountId: user.AccountId}
	} else if user.Name != "" {
		url = fmt.Sprintf(DoAssigneePathOnPremise, issueId)
		body = &onPremiseAssigneeRequestBody{}
		(body.(*onPremiseAssigneeRequestBody)).Fields.Assignee.Name = user.Name
	} else {
		return CannotPerformAssignmentErr
	}
	jsonBody, _ := json.Marshal(body)
	_, err := api.jiraRequest("PUT", url, nil, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	return nil
}
