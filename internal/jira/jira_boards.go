package jira

import (
	"encoding/json"
	"fmt"
)

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

const (
	FindAllBoardsUrl          = "/rest/agile/1.0/board"
	FindBoardConfigurationUrl = "/rest/agile/1.0/board/%d/configuration"
)

type findBoardsQueryParams struct {
	ProjectKeyOrId string `url:"projectKeyOrId"`
	StartAt        int    `url:"startAt"`
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
