package main

import (
	"fmt"
	"github.com/mk-5/fjira/cmd/fjira-cli/commands"
	"os"
)

var (
	version = "dev"
)

func main() {
	initCli()
}

func initCli() {
	rootCmd := commands.GetRootCmd()

	rootCmd.AddCommand(commands.GetIssueCmd())
	rootCmd.AddCommand(commands.GetWorkspaceCmd())
	rootCmd.AddCommand(commands.GetJqlCmd())
	rootCmd.AddCommand(commands.GetFiltersCmd())
	rootCmd.AddCommand(commands.GetVersionCmd(version))

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
