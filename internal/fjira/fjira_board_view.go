package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"strings"
)

const (
	topMargin = 2 // 1 for navigation
)

var (
	columnHeaderStyle  = TopBarItemDefault
	issueStyle         = app.DefaultStyle.Background(tcell.NewRGBColor(35, 35, 35)).Foreground(tcell.ColorWhite)
	cursorIssueStyle   = app.DefaultStyle.Foreground(tcell.ColorWhite).Background(tcell.NewRGBColor(72, 72, 72))
	selectedIssueStyle = app.DefaultStyle.Background(tcell.ColorDarkRed).Foreground(tcell.ColorWhite).Bold(true)
	titleStyle         = app.DefaultStyle.Italic(true).Foreground(tcell.NewRGBColor(236, 206, 88))
)

type boardView struct {
	app.View
	bottomBar              *app.ActionBar
	selectedIssueBottomBar *app.ActionBar
	topBar                 *app.ActionBar
	boardConfiguration     *jira.BoardConfiguration
	project                *jira.Project
	issues                 []jira.Issue
	statusesColumnsMap     map[string]int
	columnStatusesMap      map[int][]string
	columnsX               map[int]int
	issuesRow              map[string]int
	issuesX                map[string]int
	issuesColumn           map[string]int
	issuesSummaries        map[string]string
	columns                []string
	highlightedIssue       *jira.Issue
	issueSelected          bool
	tmpX                   int
	cursorX                int
	cursorY                int
	screenX, screenY       int
	scrollX                int
	columnSize             int
}

func newBoardView(project *jira.Project, boardConfiguration *jira.BoardConfiguration) *boardView {
	col := 0
	statusesColumnsMap := map[string]int{}
	columnStatusesMap := map[int][]string{}
	columns := make([]string, 0, 20)
	for _, column := range boardConfiguration.ColumnConfig.Columns {
		if len(column.Statuses) == 0 {
			continue
		}
		if columnStatusesMap[col] == nil {
			columnStatusesMap[col] = make([]string, 0, 20)
		}
		for _, status := range column.Statuses {
			statusesColumnsMap[status.Id] = col
			columnStatusesMap[col] = append(columnStatusesMap[col], status.Id)
		}
		columns = append(columns, column.Name)
		col++
	}
	bottomBar := CreateBottomLeftBar()
	bottomBar.AddItem(CreateArrowsNavigateItem())
	bottomBar.AddItem(CreateSelectItem())
	bottomBar.AddItem(NewOpenBarItem())
	bottomBar.AddItem(NewCancelBarItem())
	selectedIssueBottomBar := CreateBottomLeftBar()
	selectedIssueBottomBar.AddItem(CreateMoveArrowsItem())
	selectedIssueBottomBar.AddItem(CreateUnSelectItem())
	selectedIssueBottomBar.AddItem(NewCancelBarItem())
	topBar := CreateIssueTopBar(&jira.Issue{})
	return &boardView{
		project:                project,
		boardConfiguration:     boardConfiguration,
		statusesColumnsMap:     statusesColumnsMap,
		columnStatusesMap:      columnStatusesMap,
		columns:                columns,
		columnsX:               map[int]int{},
		issuesRow:              map[string]int{},
		issuesColumn:           map[string]int{},
		issuesSummaries:        map[string]string{},
		cursorX:                0,
		cursorY:                0,
		bottomBar:              bottomBar,
		selectedIssueBottomBar: selectedIssueBottomBar,
		topBar:                 topBar,
		scrollX:                0,
		highlightedIssue:       &jira.Issue{},
		columnSize:             28,
	}
}

