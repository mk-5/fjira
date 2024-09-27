package jira

import (
	"encoding/json"
	"fmt"
)

type IssueType struct {
	Name string `json:"name"`
}

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"Fields"`
	Id     string      `json:"id"`
}

type IssueFields struct {
	Summary     string  `json:"summary"`
	Project     Project `json:"project"`
	Description string  `json:"description,omitempty"`
	Reporter    struct {
		AccountId   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	} `json:"reporter"`
	Assignee struct {
		AccountId   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	} `json:"assignee"`
	Type struct {
		Name string `json:"name"`
	} `json:"issuetype"`
	Status  Status
	Comment struct {
		Comments   []Comment `json:"comments"`
		MaxResults int32     `json:"maxResults"`
		Total      int32     `json:"total"`
		StartAt    int32     `json:"startAt"`
	} `json:"comment"`
	Labels     []string `json:"labels"`
	Attachment []IssueAttachment
}

type IssueAttachment struct {
	Self     string `json:"self"`
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Author   struct {
		Self         string `json:"self"`
		AccountId    string `json:"accountId"`
		EmailAddress string `json:"emailAddress"`
		AvatarUrls   struct {
			X48 string `json:"48x48"`
			X24 string `json:"24x24"`
			X16 string `json:"16x16"`
			X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		DisplayName string `json:"displayName"`
		Active      bool   `json:"active"`
		TimeZone    string `json:"timeZone"`
		AccountType string `json:"accountType"`
	} `json:"author"`
	Created   string `json:"created"`
	Size      int    `json:"size"`
	MimeType  string `json:"mimeType"`
	Content   string `json:"content"`
	Thumbnail string `json:"thumbnail"`
}

const (
	GetJiraIssuePath = "/rest/api/2/issue/%s"
)

func (api *httpApi) GetIssueDetailed(id string) (*Issue, error) {
	body, err := api.jiraRequest("GET", fmt.Sprintf(GetJiraIssuePath, id), &nilParams{}, nil)
	if err != nil {
		return nil, err
	}
	var jiraIssue Issue
	if err := json.Unmarshal(body, &jiraIssue); err != nil {
		return nil, SearchDeserializeErr
	}
	return &jiraIssue, nil
}
