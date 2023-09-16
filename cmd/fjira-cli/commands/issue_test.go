package commands

import (
	"context"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetIssueCmd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create&execute IssueCmd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			cmd := GetIssueCmd()

			// then
			assert.NotNil(t, cmd)

			// and when
			cmd.SetArgs([]string{"ABC-123"})
			var err error
			go func() {
				err = cmd.ExecuteContext(context.WithValue(context.TODO(), CtxWorkspaceSettings, &workspaces.WorkspaceSettings{}))
			}() //nolint:errcheck
			for app.GetApp() == nil {
				<-time.After(50 * time.Millisecond)
			}

			// then
			assert.Nil(t, err)
		})
	}
}
