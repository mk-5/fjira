package jira

import (
	"encoding/json"
	"fmt"
)

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
