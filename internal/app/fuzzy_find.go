package app

import (
	"bytes"
	"fmt"
	"github.com/bep/debounce"
	"github.com/gdamore/tcell/v2"
	"github.com/sahilm/fuzzy"
	"strings"
	"time"
	"unicode"
)

type FuzzyFind struct {
	MarginTop        int
	MarginBottom     int
	Complete         chan FuzzyFindResult
	records          []string
	recordsProvider  func(query string) []string
	query            string
	fuzzyStatus      string
	title            string
	matches          fuzzy.Matches
	matchesAll       fuzzy.Matches
	buffer           bytes.Buffer
	dirty            bool
	selected         int
	screenX          int
	screenY          int
	supplierDebounce func(f func())
	debounceDisabled bool
}

type FuzzyFindResult struct {
	Index int
	Match string
}

const (
	ResultsMarginBottom     = 3
	WriteIndicator          = "> "
	MaxResults              = 4096
	DynamicSupplierDebounce = 50 * time.Millisecond
	SearchResultsPivot      = 6
)

var (
	boldMatchStyle   = DefaultStyle.Foreground(tcell.ColorLightGreen).Underline(true).Bold(true)
	boldRedStyle     = DefaultStyle.Foreground(tcell.ColorDarkRed).Bold(true)
	highlightDefault = DefaultStyle.Foreground(tcell.ColorWhite).Background(tcell.NewRGBColor(58, 58, 58))
	highlightBold    = highlightDefault.Foreground(tcell.ColorLightCyan).Bold(true)
	boldStyle        = DefaultStyle.Bold(true)
	titleStyle       = DefaultStyle.Italic(true).Foreground(tcell.NewRGBColor(236, 206, 88))
)

func NewFuzzyFind(title string, records []string) *FuzzyFind {
	matchesAll := make(fuzzy.Matches, 0, MaxResults)
	// TODO - not super optimize way to store results..
	for i, record := range records {
		matchesAll = append(matchesAll, fuzzy.Match{
			Str:   record,
			Index: i,
		})
	}
	return &FuzzyFind{
		Complete:         make(chan FuzzyFindResult),
		records:          records,
		query:            EmptyLine,
		fuzzyStatus:      "0/0",
		matches:          nil,
		selected:         0,
		dirty:            true,
		matchesAll:       matchesAll,
		recordsProvider:  nil,
		title:            title,
		MarginTop:        0,
		MarginBottom:     1,
		debounceDisabled: false,
	}
}

func NewFuzzyFindWithProvider(title string, recordsProvider func(query string) []string) *FuzzyFind {
	return &FuzzyFind{
		Complete:         make(chan FuzzyFindResult),
		records:          nil,
		query:            "init",
		fuzzyStatus:      "0/0",
		matches:          nil,
		selected:         0,
		dirty:            true,
		matchesAll:       make(fuzzy.Matches, 0, MaxResults),
		recordsProvider:  recordsProvider,
		supplierDebounce: debounce.New(DynamicSupplierDebounce),
		title:            title,
		MarginTop:        0,
		MarginBottom:     1,
		debounceDisabled: false,
	}
}

func (f *FuzzyFind) Draw(screen tcell.Screen) {
	if f.screenX == 0 || f.screenY == 0 {
		x, y := screen.Size()
		f.screenX = x
		f.screenY = y
	}
	f.drawRecords(screen)
	if f.title != "" {
		DrawText(screen, 2, f.screenY-ResultsMarginBottom-f.MarginBottom+1, titleStyle, f.title)
	}
	DrawText(screen, f.screenX-len(f.fuzzyStatus)-2, f.screenY-ResultsMarginBottom-f.MarginBottom+1, titleStyle, f.fuzzyStatus)
	DrawText(screen, 0, f.screenY-1-f.MarginBottom, boldStyle, WriteIndicator)
	DrawText(screen, 2, f.screenY-1-f.MarginBottom, DefaultStyle, f.query)
	screen.ShowCursor(2+len(f.query), f.screenY-1-f.MarginBottom)
}

