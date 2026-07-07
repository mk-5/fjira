package projects

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

type searchProjectsView struct {
	api       jira.Api
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
}

func NewProjectsSearchView(api jira.Api) app.View {
	bottomBar := ui.CreateBottomActionBar(ui.MessageProjectLabel, app.ActionBarLabel(""))
	topBar := ui.CreateTopActionBar(ui.MessageProjectLabel, app.ActionBarLabel(""))
	return &searchProjectsView{
		api:       api,
		bottomBar: bottomBar,
		topBar:    topBar,
	}
}

func (view *searchProjectsView) Init() {
	app.GetApp().LoadingWithText(true, ui.MessageSearchProjectsLoading)
	go view.runProjectsFuzzyFind()
}

func (view *searchProjectsView) Destroy() {
}

func (view *searchProjectsView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *searchProjectsView) Update() {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *searchProjectsView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
}

func (view *searchProjectsView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *searchProjectsView) findProjects() []jira.Project {
	projects, err := view.api.FindProjects()
	if err != nil {
		app.Error(err.Error())
	}
	if err == nil && len(projects) == 0 && isAtlassianCloud(view.api.GetApiUrl()) {
		app.Error(ui.MessageNoProjectsTokenHint)
	}
	return projects
}

func (view *searchProjectsView) reopen() {
	app.GoTo("projects", view.api)
}

func (view *searchProjectsView) runProjectsFuzzyFind() {
	defer app.GetApp().PanicRecover()
	projects := view.findProjects()
	projects = append(projects, jira.Project{Id: ui.MessageAll, Name: ui.MessageAll, Key: ui.MessageAll})
	projectsString := FormatJiraProjects(projects)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageSelectProject, projectsString)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	app.GetApp().ClearNow()
	if chosen := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if chosen.Index < 0 {
			app.GetApp().Quit()
			return
		}
		chosenProject := projects[chosen.Index]
		app.GoTo("issues-search", chosenProject.Id, view.reopen, view.api)
	}
}

// isAtlassianCloud reports whether the Jira REST URL points at an Atlassian
// Cloud site. On Cloud, an expired/wrong-scope API token returns 200 with an
// empty project list instead of 401 — so an empty result is the strongest
// available signal of an auth problem there. The token-management URL in the
// hint is Cloud-specific, so we only show it for Cloud workspaces.
func isAtlassianCloud(apiUrl string) bool {
	return strings.Contains(apiUrl, ".atlassian.net")
}
