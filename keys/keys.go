package keys

import "github.com/gdamore/tcell"

func MatchKey(input *tcell.EventKey, key string) bool {
	if input.Key() == tcell.KeyRune {
		return input.Rune() == []rune(key)[0]
	} else {
		return input.Key() == Keys[key]
	}
}

var (
	Keys = map[string]tcell.Key {
		"ctrl-a": tcell.KeyCtrlA,
		"ctrl-b": tcell.KeyCtrlB,
		"ctrl-c": tcell.KeyCtrlC,
		"ctrl-d": tcell.KeyCtrlD,
		"ctrl-e": tcell.KeyCtrlE,
		"ctrl-f": tcell.KeyCtrlF,
		"ctrl-g": tcell.KeyCtrlG,
		"ctrl-h": tcell.KeyCtrlH,
		"ctrl-i": tcell.KeyCtrlI,
		"ctrl-j": tcell.KeyCtrlJ,
		"ctrl-k": tcell.KeyCtrlK,
		"ctrl-l": tcell.KeyCtrlL,
		"ctrl-m": tcell.KeyCtrlM,
		"ctrl-n": tcell.KeyCtrlN,
		"ctrl-o": tcell.KeyCtrlO,
		"ctrl-p": tcell.KeyCtrlP,
		"ctrl-q": tcell.KeyCtrlQ,
		"ctrl-r": tcell.KeyCtrlR,
		"ctrl-s": tcell.KeyCtrlS,
		"ctrl-t": tcell.KeyCtrlT,
		"ctrl-u": tcell.KeyCtrlU,
		"ctrl-v": tcell.KeyCtrlV,
		"ctrl-w": tcell.KeyCtrlW,
		"ctrl-x": tcell.KeyCtrlX,
		"ctrl-y": tcell.KeyCtrlY,
		"ctrl-z": tcell.KeyCtrlZ,
		"f1": tcell.KeyF1,
		"f2": tcell.KeyF2,
		"f3": tcell.KeyF3,
		"f4": tcell.KeyF4,
		"f5": tcell.KeyF5,
		"f6": tcell.KeyF6,
		"f7": tcell.KeyF7,
		"f8": tcell.KeyF8,
		"f9": tcell.KeyF9,
		"f10": tcell.KeyF10,
		"f11": tcell.KeyF11,
		"f12": tcell.KeyF12,
		"home": tcell.KeyHome,
	}
)
