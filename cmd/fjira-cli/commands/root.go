package commands

import (
	"context"
	"errors"
	"github.com/mk-5/fjira/internal/fjira"
	"github.com/mk-5/fjira/internal/workspaces"
	"github.com/spf13/cobra"
	"regexp"
)

type CtxVarWorkspaceSettings string

const (
	CtxWorkspaceSettings CtxVarWorkspaceSettings = "workspace-settings"
)

var InvalidIssueKeyFormatErr = errors.New("invalid issue key format")

func GetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fjira",
		Short: "A fuzzy jira tui application",
		Long: `Fjira is a powerful terminal user interface (TUI) application designed to streamline your Jira workflow.
With its fuzzy-find capabilities, it simplifies the process of searching and accessing Jira issues, 
making it easier than ever to locate and manage your tasks and projects efficiently.
Say goodbye to manual searching and hello to increased productivity with fjira.`,
		Args: cobra.MaximumNArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// it's initializing fjira before every command
			s, err := fjira.Install("")
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), CtxWorkspaceSettings, s))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// run Issue command if issueKey provided via cli argument
			if len(args) == 1 {
				issueRegExp := regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
				issueKey := args[0]
				if !issueRegExp.MatchString(issueKey) {
					return InvalidIssueKeyFormatErr
				}
				issueCmd := GetIssueCmd()
				issueCmd.SetArgs([]string{issueKey})
				return issueCmd.ExecuteContext(cmd.Context())
			}
			projectKey, _ := cmd.Flags().GetString("project")
			s := cmd.Context().Value(CtxWorkspaceSettings).(*workspaces.WorkspaceSettings)
			f := fjira.CreateNewFjira(s)
			defer f.Close()
			f.Run(&fjira.CliArgs{
				ProjectId: projectKey,
			})
			return nil
		},
	}
	cmd.AddCommand(&cobra.Command{Use: "", Short: "Open a fuzzy finder for projects as a default action"})
	cmd.Flags().StringP("project", "p", "", "Open a project directly from CLI")
	return cmd
}
