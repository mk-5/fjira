package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

type fjiraAddLabelView struct {
	app.View
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.JiraIssue
}

func NewAddLabelView(issue *jira.JiraIssue) *fjiraAddLabelView {
	return &fjiraAddLabelView{
		issue:     issue,
		topBar:    CreateIssueTopBar(issue),
		bottomBar: CreateIssueBottomBar(),
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
	labels := view.findLabels()
	view.fuzzyFind = app.NewFuzzyFind(MessageLabelFuzzyFind, labels)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if match := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		label := view.fuzzyFind.GetQuery()
		if match.Index >= 0 {
			label = labels[match.Index]
		}
		view.fuzzyFind = nil
		view.addLabelToIssue(view.issue, label)
	}
}

func (view *fjiraAddLabelView) findLabels() []string {
	api, _ := GetApi()
	labels, err := api.FindLabels()
	if err != nil {
		app.Error(err.Error())
	}
	return labels
}

func (view *fjiraAddLabelView) addLabelToIssue(issue *jira.JiraIssue, label string) {
	if label == "" {
		app.GetApp().SetView(NewIssueView(view.issue))
		return
	}
	view.doAddLabel(issue, label)
	goIntoIssueView(view.issue.Key)
}

func (view fjiraAddLabelView) doAddLabel(issue *jira.JiraIssue, label string) {
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
