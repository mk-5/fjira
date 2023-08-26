package labels

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

func RegisterGoTo() {
	app.RegisterGoto("labels-add", func(args ...interface{}) {
		issue := args[0].(*jira.Issue)
		goBackFn := args[1].(func())
		api := args[2].(jira.Api)

		commentView := NewAddLabelView(issue, goBackFn, api)
		app.GetApp().SetView(commentView)
	})
}