func (f *FuzzyFind) Update() {
	if !f.dirty {
		return
	}
	buff := f.buffer.String()
	if f.recordsProvider != nil && (f.query != buff) {
		f.query = buff
		if f.debounceDisabled {
			f.updateRecordsFromSupplier()
		} else {
			f.supplierDebounce(f.updateRecordsFromSupplier)
			f.dirty = false
			return
		}
	}
	f.query = buff
	if len(f.query) == 0 {
		f.matches = f.matchesAll
	} else {
		f.matches = fuzzy.Find(f.query, f.records)
	}
	f.fuzzyStatus = fmt.Sprintf("%d/%d", len(f.matches), len(f.records))
	f.selected = ClampInt(f.selected, 0, f.matches.Len()-1)
	f.dirty = false
}

func (f *FuzzyFind) ForceUpdate() {
	f.dirty = true
	f.Update()
}

func (f *FuzzyFind) HandleKeyEvent(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEscape {
		f.Complete <- FuzzyFindResult{Index: -1, Match: ""}
	}
	if ev.Key() == tcell.KeyEnter {
		f.dirty = true
		f.Update()
		if len(f.matches) > 0 && f.selected >= 0 {
			match := f.matches[f.selected].Str
			index := findSelectedRecord(match, f.records)
			f.Complete <- FuzzyFindResult{Index: index, Match: match}
		} else {
			f.Complete <- FuzzyFindResult{Index: -1, Match: ""}
		}
	}
	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if f.buffer.Len() > 0 {
			f.buffer.Truncate(f.buffer.Len() - 1)
		}
		f.dirty = true
	}
	if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyTab {
		f.selected = ClampInt(f.selected+1, 0, f.matches.Len()-1)
		return
	}
	if ev.Key() == tcell.KeyDown || ev.Key() == tcell.KeyBacktab {
		f.selected = ClampInt(f.selected-1, 0, f.matches.Len()-1)
		return
	}
	if f.isEventWritable(ev) {
		f.buffer.WriteRune(ev.Rune())
		f.dirty = true
	}
}

func (f *FuzzyFind) Resize(screenX, screenY int) {
	f.screenX = screenX
	f.screenY = screenY
}

func (f *FuzzyFind) GetQuery() string {
	return f.query
}

func (f *FuzzyFind) GetSelectedItem() string {
	if len(f.records) == 0 {
		return ""
	}
	return f.records[f.selected]
}

func (f *FuzzyFind) SetDebounceDisabled(b bool) {
	f.debounceDisabled = b
}

func (f *FuzzyFind) drawRecords(screen tcell.Screen) {
	matchesLen := f.matches.Len()
	if matchesLen == 0 {
		return
	}
	var row = f.screenY - ResultsMarginBottom - f.MarginBottom
	var currentStyleDefault tcell.Style
	var currentStyleBold tcell.Style
	indexDelta := ClampInt(f.selected-row+SearchResultsPivot, 0, matchesLen-1)
	for index := indexDelta; index < matchesLen && row > f.MarginTop; index++ {
		match := f.matches[index]
		currentStyleDefault = DefaultStyle
		currentStyleBold = boldMatchStyle
		if index == f.selected {
			DrawText(screen, 0, row, boldRedStyle, WriteIndicator)
			currentStyleDefault = highlightDefault
			currentStyleBold = highlightBold
		}
		runeI := 0
		for i, s := range match.Str {
			if contains(i, match.MatchedIndexes) {
				DrawText(screen, runeI+2, row, currentStyleBold, string(s))
			} else {
				DrawText(screen, runeI+2, row, currentStyleDefault, string(s))
			}
			runeI++
		}
		row--
	}
}

func (f *FuzzyFind) updateRecordsFromSupplier() {
	f.records = f.recordsProvider(f.query)
	f.matchesAll = nil
	for i, record := range f.records {
		f.matchesAll = append(f.matchesAll, fuzzy.Match{
			Str:   record,
			Index: i,
		})
	}
	f.dirty = true
}

func (f *FuzzyFind) isEventWritable(ev *tcell.EventKey) bool {
	return unicode.IsLetter(ev.Rune()) || unicode.IsSpace(ev.Rune()) || unicode.IsDigit(ev.Rune()) ||
		ev.Rune() == '-' || ev.Rune() == '"' || ev.Rune() == '\'' || ev.Rune() == '&' ||
		ev.Rune() == ';' || ev.Rune() == '|' || ev.Rune() == '>' || ev.Rune() == '<' || ev.Rune() == '=' ||
		ev.Rune() == '!'
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func findSelectedRecord(result string, records []string) int {
	// TODO - impl faster alg
	var index int
	for i := range records {
		if strings.TrimSpace(records[i]) == result || records[i] == result {
			index = i
			break
		}
	}
	return index
}
