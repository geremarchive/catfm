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
	"os/user"
	"strings"
	"fmt"
	fu "catfm/funcs"
	co "catfm/config"
	ke "catfm/keys"
	cp "github.com/otiai10/copy"
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

	files, err := fu.GetFiles(cwd, false)

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

	width, height := s.Size()

	Views := []fu.View{}

	currView := fu.View {
		File: 0,
		Files: files,

		Buffer1: 0,
		Buffer2: height-co.YBuffBottom,

		Y: co.YBuffTop,
		Dot: false,
		Cwd: cwd,

		Width: width,
		Height: height,
	}

	for i := 0; i < 10; i++ {
		Views = append(Views, currView)
	}

	if err := fu.DrawScreen(s, currView); err != nil {
		panic(err)
	}

	s.Show()

	for {
		nw, nh := s.Size()
		if width != nw && height == nh {
			width, height = nw, nh
			currView.Width, currView.Height = nw, nh

			if err := fu.DrawScreen(s, currView); err != nil {
				panic(err)
			}
			fu.BorderPipes(s)
		} else if height != nh {
			width, height = nw, nh
			currView.Width, currView.Height = nw, nh

			if currView.Buffer1 != 0 || currView.File > currView.Height-(co.YBuffTop+co.YBuffBottom) {
				currView.Buffer1 = 0
				currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
				currView.Y = co.YBuffTop
				currView.File = 0
			}

			if err := fu.DrawScreen(s, currView); err != nil {
				panic(err)
			}

			fu.BorderPipes(s)
		}
		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				if ke.MatchKey(input, co.KeyQuit) {
					s.Fini()
					file, err := os.Create("/tmp/kitty")

					defer file.Close()

					if err != nil {
						fmt.Println("Couldn't create /tmp/kitty")
						os.Exit(0)
					} else {

						_, err := file.WriteString(currView.Cwd)

						if err != nil {
							fmt.Println("Couldn't write to /tmp/kitty")
						}
					}

					os.Exit(0)
				} else if ke.MatchKey(input, co.KeyDelete) {
					err := os.RemoveAll(currView.Files[currView.File])

					if err == nil {
						if index, in := fu.In(currView.Files[currView.File], fu.Selected); in {
							fu.Selected = append(fu.Selected[:index], fu.Selected[index+1:]...)
						}

						currView.Files = append(currView.Files[:currView.File], currView.Files[currView.File+1:]...)

						s.Clear()

						if currView.File == len(currView.Files) && currView.Buffer1 == 0 {
							currView.Y -= 1
							currView.File -= 1
						} else if currView.File == len(currView.Files) {
							currView.Buffer1 -= 1
							currView.Buffer2 -= 1
							currView.File -= 1
						}

						if err := fu.DrawScreen(s, currView); err != nil {
								panic(err)
						}

						s.Show()
					}
				} else if ke.MatchKey(input, co.KeyCopy) || ke.MatchKey(input, co.KeyMove) || ke.MatchKey(input, co.KeyBulkDelete) {
					for _, elem := range fu.Selected {
						split := strings.Split(elem, "/")
						if ke.MatchKey(input, co.KeyCopy) {
							err = cp.Copy(elem, currView.Cwd + "/" + split[len(split)-1])
						} else if ke.MatchKey(input, co.KeyMove) {
							err = cp.Copy(elem, currView.Cwd + "/" + split[len(split)-1])
							err = os.RemoveAll(elem)
						} else if ke.MatchKey(input, co.KeyBulkDelete) {
							err = os.RemoveAll(elem)
						}
					}

					if err == nil {
						fu.Selected = []string{}
						currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

						if err != nil {
							panic(err)
						}

						currView.Y = co.YBuffTop
						currView.File = 0

						currView.Buffer1, currView.Buffer2 = 0, height - co.YBuffBottom

						if err := fu.DrawScreen(s, currView); err != nil {
							panic(err)
						}
					}
				} else if ke.MatchKey(input, co.KeySelect) {
					if fu.IsSel(currView.Cwd + "/" + currView.Files[currView.File]) {
						index, _ := fu.In(currView.Cwd + "/" + currView.Files[currView.File], fu.Selected)
						fu.Selected = append(fu.Selected[:index], fu.Selected[index+1:]...)
					} else {
						fu.Selected = append(fu.Selected, currView.Cwd + "/" + currView.Files[currView.File])
					}

					formated, err := fu.FormatText(s, currView.Files[currView.File], false)

					if err == nil {
						if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
							fu.Addstr(s, tcell.StyleDefault, co.XBuff, currView.Y, formated + strings.Repeat(" ", len(co.SelectArrow)+1))
						} else {
							fu.Addstr(s, tcell.StyleDefault, co.XBuff, currView.Y, formated + " ")
						}

						if err := fu.SelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
							panic(err)
						}
						s.Show()
					}
				} else if ke.MatchKey(input, co.KeySelectAll) {
					for _, elem := range currView.Files {
						if _, in := fu.In(currView.Cwd + "/" + elem, fu.Selected); !in {
							fu.Selected = append(fu.Selected, currView.Cwd + "/" + elem)
						}
					}

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else if ke.MatchKey(input, co.KeyDeselectAll) {
					fu.Selected = []string{}

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else if ke.MatchKey(input, co.KeyDotToggle) {
					currView.Dot = !(currView.Dot)
					currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

					if err == nil {
						currView.Buffer1 = 0
						currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
						currView.File = 0
						currView.Y = co.YBuffTop
						if err := fu.DrawScreen(s, currView); err != nil {
							panic(err)
						}
					}
				} else if ke.MatchKey(input, co.KeyGoToFirst) {
					currView.File = 0
					currView.Y = co.YBuffTop

					currView.Buffer1 = 0
					currView.Buffer2 = height - co.YBuffTop

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else if ke.MatchKey(input, co.KeyGoToLast) && currView.File < len(currView.Files)-1 {
					currView.File = len(currView.Files)-1

					if len(currView.Files) + co.YBuffTop + co.YBuffBottom > height {
						currView.Y = height - (co.YBuffBottom + 1)

						currView.Buffer1 = (len(currView.Files))-(height-(co.YBuffTop+co.YBuffBottom))
						currView.Buffer2 = len(currView.Files)-1
					} else {
						currView.Y = currView.File + co.YBuffTop
					}

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else if (input.Key() == tcell.KeyDown || ke.MatchKey(input, co.KeyDown)) && len(currView.Files) != 0 {
					if currView.File == len(currView.Files)-1 {
						continue
					} else if currView.Y == (height-1)-co.YBuffBottom {
						currView.Buffer1 += 1
						if currView.Buffer2 >= len(currView.Files)-1 {
							currView.Buffer2 = len(currView.Files)-1
						} else {
							currView.Buffer2 += 1
						}
						currView.File += 1
						if err := fu.DrawScreen(s, currView); err != nil {
							panic(err)
						}
					} else {
						c := make(chan error)

						go func(c chan error) {
							c <- fu.DispBar(s, co.BarStyle, currView.Files[currView.File+1], currView.File+2, len(currView.Files))
						}(c)

						err := <-c

						if err != nil {
							panic(err)
						}

						if err := fu.DSelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
							panic(err)
						}

						currView.Y += 1
						currView.File += 1

						if err := fu.SelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
							panic(err)
						}

						s.Show()
					}

				} else if (input.Key() == tcell.KeyUp || ke.MatchKey(input, co.KeyUp)) && len(currView.Files) != 0 {
					if currView.File == 0 {
						continue
					} else if currView.Y == co.YBuffTop {
						currView.Buffer1 -= 1
						currView.Buffer2 -= 1
						currView.File -= 1

						if err := fu.DrawScreen(s, currView); err != nil {
							panic(err)
						}
					} else {
						c := make(chan error)

						go func(c chan error) {
							c <- fu.DispBar(s, co.BarStyle, currView.Files[currView.File-1], currView.File, len(currView.Files))
						}(c)

						err := <-c

						if err != nil {
							panic(err)
						}

						if err := fu.DSelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
							panic(err)
						}

						currView.Y -= 1
						currView.File -= 1

						if err := fu.SelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
							panic(err)
						}

						s.Show()
					}
				} else if (input.Key() == tcell.KeyRight || ke.MatchKey(input, co.KeyRight)) && len(currView.Files) != 0 {
					if fu.Isd(currView.Files[currView.File]) {
						err := os.Chdir(currView.Files[currView.File])

						if err == nil {

							currView.Cwd, err = os.Getwd()

							if err != nil {
								panic(err)
							}

							currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

							if err != nil {
								panic(err)
							}

							currView.Buffer1 = 0
							currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
							currView.File = 0
							currView.Y = co.YBuffTop

							s.Clear()

							if err := fu.DispFiles(s, currView.Files); err != nil {
								panic(err)
							}

							fu.BorderPipes(s)
						}

						if len(currView.Files) != 0 {
							if err := fu.SelFile(s, co.XBuff, currView.Y, currView.Files[currView.File]); err != nil {
								panic(err)
							}

							if err := fu.DispBar(s, co.BarStyle, currView.Files[currView.File], currView.File+1, len(currView.Files)); err != nil {
								panic(err)
							}
						}
						s.Show()
					} else {
						var command []string = co.FileOpen["*"]
						splitFile := strings.Split(currView.Files[currView.File], ".")
						for k, v := range co.FileOpen {
							if k == splitFile[len(splitFile)-1] {
								command = v
								break
							}
						}

						replacedString := strings.Replace(command[1], "@", currView.Files[currView.File], -1)
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
							if err := fu.DrawScreen(s, currView); err != nil {
								panic(err)
							}

							fu.BorderPipes(s)
						} else if command[0] == "g" {
							cmd.Start()
						}
					}
				} else if input.Key() == tcell.KeyLeft || ke.MatchKey(input, co.KeyLeft) {
					err := os.Chdir("..")

					if err == nil {
						currView.Cwd, err = os.Getwd()

						if err != nil {
							panic(err)
						}

						currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

						if err != nil {
							panic(err)
						}

						currView.Buffer1 = 0
						currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
						currView.File = 0
						currView.Y = co.YBuffTop

						if err := fu.DrawScreen(s, currView); err != nil {
							panic(err)
						}

					}
				} else if ke.MatchKey(input, co.KeyRefresh) && len(currView.Files) != 0 {
					currView.File = 0

					currView.Buffer1 = 0
					currView.Buffer2 = height-co.YBuffBottom

					currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

					if err != nil {
						panic(err)
					}

					currView.Y = co.YBuffTop

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else if input.Rune() >= 48 && input.Rune() <= 57 {
					Views[fu.ViewNumber] = currView

					switch input.Rune() {
					case '1':
						fu.ViewNumber = 0
					case '2':
						fu.ViewNumber = 1
					case '3':
						fu.ViewNumber = 2
					case '4':
						fu.ViewNumber = 3
					case '5':
						fu.ViewNumber = 4
					case '6':
						fu.ViewNumber = 5
					case '7':
						fu.ViewNumber = 6
					case '8':
						fu.ViewNumber = 7
					case '9':
						fu.ViewNumber = 8
					case '0':
						fu.ViewNumber = 9
					}

					currView = Views[fu.ViewNumber]

					if err := os.Chdir(currView.Cwd); err != nil {
						panic(err)
					}

					if currView.Width != width && currView.Height == height {
						currView.Width = width
					} else if currView.Height != height {
						currView.Width, currView.Height = width, height

						if currView.Buffer1 != 0 || currView.File > currView.Height-(co.YBuffTop+co.YBuffBottom) {
							currView.Buffer1 = 0
							currView.Buffer2 = height-co.YBuffBottom
							currView.File = 0
							currView.Y = co.YBuffTop
						}
					}

					if err := fu.DrawScreen(s, currView); err != nil {
						panic(err)
					}
				} else {
					for k, v := range co.Bindings {
						if ke.MatchKey(input, k) {
							replacedString := v[1]
							if len(currView.Files) != 0 {
								replacedString = strings.Replace(v[1], "@", currView.Files[currView.File], -1)
							}

							if v[0] == "cd" {
								u, err := user.Current()

								if err != nil {
									panic(err)
								}

								err = os.Chdir(strings.Replace(v[1], "~", u.HomeDir, -1))

								if err == nil {
									currView.Cwd, err = os.Getwd()

									if err != nil {
										panic(err)
									}

									currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

									if err != nil {
										panic(err)
									}

									currView.Buffer1 = 0
									currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
									currView.File = 0
									currView.Y = co.YBuffTop

									if err := fu.DrawScreen(s, currView); err != nil {
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

								if err := fu.DrawScreen(s, currView); err != nil {
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
