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
    fjira [command]
    fjira [flags]
    fjira [jira-issue] [flags]

Available Commands:
    workspace               Switch fjira workspace
    help               	    Help
    version                 Show version

Flags:
    -p, --project           Search for issues withing project, example: -p GEN.
    -w, --workspace         Use different fjira workspace, example: -w myworkspace
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
	flag.StringVar(&projectId, "project", "", "Jira Project ID")
	flag.StringVar(&projectId, "p", "", "Jira Project ID")
	flag.StringVar(&workspace, "workspace", "", "Fjira workspace")
	flag.StringVar(&workspace, "w", "", "Fjira workspace")
	flag.Parse()

	issueRegExp := regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	if len(os.Args) > 1 && issueRegExp.MatchString(os.Args[1]) {
		return fjira.CliArgs{
			IssueKey:  os.Args[1],
			Workspace: workspace,
		}
	}
	if len(os.Args) == 2 && os.Args[1] == "workspace" {
		return fjira.CliArgs{
			SwitchWorkspace: true,
		}
	}
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println(fmt.Sprintf("fjira version: %s", version))
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "help" {
		fmt.Println(usage)
		os.Exit(0)
	}
	return fjira.CliArgs{
		ProjectId:       projectId,
		Workspace:       workspace,
		SwitchWorkspace: false,
	}
}
