package funcs

import (
	"github.com/gdamore/tcell"
	"strings"
	"sort"
	"os"
	"strconv"
	"os/exec"
	"os/user"
	"code.cloudfoundry.org/bytefmt"
	co "catfm/config"
	ke "catfm/keys"
)

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		s.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func ShowFile(s tcell.Screen, x int, y int, f string, sel bool) (string, error) {
	formated, err := FormatText(s, f, sel)

	if err != nil {
		return "", err
	}

	split := strings.Split(f, ".")

	if Isd(f) {
		Addstr(s, co.FileColors["[dir]"], x, y, formated)
	} else {
		Addstr(s, co.FileColors[split[len(split)-1]], x, y, formated)
	}

	return formated, nil
}

func (v View) DispFiles(s tcell.Screen) error {
	_, height := s.Size()
	files := v.Files

	if v.Buffer1 > 0 {
		files = v.Files[v.Buffer1:v.Buffer2]
	}

	if len(files) != 0 {
		for i, f := range files {
			if i+co.YBuffTop < (height - co.YBuffBottom) {
				if _, err := ShowFile(s, co.XBuff, i+co.YBuffTop, f, false); err != nil {
					return err
				}
			} else {
				break
			}
		}
	} else {
		msg, err := FormatText(s, "The cat can't seem to find anything here...", false)

		if err != nil {
			return err
		}

		Addstr(s, tcell.StyleDefault, co.XBuff, co.YBuffTop, msg)
	}

	return nil
}

