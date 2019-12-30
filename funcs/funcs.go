package funcs

import (
	co "catfm/config"
	"github.com/gdamore/tcell"
	"io/ioutil"
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

func FormatText(s tcell.Screen, text string, sel bool) string {
	width, _ := s.Size()
	cwd, _ := os.Getwd()

	var tlen int
	var buflen int

	if sel && (co.SelectType == "arrow" || co.SelectType == "arrow-default") {
		tlen = len(text)+len(co.SelectArrow)
		buflen = len(co.SelectArrow)
	} else {
		tlen = len(text)
		buflen = 0
	}

	if Isd(text) {
		if IsSel(cwd + "/" + text) {
			if tlen > width-(co.XBuff*2)-2 {
				return "*"+text[:width-(co.XBuff*2)-buflen-5] + ".../"
			} else {
				return "*"+text+"/"
			}
		} else {
			if tlen > width-(co.XBuff*2)-1 {
				return text[:width-(co.XBuff*2)-buflen-4] + ".../"
			} else {
				return text+"/"
			}
		}
	} else {
		if IsSel(cwd + "/" + text) {
			if tlen > width-(co.XBuff*2)-1 {
				return "*"+text[:width-(co.XBuff*2)-buflen-4] + "..."
			} else {
				return "*"+text
			}
		} else {
			if tlen > width-(co.XBuff*2) {
				return text[:width-(co.XBuff*2)-buflen-3] + "..."
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
			formated := FormatText(s, f, false)
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
		Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor("#ff0000")).Bold(true), co.XBuff, co.YBuffTop, FormatText(s, "The cat can't seem to find anything here...", false))
	}
}

func SelFile(s tcell.Screen, x int, y int, file string) {
	formated := FormatText(s, file, true)
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
}

func DSelFile(s tcell.Screen, x int, y int, file string) {
	formated := FormatText(s, file, false)
	splitFile := strings.Split(file, ".")
	width, _ := s.Size()

	if co.SelectType == "full" {
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
}

func IsSel(path string) bool {
	_, in := In(path, co.Selected)
	return in
}

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
		DispFiles(s, currFs[buf1:buf2+1])
	}
	if len(currFs) > 0 {
		DispBar(s, co.BarStyle, currFs[currF])
		SelFile(s, co.XBuff, y, currFs[currF])
	}
	s.Show()
}

func Itemi(item string, slice []string) (out int) {
	for index, elem := range slice {
		if elem == item {
			out = index
			break
		}
	}

	return
}
