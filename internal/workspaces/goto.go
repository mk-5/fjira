package workspaces

import "github.com/mk-5/fjira/internal/app"

func RegisterGoTo() {
	app.RegisterGoto("workspaces-switch", func(args ...interface{}) {
		switchWorkspaceView := NewSwitchWorkspaceView()
		app.GetApp().SetView(switchWorkspaceView)
	})
}
