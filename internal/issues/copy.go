package issues

import (
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

func CopyIssue(i *jira.Issue) {
	app.CopyIssue(i.Key)
	app.Success(fmt.Sprintf(ui.MessageCopyIssueSuccess, i.Key))
	return
}
