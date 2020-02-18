/*

+--------------------------+
| Catfm Configuration File |
+--------------------------+

*/

package config

import "github.com/gdamore/tcell"

var (
	XBuff int = 0 // Blank space on sides
	YBuffTop int = 0 // Blank space on top
	YBuffBottom int = 2 // Blank space on bottom (remember to make room for the bar!)

	// Key bindings (arrow keys will always work for the movement functions)

	KeyRefresh rune = 'f'
	KeyQuit rune = 'q'
	KeyDelete rune = 'd'
	KeyBulkDelete rune = 'D'
	KeyCopy rune = 'C'
	KeyMove rune = 'M'
	KeySelect rune = ' '
	KeySelectAll rune = '*'
	KeyDeselectAll rune = '-'
	KeyDotToggle rune = '.'
	KeyGoToFirst rune = 'g'
	KeyGoToLast rune = 'G'
	KeyLeft rune = 'h'
	KeyDown rune = 'j'
	KeyUp rune = 'k'
	KeyRight rune = 'l'

	BarLocale = "bottom" // top or bottom ("" to disable the bar)
	BarBg tcell.Color = tcell.ColorWhite
	BarFg tcell.Color = tcell.GetColor("#000000")
	BarDiv string = " "
	BarStyle = map[string]tcell.Style{
		// Make sure to order your elements!
		// To get the output of a command/script encolse the script in brackets
		// To access the currently selected file use the "@" symbol, i.e. 1[ls -l @]

		"1cwd": tcell.StyleDefault.Background(BarBg).Bold(true),
		"2size": tcell.StyleDefault.Background(BarBg),
		"3mode": tcell.StyleDefault.Background(BarBg),

		// Other elements: total, tab
		// Variables: $HOST, $USER, $FILE, $TAB
	}

	SelectType string = "default" // full, default, arrow, arrow-default
	SelectStyle tcell.Style = tcell.StyleDefault.Reverse(true)
	SelectArrow string = "> "
	SelectArrowStyle tcell.Style = tcell.StyleDefault.Bold(true)

	PipeType = "" // thick, thin, hollow, round (empty string to disable)
	PipeStyle tcell.Style = tcell.StyleDefault
	PipeText = "" // user@host, catfm@host, user@catfm, or any text
	PipeTextStyle = tcell.StyleDefault

	Shell string = "sh" // Shell that will be used to execute commands
	TildeHome bool = true // replace /home/user with "~"

	FileOpen = map[string][]string {
		// Key is the file type, formatted like "jpg"
		// Value is the program,t/g. "t" for terminal, "g" for gui

		"*": []string{"t", "vi '@'"}, // vi, a terminal program will open all files. (the '@' symbol will be replaced with the currently selected file)
	}

	FileColors = map[string]tcell.Style {
		"[dir]": tcell.StyleDefault.Foreground(tcell.GetColor("#508cbe")).Bold(true),
	}

	Bindings = map[rune][]string {
		'~': []string{"cd", "~"}, // "cd" into the home directory when the user presses '~'
		'v': []string{"t", "less '@'"}, // View the selected file in less when 'v' is pressed
	}
)
