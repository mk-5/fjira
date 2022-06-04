package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/mk5/fjira/pkg/fjira"
)

const (
	usage = `Usage:
    fjira JIRA-TICKET
    fjira [OPTIONS]

Optional options:
    -p, --project               Search for issues withing project, example: GEN.
    -i, --issue                 Open Jira Issue, example: GEN-123.
    -w, --workspace             Use fjira workspace, example: my-workspace2
`
)

func main() {
	f := fjira.CreateNewFjira(nil)
	defer f.Close()
	args := parseCliArgs()
	errors := f.Install(args.Workspace)
	if errors != nil {
		for _, err := range errors {
			log.Println(err.Error())
		}
		log.Fatalln(fjira.InstallFailedErr.Error())
	}
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
		ProjectId: projectId,
		IssueKey:  issueKey,
		Workspace: workspace,
	}
}
