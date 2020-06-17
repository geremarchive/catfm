package funcs

import (
	co "catfm/config"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"strings"
	"os/user"
	"os/exec"
	"sort"
	"os"
	"fmt"
)

var (
	barLen int
)

func Space(s string, x int) (out string) {
	for i := 0; i < x; i++ {
		out += s
	}

	return
}

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

	var (
		tlen int
		buflen int
	)

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
	_, in := In(path, Selected)
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

func Errout(s tcell.Screen, msg string) {
	s.Fini()
	fmt.Println("catfm: " + msg)

	os.Exit(0)
}

func (v *View) ParseBinding(s tcell.Screen, val []string) tcell.Screen {
	replacedString := val[1]

	if len(v.Files) != 0 {
		replacedString = strings.Replace(val[1], "@", v.Files[v.File], -1)
	}

	if val[0] == "cd" {
		u, err := user.Current()

		if err != nil {
			Errout(s, "unable to get the current user")
		}

		err = os.Chdir(strings.Replace(val[1], "~", u.HomeDir, -1))

		if err == nil {
			v.Cwd, err = os.Getwd()

			if err != nil {
				Errout(s, "unable to get the working directory")
			}

			v.Files, err = GetFiles(v.Cwd, v.Dot)

			if err != nil {
				Errout(s, "couldn't read files")
			}

			_, height := s.Size()

			v.Buffer1 = 0
			v.Buffer2 = (height-co.YBuffBottom)+co.YBuffTop
			v.File = 0
			v.Y = co.YBuffTop

			if err := v.DrawScreen(s); err != nil {
				Errout(s, "couldn't draw screen")
			}
		}
	} else if val[0] == "t" {
		cmd := exec.Command(co.Shell, "-c", replacedString)

		s.Fini()

		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Run()

		var err error
		s, err = tcell.NewScreen()

		if err != nil {
			Errout(s, "couldn't initialize screen")
		}

		s.Init()

		if err := v.DrawScreen(s); err != nil {
			Errout(s, "couldn't draw screen")
		}
	} else if val[0] == "g" {
		cmd := exec.Command(co.Shell, "-c", replacedString)
		cmd.Start()
	}

	return s
}
