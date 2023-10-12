package filters

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"strings"
)

type filtersSearchView struct {
	app.View
	api       jira.Api
	bottomBar *app.ActionBar
	fuzzyFind *app.FuzzyFind
}

func NewFiltersView(api jira.Api) app.View {
	bottomBar := ui.CreateBottomLeftBar()
	bottomBar.AddItem(ui.NewCancelBarItem())
	return &filtersSearchView{
		api:       api,
		bottomBar: bottomBar,
	}
}

func (view *filtersSearchView) Init() {
	go view.startFiltersFuzzyFind()
	go view.handleBottomBarActions()
}

func (view *filtersSearchView) Destroy() {
	view.bottomBar.Destroy()
}

func (view *filtersSearchView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.bottomBar.Draw(screen)
}

func (view *filtersSearchView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *filtersSearchView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.bottomBar.Resize(screenX, screenY)
}

func (view *filtersSearchView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *filtersSearchView) startFiltersFuzzyFind() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	filters, err := view.api.GetMyFilters()
	if err != nil {
		app.Error(err.Error())
		return
	}
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageSelectFilter, FormatFilters(filters))
	view.fuzzyFind.MarginBottom = 1
	app.GetApp().Loading(false)
	if chosen := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		query := view.fuzzyFind.GetQuery()
		if chosen.Index < 0 && strings.TrimSpace(query) == "" {
			// do nothing
			return
		}
		if chosen.Index >= 0 {
			filter := filters[chosen.Index]
			app.GoTo("issues-search-jql", filter.JQL, view.reopen, view.api)
		}
	}
}

func (view *filtersSearchView) handleBottomBarActions() {
	for {
		action := <-view.bottomBar.Action
		switch action {
		case ui.ActionCancel:
			view.cancel()
			return
		}
	}
}

func (view *filtersSearchView) cancel() {
	app.GetApp().Quit()
}

func (view *filtersSearchView) reopen() {
	app.GoTo("filters", view.api)
}
