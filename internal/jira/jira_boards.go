package jira

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mk-5/fjira/internal/app"
)

type BoardItem struct {
	Id   int    `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type GetAllSprintsQueryParams struct {
	State string `url:"state,omitempty"`
}

type BoardsResponse struct {
	MaxResults int         `json:"maxResults"`
	StartAt    int         `json:"startAt"`
	Total      int         `json:"total"`
	IsLast     bool        `json:"isLast"`
	Values     []BoardItem `json:"values"`
}

type SprintsResponse struct {
	MaxResults int          `json:"maxResults"`
	StartAt    int          `json:"startAt"`
	Total      int          `json:"total"`
	IsLast     bool         `json:"isLast"`
	Values     []SprintItem `json:"values"`
}

type SprintItem struct {
	Id            int        `json:"id"`
	Self          string     `json:"self"`
	State         string     `json:"state"`
	Name          string     `json:"name"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	CompleteDate  *time.Time `json:"completeDate"`
	OriginBoardId int        `json:"originBoardId"`
	Goal          string     `json:"goal"`
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

const (
	FindAllBoardsUrl          = "/rest/agile/1.0/board"
	FindBoardConfigurationUrl = "/rest/agile/1.0/board/%d/configuration"
	FindBoardSprintsUrl       = "/rest/agile/1.0/board/%d/sprint"
	FindBoardSprintsIssuesUrl = "/rest/agile/1.0/board/%d/sprint/%d/issue"
)

type findBoardsQueryParams struct {
	ProjectKeyOrId string `url:"projectKeyOrId"`
	StartAt        int    `url:"startAt"`
}

type boardsSearchQueryParams struct {
	MaxResults int32  `url:"maxResults"`
	Fields     string `url:"fields"`
	StartAt    int32  `url:"startAt"`
}

type boardsSearchResponse struct {
	Total      int32   `json:"total"`
	MaxResults int32   `json:"maxResults"`
	Issues     []Issue `json:"issues"`
	IsLast     bool    `json:"isLast"`
}

func (api *httpApi) FindBoards(projectKeyOrId string) ([]BoardItem, error) {
	params := &findBoardsQueryParams{ProjectKeyOrId: projectKeyOrId, StartAt: 0}
	var boards []BoardItem
	for {
		resultBytes, err := api.jiraRequest("GET", FindAllBoardsUrl, params, nil)
		if err != nil {
			return nil, err
		}
		var result BoardsResponse
		err = json.Unmarshal(resultBytes, &result)
		if err != nil {
			return nil, err
		}
		if cap(boards) == 0 {
			boards = make([]BoardItem, 0, result.Total)
		}
		boards = append(boards, result.Values...)

		if result.IsLast {
			break
		}
		params.StartAt += result.MaxResults
	}
	return boards, nil
}

func (api *httpApi) GetBoardConfiguration(boardId int) (*BoardConfiguration, error) {
	resultBytes, err := api.jiraRequest("GET", fmt.Sprintf(FindBoardConfigurationUrl, boardId), &nilParams{}, nil)
	if err != nil {
		return nil, err
	}
	var result BoardConfiguration
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (api *httpApi) GetBoardSprints(boardId int) ([]SprintItem, error) {
	params := &GetAllSprintsQueryParams{State: "active,future"}
	resultBytes, err := api.jiraRequest("GET", fmt.Sprintf(FindBoardSprintsUrl, boardId), params, nil)
	if err != nil {
		return nil, err
	}
	var result SprintsResponse
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return result.Values, nil
}

func (api *httpApi) GetBoardSprintIssues(boardId int, sprintId int, page int32, pageSize int32) ([]Issue, int32, int32, error) {
	params := &boardsSearchQueryParams{
		MaxResults: pageSize,
		StartAt:    page * pageSize,
		Fields:     "id,key,summary,issuetype,project,reporter,status,assignee",
	}
	body, err := api.jiraRequest("GET", fmt.Sprintf(FindBoardSprintsIssuesUrl, boardId, sprintId), params, nil)
	if err != nil {
		return nil, -1, pageSize, err
	}
	var sResponse boardsSearchResponse
	if err := json.Unmarshal(body, &sResponse); err != nil {
		app.Error(err.Error())
		return nil, -1, pageSize, SearchDeserializeErr
	}
	return sResponse.Issues, sResponse.Total, sResponse.MaxResults, err
}
