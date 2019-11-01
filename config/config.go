package config

import "github.com/gdamore/tcell"

var (
	XBuff int = 0
	YBuffTop int = 0
	YBuffBottom int = 2

	Selected []string

	BarBg tcell.Color = tcell.GetColor("#ff0000")
	BarFg tcell.Color = tcell.GetColor("#0000ff")
	BarDiv string = " "
	BarStyle = map[string]tcell.Style{
		"1cwd": tcell.StyleDefault.Background(BarBg).Bold(true),
		"2size": tcell.StyleDefault.Background(BarBg),
		"3mode": tcell.StyleDefault.Background(BarBg),
		//"3[file @]": tcell.StyleDefault.Background(tcell.GetColor(BarBg)),
	}

	FileOpen = map[string]string {
		"*": "vi,t",
	}

	FileColors = map[string]tcell.Style {
		"[dir]": tcell.StyleDefault.Foreground(tcell.GetColor("#508cbe")).Bold(true),
	}
)
