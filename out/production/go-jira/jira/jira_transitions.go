package jira

import (
	"encoding/json"
	"log"
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
	Transitions []JiraIssueTransition `json:"transitions"`
}

type doTransitionRequest struct {
	Transition string `json:"transition"`
}

func (a *httpJiraApi) DoTransition(issueId string, transition *JiraIssueTransition) error {
	request := doTransitionRequest{Transition: transition.Id}
	requestBody, _ := json.Marshal(request)
	_, err := a.jiraRequest("POST", strings.Replace(GetTransitions, "{issue}", issueId, 1), &nilParams{}, strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}
	return nil
}

func (a *httpJiraApi) FindTransitions(issueId string) ([]JiraIssueTransition, error) {
	responseBody, _ := a.jiraRequest("GET", strings.Replace(GetTransitions, "{issue}", issueId, 1), &nilParams{}, nil)
	var sResponse transitionsResponse
	if err := json.Unmarshal(responseBody, &sResponse); err != nil {
		log.Fatalln(err)
		return nil, SearchDeserializeErr
	}
	var transitions = make([]JiraIssueTransition, 0, 1000)
	for _, tr := range sResponse.Transitions {
		transitions = append(transitions, tr)
	}
	return transitions, nil
}
