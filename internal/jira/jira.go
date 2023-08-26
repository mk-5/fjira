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
	FindProjects() ([]Project, error)
	FindLabels(issue *Issue, query string) ([]string, error)
	AddLabel(issueId string, label string) error
	FindProject(projectKey string) (*Project, error)
	FindTransitions(issueId string) ([]IssueTransition, error)
	FindProjectStatuses(projectId string) ([]IssueStatus, error)
	DoTransition(issueId string, transition *IssueTransition) error
	DoAssignee(issueId string, accountId string) error
	GetIssueDetailed(issueId string) (*Issue, error)
	DoComment(issueId string, commentBody string) error
	FindBoards(projectKeyOrId string) ([]BoardItem, error)
	GetBoardConfiguration(boardId int) (*BoardConfiguration, error)
	Close()
	GetApiUrl() string
}

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"Fields"`
	Id     string      `json:"id"`
}

type Comment struct {
	Author  User   `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type Status struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	Labels []string `json:"labels"`
}

type User struct {
	AccountId    string            `json:"accountId"`
	Active       bool              `json:"active"`
	AvatarUrls   map[string]string `json:"avatarUrls"`
	DisplayName  string            `json:"displayName"`
	EmailAddress string            `json:"emailAddress"`
	Locale       string            `json:"locale"`
	Self         string            `json:"self"`
	TimeZone     string            `json:"timeZone"`
}

type IssueTransition struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	To   struct {
		StatusUrl string `json:"self"`
		StatusId  string `json:"id"`
		Name      string `json:"name"`
	} `json:"to"`
}

type IssueStatus struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LabelsSuggestionsResponseBody struct {
	Token       string `json:"token"`
	Suggestions []struct {
		Label string `json:"label"`
		Html  string `json:"html"`
	} `json:"suggestions"`
}

type BoardItem struct {
	Id   int    `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type BoardsResponse struct {
	MaxResults int         `json:"maxResults"`
	StartAt    int         `json:"startAt"`
	Total      int         `json:"total"`
	IsLast     bool        `json:"isLast"`
	Values     []BoardItem `json:"values"`
}

type BoardConfiguration struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Self     string `json:"self"`
	Location struct {
		Type string `json:"type"`
		Key  string `json:"key"`
		Id   string `json:"id"`
		Self string `json:"self"`
		Name string `json:"name"`
	} `json:"location"`
	Filter struct {
		Id   string `json:"id"`
		Self string `json:"self"`
	} `json:"filter"`
	SubQuery struct {
		Query string `json:"query"`
	} `json:"subQuery"`
	ColumnConfig struct {
		Columns []struct {
			Name     string `json:"name"`
			Statuses []struct {
				Id   string `json:"id"`
				Self string `json:"self"`
			} `json:"statuses"`
		} `json:"columns"`
		ConstraintType string `json:"constraintType"`
	} `json:"columnConfig"`
	Ranking struct {
		RankCustomFieldId int `json:"rankCustomFieldId"`
	} `json:"ranking"`
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
	apiUrl  string
	client  *http.Client
	restUrl *url.URL
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
		apiUrl: apiUrl,
		client: &http.Client{
			Transport: &authInterceptor{core: http.DefaultTransport, token: authToken, authType: authType},
		},
		restUrl: baseUrl,
	}, nil
}

func (api *httpApi) GetApiUrl() string {
	return api.apiUrl
}

func (api *httpApi) Close() {
	api.client.CloseIdleConnections()
}