func (b *boardView) Draw(screen tcell.Screen) {
	b.topBar.Draw(screen)
	b.tmpX = 0
	for _, column := range b.columns {
		app.DrawText(screen, b.tmpX-b.scrollX, topMargin, columnHeaderStyle, centerString(column, b.columnSize))
		b.tmpX += b.columnSize + 1
	}
	if len(b.issues) == 0 {
		return
	}
	for _, issue := range b.issues {
		column := b.statusesColumnsMap[issue.Fields.Status.Id]
		x := b.columnsX[column]
		y := b.issuesRow[issue.Id]
		if b.highlightedIssue.Id == issue.Id {
			var style = &cursorIssueStyle
			if b.issueSelected {
				style = &selectedIssueStyle
			}
			app.DrawTextLimited(screen, x-b.scrollX, y+topMargin, x+b.columnSize-b.scrollX, y+1+topMargin, *style, b.issuesSummaries[issue.Id])
			continue
		}
		app.DrawTextLimited(screen, x-b.scrollX, y+topMargin, x+b.columnSize-b.scrollX, y+1+topMargin, issueStyle, b.issuesSummaries[issue.Id])
	}
	if b.highlightedIssue != nil {
		app.DrawText(screen, 0, 1, titleStyle, app.WriteIndicator)
		app.DrawText(screen, 2, 1, titleStyle, b.issuesSummaries[b.highlightedIssue.Id])
	}
	if !b.issueSelected {
		b.bottomBar.Draw(screen)
	} else {
		b.selectedIssueBottomBar.Draw(screen)
	}
	b.ensureHighlightInViewport()
}

func (b *boardView) Update() {
	b.bottomBar.Update()
	b.selectedIssueBottomBar.Update()
	b.topBar.Update()
}

func (b *boardView) Resize(screenX, screenY int) {
	b.bottomBar.Resize(screenX, screenY)
	b.selectedIssueBottomBar.Resize(screenX, screenY)
	b.topBar.Resize(screenX, screenY)
	b.screenY = screenY
	b.screenX = screenX
	for i := range b.columns {
		if i == 0 {
			b.columnsX[i] = 0
		} else {
			b.columnsX[i] = i*b.columnSize + i
		}
	}
}

func (b *boardView) Init() {
	app.GetApp().Loading(true)
	api, _ := GetApi()
	// it's not perfect - but I cannot find simple way to fetch all issues visible at given board
	issues, err := api.SearchJql(fmt.Sprintf("project=%s order by updatedDate,createdDate", b.project.Id))
	if err != nil {
		app.GetApp().Loading(false)
		app.Error(err.Error())
		return
	}
	b.issues = issues
	b.refreshIssuesSummaries()
	b.refreshIssuesRows()
	b.refreshHighlightedIssue()
	app.GetApp().Loading(false)
	go b.handleActions()
}

func (b *boardView) Destroy() {
	// ...
}

func (b *boardView) Refresh() {
	b.refreshIssuesSummaries()
	b.refreshIssuesRows()
	b.refreshIssueTopBar()
	b.refreshHighlightedIssue()
}

func (b *boardView) SetColumnSize(colSize int) {
	b.columnSize = colSize
	b.Resize(b.scrollX, b.screenY)
	b.Refresh()
}

func (b *boardView) HandleKeyEvent(ev *tcell.EventKey) {
	if app.GetApp().IsLoading() {
		return
	}
	if !b.issueSelected {
		b.bottomBar.HandleKeyEvent(ev)
	} else {
		b.selectedIssueBottomBar.HandleKeyEvent(ev)
	}
	if ev.Key() == tcell.KeyRight {
		newColumn := b.statusesColumnsMap[b.highlightedIssue.Fields.Status.Id] + 1
		if newColumn > len(b.statusesColumnsMap) {
			return
		}
		b.cursorX = app.MinInt(len(b.columns), b.cursorX+1)
		b.cursorY = 0
		if b.issueSelected {
			b.moveIssue(b.highlightedIssue, 1)
			return
		}
		b.refreshHighlightedIssue()
	}
	if ev.Key() == tcell.KeyLeft {
		newColumn := b.statusesColumnsMap[b.highlightedIssue.Fields.Status.Id] - 1
		if newColumn < 0 {
			return
		}
		b.cursorX = app.MaxInt(0, b.cursorX-1)
		b.cursorY = 0
		if b.issueSelected {
			b.moveIssue(b.highlightedIssue, -1)
			return
		}
		b.refreshHighlightedIssue()
	}
	if ev.Key() == tcell.KeyUp {
		b.cursorY = app.MaxInt(0, b.cursorY-1)
		b.refreshHighlightedIssue()
	}
	if ev.Key() == tcell.KeyDown {
		// TODO - get number of issues in column
		b.cursorY = app.MinInt(1000, b.cursorY+1)
		b.refreshHighlightedIssue()
	}
}

