package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-sixel"
	"image"
	"log"
	"math"
	"net/http"
	"os"
)

import _ "image/jpeg"
import _ "image/png"
import _ "golang.org/x/image/webp"

// https://github.com/rivo/tview/blob/master/image.go
type Image struct {
	width  int
	height int
	data   *bytes.Buffer
}

func NewImage(u string) (*Image, error) {
	i, err := readImageFromWeb(u)
	if err != nil {
		return nil, err
	}

	img := &Image{
		width:  i.Bounds().Dx(),
		height: i.Bounds().Dy(),
		data:   bytes.NewBuffer(nil),
	}

	enc := sixel.NewEncoder(img.data)
	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return img, nil
}

func (i *Image) Draw(s tcell.Screen) {
	tty, ok := s.Tty()
	if !ok {
		s.Fini()
		log.Fatal("not a terminal")
	}
	ws, err := tty.WindowSize()
	if err != nil {
		s.Fini()
		log.Fatal(err)
	}
	// Get the dimensions of a single cell
	cw, ch := ws.CellDimensions()
	if cw == 0 || ch == 0 {
		s.Fini()
		log.Fatal("terminal does not support sixel graphics")
		return
	}

	sixelWidth := int(math.Ceil(float64(i.width) / float64(cw)))
	sixelHeight := int(math.Ceil(float64(i.height) / float64(ch)))

	sixelX := ws.Width/2 - (sixelWidth / 2) // Center the image horizontally
	sixelY := ws.Height/2 - sixelHeight - 2
	if sixelY < 0 {
		sixelY = 0
	}

	s.LockRegion(sixelX, sixelY, sixelWidth, sixelHeight, true)

	// Get the terminfo for our current terminal
	ti, err := tcell.LookupTerminfo(os.Getenv("TERM"))
	if err != nil {
		s.Fini()
		log.Fatal(err)
	}

	ti.TPuts(tty, ti.TGoto(sixelX, sixelY))
	ti.TPuts(tty, i.data.String())
}

func (i *Image) Update() {
	// ...
}

func (i *Image) Resize(screenX, screenY int) {
	// ...
}

func readImageFromWeb(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, err
}
