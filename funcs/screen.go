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
)

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		s.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func DispFiles(s tcell.Screen, files []string) error {
	_, height := s.Size()

	if len(files) != 0 {
		for i, f := range files {
			formated, err := FormatText(s, f, false)

			if err != nil {
				return err
			}

			if i+co.YBuffTop < (height - co.YBuffBottom) {
				splitFile := strings.Split(f, ".")
				if Isd(f) {
					Addstr(s, co.FileColors["[dir]"], co.XBuff, i+co.YBuffTop, formated)
				} else {
					Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], co.XBuff, i+co.YBuffTop, formated)
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

		Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor("#ff0000")).Bold(true), co.XBuff, co.YBuffTop, msg)
	}

	return nil
}

func SelFile(s tcell.Screen, x int, y int, file string) error {
	formated, err := FormatText(s, file, true)

	if err != nil {
		return err
	}

	splitFile := strings.Split(file, ".")
	width, _ := s.Size()

	if co.SelectType == "full" {
		Addstr(s, co.SelectStyle, x, y, formated + strings.Repeat(" ", width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		Addstr(s, co.SelectArrowStyle, x, y, co.SelectArrow)

		if co.SelectType == "arrow-default" {
			Addstr(s, co.SelectStyle, x+len(co.SelectArrow), y, formated)
		} else {
			if Isd(file) {
				Addstr(s, co.FileColors["[dir]"], x+len(co.SelectArrow), y, formated)
			} else {
				Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x+len(co.SelectArrow), y, formated)
			}
		}
	} else if co.SelectType == "default" {
		Addstr(s, co.SelectStyle, x, y, formated)
	}

	return nil
}

func DSelFile(s tcell.Screen, x int, y int, file string) error {
	formated, err := FormatText(s, file, false)

	if err != nil {
		return err
	}

	splitFile := strings.Split(file, ".")
	width, _ := s.Size()

	if co.SelectType == "full" {
		if Isd(file) {
			Addstr(s, co.FileColors["[dir]"], x, y, formated)
		} else {
			Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x, y, formated)
		}

		Addstr(s, tcell.StyleDefault, x+len(formated), y, strings.Repeat(" ", width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		if Isd(file) {
			Addstr(s, co.FileColors["[dir]"], x, y, formated + strings.Repeat(" ", len(co.SelectArrow)))
		} else {
			Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x, y, formated + strings.Repeat(" ", len(co.SelectArrow)+1))
		}
	} else if co.SelectType == "default" {
		if Isd(file) {
			Addstr(s, co.FileColors["[dir]"], x, y, formated)
		} else {
			Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x, y, formated)
		}
	}

	return nil
}

func DispBar(s tcell.Screen, elements map[string]tcell.Style, file string, curr int, total int) error {
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

	keys := make([]string, len(elements))

	i := 0
	for k, _ := range elements {
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
			elemOutput = strconv.Itoa(curr) + "/" + strconv.Itoa(total)
		} else if k[1] == '[' && k[len(k)-1] == ']' {
			replacedString := strings.Replace(k, "@", file, -1)
			cmdOutput, _ := exec.Command("dash", "-c", replacedString[2:len(replacedString)-1]).Output()
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

		Addstr(s, elements[k], x, loc, elemOutput)
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

func DrawScreen(s tcell.Screen, v View) error {
	s.Clear()

	c := make(chan error)
	var err error

	go func(c chan error) {
		if v.Buffer1 == 0 {
			err = DispFiles(s, v.Files)
		} else {
			err = DispFiles(s, v.Files[v.Buffer1:v.Buffer2+1])
		}

		c <- err
	}(c)

	err = <-c

	if err != nil {
		return err
	}

	go BorderPipes(s)

	if len(v.Files) > 0 {
		err := DispBar(s, co.BarStyle, v.Files[v.File], v.File+1, len(v.Files))

		if err != nil {
			return err
		}

		err = SelFile(s, co.XBuff, v.Y, v.Files[v.File])

		if err != nil {
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
		v.Buffer2 = len(v.Files) - 1
	} else {
		v.Y = v.File + co.YBuffTop
	}

	if err := DrawScreen(s, *v); err != nil {
		panic(err)
	}
}

func (v *View) GoToFirst(s tcell.Screen) {
	_, height := s.Size()

	v.File = 0
	v.Y = co.YBuffTop

	v.Buffer1 = 0
	v.Buffer2 = height - co.YBuffTop

	if err := DrawScreen(s, *v); err != nil {
		panic(err)
	}
}
