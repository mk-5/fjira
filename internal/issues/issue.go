package issues

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/comments"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"math"
	"strings"
)

type issueView struct {
	app.View
	api               jira.Api
	bottomBar         *app.ActionBar
	topBar            *app.ActionBar
	fuzzyFind         *app.FuzzyFind
	issue             *jira.Issue
	goBackFn          func()
	descriptionLimitX int
	descriptionLimitY int
	scrollY           int
	descriptionLines  int
	commentsLines     int
	maxScrollY        int
	body              string
	summaryLen        int
	labels            string
	labelsLen         int
	comments          []comments.Comment
	lastY             int
	boxTitleStyle     tcell.Style
	defaultStyle      tcell.Style
}

var (
	issueNavItems = []ui.NavItemConfig{
		ui.NavItemConfig{Action: ui.ActionStatusChange, Text1: ui.MessageChangeStatus, Text2: "[s]", Rune: 's'},
		ui.NavItemConfig{Action: ui.ActionAssigneeChange, Text1: ui.MessageAssignUser, Text2: "[a]", Rune: 'a'},
		ui.NavItemConfig{Action: ui.ActionComment, Text1: ui.MessageComment, Text2: "[c]", Rune: 'c'},
		ui.NavItemConfig{Action: ui.ActionAddLabel, Text1: ui.MessageLabel, Text2: "[l]", Rune: 'l'},
		ui.NavItemConfig{Action: ui.ActionOpen, Text1: ui.MessageOpen, Text2: "[o]", Rune: 'o'},
		ui.NavItemConfig{Action: ui.ActionCopyIssue, Text1: ui.MessageCopyIssue, Text2: "[y]", Rune: 'y'},
	}
)

const (
	maxCommentLineWidth = 150
	labelsDelimiter     = " | "
)

func NewIssueView(issue *jira.Issue, goBackFn func(), api jira.Api) app.View {
	bottomBar := ui.CreateBottomActionBarWithItems(issueNavItems)
	bottomBar.AddItem(ui.CreateScrollBarItem())
	bottomBar.AddItem(ui.NewCancelBarItem())

	issueActionBar := ui.CreateIssueTopBar(issue)
	cs := comments.ParseCommentsFromIssue(issue, 1000, 1000)
	ls := strings.Join(issue.Fields.Labels, labelsDelimiter)
	labelsLen := len(ls)

	return &issueView{
		api:           api,
		bottomBar:     bottomBar,
		topBar:        issueActionBar,
		issue:         issue,
		scrollY:       0,
		body:          issue.Fields.Description,
		comments:      cs,
		labels:        ls,
		labelsLen:     labelsLen,
		summaryLen:    len(issue.Fields.Summary),
		goBackFn:      goBackFn,
		boxTitleStyle: app.DefaultStyle().Foreground(app.Color("details.foreground")),
		defaultStyle:  app.DefaultStyle(),
	}
}

func (view *issueView) Init() {
	go view.handleIssueAction()
}

func (view *issueView) Destroy() {
}

