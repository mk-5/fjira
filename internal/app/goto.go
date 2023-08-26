package app

type goToHistory struct {
	screenName string
	args       []interface{}
}

var (
	gotoRegistry = map[string]func(args ...interface{}){}
	currentGoTo  = &goToHistory{}
	previousGoTo = &goToHistory{}
)

func RegisterGoto(name string, f func(args ...interface{})) {
	gotoRegistry[name] = f
}

// GoTo it's not a perfect solution ... but it's the only one
// that works, and not lead into cycled-import errors.
// For example, you can go from projects view into issues view, and from issues view into projects view.
// Both views are in different packages therefor it leas to cyclic-import
func GoTo(name string, args ...interface{}) {
	defer GetApp().PanicRecover()
	if f, ok := gotoRegistry[name]; ok {
		f(args...)
		previousGoTo.screenName = currentGoTo.screenName
		previousGoTo.args = currentGoTo.args
		currentGoTo.screenName = name
		currentGoTo.args = args
	}
}

func CurrentScreenName() string {
	return currentGoTo.screenName
}

func PreviousScreenName() string {
	return previousGoTo.screenName
}

func GoBack() {
	if previousGoTo.screenName != "" {
		GoTo(previousGoTo.screenName, previousGoTo.args)
	}
}
