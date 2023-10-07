package commands

import (
	"github.com/mk-5/fjira/internal/fjira"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/spf13/cobra"
)

func GetIssueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "[issueKey]",
		Short: "Open a Jira issue directly from the CLI",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s := cmd.Context().Value(CtxWorkspaceSettings).(*workspaces.WorkspaceSettings)
			issueKey := args[0]
			f := fjira.CreateNewFjira(s)
			defer f.Close()
			f.Run(&fjira.CliArgs{
				IssueKey: issueKey,
			})
		},
	}
}
