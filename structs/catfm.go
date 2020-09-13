package structs

import (
	"github.com/gdamore/tcell"
	"strings"
	"os"
	fu "catfm/funcs"
	co "catfm/config"
	cp "github.com/otiai10/copy"
)

type Catfm struct {
	Views []View
	View int

	Selected []string
}

// NewCatfm
// create a new instance of catfm

func NewCatfm(s tcell.Screen) (cf Catfm, err error) {
	var v View

	if v, err = NewView(s); err != nil {
		return
	}

	cf = Catfm {[]View{v,v,v,v,v,v,v,v,v,v}, 0, []string{}}

	return
}

// IsSel
// check if a given file is currently selected

func (cf Catfm) IsSel(path string) bool {
	_, in := fu.In(path, cf.Selected)
	return in
}

// OprSelected
// operate on the selected files

func (cf *Catfm) OprSelected(s tcell.Screen, m uint8) { // 0 = copy, 1 = move, 2 = bulk delete 
	var err error
	v := cf.Views[cf.View]

	for _, elem := range cf.Selected { // loop through the selected files
		split := strings.Split(elem, "/") // split the path

		if m < 2 { // copying or moving
			err = cp.Copy(elem, v.Cwd + "/" + split[len(split)-1]) // copy the file

			if m == 1 && err == nil {
				os.RemoveAll(elem)
			}

			
		} else {
			os.RemoveAll(elem)
		}
	}

	if err == nil {
		cf.Selected = []string{}

		v.Files, err = fu.GetFiles(v.Cwd, v.Dot)

		if err != nil {
			fu.Errout(s, "unable to read files")
		}

		if m == 2 {
			v.ResetVals()
		}

		if err := v.DrawScreen(s, *cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}

		cf.Views[cf.View] = v
	}
}

// RemoveFile
// remove a file from the screen and buffer

func (cf *Catfm) RemoveFile(s tcell.Screen) {
	v := cf.Views[cf.View]

	if err := os.RemoveAll(v.Files[v.File]); err == nil {
		if index, in := fu.In(v.Files[v.File], cf.Selected); in {
			cf.Selected = append(cf.Selected[:index], cf.Selected[index+1:]...)
		}

		v.Files = append(v.Files[:v.File], v.Files[v.File+1:]...)

		if v.File == len(v.Files) { // is the file at the bottom of the screen?
			v.File -= 1
			
			if v.Buffer1 == 0 { // is the line unscrolled?
				v.Y -= 1 // move the cursor
			} else { // otherwise scroll
				v.Buffer1 -= 1
				v.Buffer2 -= 1
			}
		}

		cf.Views[cf.View] = v

		// draw to the screen

		if err := v.DrawScreen(s, *cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}
	}
}

// Recycle
// move a file to the recycle bin

func (cf *Catfm) Recycle(s tcell.Screen) {
	if co.RecycleBin != "" {
		v := cf.Views[cf.View]

		if err := cp.Copy(v.Files[v.File], co.RecycleBin + "/" + v.Files[v.File]); err == nil {
			cf.RemoveFile(s)
		} else {
			s.Fini()
			panic(err)
		}
	}
}

// Select
// add a file to the list of selected files

func (cf *Catfm) Select(s tcell.Screen) {
	v := cf.Views[cf.View]

	if cf.IsSel(v.Cwd + "/" + v.Files[v.File]) { // if the higlighted file is already selected...
		index, _ := fu.In(v.Cwd + "/" + v.Files[v.File], cf.Selected)
		cf.Selected = append(cf.Selected[:index], cf.Selected[index+1:]...) // remove it!
	} else {
		cf.Selected = append(cf.Selected, v.Cwd + "/" + v.Files[v.File]) // otherwise add it.
	}

	formated, err := v.FormatText(s, *cf, v.Files[v.File]) // format 

	if err == nil {
		// remove the selected file

		if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
			fu.Addstr(s, tcell.StyleDefault, co.XBuff, v.Y, formated + strings.Repeat(" ", len(co.SelectArrow)+1))
		} else {
			fu.Addstr(s, tcell.StyleDefault, co.XBuff, v.Y, formated + " ")
		}

		// select the file

		if err := v.SelFile(s, *cf); err != nil {
			fu.Errout(s, "unable to calculate screen width")
		}
		
		s.Show()
	}

	cf.Views[cf.View] = v
}

// SelectAll
// select all of the files in the current directory.

func (cf *Catfm) SelectAll(s tcell.Screen) {
	v := cf.Views[cf.View]

	for _, elem := range v.Files { // loop through files
		if _, in := fu.In(v.Cwd + "/" + elem, cf.Selected); !in {
			// if the file isn't already selected, add it to the list.

			cf.Selected = append(cf.Selected, v.Cwd + "/" + elem)
		}
	}

	if err := v.DrawScreen(s, *cf); err != nil {
		fu.Errout(s, "couldn't draw screen")
	}
	
	cf.Views[cf.View] = v
}

// DeselectAll
// deselect all of the files

func (cf *Catfm) DeselectAll(s tcell.Screen) {
	cf.Selected = []string{}

	if err := cf.Views[cf.View].DrawScreen(s, *cf); err != nil {
		fu.Errout(s, "couldn't draw screen")
	}
}

// TabSwitch
// switch between tabs

func (cf *Catfm) TabSwitch(s tcell.Screen, tab rune) {
	if tab == '0' {
		cf.View = 9
	} else {
		cf.View = int(tab-49)
	}

	if err := os.Chdir(cf.Views[cf.View].Cwd); err != nil {
		fu.Errout(s, "couldn't change directory")
	}

	if !(cf.Views[cf.View].Resize(s, *cf)) {
		if err := cf.Views[cf.View].DrawScreen(s, *cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}
	}
}

// ShowFile
// Display a file in its proper style

func ShowFile(s tcell.Screen, cf Catfm, x int, y int, f string) (string, error) {
	formated, err := cf.Views[cf.View].FormatText(s, cf, f)

	if err != nil {
		return "", err
	}

	if fu.Isd(f) {
		fu.Addstr(s, co.FileColors["[dir]"], x, y, formated)
	} else {
		split := strings.Split(f, ".")
		fu.Addstr(s, co.FileColors[split[len(split)-1]], x, y, formated)
	}

	return formated, nil
}

