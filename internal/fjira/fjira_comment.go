package fjira

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"unicode"
)

type fjiraCommentView struct {
	app.View
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	issue     *jira.JiraIssue
	buffer    bytes.Buffer
	text      string
}

const (
	MaxCommentLineWidth = 150
)

var (
	headerStyle = app.DefaultStyle.Foreground(tcell.ColorWhite).Underline(true)
)

func NewCommentView(issue *jira.JiraIssue) *fjiraCommentView {
	bottomBar := CreateIssueBottomBar(issue)
	bottomBar.AddItem(NewSaveBarItem())
	bottomBar.AddItem(NewCancelBarItem())
	return &fjiraCommentView{
		issue:     issue,
		topBar:    CreateIssueTopBar(issue),
		bottomBar: bottomBar,
		text:      "",
	}
}

func (view *fjiraCommentView) Init() {
	go view.handleBottomBarActions()
}

func (view *fjiraCommentView) Destroy() {
	// do nothing
}

func (view *fjiraCommentView) Draw(screen tcell.Screen) {
	app.DrawText(screen, 1, 2, headerStyle, MessageTypeCommentAndSave)
	app.DrawTextLimited(screen, 1, 4, MaxCommentLineWidth, 100, tcell.StyleDefault, view.text)
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *fjiraCommentView) Update() {
	view.topBar.Update()
	view.bottomBar.Update()
}

func (view *fjiraCommentView) Resize(screenX, screenY int) {
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraCommentView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if unicode.IsLetter(ev.Rune()) || unicode.IsDigit(ev.Rune()) || unicode.IsSpace(ev.Rune()) ||
		unicode.IsPunct(ev.Rune()) || unicode.IsSymbol(ev.Rune()) {
		view.buffer.WriteRune(ev.Rune())
	}
	if ev.Key() == tcell.KeyEnter {
		view.buffer.WriteRune('\n')
	}
	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if view.buffer.Len() > 0 {
			view.buffer.Truncate(view.buffer.Len() - 1)
		}
	}
	view.text = view.buffer.String()
}

func (view *fjiraCommentView) handleBottomBarActions() {
	action := <-view.bottomBar.Action
	switch action {
	case ActionYes:
		view.doComment(view.issue, view.buffer.String())
	}
	go goIntoIssueView(view.issue.Key)
}

func (view fjiraCommentView) doComment(issue *jira.JiraIssue, comment string) {
	app.GetApp().LoadingWithText(true, MessageAddingComment)
	api, _ := GetApi()
	err := api.DoComment(issue.Key, comment)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(MessageCannotAddComment, issue.Key, err))
	}
	app.Success(fmt.Sprintf(MessageCommentSuccess, issue.Key))
}
