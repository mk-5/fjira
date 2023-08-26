package boards

import "github.com/mk-5/fjira/internal/jira"

func FormatJiraBoards(boards []jira.BoardItem) []string {
	formatted := make([]string, 0, len(boards))
	for _, board := range boards {
		formatted = append(formatted, board.Name)
	}
	return formatted
}
