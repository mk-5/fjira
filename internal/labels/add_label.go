package labels

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

type addLabelView struct {
	app.View
	api       jira.Api
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
	goBackFn  func()
	labels    []string
}

func NewAddLabelView(issue *jira.Issue, goBackFn func(), api jira.Api) app.View {
	return &addLabelView{
		api:       api,
		issue:     issue,
		goBackFn:  goBackFn,
		topBar:    ui.CreateIssueTopBar(issue),
		bottomBar: ui.CreateBottomLeftBar(),
	}
}

func (view *addLabelView) Init() {
	go view.startLabelSearching()
}

func (*addLabelView) Destroy() {
}

func (view *addLabelView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *addLabelView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *addLabelView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *addLabelView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *addLabelView) startLabelSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	view.fuzzyFind = app.NewFuzzyFindWithProvider(ui.MessageLabelFuzzyFind, view.findLabels)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if match := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		label := view.fuzzyFind.GetQuery()
		if match.Index >= 0 {
			label = view.labels[match.Index]
		}
		view.fuzzyFind = nil
		view.addLabelToIssue(view.issue, label)
	}
}

func (view *addLabelView) findLabels(query string) []string {
	app.GetApp().LoadingWithText(true, ui.MessageSearchLabelsLoading)
	labels, err := view.api.FindLabels(view.issue, query)
	if err != nil {
		app.Error(err.Error())
	}
	app.GetApp().Loading(false)
	view.labels = labels
	return labels
}

func (view *addLabelView) addLabelToIssue(issue *jira.Issue, label string) {
	if label == "" {
		view.goBackFn()
		return
	}
	view.doAddLabel(issue, label)
	view.goBackFn()
}

func (view *addLabelView) doAddLabel(issue *jira.Issue, label string) {
	app.GetApp().LoadingWithText(true, ui.MessageAddingLabel)
	err := view.api.AddLabel(issue.Key, label)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(ui.MessageCannotAddLabel, label, issue.Key, err))
		return
	}
	app.Success(fmt.Sprintf(ui.MessageAddLabelSuccess, label, issue.Key))
}
