package filters

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

func RegisterGoTo() {
	app.RegisterGoto("filters", func(args ...interface{}) {
		defer app.GetApp().PanicRecover()
		api := args[0].(jira.Api)
		view := NewFiltersView(api)
		app.GetApp().SetView(view)
	})
}