func (v View) SelFile(s tcell.Screen) error {
	formated, err := FormatText(s, v.Files[v.File], true)

	if err != nil {
		return err
	}

	width, _ := s.Size()

	if co.SelectType == "full" {
		Addstr(s, co.SelectStyle, co.XBuff, v.Y, formated + strings.Repeat(" ", width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		Addstr(s, co.SelectArrowStyle, co.XBuff, v.Y, co.SelectArrow)

		if co.SelectType == "arrow-default" {
			Addstr(s, co.SelectStyle, co.XBuff+len(co.SelectArrow), v.Y, formated)
		} else {
			if _, err = ShowFile(s, co.XBuff+len(co.SelectArrow), v.Y, v.Files[v.File], false); err != nil {
				return err
			}
		}
	} else if co.SelectType == "default" {
		Addstr(s, co.SelectStyle, co.XBuff, v.Y, formated)
	}

	return nil
}

func (v View) DSelFile(s tcell.Screen) error {
	width, _ := s.Size()

	formated, err := ShowFile(s, co.XBuff, v.Y, v.Files[v.File], false)

	if err != nil {
		return err
	}

	if co.SelectType == "full" {
		Addstr(s, tcell.StyleDefault, co.XBuff+len(formated), v.Y, strings.Repeat(" ", width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		Addstr(s, tcell.StyleDefault, co.XBuff+len(formated), v.Y, strings.Repeat(" ", len(co.SelectArrow)+1))
	}

	return nil
}

func (v View) DispBar(s tcell.Screen, file string, curr int) error {
	var x int = co.XBuff
	var elemOutput string
	var loc int

	width, y := s.Size()

	if co.BarLocale == "bottom" {
		loc = y-(co.YBuffBottom)+1
	} else if co.BarLocale == "top" {
		loc = co.YBuffTop-2
	} else if co.BarLocale == "" {
		loc = y+5
	}

	Addstr(s, tcell.StyleDefault, co.XBuff, loc, strings.Repeat(" ", barLen))

	keys := make([]string, len(co.BarStyle))

	i := 0
	for k, _ := range co.BarStyle {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var err error

	for _, k := range keys {
		if k[1:] == "cwd" {
			if co.TildeHome {
				u, err := user.Current()

				if err != nil {
					return err
				}

				ucwd, err := os.Getwd()

				if err != nil {
					return err
				}

				elemOutput = strings.Replace(ucwd, u.HomeDir, "~", -1)
			} else {
				elemOutput, err = os.Getwd()

				if err != nil {
					return err
				}
			}
		} else if k[1:] == "size" {
			f, err := os.Stat(file)

			if err != nil {
				return err
			}

			elemOutput = bytefmt.ByteSize(uint64(f.Size()))
		} else if k[1:] == "mode" {
			f, err := os.Stat(file)

			if err != nil {
				return err
			}

			elemOutput = f.Mode().String()
		} else if k[1:] == "total" {
			elemOutput = strconv.Itoa(curr) + "/" + strconv.Itoa(len(v.Files))
		} else if k[1] == '[' && k[len(k)-1] == ']' {
			replacedString := strings.Replace(k, "@", file, -1)
			cmdOutput, _ := exec.Command(co.Shell, "-c", replacedString[2:len(replacedString)-1]).Output()
			elemOutput = string(cmdOutput)
		} else if k[1:] == "tab" {
			elemOutput = "[" + strconv.Itoa(ViewNumber+1) + "]"
		} else {
			elemOutput = k[1:]
		}

		if strings.Contains(k[1:], "$HOST") {
			host, err := os.Hostname()

			if err != nil {
				return err
			}

			elemOutput = strings.Replace(elemOutput, "$HOST", host, -1)
		}

		if strings.Contains(k[1:], "$USER") {
			elemOutput = strings.Replace(elemOutput, "$USER", os.Getenv("USER"), -1)
		}

		if strings.Contains(k[1:], "$FILE") {
			elemOutput = strings.Replace(elemOutput, "$FILE", file, -1)
		}

		if strings.Contains(k[1:], "$TAB") {
			elemOutput = strings.Replace(elemOutput, "$TAB", strconv.Itoa(ViewNumber+1), -1)
		}

		Addstr(s, co.BarStyle[k], x, loc, elemOutput)
		if num, _ := strconv.Atoi(string(k[0])); num != len(keys) {
			Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), x+len(elemOutput), loc, co.BarDiv)
			x += len(elemOutput + co.BarDiv)
		} else {
			x += len(elemOutput)
		}
	}
	if x > width {
		Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), width-3, loc, "...")
	}
	barLen = x

	return nil
}

func (v View) DrawScreen(s tcell.Screen) error {
	s.Clear()

	if err := v.DispFiles(s); err != nil {
		return err
	}

	go BorderPipes(s)

	if len(v.Files) > 0 {
		if err := v.DispBar(s, v.Files[v.File], v.File+1); err != nil {
			return err
		}

		if err := v.SelFile(s); err != nil {
			return err
		}
	}

	s.Show()

	return nil
}

func BorderPipes(s tcell.Screen) {
	if co.PipeType != "" {
		thin := []rune{'┌', '┐', '└', '┘', '─', '│'}
		thick := []rune{'┏', '┓', '┗', '┛', '━', '┃'}
		hollow := []rune{'╔', '╗', '╚', '╝', '═', '║'}
		round := []rune{'╭', '╮', '╰', '╯', '─', '│'}
		used := []rune{}

		width, height := s.Size()

		if co.PipeType == "thick" {
			used = thick
		} else if co.PipeType == "thin" {
			used = thin
		} else if co.PipeType == "hollow" {
			used = hollow
		} else if co.PipeType == "round" {
			used = round
		}

		s.SetContent(0, 0, used[0], []rune(""), co.PipeStyle)
		s.SetContent(width-1, 0, used[1], []rune(""), co.PipeStyle)

		s.SetContent(0, height-1, used[2], []rune(""), co.PipeStyle)
		s.SetContent(width-1, height-1, used[3], []rune(""), co.PipeStyle)

		for i := 1; i < width-1; i++ {
			s.SetContent(i, 0, used[4], []rune(""), co.PipeStyle)
			s.SetContent(i, height-1, used[4], []rune(""), co.PipeStyle)
		}

		for i := 1; i < height-1; i++ {
			s.SetContent(0, i, used[5], []rune(""), co.PipeStyle)
			s.SetContent(width-1, i, used[5], []rune(""), co.PipeStyle)
		}

		s.Show()

		if co.PipeText != "" {
			user := os.Getenv("USER")
			host, err := os.Hostname()

			var text string

			if co.PipeText == "user@host" {
				if err == nil {
					text = user + "@" + host
				}
			} else if co.PipeText == "catfm@host" {
				if err == nil {
					text = "catfm@" + host
				}
			} else if co.PipeText == "user@catfm" {
				text = user + "@catfm"
			} else {
				text = co.PipeText
			}

			Addstr(s, co.PipeTextStyle, 1, 0, text)
			s.Show()
		}
	}
}

func (v *View) GoToLast(s tcell.Screen) {
	_, height := s.Size()
	v.File = len(v.Files) - 1

	if len(v.Files) + co.YBuffTop + co.YBuffBottom > height {
		v.Y = height - (co.YBuffBottom + 1)

		v.Buffer1 = (len(v.Files))-(height-(co.YBuffTop+co.YBuffBottom))
		v.Buffer2 = len(v.Files)
	} else {
		v.Y = v.File + co.YBuffTop
	}

	if err := v.DrawScreen(s); err != nil {
		Errout(s, "couldn't draw screen")
	}
}

func (v *View) GoToFirst(s tcell.Screen) {
	_, height := s.Size()

	v.File = 0
	v.Y = co.YBuffTop

	v.Buffer1 = 0
	v.Buffer2 = height-co.YBuffBottom

	if err := v.DrawScreen(s); err != nil {
		Errout(s, "couldn't draww screen")
	}
}

func (v *View) Move(s tcell.Screen, n int) {
	_, h := s.Size()

	if v.File == 0 && n == -1 {
		v.GoToLast(s)
	} else if v.File == len(v.Files)-1 && n == 1 {
		v.GoToFirst(s)
	} else if (v.Y == (h-1)-co.YBuffBottom && n == 1) || (v.Y == co.YBuffTop && n == -1) {
		v.Buffer1 += n
		v.Buffer2 += n
		v.File += n

		if err := v.DrawScreen(s); err != nil {
			Errout(s, "couldn't draw screen")
		}
	} else {
		c := make(chan error)

		go func(c chan error) {
			cf := 0

			if n == 1 {
				cf = 2
			}

			c <- v.DispBar(s, v.Files[v.File+n], v.File+cf)
		}(c)

		err := <-c

		if err != nil {
			Errout(s, "unable to display the infobar")
		}

		if err := v.DSelFile(s); err != nil {
			Errout(s, "unable to deselect file")
		}

		v.Y += n
		v.File += n

		if err := v.SelFile(s); err != nil {
			Errout(s, "unable to select file")
		}

		s.Show()

	}
}

func (v *View) Resize(s tcell.Screen) {
	nw, nh := s.Size()

	if v.Width != nw && v.Height == nh {
		v.Width, v.Height = nw, nh

		if err := v.DrawScreen(s); err != nil {
			Errout(s, "unable to draw screen")
		}

	} else if v.Height != nh {
		v.Width, v.Height = nw, nh

		if v.Buffer1 != 0 || v.File > v.Height-(co.YBuffTop+co.YBuffBottom) {
			v.Buffer1 = 0
			v.Buffer2 = (v.Height-co.YBuffBottom)+co.YBuffTop
			v.Y = co.YBuffTop
			v.File = 0
		}

		if err := v.DrawScreen(s); err != nil {
			Errout(s, "unable to draw screen")
		}

	}
}

func (v *View) Search(s tcell.Screen) {
	var (
		original []string = v.Files
		searching bool = true
		search string
	)

	for searching {
		v.Resize(s)

		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				if ke.MatchKey(input, co.KeyToggleSearch) {
					searching = false
				} else if input.Key() == tcell.KeyDEL || input.Key() == tcell.KeyBS {
					if len(search) > 0 {
						if len(search) == 1 {
							search = ""
						} else {
							search = search[:len(search)-1]
						}
					}
				} else {
					search += string(input.Rune())
				}

				v.Files = []string{}

				for _, e := range original {
					if strings.Contains(e, search) {
						v.Files = append(v.Files, e)
					}
				}

				if v.File != 0 {
					v.File = 0
					v.Y = co.YBuffTop

					v.Buffer1 = 0
					v.Buffer2 = v.Height-co.YBuffBottom
				}

				if err := v.DrawScreen(s); err != nil {
					Errout(s, "unable to draw screen")
				}

				BorderPipes(s)


		}
	}
}
