/*

æ Welcome To The Lunae Source Code æ

------------------------------------

Info:

∙ This project was created and is maintained by geremachek (gmk).
∙ Lunae is registered under the GPL v. 3 witch means you are free to modify and distribute it and its source code
∙ This project was based on the now deprecated pluto file-manager and shares a lot of the core design ideas of that file-manager

More Info:

∙ <http://github.com/geremachek/lunae/>
∙ <http://geremachek.io/>

*/

package main

import (
	"github.com/gdamore/tcell"
	"os"
	"os/exec"
	"strings"
	"fmt"
	fu "lunae/funcs"
	co "lunae/config"
	cp "github.com/otiai10/copy"
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

	b2 = (height-co.YBuffBottom)
	fu.DispFiles(s, currFiles)
	if len(currFiles) != 0 {
		fu.DispBar(s, co.BarStyle, currFiles[currFile])
		fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
	}
	s.Show()

	for {
		nw, nh := s.Size()
		if width != nw && height == nh {
			width, height = nw, nh
			fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
		} else if height != nh && width != nh {
			width, height = nw, nh
			b1 = 0
			b2 = (height-co.YBuffBottom)+co.YBuffTop
			currY = co.YBuffTop
			currFile = 0
			fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
		}
		input := s.PollEvent()
		switch input := input.(type) {
			case *tcell.EventKey:
				if input.Rune() == co.KeyQuit {
					s.Fini()

					file, _ := os.Create("/tmp/lunar")
					file.WriteString(cwd)
					file.Close()

					fmt.Println(b1, b2, currFile)

					os.Exit(0)
				} else if input.Rune() == co.KeyDelete {
					os.RemoveAll(currFiles[currFile])

					if index, in := fu.In(currFiles[currFile], co.Selected); in {
						co.Selected = append(co.Selected[:index], co.Selected[index+1:]...)
					}

					currFiles = append(currFiles[:currFile], currFiles[currFile+1:]...)

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
				} else if input.Rune() == co.KeyCopy || input.Rune() == co.KeyMove || input.Rune() == co.KeyBulkDelete {
					for _, elem := range co.Selected {
						split := strings.Split(elem, "/")
						if input.Rune() == co.KeyCopy {
							cp.Copy(elem, cwd + "/" + split[len(split)-1])
						} else if input.Rune() == co.KeyMove {
							cp.Copy(elem, cwd + "/" + split[len(split)-1])
							os.RemoveAll(elem)
						} else if input.Rune() == co.KeyBulkDelete {
							os.RemoveAll(elem)
						}
					}

					co.Selected = []string{}
					currFiles = fu.GetFiles(cwd, dot)

					currY = co.YBuffTop
					currFile = 0

					b1, b2 = 0, height - co.YBuffBottom

					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeySelect {
					if fu.IsSel(cwd + "/" + currFiles[currFile]) {
						index, _ := fu.In(cwd + "/" + currFiles[currFile], co.Selected)
						co.Selected = append(co.Selected[:index], co.Selected[index+1:]...)
					} else {
						co.Selected = append(co.Selected, cwd + "/" + currFiles[currFile])
					}
					if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
						fu.Addstr(s, tcell.StyleDefault, co.XBuff, currY, fu.FormatText(s, currFiles[currFile]) + strings.Repeat(" ", len(co.SelectArrow)+1))
					} else {
						fu.Addstr(s, tcell.StyleDefault, co.XBuff, currY, fu.FormatText(s, currFiles[currFile]) + " ")
					}
					fu.SelFile(s, co.XBuff, currY, currFiles[currFile])
					s.Show()
				} else if input.Rune() == co.KeySelectAll {
					for _, elem := range currFiles {
						if _, in := fu.In(cwd + "/" + elem, co.Selected); !in {
							co.Selected = append(co.Selected, cwd + "/" + elem)
						}
					}

					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeyDeselectAll {
					co.Selected = []string{}

					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeyDotToggle {
					dot = !(dot)
					currFiles = fu.GetFiles(cwd, dot)
					b1 = 0
					b2 = (height-co.YBuffBottom)+co.YBuffTop
					currFile = 0
					currY = co.YBuffTop
					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeyGoToFirst {
					currFile = 0
					currY = co.YBuffTop

					b1 = 0
					b2 = height - co.YBuffTop

					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeyGoToLast {
					currFile = len(currFiles)-1

					if len(currFiles) + co.YBuffTop + co.YBuffBottom > height {
						currY = height - (co.YBuffBottom + 1)

						b1 = (len(currFiles) - len(currFiles[b1:b2])) + (co.YBuffTop*2)
						b2 = len(currFiles)
					} else {
						currY = currFile + co.YBuffTop
					}

					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
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
						fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
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
						fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
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
						var command []string = co.FileOpen["*"]
						splitFile := strings.Split(currFiles[currFile], ".")
						for k, v := range co.FileOpen {
							if k == splitFile[len(splitFile)-1] {
								command = v
								break
							}
						}
						replacedString := strings.Replace(command[1], "@", currFiles[currFile], -1)
						cmd := exec.Command(co.Shell, "-c", replacedString)
						if command[0] == "t" {
							cmd.Stdout = os.Stdout
							cmd.Stdin = os.Stdin
							s.Fini()
							cmd.Run()
							s, _ = tcell.NewScreen()
							s.Init()
							fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
						} else if command[0] == "g" {
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
					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else if input.Rune() == co.KeyRefresh && len(currFiles) != 0 {
					currFiles = fu.GetFiles(cwd, dot)
					fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
				} else {
					for k, v := range co.Bindings {
						if k == input.Rune() {
							replacedString := strings.Replace(v[1], "@", currFiles[currFile], -1)
							if v[0] == "cd" {
								os.Chdir(strings.Replace(string(v[1]), "~", os.Getenv("HOME"), -1))
								cwd, _ = os.Getwd()
								currFiles = fu.GetFiles(cwd, dot)
								b1 = 0
								b2 = (height-co.YBuffBottom)+co.YBuffTop
								currFile = 0
								currY = co.YBuffTop
								fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
							} else if v[0] == "t" {
								cmd := exec.Command(co.Shell, "-c", replacedString)
								cmd.Stdout = os.Stdout
								cmd.Stdin = os.Stdin
								s.Fini()
								cmd.Run()
								s, _ = tcell.NewScreen()
								s.Init()
								fu.DrawScreen(s, currFiles, currFile, currY, b1, b2)
							} else if v[0] == "g" {
								cmd := exec.Command(co.Shell, "-c", replacedString)
								cmd.Start()
							}
						}
					}
				}
		}
	}
}
