package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"regexp"
)

const (
	EmptyLine        = ""
	BoldRune         = rune(999912)
	ItalicRune       = rune(999913)
	UnderlineRune    = rune(999914)
	StrikeTroughRune = rune(999915)
)

var (
	emptyStyle          tcell.Style
	regexpBold          = regexp.MustCompile("\\*(\\S.*?\\S)\\*")
	regexpItalic        = regexp.MustCompile("_(\\S.*?\\S)_")
	regexpUnderline     = regexp.MustCompile("\\+(\\S.*?\\S)\\+")
	regexpStrikeThrough = regexp.MustCompile("-(\\S.*?\\S)-")
	regexpListItem      = regexp.MustCompile("(?m)^\\*(.*?)$")
)

func DrawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
	for _, r := range text {
		if r == '\n' {
			row++
			col = x
			continue
		}
		screen.SetContent(col, row, r, nil, style)
		col++
	}
}

func DrawTextLimited(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) int {
	row := y1
	col := x1
	for _, r := range text {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		if s != nil {
			s.SetContent(col, row, r, nil, style)
		}
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
	return row
}

func DrawTextLimitedMarkdown(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) int {
	row := y1
	col := x1
	richStyle := emptyStyle
	richStartedRune := rune(0)
	for _, r := range text {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		richTagStyle := richTextRune(style, r)
		if richTagStyle != emptyStyle && richStartedRune == 0 {
			richStartedRune = r
			richStyle = richTagStyle
			continue
		} else if richTagStyle != emptyStyle && richStartedRune != 0 {
			richStartedRune = 0
			richStyle = emptyStyle
			continue
		}
		if s != nil {
			if richStyle == emptyStyle {
				s.SetContent(col, row, r, nil, style)
			} else {
				s.SetContent(col, row, r, nil, richStyle)
			}
		}
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
	return row
}

func DrawBox(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		screen.SetContent(col, y1, tcell.RuneHLine, nil, style)
		screen.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		screen.SetContent(x1, row, tcell.RuneVLine, nil, style)
		screen.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		screen.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		screen.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		screen.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		screen.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

func PrepareRichText(t string) string {
	t = regexpBold.ReplaceAllString(t, fmt.Sprintf("%s$1%s", string(BoldRune), string(BoldRune)))
	t = regexpItalic.ReplaceAllString(t, fmt.Sprintf("%s$1%s", string(ItalicRune), string(ItalicRune)))
	t = regexpUnderline.ReplaceAllString(t, fmt.Sprintf("%s$1%s", string(UnderlineRune), string(UnderlineRune)))
	t = regexpStrikeThrough.ReplaceAllString(t, fmt.Sprintf("%s$1%s", string(StrikeTroughRune), string(StrikeTroughRune)))
	t = regexpListItem.ReplaceAllString(t, " - $1")
	return t
}

func richTextRune(orgStyle tcell.Style, r rune) tcell.Style {
	switch r {
	case BoldRune:
		return orgStyle.Bold(true)
	case ItalicRune:
		return orgStyle.Italic(true)
	case UnderlineRune:
		return orgStyle.Underline(true)
	case StrikeTroughRune:
		return orgStyle.StrikeThrough(true)
	default:
		return emptyStyle
	}
}
