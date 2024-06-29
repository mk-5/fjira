package app

import "testing"

func TestCopyIssue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run CopyIssue func without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyIssue("ABC-123")
		})
	}
}
