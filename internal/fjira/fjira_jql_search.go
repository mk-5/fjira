package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"strings"
)

const (
	MaxJqlLines     = 1000
	DefaultJqlQuery = "created >= -30d order by created DESC"
	MaxJqlLength    = 2000
)

type fjiraJqlSearchView struct {
	app.View
	bottomBar  *app.ActionBar
	fuzzyFind  *app.FuzzyFind
	jqlStorage *jqlStorage
}

func NewJqlSearchView() *fjiraJqlSearchView {
	bottomBar := CreateBottomLeftBar()
	bottomBar.AddItem(NewNewJqlItem())
	bottomBar.AddItem(NewDeleteItem())
	bottomBar.AddItem(NewCancelBarItem())
	return &fjiraJqlSearchView{
		bottomBar:  bottomBar,
		jqlStorage: &jqlStorage{},
	}
}

func (view *fjiraJqlSearchView) Init() {
	go view.startJqlFuzzyFind()
	go view.handleBottomBarActions()
}

func (view *fjiraJqlSearchView) Destroy() {
	view.bottomBar.Destroy()
}

func (view *fjiraJqlSearchView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.bottomBar.Draw(screen)
}

func (view *fjiraJqlSearchView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraJqlSearchView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraJqlSearchView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraJqlSearchView) startJqlFuzzyFind() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	jqls, err := view.jqlStorage.readAll()
	if err != nil {
		app.Error(err.Error())
		return
	}
	jqls = append(jqls, DefaultJqlQuery)
	view.fuzzyFind = app.NewFuzzyFind(MessageJqlFuzzyFind, jqls)
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
			app.GetApp().SetView(NewJqlSearchView())
			return
		}
		view.fuzzyFind = nil
		goIntoIssuesSearchForJql(jqls[jql.Index])
	}
}

func (view *fjiraJqlSearchView) handleBottomBarActions() {
	for {
		action := <-view.bottomBar.Action
		switch action {
		case ActionCancel:
			if view.fuzzyFind != nil {
				app.GetApp().Quit()
			}
			if view.fuzzyFind == nil {
				app.GetApp().SetView(NewJqlSearchView())
			}
			return
		case ActionNew:
			app.GetApp().SetView(newTextWriterView(&textWriterArgs{
				header: MessageTypeJqlAndSave,
				goBack: func() {
					app.GetApp().SetView(NewJqlSearchView())
				},
				textConsumer: func(s string) {
					err := view.jqlStorage.addNew(s)
					if err != nil {
						app.Error(err.Error())
						return
					}
					app.Success(MessageJqlAddSuccess)
				},
				maxLength: MaxJqlLength,
			}))
			return
		case ActionDelete:
			go view.confirmJqlDelete(view.fuzzyFind.GetSelectedItem())
		}
	}
}

func (view *fjiraJqlSearchView) confirmJqlDelete(jql string) {
	message := fmt.Sprintf(MessageJqlRemoveConfirm, jql)
	view.fuzzyFind = nil
	app.GetApp().ClearNow()
	view.bottomBar.Clear()
	view.bottomBar.AddItem(NewYesBarItem())
	view.bottomBar.AddItem(NewCancelBarItem())
	createNewJql := app.Confirm(app.GetApp(), message)
	switch createNewJql {
	case true:
		err := view.jqlStorage.remove(jql)
		if err != nil {
			app.Error(fmt.Sprintf(MessageCannotAddNewJql, err))
			return
		}
		app.Success(MessageJqlRemoveSuccess)
		app.GetApp().SetView(NewJqlSearchView())
	case false:
		app.GetApp().SetView(NewJqlSearchView())
	}
}