func (b *boardView) handleActions() {
	defer app.GetApp().PanicRecover()
	for {
		select {
		case action := <-b.bottomBar.Action:
			switch action {
			case ActionSelect:
				b.issueSelected = true
			case ActionCancel:
				goIntoIssuesSearchForProject(b.project.Id)
			case ActionOpen:
				goIntoIssueView(b.highlightedIssue.Key)
			}
		case action := <-b.selectedIssueBottomBar.Action:
			switch action {
			case ActionUnselect:
				b.issueSelected = false
			case ActionCancel:
				b.issueSelected = false
			case ActionOpen:
				goIntoIssueView(b.highlightedIssue.Key)
			}
		default:
		}
	}
}

func (b *boardView) refreshIssueTopBar() {
	b.topBar.GetItem(0).ChangeText2(b.highlightedIssue.Key)
	b.topBar.GetItem(1).ChangeText2(b.highlightedIssue.Fields.Reporter.DisplayName)
	b.topBar.GetItem(2).ChangeText2(b.highlightedIssue.Fields.Assignee.DisplayName)
	b.topBar.GetItem(3).ChangeText2(b.highlightedIssue.Fields.Type.Name)
	b.topBar.GetItem(4).ChangeText2(b.highlightedIssue.Fields.Status.Name)
	b.topBar.Resize(b.screenX, b.screenY)
}

func (b *boardView) refreshHighlightedIssue() {
	for i, issue := range b.issues {
		y := b.issuesRow[issue.Id]
		if b.issuesColumn[issue.Id] == b.cursorX && y-1 == b.cursorY {
			if b.highlightedIssue.Key != issue.Key {
				b.highlightedIssue = &b.issues[i]
				b.refreshIssueTopBar()
			}
		}
	}
}

func (b *boardView) pointCursorTo(issueId string) {
	b.cursorX = b.issuesColumn[issueId]
	b.cursorY = b.issuesRow[issueId] - 1
}

func (b *boardView) refreshIssuesRows() {
	rows := map[int]int{}
	for _, issue := range b.issues {
		column := b.statusesColumnsMap[issue.Fields.Status.Id]
		y := rows[column] + 1
		b.issuesRow[issue.Id] = y
		b.issuesColumn[issue.Id] = column
		rows[column] = y
	}
}

func (b *boardView) refreshIssuesSummaries() {
	for _, issue := range b.issues {
		b.issuesSummaries[issue.Id] = fmt.Sprintf("%s %s", issue.Key, issue.Fields.Summary)
	}
}

func (b *boardView) moveIssue(issue *jira.Issue, direction int) {
	app.GetApp().Loading(true)
	inc := 1
	if direction < 0 {
		inc = -1
	}
	column := b.statusesColumnsMap[issue.Fields.Status.Id] + inc
	targetColumnStatuses := b.columnStatusesMap[column]
	api, _ := GetApi()
	transitions, err := api.FindTransitions(issue.Id)
	if err != nil {
		app.GetApp().Loading(false)
		app.Error(err.Error())
		return
	}
	var targetTransition *jira.IssueTransition
	for i, transition := range transitions {
		for _, targetStatus := range targetColumnStatuses {
			if transition.To.StatusId == targetStatus {
				targetTransition = &transitions[i]
				break
			}
		}
	}
	if targetTransition == nil {
		app.GetApp().Loading(false)
		app.Error(MessageCannotFindStatusForColumn)
		return
	}
	err = api.DoTransition(issue.Id, targetTransition)
	if err != nil {
		app.GetApp().Loading(false)
		app.Error(err.Error())
		return
	}
	app.GetApp().Loading(false)
	app.Success(fmt.Sprintf(MessageChangeStatusSuccess, issue.Key, targetTransition.To.Name))
	issue.Fields.Status.Id = targetTransition.To.StatusId
	issue.Fields.Status.Name = targetTransition.To.Name
	b.issueSelected = false
	b.refreshIssuesRows()
	b.pointCursorTo(issue.Id)
	b.refreshHighlightedIssue()
}

func (b *boardView) ensureHighlightInViewport() {
	if b.highlightedIssue == nil {
		return
	}
	if b.scrollX+(b.cursorX*b.columnSize)+b.columnSize > b.screenX { // highlighted issue out of screen
		b.scrollX = app.MaxInt(0, (b.cursorX-2)*b.columnSize)
	}
}

func centerString(str string, width int) string {
	if len(str) > width {
		str = str[:width]
	}
	spaces := int(float64(width-len(str)) / 2)
	return strings.Repeat(" ", spaces) + str + strings.Repeat(" ", width-(spaces+len(str)))
}
