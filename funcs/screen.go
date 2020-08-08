package funcs

import (
	"github.com/gdamore/tcell"
)

// Addstr
// add text to the screen

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	arr := []rune(text)

	for i := x; i < len(arr)+x; i++ {
		s.SetContent(i, y, arr[i-x], []rune(""), style)
	}
}
