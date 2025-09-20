package app

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

type App struct {
	ScreenX      int
	ScreenY      int
	screen       tcell.Screen
	spinnerIndex int32
	keyEvent     chan *tcell.EventKey
	drawables    []Drawable
	flash        []Drawable
	systems      []System
	// clear/add/remove is less accurate execution than clear.
	// so it makes sense to store keep-alive stuff like this, instead of having
	// separate arrays to iterate through
	keepAlive    map[interface{}]bool
	changeMutex  sync.Mutex
	changeMutex2 sync.Mutex
	viewMutex    sync.Mutex
	quit         bool
	// re-render screen if true
	dirty           bool
	loading         bool
	closed          bool
	runOnAppRoutine []func()
	spinner         *SpinnerTCell
	view            View
	style           tcell.Style
}

const (
	FPS     = 30
	FPSMill = time.Second / FPS
)

var (
	appInstance *App
	once        sync.Once
)

func CreateNewApp() *App {
	once.Do(initApp)
	return appInstance
}

// CreateNewAppWithScreen accessible for testing
func CreateNewAppWithScreen(screen tcell.Screen) *App {
	initAppWithScreen(screen)
	once.Do(func() {
		// ... do nothing, complete 'once'
	})
	return appInstance
}

func GetApp() *App {
	return appInstance
}

func DefaultStyle() tcell.Style {
	return tcell.StyleDefault.Background(Color("default.background")).Foreground(Color("default.foreground"))
}

func initApp() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalln(err)
	}
	initAppWithScreen(screen)
}

func initAppWithScreen(screen tcell.Screen) {
	if os.Getenv("TERM") == "cygwin" {
		os.Setenv("TERM", "")
	}
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	MustLoadColorScheme()
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	screen.SetStyle(DefaultStyle())
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	s := NewSimpleSpinner()
	x, y := screen.Size()
	appInstance = &App{
		screen:          screen,
		ScreenX:         x,
		ScreenY:         y,
		spinnerIndex:    0,
		keyEvent:        make(chan *tcell.EventKey),
		runOnAppRoutine: make([]func(), 0, 64),
		drawables:       make([]Drawable, 0, 256),
		systems:         make([]System, 0, 128),
		flash:           make([]Drawable, 0, 5),
		keepAlive:       make(map[interface{}]bool),
		dirty:           true,
		spinner:         s,
		style:           DefaultStyle(),
	}
}

func (a *App) Start() {
	defer a.Close()
	go a.processTerminalEvents()
	go a.processOsSignals()
	defer a.PanicRecover()

	for {
		if a.quit {
			return
		}
		a.Render()
		if len(a.runOnAppRoutine) == 0 {
			time.Sleep(FPSMill)
			continue
		}
		funcsToRun := len(a.runOnAppRoutine) - 1
		for i := funcsToRun; i >= 0; i-- {
			a.runOnAppRoutine[i]()
		}
		if len(a.runOnAppRoutine) > funcsToRun {
			a.runOnAppRoutine = a.runOnAppRoutine[funcsToRun+1:]
			continue
		}
		a.runOnAppRoutine = nil
	}
}

func (a *App) Render() {
	a.screen.Show()
	for _, system := range a.systems {
		system.Update()
	}
	if !a.dirty && !a.loading {
		time.Sleep(FPSMill)
		return
	}
	a.screen.Fill(' ', a.style)
	if a.loading {
		a.spinner.Draw(a.screen)
	}
	for _, drawable := range a.drawables {
		drawable.Draw(a.screen)
	}
	for _, flash := range a.flash {
		flash.Draw(a.screen)
	}
	a.dirty = false
}

func (a *App) Close() {
	if a.closed {
		return
	}
	a.closed = true
	a.screen.DisableMouse()
	a.screen.Fill(' ', a.style)
	a.screen.Show()
	a.screen.Fini()
	close(a.keyEvent)
}

func (a *App) Loading(flag bool) {
	a.spinner.text = "Fetching"
	a.loading = flag
	a.setDirty()
}

func (a *App) IsLoading() bool {
	return a.loading
}

func (a *App) IsQuit() bool {
	return a.quit
}

func (a *App) LoadingWithText(flag bool, text string) {
	a.spinner.text = text
	a.loading = flag
}

func (a *App) SetView(view View) {
	a.viewMutex.Lock()
	a.setDirty()
	if a.view != nil {
		a.view.Destroy()
		delete(a.keepAlive, a.view)
		a.RemoveDrawable(a.view.(Drawable))
		a.RemoveSystem(a.view.(System))
	}
	a.view = view
	a.ClearNow()
	a.AddDrawable(view.(Drawable))
	a.AddSystem(view.(System))
	a.keepAlive[view] = true
	view.Init()
	a.viewMutex.Unlock()
}

func (a *App) CurrentView() interface{} {
	return a.view
}

