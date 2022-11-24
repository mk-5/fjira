package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

type fjiraAddLabelView struct {
	app.View
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
	labels    []string
}

func NewAddLabelView(issue *jira.Issue) *fjiraAddLabelView {
	return &fjiraAddLabelView{
		issue:     issue,
		topBar:    CreateIssueTopBar(issue),
		bottomBar: CreateBottomLeftBar(),
	}
}

func (view *fjiraAddLabelView) Init() {
	go view.startLabelSearching()
}

func (*fjiraAddLabelView) Destroy() {
}

func (view *fjiraAddLabelView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *fjiraAddLabelView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraAddLabelView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraAddLabelView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraAddLabelView) startLabelSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	view.fuzzyFind = app.NewFuzzyFindWithProvider(MessageLabelFuzzyFind, view.findLabels)
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

func (view *fjiraAddLabelView) findLabels(query string) []string {
	api, _ := GetApi()
	app.GetApp().LoadingWithText(true, MessageSearchLabelsLoading)
	labels, err := api.FindLabels(view.issue, query)
	if err != nil {
		app.Error(err.Error())
	}
	app.GetApp().Loading(false)
	view.labels = labels
	return labels
}

func (view *fjiraAddLabelView) addLabelToIssue(issue *jira.Issue, label string) {
	if label == "" {
		app.GetApp().SetView(NewIssueView(view.issue))
		return
	}
	view.doAddLabel(issue, label)
	goIntoIssueView(view.issue.Key)
}

func (view fjiraAddLabelView) doAddLabel(issue *jira.Issue, label string) {
	app.GetApp().LoadingWithText(true, MessageAddingLabel)
	api, _ := GetApi()
	err := api.AddLabel(issue.Key, label)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(MessageCannotAddLabel, label, issue.Key, err))
		return
	}
	app.Success(fmt.Sprintf(MessageAddLabelSuccess, label, issue.Key))
}
