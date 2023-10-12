package filters

import (
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFilters(t *testing.T) {
	type args struct {
		filters []jira.Filter
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"should format filters",
			args{filters: []jira.Filter{{Name: "Test ABC"}}},
			[]string{"Test ABC"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FormatFilters(tt.args.filters), "FormatFilters(%v)", tt.args.filters)
		})
	}
}
