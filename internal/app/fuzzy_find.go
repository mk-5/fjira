package app

import (
	"bytes"
	"fmt"
	"github.com/bep/debounce"
	"github.com/gdamore/tcell"
	"github.com/sahilm/fuzzy"
	"strings"
	"time"
	"unicode"
)

type FuzzyFind struct {
	MarginTop        int
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
}

type FuzzyFindResult struct {
	Index int
	Match string
}

const (
	FuzzyFindMarginBottom   = 3
	ResultsMarginBottom     = 3
	WriteIndicator          = "> "
	MaxResults              = 4096
	DynamicSupplierDebounce = 150 * time.Millisecond
	SearchResultsPivot      = 6
)

var (
	boldMatchStyle   = tcell.StyleDefault.Foreground(tcell.ColorLightGreen).Background(tcell.ColorDefault).Underline(true).Bold(true)
	boldRedStyle     = tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Background(tcell.ColorDefault).Bold(true)
	highlightDefault = tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorLightSlateGray)
	highlightBold    = tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Background(tcell.ColorLightSlateGray).Bold(true)
	boldStyle        = tcell.StyleDefault.Bold(true)
	titleStyle       = tcell.StyleDefault.Italic(true).Foreground(tcell.ColorOlive)
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
		Complete:        make(chan FuzzyFindResult),
		records:         records,
		query:           EmptyLine,
		fuzzyStatus:     "0/0",
		matches:         nil,
		selected:        0,
		dirty:           true,
		matchesAll:      matchesAll,
		recordsProvider: nil,
		title:           title,
		MarginTop:       0,
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
		DrawText(screen, 2, f.screenY-ResultsMarginBottom-FuzzyFindMarginBottom+1, titleStyle, f.title)
	}
	DrawText(screen, f.screenX-len(f.fuzzyStatus)-2, f.screenY-ResultsMarginBottom-FuzzyFindMarginBottom+1, titleStyle, f.fuzzyStatus)
	DrawText(screen, 0, f.screenY-1-FuzzyFindMarginBottom, boldStyle, WriteIndicator)
	DrawText(screen, 2, f.screenY-1-FuzzyFindMarginBottom, tcell.StyleDefault, f.query)
}

func (f *FuzzyFind) Update() {
	if !f.dirty {
		return
	}
	if f.query != f.buffer.String() && f.recordsProvider != nil {
		f.query = f.buffer.String()
		f.supplierDebounce(f.updateRecordsFromSupplier)
		f.dirty = false
		return
	}
	f.query = strings.TrimSpace(f.buffer.String())
	if len(f.query) == 0 {
		f.matches = f.matchesAll
	} else {
		f.matches = fuzzy.Find(f.query, f.records)
	}
	f.fuzzyStatus = fmt.Sprintf("%d/%d", len(f.matches), len(f.records))
	f.dirty = false
}

func (f *FuzzyFind) HandleKeyEvent(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEscape {
		f.Complete <- FuzzyFindResult{Index: -1, Match: ""}
	}
	if ev.Key() == tcell.KeyEnter {
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
			f.dirty = true
		}
	}
	if ev.Key() == tcell.KeyUp {
		f.selected = ClampInt(f.selected+1, 0, f.matches.Len()-1)
	}
	if ev.Key() == tcell.KeyDown {
		f.selected = ClampInt(f.selected-1, 0, f.matches.Len()-1)
	}
	if unicode.IsLetter(ev.Rune()) || unicode.IsSpace(ev.Rune()) || unicode.IsDigit(ev.Rune()) || ev.Rune() == '-' {
		f.buffer.WriteRune(ev.Rune())
		f.dirty = true
	}
}

func (f *FuzzyFind) Resize(screenX, screenY int) {
	f.screenX = screenX
	f.screenY = screenY
}

func (f *FuzzyFind) drawRecords(screen tcell.Screen) {
	var row = f.screenY - ResultsMarginBottom - FuzzyFindMarginBottom
	var currentStyleDefault tcell.Style
	var currentStyleBold tcell.Style
	indexDelta := ClampInt(f.selected-row+SearchResultsPivot, 0, f.matches.Len()-1)
	for index := indexDelta; index < f.matches.Len() && row > f.MarginTop; index++ {
		match := f.matches[index]
		currentStyleDefault = tcell.StyleDefault
		currentStyleBold = boldMatchStyle
		if index == f.selected {
			DrawText(screen, 0, row, boldRedStyle, WriteIndicator)
			currentStyleDefault = highlightDefault
			currentStyleBold = highlightBold
		}
		for i, s := range match.Str {
			if contains(i, match.MatchedIndexes) {
				DrawText(screen, i+2, row, currentStyleBold, string(s))
			} else {
				DrawText(screen, i+2, row, currentStyleDefault, string(s))
			}
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
