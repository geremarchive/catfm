package funcs

import (
	co "catfm/config"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"sort"
	"os"
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

func GetFiles(path string, dot bool) ([]string, error) {
	var d []string
	var f []string
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return []string{}, err
	}

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
	return append(d, f...), nil
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

func FormatText(s tcell.Screen, text string, sel bool) (string, error) {
	width, _ := s.Size()
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

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
				return "*"+text[:width-(co.XBuff*2)-buflen-5] + ".../", nil
			} else {
				return "*"+text+"/", nil
			}
		} else {
			if tlen > width-(co.XBuff*2)-1 {
				return text[:width-(co.XBuff*2)-buflen-4] + ".../", nil
			} else {
				return text+"/", nil
			}
		}
	} else {
		if IsSel(cwd + "/" + text) {
			if tlen > width-(co.XBuff*2)-1 {
				return "*"+text[:width-(co.XBuff*2)-buflen-4] + "...", nil
			} else {
				return "*"+text, nil
			}
		} else {
			if tlen > width-(co.XBuff*2) {
				return text[:width-(co.XBuff*2)-buflen-3] + "...", nil
			} else {
				return text, nil
			}
		}
	}
}

func IsSel(path string) bool {
	_, in := In(path, co.Selected)
	return in
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
