package app

import (
	"fmt"
	"github.com/gdamore/tcell"
)

const (
	EmptyLine = ""
)

func StartCli() {
	//fmt.Fprint(os.Stderr, "\x1b[?1049h")
	//fmt.Fprint(os.Stderr, "\x1b[?25l"+"\x1b[?1049h"+"\x1b[?25h")

	//ClearCli()
}

func ClearCli() {
	//fmt.Printf("\x1bc")   // bottom
	//fmt.Printf("\x1b[2J") // top
	//fmt.Printf("\\033c")
	//var cmd *exec.Cmd
	//switch runtime.GOOS {
	//case "windows":
	//	cmd = exec.Command("cls")
	//	break
	//default:
	//	cmd = exec.Command("clear")
	//	break
	//}
	//cmd.Stdout = os.Stdout
	//cmd.Run()
}

func ClearLine() {
	fmt.Printf("\n\033[1A\033[K")
	fmt.Printf("\r")
}

func DrawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = x
			continue
		}
		screen.SetContent(col, row, r, nil, style)
		col++
	}
}

func DrawTextLimited(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