func (view *issueView) Draw(screen tcell.Screen) {
	if view.fuzzyFind == nil {
		app.DrawBox(screen, 1, 2-view.scrollY, view.summaryLen+4, 4-view.scrollY, view.boxTitleStyle)
		app.DrawText(screen, 2, 2-view.scrollY, view.boxTitleStyle, ui.MessageSummary)
		app.DrawText(screen, 3, 3-view.scrollY, view.defaultStyle, view.issue.Fields.Summary)

		view.lastY = 2 - view.scrollY + 2

		if view.labels != "" {
			app.DrawBox(screen, 1, view.lastY+1, view.labelsLen+4, view.lastY+3, view.boxTitleStyle)
			app.DrawText(screen, 2, view.lastY+1, view.boxTitleStyle, ui.MessageLabels)
			app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.lastY+2, view.defaultStyle, view.labels)
			view.lastY = view.lastY + 3
		}

		app.DrawBox(screen, 1, view.lastY+1, view.descriptionLimitX+4, view.lastY+1+view.descriptionLines+4, view.boxTitleStyle)
		app.DrawText(screen, 2, view.lastY+1, view.boxTitleStyle, ui.MessageDescription)
		app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.descriptionLimitY, view.defaultStyle, view.body)

		view.lastY = view.lastY + view.descriptionLines + 6

		for _, comment := range view.comments {
			app.DrawBox(screen, 1, view.lastY+1, view.descriptionLimitX+4, view.lastY+1+comment.Lines+2, view.boxTitleStyle)
			app.DrawText(screen, 2, view.lastY+1, view.boxTitleStyle, comment.Title)
			app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.descriptionLimitY, view.defaultStyle, comment.Body)
			view.lastY = view.lastY + 1 + comment.Lines + 3
		}
	}
	view.bottomBar.Draw(screen)
	view.topBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *issueView) Update() {
	view.bottomBar.Update()
	view.topBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *issueView) Resize(screenX, screenY int) {
	view.descriptionLimitX = app.ClampInt(int(math.Floor(float64(screenX)*0.9)), 1, 10000)
	view.descriptionLimitY = 1000
	view.descriptionLines = app.DrawTextLimited(nil, 0, 0, view.descriptionLimitX, view.descriptionLimitY, view.defaultStyle, view.body) + 1
	commentsLines := 0
	view.comments = comments.ParseCommentsFromIssue(view.issue, view.descriptionLimitX, view.descriptionLimitY)
	for _, comment := range view.comments {
		commentsLines = commentsLines + comment.Lines + 3
	}
	view.commentsLines = commentsLines + len(view.comments) + 1
	topAndBottomBarSize := 12
	view.maxScrollY = app.ClampInt(int(math.Abs(float64(screenY-topAndBottomBarSize-view.descriptionLines-view.commentsLines-10))), 0, 2000)
	view.bottomBar.Resize(screenX, screenY)
	view.topBar.Resize(screenX, screenY)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
}

func (view *issueView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	view.topBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
	if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyTab {
		view.scrollY = app.ClampInt(view.scrollY-1, 0, view.maxScrollY)
	}
	if ev.Key() == tcell.KeyDown || ev.Key() == tcell.KeyBacktab {
		view.scrollY = app.ClampInt(view.scrollY+1, 0, view.maxScrollY)
	}
}

func (view *issueView) goBack() {
	if view.goBackFn != nil {
		view.goBackFn()
	}
}

func (view *issueView) handleIssueAction() {
	if selectedAction := <-view.bottomBar.Action; true {
		switch selectedAction {
		case ui.ActionCancel:
			view.goBack()
			return
		case ui.ActionStatusChange:
			app.GoTo("status-change", view.issue, view.reopen, view.api)
			return
		case ui.ActionAssigneeChange:
			app.GoTo("users-assign", view.issue, view.reopen, view.api)
			return
		case ui.ActionComment:
			app.GoTo("text-writer", &ui.TextWriterArgs{
				Header: ui.MessageTypeCommentAndSave,
				GoBack: func() {
					view.reopen()
				},
				TextConsumer: func(s string) {
					view.doComment(view.issue, s)
				},
				MaxLength: maxCommentLineWidth,
			})
			return
		case ui.ActionAddLabel:
			app.GoTo("labels-add", view.issue, view.reopen, view.api)
			return
		case ui.ActionOpen:
			OpenIssueInBrowser(view.issue, view.api)
			go view.handleIssueAction()
			return
		case ui.ActionCopyIssue:
			CopyIssue(view.issue)
			go view.handleIssueAction()
			return
		}
	}
}

func (view *issueView) reopen() {
	app.GoTo("issue", view.issue.Key, view.goBackFn, view.api)
}

func (view *issueView) doComment(issue *jira.Issue, comment string) {
	app.GetApp().LoadingWithText(true, ui.MessageAddingComment)
	err := view.api.DoComment(issue.Key, comment)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(ui.MessageCannotAddComment, issue.Key, err))
	}
	app.Success(fmt.Sprintf(ui.MessageCommentSuccess, issue.Key))
}
