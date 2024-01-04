package jira

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
)

type Api interface {
	Search(query string) ([]Issue, int32, error)
	SearchJql(query string) ([]Issue, error)
	SearchJqlPageable(query string, page int32, pageSize int32) ([]Issue, int32, int32, error)
	FindUsers(project string) ([]User, error)
	FindUsersWithQuery(project string, query string) ([]User, error)
	FindProjects() ([]Project, error)
	FindLabels(issue *Issue, query string) ([]string, error)
	AddLabel(issueId string, label string) error
	FindProject(projectKey string) (*Project, error)
	FindTransitions(issueId string) ([]IssueTransition, error)
	FindProjectStatuses(projectId string) ([]IssueStatus, error)
	DoTransition(issueId string, transition *IssueTransition) error
	DoAssignee(issueId string, user *User) error
	GetIssueDetailed(issueId string) (*Issue, error)
	DoComment(issueId string, commentBody string) error
	FindBoards(projectKeyOrId string) ([]BoardItem, error)
	GetBoardConfiguration(boardId int) (*BoardConfiguration, error)
	GetFilter(filterId string) (*Filter, error)
	GetMyFilters() ([]Filter, error)
	Close()
	GetApiUrl() string

	IsJiraServer() bool
}

type ApiCredentials struct {
	Host   string
	ApiKey string
}

type JiraTokenType string

const (
	ApiToken      JiraTokenType = "api token"
	PersonalToken JiraTokenType = "personal token"
)

type httpApi struct {
	apiUrl    string
	tokenType JiraTokenType
	client    *http.Client
	restUrl   *url.URL
}

func NewApi(apiUrl string, username string, token string, tokenType JiraTokenType) (Api, error) {
	baseUrl, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatalln(err)
	}
	var authToken string
	var authType AuthType
	switch tokenType {
	case PersonalToken:
		authToken = token
		authType = Bearer
	default:
		authToken = base64.StdEncoding.EncodeToString([]byte(username + ":" + token))
		authType = Basic
	}
	return &httpApi{
		apiUrl:    apiUrl,
		tokenType: tokenType,
		client: &http.Client{
			Transport: &authInterceptor{core: defaultHttpTransport, token: authToken, authType: authType},
		},
		restUrl: baseUrl,
	}, nil
}

func (api *httpApi) GetApiUrl() string {
	return api.apiUrl
}

func (api *httpApi) IsJiraServer() bool {
	// for now - just a stupid impl like this
	return api.tokenType == PersonalToken
}

func (api *httpApi) Close() {
	api.client.CloseIdleConnections()
}
