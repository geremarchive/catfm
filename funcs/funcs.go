package funcs

import (
	co "../config"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"io"
	"sort"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"code.cloudfoundry.org/bytefmt"
)

var (
	barLen int
)

func In(item string, array []string) (int, bool) {
	for i, elem := range array {
		if elem == item {
			return i, true
		}
	}
	return 0, false
}

func GetFiles(path string, dot bool) ([]string) {
	var d []string
	var f []string
	files, _ := ioutil.ReadDir(path)
	for _, elem := range files {
		if elem.Name()[0] == '.' && !(dot) {
			continue
		} else if elem.IsDir() {
			d = append(d, string(elem.Name()))
		} else {
			f = append(f, string(elem.Name()))
		}
	}
	sort.Strings(d)
	sort.Strings(f)
	return append(d, f...)
}

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		s.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func Isd(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return true
	} else {
		return false
	}
}

func FormatText(s tcell.Screen, text string) string {
	width, _ := s.Size()
	cwd, _ := os.Getwd()
	if Isd(text) {
		if IsSel(cwd + "/" + text) {
			if len(text) > width-(co.XBuff*2)-2 {
				return "*"+text[:width-(co.XBuff*2)-5] + ".../"
			} else {
				return "*"+text+"/"
			}
		} else {
			if len(text) > width-(co.XBuff*2)-1 {
				return text[:width-(co.XBuff*2)-4] + ".../"
			} else {
				return text+"/"
			}
		}
	} else {
		if IsSel(cwd + "/" + text) {
			if len(text) > width-(co.XBuff*2)-1 {
				return "*"+text[:width-(co.XBuff*2)-4] + "..."
			} else {
				return "*"+text
			}
		} else {
			if len(text) > width-(co.XBuff*2) {
				return text[:width-(co.XBuff*2)-3] + "..."
			} else {
				return text
			}
		}
	}
}

func DispFiles(s tcell.Screen, files []string) {
	_, height := s.Size()
	if len(files) != 0 {
		for i, f := range files {
			if i+co.YBuffTop != (height - co.YBuffBottom) {
				splitFile := strings.Split(f, ".")
				if Isd(f) {
					Addstr(s, co.FileColors["[dir]"], co.XBuff, i+co.YBuffTop, FormatText(s, f))
				} else {
					Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], co.XBuff, i+co.YBuffTop, FormatText(s, f))
				}
			} else {
				break
			}
		}
	} else {
		Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor("#ff0000")).Bold(true), co.XBuff, co.YBuffTop, FormatText(s, "[EMPTY]"))
	}
}

func SelFile(s tcell.Screen, x int, y int, file string) {
	formated := FormatText(s, file)
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
				Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x+len(co.SelectArrow), y, FormatText(s, file))
			}
		}
	} else if co.SelectType == "default" {
		Addstr(s, co.SelectStyle, x, y, formated)
	}
}

func DSelFile(s tcell.Screen, x int, y int, file string) {
	formated := FormatText(s, file)
	splitFile := strings.Split(file, ".")
	width, _ := s.Size()

	if co.SelectType == "full" {
		Addstr(s, co.FileColors["[dir]"], x, y, formated)
		Addstr(s, tcell.StyleDefault, x+len(formated), y, strings.Repeat(" ", width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		if Isd(file) {
			Addstr(s, co.FileColors["[dir]"], x, y, formated + strings.Repeat(" ", len(co.SelectArrow)))
		} else {
			Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x, y, formated + strings.Repeat(" ", len(co.SelectArrow)))
		}
	} else if co.SelectType == "default" {
		Addstr(s, co.FileColors[splitFile[len(splitFile)-1]], x, y, formated)
	}
}

func IsSel(path string) bool {
	_, in := In(path, co.Selected)
	return in
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func Move(src, dst string) {
	Copy(src, dst)
	os.RemoveAll(src)
}

/*func Msgnread(s tcell.Screen, message string) {
	_, y := s.Size()

	Addstr(s, tcell.StyleDefault, co.XBuff, y-(co.YBuffBottom)+1, strings.Repeat(" ", barLen))
	Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), co.XBuff, y-(co.YBuffBottom)+1, message)
	s.Show()

	barLen = co.XBuff + len(message)

	//var text string
	var currX int = co.XBuff + len(message) + 1

	s.ShowCursor(currX, y-(YBuffBottom)+1)
	s.Show()

	for {
		input := s.PollEvent()
		switch input := input.(type) {
			case *tcell.EventKey:
				if input.Key() == 
				Addstr(s ,tcell.StyleDefault, currX, y-(co.YBuffBottom)+1, string(input.Rune()))
				currX += 1
				s.ShowCursor(currX, y-(co.YBuffBottom)+1)
				s.Show()
		}
	}
}*/

func DispBar(s tcell.Screen, elements map[string]tcell.Style, file string) {
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
		} else if k[1] == '[' && k[len(k)-1] == ']' {
			replacedString := strings.Replace(k, "@", file, -1)
			cmdOutput, _ := exec.Command("dash", "-c", replacedString[2:len(replacedString)-1]).Output()
			elemOutput = string(cmdOutput)
		} else {
			elemOutput = k[1:]
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
}

func DrawScreen(s tcell.Screen, currFs []string, currF int, y int, buf1 int, buf2 int) {
	s.Clear()
	if buf1 == 0 {
		DispFiles(s, currFs)
	} else {
		DispFiles(s, currFs[buf1:buf2])
	}
	DispBar(s, co.BarStyle, currFs[currF])
	SelFile(s, co.XBuff, y, currFs[currF])
	s.Show()
}
