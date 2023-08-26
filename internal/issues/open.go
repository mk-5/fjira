package issues

import (
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

func OpenIssueInBrowser(i *jira.Issue, api jira.Api) {
	jiraUrl := api.GetApiUrl()
	app.OpenLink(fmt.Sprintf("%s/browse/%s", jiraUrl, i.Key))
}
