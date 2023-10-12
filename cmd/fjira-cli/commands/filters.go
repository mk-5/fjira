package commands

import (
	"github.com/mk-5/fjira/internal/fjira"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/spf13/cobra"
)

func GetFiltersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "filters",
		Short: "Search using Jira filters",
		Run: func(cmd *cobra.Command, args []string) {
			s := cmd.Context().Value(CtxWorkspaceSettings).(*workspaces.WorkspaceSettings)
			f := fjira.CreateNewFjira(s)
			defer f.Close()
			f.Run(&fjira.CliArgs{
				FiltersMode: true,
			})
		},
	}
}
