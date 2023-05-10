package fjira

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"unicode"
)

type fjiraTextWriterView struct {
	app.View
	bottomBar   *app.ActionBar
	buffer      bytes.Buffer
	text        string
	headerStyle tcell.Style
	args        textWriterArgs
}

type textWriterArgs struct {
	maxLength    int
	textConsumer func(string)
	goBack       func()
	header       string
}

func newTextWriterView(args *textWriterArgs) *fjiraTextWriterView {
	bottomBar := CreateBottomLeftBar()
	bottomBar.AddItem(NewSaveBarItem())
	bottomBar.AddItem(NewCancelBarItem())
	if args.maxLength == 0 {
		args.maxLength = 150
	}
	if args.goBack == nil {
		args.goBack = func() {
		}
	}
	if args.textConsumer == nil {
		args.textConsumer = func(str string) {
		}
	}
	return &fjiraTextWriterView{
		bottomBar:   bottomBar,
		text:        "",
		args:        *args,
		headerStyle: app.DefaultStyle.Foreground(tcell.ColorWhite).Underline(true),
	}
}

func (view *fjiraTextWriterView) Init() {
	go view.handleBottomBarActions()
}

func (view *fjiraTextWriterView) Destroy() {
	view.bottomBar.Destroy()
}

func (view *fjiraTextWriterView) Draw(screen tcell.Screen) {
	app.DrawText(screen, 1, 2, view.headerStyle, view.args.header)
	app.DrawTextLimited(screen, 1, 4, view.args.maxLength, 100, app.DefaultStyle, view.text)
	view.bottomBar.Draw(screen)
}

func (view *fjiraTextWriterView) Update() {
	view.bottomBar.Update()
}

func (view *fjiraTextWriterView) Resize(screenX, screenY int) {
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraTextWriterView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if (unicode.IsLetter(ev.Rune()) || unicode.IsDigit(ev.Rune()) || unicode.IsSpace(ev.Rune()) ||
		unicode.IsPunct(ev.Rune()) || unicode.IsSymbol(ev.Rune())) && !isBackspace(ev) {
		view.buffer.WriteRune(ev.Rune())
	}
	if ev.Key() == tcell.KeyEnter {
		view.buffer.WriteRune('\n')
	}
	if isBackspace(ev) {
		if view.buffer.Len() > 0 {
			view.buffer.Truncate(view.buffer.Len() - 1)
		}
	}
	view.text = view.buffer.String()
}

func (view *fjiraTextWriterView) handleBottomBarActions() {
	action := <-view.bottomBar.Action
	switch action {
	case ActionYes:
		view.args.textConsumer(view.buffer.String())
	}
	go view.args.goBack()
}

func isBackspace(ev *tcell.EventKey) bool {
	return ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBS
}
