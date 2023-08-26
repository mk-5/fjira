package ui

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"unicode"
)

// TODO - should be here?
type TextWriterView struct {
	app.View
	bottomBar   *app.ActionBar
	buffer      bytes.Buffer
	text        string
	headerStyle tcell.Style
	args        TextWriterArgs
}

type TextWriterArgs struct {
	MaxLength    int
	TextConsumer func(string)
	GoBack       func()
	Header       string
}

func NewTextWriterView(args *TextWriterArgs) app.View {
	bottomBar := CreateBottomLeftBar()
	bottomBar.AddItem(NewSaveBarItem())
	bottomBar.AddItem(NewCancelBarItem())
	if args.MaxLength == 0 {
		args.MaxLength = 150
	}
	if args.GoBack == nil {
		args.GoBack = func() {
		}
	}
	if args.TextConsumer == nil {
		args.TextConsumer = func(str string) {
		}
	}
	return &TextWriterView{
		bottomBar:   bottomBar,
		text:        "",
		args:        *args,
		headerStyle: app.DefaultStyle.Foreground(tcell.ColorWhite).Underline(true),
	}
}

func (view *TextWriterView) Init() {
	go view.handleBottomBarActions()
}

func (view *TextWriterView) Destroy() {
	view.bottomBar.Destroy()
}

func (view *TextWriterView) Draw(screen tcell.Screen) {
	app.DrawText(screen, 1, 2, view.headerStyle, view.args.Header)
	app.DrawTextLimited(screen, 1, 4, view.args.MaxLength, 100, app.DefaultStyle, view.text)
	view.bottomBar.Draw(screen)
}

func (view *TextWriterView) Update() {
	view.bottomBar.Update()
}

func (view *TextWriterView) Resize(screenX, screenY int) {
	view.bottomBar.Resize(screenX, screenY)
}

func (view *TextWriterView) HandleKeyEvent(ev *tcell.EventKey) {
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

func (view *TextWriterView) handleBottomBarActions() {
	action := <-view.bottomBar.Action
	switch action {
	case ActionYes:
		view.args.TextConsumer(view.buffer.String())
	}
	go view.args.GoBack()
}

func isBackspace(ev *tcell.EventKey) bool {
	return ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBS
}
