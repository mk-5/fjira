package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"unicode/utf8"
)

type ActionBarAction int

type ActionBar struct {
	Action           chan ActionBarAction
	Y                int
	X                int
	screenX, screenY int
	items            []*ActionBarItem
	vAlign           int
	hAlign           int
}

const (
	Bottom               = -1
	Top                  = 1
	Left                 = 1
	Right                = -1
	ActionBarMaxItems    = 10
	ActionBarItemPadding = 1
	MessageLabelNone     = "-"
)

func NewActionBar(vAlign int, hAlign int) *ActionBar {
	return &ActionBar{
		Action:  make(chan ActionBarAction),
		Y:       0,
		screenX: 0,
		screenY: 0,
		vAlign:  vAlign,
		hAlign:  hAlign,
		items:   make([]*ActionBarItem, 0, ActionBarMaxItems),
	}
}

func (b *ActionBar) AddTextItem(id string, text string) {
	item := NewActionBarItem(len(b.items)+1, text, 0, 0)
	b.AddItem(item)
}

func (b *ActionBar) AddItem(item *ActionBarItem) {
	item.text = fmt.Sprintf("%s%s", item.Text1, item.Text2)
	item.x = b.getNextItemX()
	item.y = b.Y
	b.items = append(b.items, item)
	b.Resize(b.screenX, b.screenY)
}

func (b *ActionBar) AddItemWithStyles(text1 string, text2 string, text1Style tcell.Style, text2Style tcell.Style) {
	item := &ActionBarItem{
		Id:         len(b.items) + 1,
		Text1:      text1,
		Text2:      text2,
		Text1Style: text1Style,
		Text2Style: text2Style,
	}
	b.AddItem(item)
	b.Resize(b.screenX, b.screenY)
}

func (b *ActionBar) GetItem(index int) *ActionBarItem {
	if index > len(b.items)-1 {
		return nil
	}
	return b.items[index]
}

func (b *ActionBar) RemoveItem(id int) {
	for i, item := range b.items {
		if item.Id == id {
			b.RemoveItemAtIndex(i)
			return
		}
	}
}

func (b *ActionBar) RemoveItemAtIndex(index int) {
	if index > 0 {
		b.items = append(b.items[:index], b.items[index+1:]...)
		b.Resize(b.screenX, b.screenY)
	}
}

func (b *ActionBar) TrimItemsTo(index int) {
	if index > 0 {
		b.items = b.items[:index]
		b.Resize(b.screenX, b.screenY)
	}
}

func (b *ActionBar) Draw(screen tcell.Screen) {
	for _, item := range b.items {
		DrawText(screen, item.x, item.y, item.Text1Style, item.Text1)
		DrawText(screen, item.x+len(item.Text1), item.y, item.Text2Style, item.Text2)
		DrawText(screen, item.x+len(item.Text1)+len(item.Text2), item.y, DefaultStyle, " ")
	}
}

func (b *ActionBar) Update() {
	// do nothing
}

func (b *ActionBar) HandleKeyEvent(ev *tcell.EventKey) {
	for _, item := range b.items {
		if item.TriggerRune > 0 && item.TriggerRune == ev.Rune() {
			b.Action <- ActionBarAction(item.Id)
			return
		} else if item.TriggerKey > 0 && item.TriggerKey == ev.Key() {
			b.Action <- ActionBarAction(item.Id)
			return
		}
	}
}

func (b *ActionBar) Resize(screenX, screenY int) {
	b.screenX = screenX
	b.screenY = screenY
	switch b.vAlign {
	case Bottom:
		b.Y = screenY - 1
	case Top:
		b.Y = 0
	}
	for _, item := range b.items {
		item.y = b.Y
	}
}

func (b *ActionBar) getNextItemX() int {
	switch b.hAlign {
	case Right:
		if len(b.items) == 0 {
			return b.screenX
		}
		item := b.items[len(b.items)-1]
		return item.x - utf8.RuneCountInString(item.text) - ActionBarItemPadding
	default:
		if len(b.items) == 0 {
			return 0
		}
		item := b.items[len(b.items)-1]
		return item.x + utf8.RuneCountInString(item.text) + ActionBarItemPadding
	}
}
