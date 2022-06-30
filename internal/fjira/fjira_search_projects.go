package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

type fjiraSearchProjectsView struct {
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
}

func NewProjectsSearchView() *fjiraSearchProjectsView {
	bottomBar := CreateProjectBottomBar()
	return &fjiraSearchProjectsView{
		bottomBar: bottomBar,
		topBar:    CreateProjectsTopBar(),
	}
}

func (view *fjiraSearchProjectsView) Init() {
	app.GetApp().LoadingWithText(true, MessageSearchProjectsLoading)
	go view.runProjectsFuzzyFind()
}

func (view *fjiraSearchProjectsView) Destroy() {
}

func (view *fjiraSearchProjectsView) Draw(screen tcell.Screen) {
	//view.bottomBar.Draw(screen)
	//view.topBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *fjiraSearchProjectsView) Update() {
	//view.bottomBar.Update()
	//view.topBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraSearchProjectsView) Resize(screenX, screenY int) {
	//view.bottomBar.Resize(screenX, screenY)
	//view.topBar.Resize(screenX, screenY)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
}

func (view *fjiraSearchProjectsView) HandleKeyEvent(ev *tcell.EventKey) {
	//view.topBar.HandleKeyEvent(ev)
	//view.bottomBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraSearchProjectsView) findProjects() []jira.JiraProject {
	api, _ := GetApi()
	projects, err := api.FindProjects()
	if err != nil {
		app.Error(err.Error())
	}
	return projects
}

func (view *fjiraSearchProjectsView) runProjectsFuzzyFind() {
	projects := view.findProjects()
	formatter, _ := GetFormatter()
	projectsString := formatter.formatJiraProjects(projects)
	view.fuzzyFind = app.NewFuzzyFind(MessageSelectProject, projectsString)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	app.GetApp().ClearNow()
	select {
	case chosen := <-view.fuzzyFind.Complete:
		app.GetApp().ClearNow()
		if chosen.Index < 0 {
			app.GetApp().Quit()
			return
		}
		chosenProject := projects[chosen.Index]
		go goIntoIssuesSearch(&chosenProject)
	}
}
