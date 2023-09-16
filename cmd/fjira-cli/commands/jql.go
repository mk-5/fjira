package commands

import (
	"github.com/mk-5/fjira/internal/fjira"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/spf13/cobra"
)

func GetJqlCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "jql",
		Short: "Search using custom-jql",
		Run: func(cmd *cobra.Command, args []string) {
			s := cmd.Context().Value(CtxWorkspaceSettings).(*workspaces.WorkspaceSettings)
			f := fjira.CreateNewFjira(s)
			defer f.Close()
			f.Run(&fjira.CliArgs{
				JqlMode: true,
			})
		},
	}
}
