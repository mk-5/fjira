package app

import "github.com/gdamore/tcell/v2"

type View interface {
	Init()
	Destroy()
}

type Drawable interface {
	Draw(screen tcell.Screen)
}

type Resizable interface {
	Resize(screenX, screenY int)
}

type System interface {
	Update()
}

type KeyListener interface {
	HandleKeyEvent(keyEvent *tcell.EventKey)
}

type Dirtyable struct {
	dirty bool
}
