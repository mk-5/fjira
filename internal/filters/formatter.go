package filters

import "github.com/mk-5/fjira/internal/jira"

func FormatFilters(filters []jira.Filter) []string {
	s := make([]string, 0, len(filters))
	for _, filter := range filters {
		s = append(s, filter.Name)
	}
	return s
}
