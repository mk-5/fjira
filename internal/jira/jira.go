package jira

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
)

type JiraApi interface {
	Search(query string) ([]JiraIssue, int32, error)
	SearchJql(query string) ([]JiraIssue, error)
	SearchJqlPageable(query string, page int32, pageSize int32) ([]JiraIssue, int32, int32, error)
	FindUsers(project string) ([]JiraUser, error)
	FindProjects() ([]JiraProject, error)
	FindLabels() ([]string, error)
	FindProject(projectKey string) (*JiraProject, error)
	FindTransitions(issueId string) ([]JiraIssueTransition, error)
	FindProjectStatuses(projectId string) ([]JiraIssueStatus, error)
	DoTransition(issueId string, transition *JiraIssueTransition) error
	DoAssignee(issueId string, accountId string) error
	GetIssueDetailed(issueId string) (*JiraIssue, error)
	DoComment(issueId string, commentBody string) error
	Close()
}

type JiraProject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type JiraIssueType struct {
	Name string `json:"name"`
}

type JiraIssue struct {
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"Fields"`
	Id     string          `json:"id"`
}

type JiraComment struct {
	Author  JiraUser `json:"author"`
	Body    string   `json:"body"`
	Created string   `json:"created"`
}

type JiraIssueFields struct {
	Summary     string      `json:"summary"`
	Project     JiraProject `json:"project"`
	Description string      `json:"description,omitempty"`
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
	Status struct {
		Name string `json:"name"`
	} `json:"status"`
	Comment struct {
		Comments   []JiraComment `json:"comments"`
		MaxResults int32         `json:"maxResults"`
		Total      int32         `json:"total"`
		StartAt    int32         `json:"startAt"`
	} `json:"comment"`
	Labels []string `json:"labels"`
}

type JiraUser struct {
	AccountId    string            `json:"accountId"`
	Active       bool              `json:"active"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	DisplayName  string            `json:"displayName"`
	EmailAddress string            `json:"emailAddress"`
	Locale       string            `json:"locale"`
	Self         string            `json:"self"`
	TimeZone     string            `json:"timeZone"`
}

type JiraIssueTransition struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	To   struct {
		StatusUrl string `json:"self"`
		StatusId  string `json:"id"`
		Name      string `json:"name"`
	} `json:"to"`
}

type JiraIssueStatus struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JiraLabelsResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []string `json:"values"`
}

type JiraApiCredentials struct {
	Host   string
	ApiKey string
}

type httpJiraApi struct {
	client  *http.Client
	restUrl *url.URL
}

func NewJiraApi(apiUrl string, username string, token string) (JiraApi, error) {
	baseUrl, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatalln(err)
	}
	basicToken := base64.StdEncoding.EncodeToString([]byte(username + ":" + token))
	return &httpJiraApi{
		client: &http.Client{
			Transport: &authInterceptor{core: http.DefaultTransport, token: basicToken},
		},
		restUrl: baseUrl,
	}, nil
}

func (api *httpJiraApi) Close() {
	api.client.CloseIdleConnections()
}
