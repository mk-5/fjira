package fjira

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"math"
	"strings"
)

type fjiraIssueView struct {
	app.View
	bottomBar         *app.ActionBar
	topBar            *app.ActionBar
	fuzzyFind         *app.FuzzyFind
	issue             *jira.Issue
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
	comments          []struct {
		body  string
		title string
		lines int
	}
	lastY int
}

var (
	boxTitleStyle = app.DefaultStyle.Foreground(tcell.ColorDimGrey)
)

const (
	labelsDelimiter = " | "
)

func NewIssueView(issue *jira.Issue) *fjiraIssueView {
	bottomBar := CreateIssueBottomBar()
	bottomBar.AddItem(NewStatusChangeBarItem())
	bottomBar.AddItem(NewAssigneeChangeBarItem())
	bottomBar.AddItem(CreateCommentBarItem())
	bottomBar.AddItem(CreateAddLabelBarItem())
	bottomBar.AddItem(NewOpenBarItem())
	bottomBar.AddItem(CreateScrollBarItem())
	bottomBar.AddItem(NewCancelBarItem())

	issueActionBar := CreateIssueTopBar(issue)
	comments := parseComments(issue, 1000, 1000)
	labels := strings.Join(issue.Fields.Labels, labelsDelimiter)
	labelsLen := len(labels)

	return &fjiraIssueView{
		bottomBar:  bottomBar,
		topBar:     issueActionBar,
		issue:      issue,
		scrollY:    0,
		body:       issue.Fields.Description,
		comments:   comments,
		labels:     labels,
		labelsLen:  labelsLen,
		summaryLen: len(issue.Fields.Summary),
	}
}

func (view *fjiraIssueView) Init() {
	go view.handleIssueAction()
}

func (view *fjiraIssueView) Destroy() {
}

func (view *fjiraIssueView) Draw(screen tcell.Screen) {
	if view.fuzzyFind == nil {
		app.DrawBox(screen, 1, 2-view.scrollY, view.summaryLen+4, 4-view.scrollY, boxTitleStyle)
		app.DrawText(screen, 2, 2-view.scrollY, boxTitleStyle, MessageSummary)
		app.DrawText(screen, 3, 3-view.scrollY, app.DefaultStyle, view.issue.Fields.Summary)

		view.lastY = 2 - view.scrollY + 2

		if view.labels != "" {
			app.DrawBox(screen, 1, view.lastY+1, view.labelsLen+4, view.lastY+3, boxTitleStyle)
			app.DrawText(screen, 2, view.lastY+1, boxTitleStyle, MessageLabels)
			app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.lastY+2, app.DefaultStyle, view.labels)
			view.lastY = view.lastY + 3
		}

		app.DrawBox(screen, 1, view.lastY+1, view.descriptionLimitX+4, view.lastY+1+view.descriptionLines+4, boxTitleStyle)
		app.DrawText(screen, 2, view.lastY+1, boxTitleStyle, MessageDescription)
		app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.descriptionLimitY, app.DefaultStyle, view.body)

		view.lastY = view.lastY + view.descriptionLines + 6

		for _, comment := range view.comments {
			app.DrawBox(screen, 1, view.lastY+1, view.descriptionLimitX+4, view.lastY+1+comment.lines+2, boxTitleStyle)
			app.DrawText(screen, 2, view.lastY+1, boxTitleStyle, comment.title)
			app.DrawTextLimited(screen, 3, view.lastY+2, view.descriptionLimitX, view.descriptionLimitY, app.DefaultStyle, comment.body)
			view.lastY = view.lastY + 1 + comment.lines + 3
		}
	}
	view.bottomBar.Draw(screen)
	view.topBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *fjiraIssueView) Update() {
	view.bottomBar.Update()
	view.topBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraIssueView) Resize(screenX, screenY int) {
	view.descriptionLimitX = app.ClampInt(int(math.Floor(float64(screenX)*0.9)), 1, 10000)
	view.descriptionLimitY = 1000
	view.descriptionLines = app.DrawTextLimited(nil, 0, 0, view.descriptionLimitX, view.descriptionLimitY, app.DefaultStyle, view.body) + 1
	commentsLines := 0
	view.comments = parseComments(view.issue, view.descriptionLimitX, view.descriptionLimitY)
	for _, comment := range view.comments {
		commentsLines = commentsLines + comment.lines + 3
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

func (view *fjiraIssueView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	view.topBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
	if ev.Key() == tcell.KeyUp {
		view.scrollY = app.ClampInt(view.scrollY-1, 0, view.maxScrollY)
	}
	if ev.Key() == tcell.KeyDown {
		view.scrollY = app.ClampInt(view.scrollY+1, 0, view.maxScrollY)
	}
}

func (view *fjiraIssueView) handleIssueAction() {
	if selectedAction := <-view.bottomBar.Action; true {
		switch selectedAction {
		case ActionEscape:
			app.GetApp().SetView(NewIssuesSearchView(&view.issue.Fields.Project))
			return
		case ActionStatusChange:
			goIntoChangeStatus(view.issue)
			return
		case ActionAssigneeChange:
			goIntoChangeAssignment(view.issue)
			return
		case ActionComment:
			goIntoCommentView(view.issue)
			return
		case ActionAddLabel:
			goIntoAddLabelView(view.issue)
			return
		case ActionOpen:
			jiraUrl, _ := GetJiraUrl()
			app.OpenLink(fmt.Sprintf("%s/browse/%s", jiraUrl, view.issue.Key))
			go view.handleIssueAction()
			return
		}
	}
}

// TODO - could be optimized a bit
func parseComments(issue *jira.Issue, limitX, limitY int) []struct {
	body  string
	title string
	lines int
} {
	comments := make([]struct {
		body  string
		title string
		lines int
	}, 0, 100)
	var commentsBuffer bytes.Buffer
	if len(issue.Fields.Comment.Comments) > 0 {
		for _, comment := range issue.Fields.Comment.Comments {
			title := fmt.Sprintf("%s, %s", comment.Created, comment.Author.DisplayName)
			body := fmt.Sprintf("\n%s", comment.Body)
			lines := app.DrawTextLimited(nil, 0, 0, limitX, limitY, app.DefaultStyle, comment.Body) + 2
			comments = append(comments, struct {
				body  string
				title string
				lines int
			}{body, title, lines})
			commentsBuffer.Reset()
		}
	}
	return comments
}
