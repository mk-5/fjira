package app

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

type ActionBarAction int

type ActionBar struct {
	Action           chan ActionBarAction
	Y                int
	X                int
	screenX, screenY int
	items            []*ActionBarItem
	lastItemX        int
	vAlign           int
	hAlign           int
}

type ActionBarItem struct {
	TextBox
	Id          int
	Text1       string
	Text2       string
	Text1Style  tcell.Style
	Text2Style  tcell.Style
	TriggerRune rune
	TriggerKey  tcell.Key
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

func NewActionBarItem(id int, text string, triggerRune rune, triggerKey tcell.Key) *ActionBarItem {
	item := &ActionBarItem{
		Id:          id,
		TriggerRune: triggerRune,
		TriggerKey:  triggerKey,
	}
	item.text = text
	item.bgStyle = DefaultStyle
	item.borderStyle = DefaultStyle.Foreground(tcell.ColorWhite).Background(AppBackground)
	return item
}

func ActionBarLabel(str string) string {
	if str == "" {
		return MessageLabelNone
	}
	return str
}

func (b *ActionBar) AddTextItem(id string, text string) {
	item := NewActionBarItem(len(b.items)+1, text, 0, 0)
	b.AddItem(item)
}

// TODO - refactor
func (b *ActionBar) AddItem(item *ActionBarItem) {
	item.bgStyle = DefaultStyle
	item.borderStyle = DefaultStyle.Foreground(tcell.ColorWhite).Background(AppBackground)
	// It's making problems if we want to have different styles.
	// TODO - We could think about introducing textMap argument where you can put styles and text together (or second array)
	item.text = fmt.Sprintf("%s%s", item.Text1, item.Text2)
	item.x = b.lastItemX
	item.x2 = item.x + len(item.text) + ActionBarItemPadding*2 + 1
	item.y = b.Y
	item.y2 = b.Y
	item.borderTop = b.vAlign == Bottom
	item.borderBottom = b.vAlign != Bottom
	b.lastItemX = b.getLastItemX()
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
	item.text = strings.Repeat(" ", len(item.Text1)+len(item.Text2))
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
		b.lastItemX = b.getLastItemX()
		b.Resize(b.screenX, b.screenY)
	}
}

func (b *ActionBar) TrimItemsTo(index int) {
	if index > 0 {
		b.items = b.items[:index]
		b.lastItemX = b.getLastItemX()
		b.Resize(b.screenX, b.screenY)
	}
}

func (b *ActionBar) Draw(screen tcell.Screen) {
	for _, item := range b.items {
		item.Draw(screen)
		DrawText(screen, item.x+2, item.y+1, item.Text1Style, item.Text1)
		DrawText(screen, item.x+2+len(item.Text1), item.y+1, item.Text2Style, item.Text2)
	}
}

func (b *ActionBar) Update() {
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
		break
	case Top:
		b.Y = 1
	}
	b.lastItemX = 0
	if b.hAlign == Right {
		b.lastItemX = screenX
	}
	for _, item := range b.items {
		item.y = b.Y
		//item.y2 = b.Y - ActionBarItemPadding + 1
		item.y2 = b.Y - 1 - ActionBarItemPadding
		item.x = b.lastItemX
		item.x2 = item.x + (len(item.text)+ActionBarItemPadding*2+1)*b.hAlign
		b.lastItemX = item.x2
	}
}

func (b *ActionBarItem) ChangeText(text1 string, text2 string) {
	b.Text2 = text2
	b.Text1 = text1
	b.text = fmt.Sprintf("%s%s", b.Text1, b.Text2)
}

func (b *ActionBar) getLastItemX() int {
	if len(b.items) == 0 {
		return 0
	}
	item := b.items[len(b.items)-1]
	return b.lastItemX + len(item.text) + ActionBarItemPadding*2
}