func (a *App) KeepAlive(component interface{}) {
	a.changeMutex.Lock()
	a.keepAlive[component] = true
	a.changeMutex.Unlock()
}

func (a *App) UnKeepAlive(component interface{}) {
	a.changeMutex.Lock()
	delete(a.keepAlive, component)
	a.changeMutex.Unlock()
}

func (a *App) AddDrawable(drawable Drawable) {
	a.changeMutex.Lock()
	a.drawables = append(a.drawables, drawable)
	a.changeMutex.Unlock()
	if resizable, ok := drawable.(Resizable); ok {
		resizable.Resize(a.ScreenX, a.ScreenY)
	}
}

func (a *App) RemoveDrawable(drawable Drawable) {
	if a.keepAlive[drawable] {
		return
	}
	a.changeMutex.Lock()
	index := -1
	for i := range a.drawables {
		if a.drawables[i] == drawable {
			index = i
			break
		}
	}
	if index >= 0 {
		a.drawables = append(a.drawables[:index], a.drawables[index+1:]...)
	}
	a.changeMutex.Unlock()
}

func (a *App) AddFlash(flash Drawable, duration time.Duration) {
	defer a.PanicRecover()
	a.changeMutex.Lock()
	a.flash = append(a.flash, flash)
	a.changeMutex.Unlock()
	if resizable, ok := flash.(Resizable); ok {
		resizable.Resize(a.ScreenX, a.ScreenY)
	}
	timer := time.NewTimer(duration)
	a.setDirty()
	go func() {
		defer a.PanicRecover()
		<-timer.C
		a.changeMutex.Lock()
		a.flash = nil // it could lead to removing just-added flash message. For now, it's a good-enough solution
		a.changeMutex.Unlock()
		a.setDirty()
	}()
}

func (a *App) AddSystem(system System) {
	a.changeMutex.Lock()
	a.systems = append(a.systems, system)
	a.changeMutex.Unlock()
}

func (a *App) RemoveSystem(system System) {
	if a.keepAlive[system] {
		return
	}
	a.changeMutex.Lock()
	index := -1
	for i := range a.systems {
		if a.systems[i] == system {
			index = i
			break
		}
	}
	if index >= 0 {
		a.systems = append(a.systems[:index], a.systems[index+1:]...)
	}
	a.changeMutex.Unlock()
}

func (a *App) LastDrawable() Drawable {
	if len(a.drawables) == 0 {
		return nil
	}
	return a.drawables[len(a.drawables)-1]
}

func (a *App) SetDirty() {
	a.dirty = true
}

func (a *App) ClearNow() {
	a.setDirty()
	a.clear()
	// a.screen.Clear() is preserving terminal buffer (not alternate screen buffer) :/ different then in 1.3
	//a.screen.Clear()
	a.screen.Fill(' ', a.style)
	a.screen.HideCursor()
}

func (a *App) RunOnAppRoutine(f func()) {
	a.runOnAppRoutine = append(a.runOnAppRoutine, f)
}

func (a *App) Quit() {
	a.quit = true
}

func (a *App) PanicRecover() {
	rec := recover()
	if rec != nil {
		a.Close()
		panic(rec)
	}
}

func (a *App) setDirty() {
	a.dirty = true
	a.RunOnAppRoutine(func() {
		a.dirty = true
	})
}

func (a *App) clear() {
	a.changeMutex.Lock()
	a.drawables = nil
	a.systems = nil
	a.changeMutex.Unlock()
	if len(a.keepAlive) > 0 {
		a.changeMutex2.Lock()
		for s := range a.keepAlive {
			if _, ok := s.(System); ok {
				a.AddSystem(s.(System))
			}
			if _, ok := s.(Drawable); ok {
				a.AddDrawable(s.(Drawable))
			}
		}
		a.changeMutex2.Unlock()
	}
}

func (a *App) processTerminalEvents() {
	defer a.PanicRecover()
	for {
		if a.quit {
			return
		}
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.setDirty()
			a.screen.Sync()
			x, y := a.screen.Size()
			a.ScreenX = x
			a.ScreenY = y
			for _, s := range a.drawables {
				if ft, ok := (s).(Resizable); ok {
					ft.Resize(x, y)
				}
			}
		case *tcell.EventKey:
			a.setDirty()
			if ev.Key() == tcell.KeyCtrlC {
				a.Quit()
				return
			}
			if len(a.systems) == 0 && ev.Key() == tcell.KeyEscape {
				a.quit = true
			}
			// TODO - should keep only one array with components?
			for _, s := range a.systems {
				if ft, ok := (s).(KeyListener); ok {
					ft.HandleKeyEvent(ev)
				}
			}
		default:
			continue
		}
	}
}

func (a *App) processOsSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		<-signals
		a.quit = true
	}()
}
