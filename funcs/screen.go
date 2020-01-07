package funcs

import (
	"github.com/gdamore/tcell"
	"strings"
	"sort"
	"os"
	"strconv"
	"os/exec"
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
	width, y := s.Size()
	Addstr(s, tcell.StyleDefault, co.XBuff, y-(co.YBuffBottom)+1, strings.Repeat(" ", barLen))

	keys := make([]string, len(elements))

	i := 0
	for k, _ := range elements {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		if k[1:] == "cwd" {
			elemOutput, _ = os.Getwd()
		} else if k[1:] == "size" {
			f, _ := os.Stat(file)
			elemOutput = bytefmt.ByteSize(uint64(f.Size()))
		} else if k[1:] == "mode" {
			f, _ := os.Stat(file)
			elemOutput = f.Mode().String()
		} else if k[1:] == "total" {
			elemOutput = strconv.Itoa(curr) + "/" + strconv.Itoa(total)
		} else if k[1] == '[' && k[len(k)-1] == ']' {
			replacedString := strings.Replace(k, "@", file, -1)
			cmdOutput, _ := exec.Command("dash", "-c", replacedString[2:len(replacedString)-1]).Output()
			elemOutput = string(cmdOutput)
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

		Addstr(s, elements[k], x, y-(co.YBuffBottom)+1, elemOutput)
		if num, _ := strconv.Atoi(string(k[0])); num != len(keys) {
			Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), x+len(elemOutput), y-(co.YBuffBottom)+1, co.BarDiv)
			x += len(elemOutput + co.BarDiv)
		} else {
			x += len(elemOutput)
		}
	}
	if x > width {
		Addstr(s, tcell.StyleDefault.Background(co.BarBg), width-3, y-(co.YBuffBottom)+1, "...")
	}
	barLen = x

	return nil
}

func DrawScreen(s tcell.Screen, currFs []string, currF int, y int, buf1 int, buf2 int) error {
	s.Clear()
	if buf1 == 0 {
		err := DispFiles(s, currFs)

		if err != nil {
			return err
		}
	} else {
		err := DispFiles(s, currFs[buf1:buf2+1])

		if err != nil {
			return err
		}
	}
	if len(currFs) > 0 {
		err := DispBar(s, co.BarStyle, currFs[currF], currF+1, len(currFs))

		if err != nil {
			return err
		}

		err = SelFile(s, co.XBuff, y, currFs[currF])

		if err != nil {
			return err
		}
	}

	BorderPipes(s)
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
	}

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