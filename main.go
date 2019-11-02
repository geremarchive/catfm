/*

æ Welcome To The Lunae Source Code æ

------------------------------------

Info:

∙ This project was created and is maintained by geremachek (gmk).
∙ Lunae is registered under the GPL v. 3 witch means you are free to modify and distribute it and its source code
∙ This project was based on the now deprecated pluto project and shares a lot of the core design ideas from that project

More Info:

∙ <http://github.com/geremachek/lunae/>
∙ <http://geremachek.io/>

*/

package main

import (
	"github.com/gdamore/tcell"
	"os"
	"os/exec"
	"fmt"
	"strings"
	fu "./funcs"
	co "./config"
)

var (
	currFile int
	currY int = co.YBuffTop

	b1 int
	b2 int

	width int
	height int

	dot bool = false
)

func main() {
	cwd, _ := os.Getwd()
	currFiles := fu.GetFiles(cwd, dot)

	s, _ := tcell.NewScreen()
	s.Init()
	width, height = s.Size()

	b2 = (height-co.YBuffBottom)//+co.YBuffTop
	fu.DispFiles(s, currFiles)
	fu.DispBar(s, co.BarStyle, currFiles[currFile])
	fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
	s.Show()

	for {
		nw, nh := s.Size()
		if width != nw || height != nh {
			width, height = nw, nh
			b1 = 0
			b2 = (height-co.YBuffBottom)+co.YBuffTop
			currY = co.YBuffTop
			currFile = 0
			fu.DrawScreen(s, currFiles, currFile, currY)
		}
		input := s.PollEvent()
		switch input := input.(type) {
			case *tcell.EventKey:
				if input.Rune() == co.KeyQuit {
					s.Fini()
					fmt.Print()
					os.Exit(0)
				} else if input.Rune() == co.KeyDelete {
					os.RemoveAll(currFiles[currFile])
					currFiles = fu.GetFiles(cwd, dot)
					s.Clear()
					if currFile == len(currFiles) && b1 == 0 {
						currY -= 1
						currFile -= 1
						fu.DispFiles(s, currFiles)
					} else if b1 == 0 {
						fu.DispFiles(s, currFiles)
					} else if currFile == len(currFiles) {
						b1 -= 1
						b2 -= 1
						currFile -= 1
						fu.DispFiles(s, currFiles[b1:b2])
					} else {
						fu.DispFiles(s, currFiles[b1:b2])
					}
					if len(currFiles) != 0 {
						fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
						fu.DispBar(s, co.BarStyle, currFiles[currFile])
					}
					s.Show()
				} else if input.Rune() == co.KeyGroupDelete {
					for _, elem := range co.Selected {
						os.RemoveAll(elem)
						if strings.Contains(elem, cwd) {
							currFile -= 1
							if b1 != 0 {
								b1 -= 1
								b2 -= 1
							} else {
								currY -= 1
							}
						}
					}
					co.Selected = []string{}
					currFiles = fu.GetFiles(cwd, dot)
					s.Clear()
					if b1 == 0 {
						fu.DispFiles(s, currFiles)
					} else {
						fu.DispFiles(s, currFiles[b1:b2])
					}
					if len(currFiles) != 0 {
						fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
						fu.DispBar(s, co.BarStyle, currFiles[currFile])
					}
					s.Show()
					// This may be dangerous, but i don't really care :)
				} else if input.Rune() == co.KeyCopy || input.Rune() == co.KeyMove {
					for _, elem := range co.Selected {
						split := strings.Split(elem, "/")
						if input.Rune() == co.KeyCopy {
							fu.Copy(elem, cwd + "/" + split[len(split)-1])
						} else if input.Rune() == co.KeyMove {
							fu.Move(elem, cwd + "/" + split[len(split)-1])
						}
					}
					co.Selected = []string{}
					currFiles = fu.GetFiles(cwd, dot)
					s.Clear()
					if b1 == 0 {
						fu.DispFiles(s, currFiles)
					} else {
							fu.DispFiles(s, currFiles[b1:b2])
					}
					fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
					fu.DispBar(s, co.BarStyle, currFiles[currFile])
					s.Show()
				} else if input.Rune() == co.KeySelect {
					if fu.IsSel(cwd + "/" + currFiles[currFile]) {
						index, _ := fu.In(cwd + "/" + currFiles[currFile], co.Selected)
						co.Selected = append(co.Selected[:index], co.Selected[index+1:]...)
					} else {
						co.Selected = append(co.Selected, cwd + "/" + currFiles[currFile])
					}
					fu.Addstr(s, tcell.StyleDefault, co.XBuff, currY, fu.FormatText(s, currFiles[currFile]) + "  ")
					fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
					s.Show()
				} else if input.Rune() == co.KeyDotToggle {
					dot = !(dot)
					currFiles = fu.GetFiles(cwd, dot)
					b1 = 0
					b2 = (height-co.YBuffBottom)+co.YBuffTop
					currFile = 0
					currY = co.YBuffTop
					fu.DrawScreen(s, currFiles, currFile, currY)
				} else if (input.Key() == tcell.KeyDown || input.Rune() == co.KeyDown) && len(currFiles) != 0 {
					if currFile == len(currFiles)-1 {
						continue
					} else if currY == (height-1)-co.YBuffBottom {
						b1 += 1
						if b2 >= len(currFiles)-1 {
							b2 = len(currFiles)-1
						} else {
							b2 += 1
						}
						currFile += 1
						fu.DrawScreen(s, currFiles, currFile, currY)
					} else {
						fu.DSelFile(s, co.XBuff, currY, currFiles[currFile])
						currY += 1
						currFile += 1
						fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
						fu.DispBar(s, co.BarStyle, currFiles[currFile])
						s.Show()
					}

				} else if (input.Key() == tcell.KeyUp || input.Rune() == co.KeyUp) && len(currFiles) != 0 {
					if currFile == 0 {
						continue
					} else if currY == co.YBuffTop {
						b1 -= 1
						b2 -= 1
						currFile -= 1
						fu.DrawScreen(s, currFiles, currFile, currY)
					} else {
						fu.DSelFile(s, co.XBuff, currY, currFiles[currFile])
						currY -= 1
						currFile -= 1
						fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
						fu.DispBar(s, co.BarStyle, currFiles[currFile])
						s.Show()

					}
				} else if (input.Key() == tcell.KeyRight || input.Rune() == co.KeyRight) && len(currFiles) != 0 {
					if fu.Isd(currFiles[currFile]) {
						os.Chdir(currFiles[currFile])
						cwd, _ = os.Getwd()
						currFiles = fu.GetFiles(cwd, dot)
						b1 = 0
						b2 = (height-co.YBuffBottom)+co.YBuffTop
						currFile = 0
						currY = co.YBuffTop
						s.Clear()
						fu.DispFiles(s, currFiles)
						if len(currFiles) != 0 {
							fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
							fu.DispBar(s, co.BarStyle, currFiles[currFile])
						}
						s.Show()
					} else {
						var command []string = strings.Split(co.FileOpen["*"], ",")
						splitFile := strings.Split(currFiles[currFile], ".")
						for k, v := range co.FileOpen {
							if k == splitFile[len(splitFile)-1] {
								command = strings.Split(v, ",")
								break
							}
						}
						cmd := exec.Command(command[0], currFiles[currFile])
						if command[1] == "t" {
							cmd.Stdout = os.Stdout
							cmd.Stdin = os.Stdin
							s.Fini()
							cmd.Run()
							s, _ = tcell.NewScreen()
							s.Init()
							fu.DispFiles(s, currFiles)
							fu.DispBar(s, co.BarStyle, currFiles[currFile])
							fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
							s.Show()
						} else if command[1] == "g" {
							cmd.Start()
						}
					}
				} else if input.Key() == tcell.KeyLeft || input.Rune() == co.KeyLeft {
					os.Chdir("..")
					cwd, _ = os.Getwd()
					currFiles = fu.GetFiles(cwd, dot)
					b1 = 0
					b2 = (height-co.YBuffBottom)+co.YBuffTop
					currFile = 0
					currY = co.YBuffTop
					fu.DrawScreen(s, currFiles, currFile, currY)
				} else if input.Rune() == co.KeyRefresh && len(currFiles) != 0 {
					currFiles = fu.GetFiles(cwd, dot)
					fu.DrawScreen(s, currFiles, currFile, currY)
				} else {
					for k, v := range co.Bindings {
						if k == input.Rune() {
							if v[0] == "cd" {
								os.Chdir(strings.Replace(string(v[1]), "~", os.Getenv("HOME"), -1))
								cwd, _ = os.Getwd()
								currFiles = fu.GetFiles(cwd, dot)
								b1 = 0
								b2 = (height-co.YBuffBottom)+co.YBuffTop
								currFile = 0
								currY = co.YBuffTop
								fu.DrawScreen(s, currFiles, currFile, currY)
							}
						}
					}
				}
		}
	}
}
