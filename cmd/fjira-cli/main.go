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
    fjira JIRA-TICKET
    fjira workspace
    fjira version
    fjira [OPTIONS]

Optional options:
    -p, --project               Search for issues withing project, example: GEN.
    -i, --issue                 Open Jira Issue, example: GEN-123.
    -w, --workspace             Use fjira workspace, example: myworkspace
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
	issueRegExp := regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	if len(os.Args) > 1 && issueRegExp.MatchString(os.Args[1]) {
		return fjira.CliArgs{
			IssueKey: os.Args[1],
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
	var projectId string
	var issueKey string
	var workspace string
	//var help bool
	flag.StringVar(&projectId, "project", "", "Jira Project ID")
	flag.StringVar(&projectId, "p", "", "Jira Project ID")
	flag.StringVar(&issueKey, "issue", "", "Jira Issue Key")
	flag.StringVar(&issueKey, "i", "", "Jira Issue Key")
	flag.StringVar(&workspace, "workspace", "", "Fjira workspace")
	flag.StringVar(&workspace, "w", "", "Fjira workspace")
	flag.Parse()
	return fjira.CliArgs{
		ProjectId:       projectId,
		IssueKey:        issueKey,
		Workspace:       workspace,
		SwitchWorkspace: false,
	}
}
