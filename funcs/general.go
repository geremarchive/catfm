package funcs

import (
	"github.com/gdamore/tcell"
	"io/ioutil"
	"sort"
	"os"
	"fmt"
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

func Perr(msg string) {
	fmt.Fprintf(os.Stderr, "catfm: %s\n", msg)
	os.Exit(1)
}

func Errout(s tcell.Screen, msg string) {
	s.Fini()
	Perr(msg)
}
