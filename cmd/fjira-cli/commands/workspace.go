package commands

import (
	"github.com/mk-5/fjira/internal/fjira"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func GetWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Switch to a different workspace",
		Run: func(cmd *cobra.Command, args []string) {
			edit, _ := cmd.Flags().GetString("edit")
			n, _ := cmd.Flags().GetString("new")
			var s *workspaces.WorkspaceSettings
			var err error
			if edit != "" {
				s, err = fjira.EditWorkspaceAndReadSettings(os.Stdin, edit)
				if err != nil {
					log.Println(err)
					log.Fatalln(fjira.ErrInstallFailed.Error())
				}
			} else if n != "" {
				s, err = fjira.Install(n)
				if err != nil {
					log.Println(err)
					log.Fatalln(fjira.ErrInstallFailed.Error())
				}
			} else {
				s = cmd.Context().Value(CtxWorkspaceSettings).(*workspaces.WorkspaceSettings)
			}
			f := fjira.CreateNewFjira(s)
			defer f.Close()
			f.Run(&fjira.CliArgs{
				WorkspaceSwitch: true,
			})
		},
	}
	cmd.Flags().String("edit", "", "Edit workspace")
	cmd.Flags().String("new", "", "Create a new workspace")
	return cmd
}
