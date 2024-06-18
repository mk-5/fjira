package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	os2 "github.com/mk-5/fjira/internal/os"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	schemeMap map[string]interface{}
	colorsMap = map[string]tcell.Color{}
)

func Color(c string) tcell.Color {
	if len(colorsMap) == 0 {
		MustLoadColorScheme()
	}
	if color, ok := colorsMap[c]; ok {
		return color
	}
	panic("unknown color " + c)
}

func MustLoadColorScheme() map[string]interface{} {
	d := os2.MustGetUserHomeDir()
	p := fmt.Sprintf("%s/.fjira/colors.yml", d)
	b, err := os.ReadFile(p)
	if err != nil {
		schemeMap = parseYMLStr(defaultColorsYML())
	} else {
		schemeMap = parseYMLStr(string(b))
	}
	colorsMap = parseYamlToDotNotationMap("", schemeMap, colorsMap)
	return schemeMap
}

func parseYamlToDotNotationMap(prefix string, yml map[string]interface{}, targetMap map[string]tcell.Color) map[string]tcell.Color {
	var key string
	for k, v := range yml {
		if prefix == "" {
			key = k
		} else {
			key = fmt.Sprintf("%s.%s", prefix, k)
		}
		if m, ok := v.(map[string]interface{}); ok {
			targetMap = parseYamlToDotNotationMap(key, m, targetMap)
		}
		if h, ok := v.(string); ok {
			targetMap[key] = tcell.GetColor(h)
		}
	}
	return targetMap
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
