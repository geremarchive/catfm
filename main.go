/*

Welcome To The Catfm Source Code

------------------------------------

Info:

∙ This project was created and is maintained by geremachek (gmk).
∙ Catfm is registered under the GPL v. 3 witch means you are free to modify and distribute it and its source code
∙ This project was based on the now deprecated pluto file-manager and shares a lot of the core design ideas of that file-manager

More Info:

∙ <http://github.com/geremachek/catfm/>
∙ <http://geremachek.io/software>

*/

package main

import (
	"github.com/gdamore/tcell"
	"os"
	"os/exec"
	"strings"
	"fmt"
	fu "catfm/funcs"
	co "catfm/config"
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
	if len(os.Args) > 1 {
		err := os.Chdir(os.Args[1])

		if err != nil {
			fmt.Println("Couldn't open '" + os.Args[1] + "'")
		}
	}

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	currFiles, err := fu.GetFiles(cwd, dot)

	if err != nil {
		panic(err)
	}

	s, err := tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	s.Init()

	s.Fini()

	s, err = tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	s.Init()

	defer s.Fini()

	width, height = s.Size()
	b2 = (height-co.YBuffBottom)

	if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
		panic(err)
	}

	s.Show()

	for {
		nw, nh := s.Size()
		if width != nw && height == nh {
			width, height = nw, nh

			if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
				panic(err)
			}
			fu.BorderPipes(s)
		} else if height != nh && width != nh {
			width, height = nw, nh
			b1 = 0
			b2 = (height-co.YBuffBottom)+co.YBuffTop
			currY = co.YBuffTop
			currFile = 0

			if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
				panic(err)
			}

			fu.BorderPipes(s)
		}
		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				if input.Rune() == co.KeyQuit {
					s.Fini()

					file, err := os.Create("/tmp/kitty")

					defer file.Close()

					if err != nil {
						fmt.Println("Couldn't create /tmp/kitty")
						os.Exit(0)
					} else {

						_, err := file.WriteString(cwd)

						if err != nil {
							fmt.Println("Couldn't write to /tmp/kitty")
						}
					}

					os.Exit(0)
				} else if input.Rune() == co.KeyDelete {
					err := os.RemoveAll(currFiles[currFile])

					if err == nil {
						if index, in := fu.In(currFiles[currFile], co.Selected); in {
							co.Selected = append(co.Selected[:index], co.Selected[index+1:]...)
						}

						currFiles = append(currFiles[:currFile], currFiles[currFile+1:]...)

						s.Clear()

						if currFile == len(currFiles) && b1 == 0 {
							currY -= 1
							currFile -= 1
						} else if currFile == len(currFiles) {
							b1 -= 1
							b2 -= 1
							currFile -= 1
						}

						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
								panic(err)
						}

						s.Show()
					}
				} else if input.Rune() == co.KeyCopy || input.Rune() == co.KeyMove || input.Rune() == co.KeyBulkDelete {
					for _, elem := range co.Selected {
						split := strings.Split(elem, "/")
						if input.Rune() == co.KeyCopy {
							err = cp.Copy(elem, cwd + "/" + split[len(split)-1])
						} else if input.Rune() == co.KeyMove {
							err = cp.Copy(elem, cwd + "/" + split[len(split)-1])
							err = os.RemoveAll(elem)
						} else if input.Rune() == co.KeyBulkDelete {
							err = os.RemoveAll(elem)
						}
					}

					if err == nil {
						co.Selected = []string{}
						currFiles, err = fu.GetFiles(cwd, dot)

						if err != nil {
							panic(err)
						}

						currY = co.YBuffTop
						currFile = 0

						b1, b2 = 0, height - co.YBuffBottom

						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
							panic(err)
						}
					}
				} else if input.Rune() == co.KeySelect {
					if fu.IsSel(cwd + "/" + currFiles[currFile]) {
						index, _ := fu.In(cwd + "/" + currFiles[currFile], co.Selected)
						co.Selected = append(co.Selected[:index], co.Selected[index+1:]...)
					} else {
						co.Selected = append(co.Selected, cwd + "/" + currFiles[currFile])
					}

					formated, err := fu.FormatText(s, currFiles[currFile], false)

					if err == nil {
						if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
							fu.Addstr(s, tcell.StyleDefault, co.XBuff, currY, formated + strings.Repeat(" ", len(co.SelectArrow)+1))
						} else {
							fu.Addstr(s, tcell.StyleDefault, co.XBuff, currY, formated + " ")
						}

						if err := fu.SelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
							panic(err)
						}
						s.Show()
					}
				} else if input.Rune() == co.KeySelectAll {
					for _, elem := range currFiles {
						if _, in := fu.In(cwd + "/" + elem, co.Selected); !in {
							co.Selected = append(co.Selected, cwd + "/" + elem)
						}
					}

					if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
						panic(err)
					}
				} else if input.Rune() == co.KeyDeselectAll {
					co.Selected = []string{}

					if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
						panic(err)
					}
				} else if input.Rune() == co.KeyDotToggle {
					dot = !(dot)
					currFiles, err = fu.GetFiles(cwd, dot)

					if err == nil {
						b1 = 0
						b2 = (height-co.YBuffBottom)+co.YBuffTop
						currFile = 0
						currY = co.YBuffTop
						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
							panic(err)
						}
					}
				} else if input.Rune() == co.KeyGoToFirst {
					currFile = 0
					currY = co.YBuffTop

					b1 = 0
					b2 = height - co.YBuffTop

					if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
						panic(err)
					}
				} else if input.Rune() == co.KeyGoToLast && currFile < len(currFiles)-1 {
					currFile = len(currFiles)-1

					if len(currFiles) + co.YBuffTop + co.YBuffBottom > height {
						currY = height - (co.YBuffBottom + 1)

						b1 = (len(currFiles))-(height-(co.YBuffTop+co.YBuffBottom))
						b2 = len(currFiles)-1
					} else {
						currY = currFile + co.YBuffTop
					}

					if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
						panic(err)
					}
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
						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
							panic(err)
						}
					} else {
						c := make(chan error)

						go func(c chan error) {
							c <- fu.DispBar(s, co.BarStyle, currFiles[currFile], currFile+2, len(currFiles))
						}(c)

						err := <-c

						if err != nil {
							panic(err)
						}

						if err := fu.DSelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
							panic(err)
						}

						currY += 1
						currFile += 1

						if err := fu.SelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
							panic(err)
						}

						s.Show()
					}

				} else if (input.Key() == tcell.KeyUp || input.Rune() == co.KeyUp) && len(currFiles) != 0 {
					if currFile == 0 {
						continue
					} else if currY == co.YBuffTop {
						b1 -= 1
						b2 -= 1
						currFile -= 1

						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
							panic(err)
						}
					} else {
						c := make(chan error)

						go func(c chan error) {
							c <- fu.DispBar(s, co.BarStyle, currFiles[currFile], currFile, len(currFiles))
						}(c)

						err := <-c

						if err != nil {
							panic(err)
						}

						if err := fu.DSelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
							panic(err)
						}

						currY -= 1
						currFile -= 1

						if err := fu.SelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
							panic(err)
						}

						s.Show()
					}
				} else if (input.Key() == tcell.KeyRight || input.Rune() == co.KeyRight) && len(currFiles) != 0 {
					if fu.Isd(currFiles[currFile]) {
						err := os.Chdir(currFiles[currFile])

						if err == nil {

							cwd, err = os.Getwd()

							if err != nil {
								panic(err)
							}

							currFiles, err = fu.GetFiles(cwd, dot)

							if err != nil {
								panic(err)
							}

							b1 = 0
							b2 = (height-co.YBuffBottom)+co.YBuffTop
							currFile = 0
							currY = co.YBuffTop

							s.Clear()

							if err := fu.DispFiles(s, currFiles); err != nil {
								panic(err)
							}

							fu.BorderPipes(s)
						}

						if len(currFiles) != 0 {
							if err := fu.SelFile(s, co.XBuff, currY, currFiles[currFile]); err != nil {
								panic(err)
							}

							if err := fu.DispBar(s, co.BarStyle, currFiles[currFile], currFile+1, len(currFiles)); err != nil {
								panic(err)
							}
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
							s, err = tcell.NewScreen()

							if err != nil {
								panic(err)
							}

							s.Init()
							if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
								panic(err)
							}

							fu.BorderPipes(s)
						} else if command[0] == "g" {
							cmd.Start()
						}
					}
				} else if input.Key() == tcell.KeyLeft || input.Rune() == co.KeyLeft {
					err := os.Chdir("..")

					if err == nil {
						cwd, err = os.Getwd()

						if err != nil {
							panic(err)
						}

						currFiles, err = fu.GetFiles(cwd, dot)

						if err != nil {
							panic(err)
						}

						b1 = 0
						b2 = (height-co.YBuffBottom)+co.YBuffTop
						currFile = 0
						currY = co.YBuffTop

						if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
							panic(err)
						}

					}
				} else if input.Rune() == co.KeyRefresh && len(currFiles) != 0 {
					currFile = 0

					b1 = 0
					b2 = height-co.YBuffBottom

					currFiles, err = fu.GetFiles(cwd, dot)

					if err != nil {
						panic(err)
					}

					currY = co.YBuffTop

					if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
						panic(err)
					}
				} else {
					for k, v := range co.Bindings {
						if k == input.Rune() {
							replacedString := strings.Replace(v[1], "@", currFiles[currFile], -1)

							if v[0] == "cd" {
								err := os.Chdir(strings.Replace(string(v[1]), "~", os.Getenv("HOME"), -1))

								if err == nil {
									cwd, err = os.Getwd()

									if err != nil {
										panic(err)
									}

									currFiles, err = fu.GetFiles(cwd, dot)

									if err != nil {
										panic(err)
									}

									b1 = 0
									b2 = (height-co.YBuffBottom)+co.YBuffTop
									currFile = 0
									currY = co.YBuffTop

									if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
										panic(err)
									}
								}
							} else if v[0] == "t" {
								cmd := exec.Command(co.Shell, "-c", replacedString)

								s.Fini()

								cmd.Stdout = os.Stdout
								cmd.Stdin = os.Stdin
								cmd.Run()

								fmt.Print()

								s, err = tcell.NewScreen()

								if err != nil {
									panic(err)
								}

								s.Init()

								if err := fu.DrawScreen(s, currFiles, currFile, currY, b1, b2); err != nil {
									panic(err)
								}
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
