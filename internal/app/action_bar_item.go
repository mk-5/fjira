package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type ActionBarItem struct {
	Id          int
	Text1       string
	Text2       string
	Text1Style  tcell.Style
	Text2Style  tcell.Style
	TriggerRune rune
	TriggerKey  tcell.Key
	text        string
	x           int
	y           int
}

func NewActionBarItem(id int, text string, triggerRune rune, triggerKey tcell.Key) *ActionBarItem {
	item := &ActionBarItem{
		Id:          id,
		TriggerRune: triggerRune,
		TriggerKey:  triggerKey,
	}
	item.Text1 = text
	return item
}

func ActionBarLabel(str string) string {
	if str == "" {
		return MessageLabelNone
	}
	return str
}

func (b *ActionBarItem) ChangeText(text1 string, text2 string) {
	b.Text2 = text2
	b.Text1 = text1
	b.text = fmt.Sprintf("%s%s", b.Text1, b.Text2)
}
