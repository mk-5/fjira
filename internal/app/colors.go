package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var (
	schemeMap map[string]interface{}
	colorsMap = map[string]tcell.Color{}
)

func Color(c string) tcell.Color {
	if color, ok := colorsMap[c]; ok {
		return color
	}
	parts := strings.Split(c, ".")
	var t interface{}
	var hex string
	t = schemeMap
	for _, p := range parts {
		if m, ok := t.(map[string]interface{}); ok {
			t = m[p]
		}
		if h, ok := t.(string); ok {
			hex = h
		}
	}
	color := tcell.GetColor(hex)
	colorsMap[c] = color
	return color
}

func MustLoadColorScheme() map[string]interface{} {
	d, _ := os.UserHomeDir()
	p := fmt.Sprintf("%s/.fjira/colors.yml", d)
	b, err := os.ReadFile(p)
	if err != nil {
		schemeMap = parseYMLStr(defaultColorsYML())
		return schemeMap
	}
	schemeMap = parseYMLStr(string(b))
	return schemeMap
}

func parseYMLStr(y string) map[string]interface{} {
	var yml map[string]interface{}
	err := yaml.Unmarshal([]byte(y), &yml)
	if err != nil {
		Error(err.Error())
		return map[string]interface{}{}
	}
	return yml
}

func defaultColorsYML() string {
	return `
default:
  background: "#161616"
  foreground: "#c7c7c7"
  foreground2: "#FFFFFF"

finder:
  cursor: "#8B0000"
  title: "#E9CE58"
  match: "#90EE90"
  highlight:
    background: "#3A3A3A"
    foreground: "#FFFFFF"
    match: "#E0FFFF"

navigation:
  top:
    background: "#5F875f"
    foreground1: "#FFFFFF"
    foreground2: "#151515"
  bottom:
    background: "#5F87AF"
    foreground1: "#FFFFFF"
    foreground2: "#151515"

details:
  foreground: "#696969"

boards:
  title:
    foreground: "#ecce58"
  headers:
    background: "#5F875f"
    foreground: "#FFFFFF"
  column:
    background: "#232323"
    foreground: "#ffffff"
  highlight:
    background: "#484848"
    foreground: "#ffffff"
  selection:
    background: "#8B0000"
    foreground: "#ffffff"

spinner:
  accent: "#FF0000"

alerts:
  success:
    background: "#F5F5F5"
    foreground: "#006400"
  error:
    background: "#F5F5F5"
    foreground: "#8B0000"
`
}
