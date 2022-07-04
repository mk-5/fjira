package main

import (
	"flag"
	"fmt"
	"github.com/mk5/fjira/internal/fjira"
	"log"
	"os"
	"regexp"
)

const (
	usage = `Usage:
    fjira
    fjira [command]
    fjira [command] [flags]
    fjira [flags]
    fjira [jira-issue] [flags]

Available Commands:
    workspace               Switch fjira workspace
    help               	    Help
    version                 Show version

Flags:
    -p, --project             Open project issues search directly from CLI, example: -p GEN.
    -w, --workspace           Use different fjira workspace without switching it globally, example: -w myworkspace
    -nw, --new-workspace      Create new workspace, example: fjira --new-workspace=abc
`
)

var (
	version = "dev"
)

func main() {
	args := parseCliArgs()
	settings, err := fjira.Install(args.Workspace)
	if err != nil {
		log.Println(err)
		log.Fatalln(fjira.InstallFailedErr.Error())
	}
	f := fjira.CreateNewFjira(settings)
	defer f.Close()
	f.Run(&args)
}

func parseCliArgs() fjira.CliArgs {
	flag.Usage = func() {
		fmt.Print(usage)
	}
	var projectId string
	var workspace string
	var newWorkspace string
	flag.StringVar(&projectId, "project", "", "Jira Project Key")
	flag.StringVar(&projectId, "p", "", "Jira Project Key")
	flag.StringVar(&workspace, "workspace", "", "Fjira workspace")
	flag.StringVar(&workspace, "w", "", "Fjira workspace")
	flag.StringVar(&newWorkspace, "new-workspace", "", "New workspace name")
	flag.StringVar(&newWorkspace, "nw", "", "New workspace name")
	flag.Parse()

	issueRegExp := regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	if len(os.Args) > 1 && issueRegExp.MatchString(os.Args[1]) {
		return fjira.CliArgs{
			IssueKey:  os.Args[1],
			Workspace: workspace,
		}
	}
	if newWorkspace != "" {
		return fjira.CliArgs{
			Workspace:       newWorkspace,
			SwitchWorkspace: false,
		}
	}
	if len(os.Args) >= 2 && os.Args[1] == "workspace" {
		return fjira.CliArgs{
			SwitchWorkspace: true,
		}
	}
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Printf("fjira version: %s", version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "help" {
		fmt.Print(usage)
		os.Exit(0)
	}
	return fjira.CliArgs{
		ProjectId:       projectId,
		Workspace:       workspace,
		SwitchWorkspace: false,
	}
}
