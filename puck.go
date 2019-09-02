package main

import (

	"fmt"
	"os"
	"sort"
	"strings"
	"io/ioutil"
        "github.com/gdamore/tcell"
)

type Format struct {
	What string
	Special string
	Fg string
	Bg string
	Bold bool
	Dim bool
	Reverse bool
	Underline bool
}

func format2style(f Format) (out tcell.Style) {
	out = out.Foreground(tcell.GetColor(f.Fg)).Background(tcell.GetColor(f.Bg)).Bold(f.Bold).Dim(f.Dim).Reverse(f.Reverse).Underline(f.Underline)
	return
}

func getFiles(path string) (out []string) {
	var dirs []string
	var file []string
	files, _ := ioutil.ReadDir(path)
	for _, elem := range files {
		if elem.IsDir() {
			dirs = append(dirs, string(elem.Name()))
		} else {
			file = append(file, string(elem.Name()))
		}
	}
	sort.Strings(dirs)
	sort.Strings(file)
	out = append(dirs, file...)
	return
}

func addstr(screen tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		screen.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func isd(path string) bool {
	f, _ := os.Stat(path)
	if f.IsDir() {
		return true
	} else {
		return false
	}
}

func sizeText(screen tcell.Screen, dir bool, buffer int, text string) string {
	width, _ := screen.Size()
	if dir {
		if len(text) > width-(buffer*2)-1 {
			return text[:width-(buffer*2)-4] + ".../"
		} else {
			return text+"/"
		}
	} else {
		if len(text) > width-(buffer*2) {
			return text[:width-(buffer*2)-3] + "..."
		} else {
			return text
		}
	}
}

func getSchemeIndex(scheme []Format, normal int, file string) (out int) {
	for i, elem := range scheme {
		if strings.Contains(file, ".") {
			if elem.What == strings.Split(file, ".")[1] {
				out = i
			} else {
				out = normal
			}

		} else {
			out = normal
		}
	}
	return
}

func dispFiles(screen tcell.Screen, scheme []Format, files []string, top int, bottom int, side int, dir int, normal int) {
	_, height := screen.Size()
	if len(files) == 0 {
		addstr(screen, tcell.StyleDefault.Background(tcell.GetColor("#ff0000")).Bold(true), side, top, sizeText(screen, false, side, "There's nothing here!"))
	} else {
		for i, elem := range files {
			if i != height - (bottom) {
				if isd(elem) {
					addstr(screen, format2style(scheme[dir]), side, i+top, sizeText(screen, true, side, elem))
				} else {

					addstr(screen, format2style(scheme[getSchemeIndex(scheme, normal, elem)]), side, i+top, sizeText(screen, false, side, elem))
				}
			} else {
				break
			}
		}
	}
}

func rstrx(str string, x int) (out string) {
	for i := 0; i < x; i++ {
		out += str
	}
	return
}

func selFile(screen tcell.Screen, selStyle Format, x int, y int, file string) {
	style := format2style(selStyle)
	width, _ := screen.Size()
	if isd(file) {
		text := sizeText(screen, true, x, file)
		if selStyle.Special == "full" {
			addstr(screen, style, x, y, text+rstrx(" ", width-len(text)-(x*2)))
		} else {
			addstr(screen, style, x, y, text)
		}
	} else {
		text := sizeText(screen, false, x, file)
		if selStyle.Special == "full" {
			addstr(screen, style, x, y, text+rstrx(" ", width-len(text)-(x*2)))
		} else {
			addstr(screen, style, x, y, text)
		}
	}
}

func dSelFile(screen tcell.Screen, y int, index int, side int, dir int, normal int, array []string, scheme []Format) {
	width, _ := screen.Size()
	addstr(screen, tcell.StyleDefault, 0, y, rstrx(" ", width))
	if isd(array[index]) {
		addstr(screen, format2style(scheme[getSchemeIndex(scheme, dir, array[index])]), 0, y, sizeText(screen, true, side, array[index]))
	} else {
		addstr(screen, format2style(scheme[getSchemeIndex(scheme, normal, array[index])]), 0, y, sizeText(screen, false, side, array[index]))
	}

}

func main() {
	s, _ := tcell.NewScreen()
	s.Init()
	s.Clear()
	_, height := s.Size()
	var (
		dir int
		normal int
		top int = 0
		bottom int = 0
		side int = 0
	)
	var (
		currY int = top
		currFile int
		b1 int
		b2 int = (height - bottom)-top
	)
	cwd, _ := os.Getwd()
	mainList := getFiles(cwd)
	selStyle := Format{"", "full", "#ffffff", "#878080", false, false, false, false}
	defScheme := []Format{ Format{"[dir]", "", "#94bff3", "", true, false, false, false}, Format{"[normal]", "", "", "", false, false, false, false}, Format{"jpg", "", "#ff00ff", "", false, false, false, false} }
	for i, elem := range defScheme {
		if elem.What == "[dir]" {
			dir = i
		} else if elem.What == "[normal]" {
			normal = i
		} else {
			continue
		}
	}
	dispFiles(s, defScheme, mainList, top, bottom, side, dir, normal)
	selFile(s, selStyle, side, currY, mainList[currFile])
	s.Show()

	for {
		_, height := s.Size()
		input := s.PollEvent()
		switch input := input.(type) {
			case *tcell.EventKey:
				if input.Rune() == 'q' {
					s.Fini()
					fmt.Println(currY, currFile)
					os.Exit(0)
				} else if input.Key() == tcell.KeyDown {
					if currFile == len(mainList)-1 {
						continue
					} else if currY == (height-1)-bottom {
						b1 += 1
						b2 += 1
						currFile += 1
						s.Clear()
						dispFiles(s, defScheme, mainList[b1:b2], top, bottom, side, dir, normal)
						selFile(s, selStyle, side, currY, mainList[currFile])
						s.Show()
					} else {
						dSelFile(s, currY, currFile, side, dir, normal, mainList, defScheme)
						currY += 1
						currFile += 1
						selFile(s, selStyle, side, currY, mainList[currFile])
						s.Show()
					}
				} else if input.Key() == tcell.KeyUp {
					if currFile == 0 {
						continue
					} else if currY == top && currFile != 0 {
						b1 -= 1
						b2 -= 1
						currFile -= 1
						s.Clear()
						dispFiles(s, defScheme, mainList[b1:b2], top, bottom, side, dir, normal)
						selFile(s, selStyle, side, currY, mainList[currFile])
						s.Show()
					} else {
						dSelFile(s, currY, currFile, side, dir, normal, mainList, defScheme)
						currY -= 1
						currFile -= 1
						selFile(s, selStyle, side, currY, mainList[currFile])
						s.Show()
					}
				} else if input.Key() == tcell.KeyRight {
					if isd(mainList[currFile]) {
						os.Chdir(mainList[currFile])
						cwd, _ = os.Getwd()
						mainList = getFiles(cwd)
						b1, b2, currFile, currY = 0, (height - bottom)-top, 0, top
						s.Clear()
						dispFiles(s, defScheme, mainList, top, bottom, side, dir, normal)
						selFile(s, selStyle, side, currY, mainList[currFile])
						s.Show()
					} else {
						continue
					}
				}
		}
	}
}
