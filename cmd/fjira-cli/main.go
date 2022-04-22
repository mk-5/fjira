package main

import (
	"flag"
	"fmt"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/pkg/fjira"
	"log"
	"os"
	"regexp"
)

const (
	usage = `Usage:
    fjira JIRA-TICKET
    fjira [OPTIONS]

Optional options:
    -p, --project               Search for issues withing project, example: GEN.
    -i, --issue                 Open Jira Issue, example: GEN-123.

`
)

func main() {
	f := fjira.CreateNewFjira()
	defer f.Close()
	errors := f.Install()
	if errors != nil {
		for _, err := range errors {
			log.Println(err.Error())
		}
		log.Fatalln(fjira.InstallFailedErr.Error())
	}
	app.StartCli()
	args := parseCliArgs()
	f.Run(&args)
}

func parseCliArgs() fjira.CliArgs {
	flag.Usage = func() {
		fmt.Println(usage)
	}
	issueRegExp := regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	if len(os.Args) > 1 && issueRegExp.MatchString(os.Args[1]) {
		return fjira.CliArgs{
			IssueKey: os.Args[1],
		}
	}
	var projectId string
	var issueKey string
	//var help bool
	flag.StringVar(&projectId, "project", "", "Jira Project ID")
	flag.StringVar(&projectId, "p", "", "Jira Project ID")
	flag.StringVar(&issueKey, "issue", "", "Jira Issue Key")
	flag.StringVar(&issueKey, "i", "", "Jira Issue Key")
	flag.Parse()
	return fjira.CliArgs{
		ProjectId: projectId,
		IssueKey:  issueKey,
	}
}
