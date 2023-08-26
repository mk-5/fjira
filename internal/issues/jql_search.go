package issues

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"strings"
)

const (
	MaxJqlLines     = 1000
	DefaultJqlQuery = "created >= -30d order by created DESC"
	MaxJqlLength    = 2000
)

type jqlSearchView struct {
	app.View
	api        jira.Api
	bottomBar  *app.ActionBar
	fuzzyFind  *app.FuzzyFind
	jqlStorage *jqlStorage
}

func NewJqlSearchView(api jira.Api) app.View {
	bottomBar := ui.CreateBottomActionBarWithItems([]ui.NavItemConfig{
		ui.NavItemConfig{Action: ui.ActionNew, Text1: ui.MessageNew, Text2: "[F1]", Key: tcell.KeyF1},
	})
	bottomBar.AddItem(ui.NewDeleteItem())
	bottomBar.AddItem(ui.NewCancelBarItem())
	return &jqlSearchView{
		api:        api,
		bottomBar:  bottomBar,
		jqlStorage: &jqlStorage{},
	}
}

func (view *jqlSearchView) Init() {
	go view.startJqlFuzzyFind()
	go view.handleBottomBarActions()
}

func (view *jqlSearchView) Destroy() {
	view.bottomBar.Destroy()
}

func (view *jqlSearchView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.bottomBar.Draw(screen)
}

func (view *jqlSearchView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *jqlSearchView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.bottomBar.Resize(screenX, screenY)
}

func (view *jqlSearchView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *jqlSearchView) startJqlFuzzyFind() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	jqls, err := view.jqlStorage.readAll()
	if err != nil {
		app.Error(err.Error())
		return
	}
	jqls = append(jqls, DefaultJqlQuery)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageJqlFuzzyFind, jqls)
	view.fuzzyFind.MarginBottom = 1
	app.GetApp().Loading(false)
	if jql := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		query := view.fuzzyFind.GetQuery()
		if jql.Index < 0 && strings.TrimSpace(query) == "" {
			// do nothing
			return
		}
		if jql.Index < 0 {
			// do nothing, restart view
			app.GetApp().SetView(NewJqlSearchView(view.api))
			return
		}
		view.fuzzyFind = nil
		app.GoTo("issues-search-jql", jqls[jql.Index], view.api)
		//GoIntoIssuesSearchForJql(jqls[jql.Index], view.api)
	}
}

func (view *jqlSearchView) handleBottomBarActions() {
	for {
		action := <-view.bottomBar.Action
		switch action {
		case ui.ActionCancel:
			view.cancel()
			return
		case ui.ActionNew:
			view.newJql()
			return
		case ui.ActionDelete:
			go view.confirmJqlDelete(view.fuzzyFind.GetSelectedItem())
		}
	}
}

func (view *jqlSearchView) cancel() {
	if view.fuzzyFind != nil {
		app.GetApp().Quit()
	}
	if view.fuzzyFind == nil {
		app.GetApp().SetView(NewJqlSearchView(view.api))
	}
}

func (view *jqlSearchView) newJql() {
	app.GetApp().SetView(ui.NewTextWriterView(&ui.TextWriterArgs{
		Header: ui.MessageTypeJqlAndSave,
		GoBack: func() {
			app.GetApp().SetView(NewJqlSearchView(view.api))
		},
		TextConsumer: func(s string) {
			err := view.jqlStorage.addNew(s)
			if err != nil {
				app.Error(err.Error())
				return
			}
			app.Success(ui.MessageJqlAddSuccess)
		},
		MaxLength: MaxJqlLength,
	}))
}

func (view *jqlSearchView) confirmJqlDelete(jql string) {
	message := fmt.Sprintf(ui.MessageJqlRemoveConfirm, jql)
	view.fuzzyFind = nil
	app.GetApp().ClearNow()
	view.bottomBar.Clear()
	view.bottomBar.AddItem(ui.NewYesBarItem())
	view.bottomBar.AddItem(ui.NewCancelBarItem())
	createNewJql := app.Confirm(app.GetApp(), message)
	switch createNewJql {
	case true:
		err := view.jqlStorage.remove(jql)
		if err != nil {
			app.Error(fmt.Sprintf(ui.MessageCannotAddNewJql, err))
			return
		}
		app.Success(ui.MessageJqlRemoveSuccess)
		app.GetApp().SetView(NewJqlSearchView(view.api))
	case false:
		app.GetApp().SetView(NewJqlSearchView(view.api))
	}
}
