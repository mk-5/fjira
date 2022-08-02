package jira

import (
	"encoding/json"
	"github.com/mk5/fjira/internal/app"
	"strings"
)

//
// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-transitions-post
//

const (
	GetTransitions = "/rest/api/2/issue/{issue}/transitions"
	DoTransition   = "/rest/api/2/transitions"
)

type nilParams struct{}

type transitionsResponse struct {
	Transitions []IssueTransition `json:"transitions"`
}

type doTransitionRequest struct {
	Transition string `json:"transition"`
}

func (a *httpApi) DoTransition(issueId string, transition *IssueTransition) error {
	request := doTransitionRequest{Transition: transition.Id}
	requestBody, _ := json.Marshal(request)
	_, err := a.jiraRequest("POST", strings.Replace(GetTransitions, "{issue}", issueId, 1), &nilParams{}, strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}
	return nil
}

func (a *httpApi) FindTransitions(issueId string) ([]IssueTransition, error) {
	responseBody, _ := a.jiraRequest("GET", strings.Replace(GetTransitions, "{issue}", issueId, 1), &nilParams{}, nil)
	var sResponse transitionsResponse
	if err := json.Unmarshal(responseBody, &sResponse); err != nil {
		app.Error(err.Error())
		return nil, SearchDeserializeErr
	}
	var transitions = make([]IssueTransition, 0, 1000)
	transitions = append(transitions, sResponse.Transitions...)
	return transitions, nil
}
