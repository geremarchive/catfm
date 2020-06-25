/*

Welcome To The Catfm Source Code

------------------------------------

Info:

∙ This project was created and is maintained by Jonah G. Rongstad (geremachek).
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
	"strings"
	"fmt"
	"io/ioutil"
	fu "catfm/funcs"
	co "catfm/config"
	ke "catfm/keys"
	he "catfm/help"
	cp "github.com/otiai10/copy"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			he.Help()

			return
		} else {
			err := os.Chdir(os.Args[1])

			if err != nil {
				fmt.Println("Couldn't open '" + os.Args[1] + "'")
			}
		}
	}

	co.Init()

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

	if err := currView.DrawScreen(s); err != nil {
		fu.Errout(s, "unable to draw screen")
	}

	s.Show()

	for {
		currView.Resize(s)

		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				if ke.MatchKey(input, co.KeyQuit) {
					s.Fini()
					file, err := os.Create("/tmp/kitty")

					defer file.Close()

					if err != nil {
						fmt.Println("catfm: couldn't create /tmp/kitty")
						os.Exit(0)
					} else {

						_, err := file.WriteString(currView.Cwd)

						if err != nil {
							fmt.Println("catfm: couldn't write to /tmp/kitty")
						}
					}

					os.Exit(0)
				} else if ke.MatchKey(input, co.KeyDelete) || ke.MatchKey(input, co.KeyRecycle) {
					if ke.MatchKey(input, co.KeyRecycle) && co.RecycleBin != "" {
						if err := cp.Copy(currView.Files[currView.File], co.RecycleBin + "/" + currView.Files[currView.File]); err != nil {
							fu.Errout(s, "unable to copy file")
						}
					}

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

						if err := currView.DrawScreen(s); err != nil {
							fu.Errout(s, "couldn't draw screen")
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
							fu.Errout(s, "unable to read files")
						}

						if ke.MatchKey(input, co.KeyBulkDelete) {
							currView.Y = co.YBuffTop
							currView.File = 0

							currView.Buffer1, currView.Buffer2 = 0, height - co.YBuffBottom
						}

						if err := currView.DrawScreen(s); err != nil {
							fu.Errout(s, "couldn't draw screen")
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

						if err := currView.SelFile(s); err != nil {
							fu.Errout(s, "unable to calculate screen width")
						}
						s.Show()
					}
				} else if ke.MatchKey(input, co.KeySelectAll) {
					for _, elem := range currView.Files {
						if _, in := fu.In(currView.Cwd + "/" + elem, fu.Selected); !in {
							fu.Selected = append(fu.Selected, currView.Cwd + "/" + elem)
						}
					}

					if err := currView.DrawScreen(s); err != nil {
						fu.Errout(s, "couldn't draw screen")
					}
				} else if ke.MatchKey(input, co.KeyDeselectAll) {
					fu.Selected = []string{}

					if err := currView.DrawScreen(s); err != nil {
						fu.Errout(s, "couldn't draw screen")
					}
				} else if ke.MatchKey(input, co.KeyDotToggle) {
					currView.Dot = !(currView.Dot)
					currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

					if err == nil {
						currView.Buffer1 = 0
						currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
						currView.File = 0
						currView.Y = co.YBuffTop

						if err := currView.DrawScreen(s); err != nil {
							fu.Errout(s, "couldn't draw screen")
						}
					}
				} else if ke.MatchKey(input, co.KeyGoToFirst) {
					currView.GoToFirst(s)
				} else if ke.MatchKey(input, co.KeyGoToLast) && currView.File < len(currView.Files)-1 {
					currView.GoToLast(s)
				} else if ke.MatchKey(input, co.KeyRename) {
					err := ioutil.WriteFile("/tmp/rename.catfm", []byte(currView.Files[currView.File]), 0644)

					if err == nil {
						replaced := strings.Replace(co.FileOpen["*"][1], "@", "/tmp/rename.catfm", -1)

						s = currView.ParseBinding(s, []string{co.FileOpen["*"][0],  replaced})

						if _, err := os.Stat("/tmp/rename.catfm"); err == nil {
							file, err := os.Open("/tmp/rename.catfm")

							defer file.Close()

							if err == nil {
								b, err := ioutil.ReadAll(file)

								if err == nil {
									newName := strings.Split(string(b), "\n")[0]
									_, in := fu.In(newName, currView.Files)

									if newName != "" && !(in) {
										os.Rename(currView.Files[currView.File], newName)
										currView.Files[currView.File] = newName

										if err := currView.DrawScreen(s); err != nil {
											fu.Errout(s, "couldn't draw screen")
										}
									}
								}
							}
						}
					}
				} else if (input.Key() == tcell.KeyDown || ke.MatchKey(input, co.KeyDown)) && len(currView.Files) != 0 {
					currView.Move(s, 1)
				} else if (input.Key() == tcell.KeyUp || ke.MatchKey(input, co.KeyUp)) && len(currView.Files) != 0 {
					currView.Move(s, -1)
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

							if err := currView.DispFiles(s); err != nil {
								fu.Errout(s, "unable to display files")
							}

							fu.BorderPipes(s)
						}

						if len(currView.Files) != 0 {
							if err := currView.SelFile(s); err != nil {
								fu.Errout(s, "unable to select file")
							}

							if err := currView.DispBar(s, currView.Files[currView.File], currView.File+1); err != nil {
								fu.Errout(s, "couldn't display the infobar")
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

						s = currView.ParseBinding(s, command)
					}
				} else if input.Key() == tcell.KeyLeft || ke.MatchKey(input, co.KeyLeft) {
					err := os.Chdir("..")

					if err == nil {
						currView.Cwd, err = os.Getwd()

						if err != nil {
							fu.Errout(s, "unable to get the working directory")
						}

						currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

						if err != nil {
							fu.Errout(s, "unable to read files")
						}

						currView.Buffer1 = 0
						currView.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
						currView.File = 0
						currView.Y = co.YBuffTop

						if err := currView.DrawScreen(s); err != nil {
							fu.Errout(s, "couldn't draw screen")
						}

					}
				} else if ke.MatchKey(input, co.KeyRefresh) {
					currView.File = 0

					currView.Buffer1 = 0
					currView.Buffer2 = height-co.YBuffBottom

					currView.Files, err = fu.GetFiles(currView.Cwd, currView.Dot)

					if err != nil {
						fu.Errout(s, "unable to read files")
					}

					currView.Y = co.YBuffTop

					if err := currView.DrawScreen(s); err != nil {
						fu.Errout(s, "couldn't draw screen")
					}
				} else if ke.MatchKey(input, co.KeyToggleSearch) {
					currView.Search(s)
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
						fu.Errout(s, "couldn't change directory")
					}

					currView.Resize(s)
				} else {
					for k, v := range co.Bindings {
						if ke.MatchKey(input, k) {
							s = currView.ParseBinding(s, v)
							break
						}
					}
				}
		}
	}
}
