package projects

import (
	"fmt"
	"github.com/mk-5/fjira/internal/jira"
)

func FormatJiraProject(project *jira.Project) string {
	return fmt.Sprintf("[%s] %s", project.Key, project.Name)
}

func FormatJiraProjects(projects []jira.Project) []string {
	formatted := make([]string, 0, len(projects))
	for _, project := range projects {
		formatted = append(formatted, FormatJiraProject(&project))
	}
	return formatted
}
