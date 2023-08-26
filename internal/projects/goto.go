package projects

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

func RegisterGoto() {
	app.RegisterGoto("projects", func(args ...interface{}) {
		api := args[0].(jira.Api)
		projectsView := NewProjectsSearchView(api)
		app.GetApp().SetView(projectsView)
	})
}
