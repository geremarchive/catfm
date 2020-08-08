package structs

import (
	"github.com/gdamore/tcell"
	"os"
	"os/user"
	"os/exec"
	"strings"
	"sort"
	"fmt"
	"io/ioutil"
	"strconv"
	"code.cloudfoundry.org/bytefmt"
	co "catfm/config"
	fu "catfm/funcs"
	ke "catfm/keys"
)

type View struct {
	File int
	Files []string

	Buffer1 int
	Buffer2 int

	Y int
	Dot bool
	Cwd string

	Width int
	Height int
}

// NewView
// create a new view

func NewView(s tcell.Screen) (View, error) {
	w, h := s.Size()

	var (
		nv View = View {0, []string{}, 0, h-(co.YBuffBottom+co.YBuffTop), co.YBuffTop, false, "", w, h}
		err error
	)

	if nv.Cwd, err = os.Getwd(); err != nil {
		return nv, err
	}

	if nv.Files, err = fu.GetFiles(nv.Cwd, false); err != nil {
		return nv, err
	}

	return nv, nil
}

// ResetVals
// reset positional variables

func (v *View) ResetVals() {
	v.Buffer1 = 0

	v.Buffer2 = v.Height-(co.YBuffBottom+co.YBuffTop)
	v.File = 0
	v.Y = co.YBuffTop
}

func (v *View) ChangeDir(s tcell.Screen, cf Catfm, dir string) {
	err := os.Chdir(dir)
	
	if err == nil {
		v.Cwd, err = os.Getwd()

		if err != nil {
			fu.Errout(s, "unable to get the working directory")
		}

		v.Files, err = fu.GetFiles(v.Cwd, v.Dot) // update the files

		if err != nil {
			fu.Errout(s, "couldn't read files")
		}

		v.ResetVals()

		if err := v.DrawScreen(s, cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}
	}
}

// Quit
// quit catfm

func (cv View) Quit(s tcell.Screen) {
	file, err := os.Create("/tmp/kitty") // create the file to store the CWD

	defer file.Close()

	if err != nil {
		fu.Errout(s, "couldn't create /tmp/kitty") // error out if it wasn't able to be created
	} else {
		s.Fini() // close the screen

		_, err := file.WriteString(cv.Cwd)

		if err != nil {
			fu.Errout(s, "couldn't write to /tmp/kitty")
		}
	}

	os.Exit(0)
}

func (v *View) ParseBinding(s tcell.Screen, cf Catfm, val []string) tcell.Screen {
	replacedString := val[1] // store the actual command

	if len(v.Files) != 0 {
		// if the current directory isn't empty, replace the 
		// '@' symbol with the currently selected file

		replacedString = strings.Replace(val[1], "@", v.Files[v.File], -1)
	}

	if val[0] == "cd" { // if the command type is 'cd' (we move to a new directory)
		u, err := user.Current()

		if err != nil {
			fu.Errout(s, "unable to get the current user")
		}

		v.ChangeDir(s, cf, strings.Replace(val[1], "~", u.HomeDir, -1)) // go to the specified directory (replace ~ with $HOME)

	} else if val[0] == "t" {
		cmd := exec.Command(co.Shell, "-c", replacedString)

		s.Fini()

		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Run()

		var err error
		s, err = tcell.NewScreen()

		if err != nil {
			fu.Errout(s, "couldn't initialize screen")
		}

		s.Init()

		if err := v.DrawScreen(s, cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}
	} else if val[0] == "g" {
		cmd := exec.Command(co.Shell, "-c", replacedString)
		cmd.Start()
	}

	return s
}

func (v View) DispFiles(s tcell.Screen, cf Catfm) error {
	if len(v.Files) != 0 {
		buff := len(v.Files)

		if v.Buffer1 > 0 {
			buff = v.Buffer2
		} else if len(v.Files) > v.Buffer2 {
			buff = v.Height-(co.YBuffTop+co.YBuffBottom)
		}

		for i, f := range v.Files[v.Buffer1:buff] {
			if _, err := ShowFile(s, cf, co.XBuff, i+co.YBuffTop, f); err != nil {
				return err
			}
		}
	} else {
		msg, err := v.FormatText(s, cf, "The cat can't seem to find anything here...")

		if err != nil {
			return err
		}

		fu.Addstr(s, tcell.StyleDefault, co.XBuff, co.YBuffTop, msg)
	}

	return nil
}

func (v View) SelFile(s tcell.Screen, cf Catfm) error {
	formated, err := v.FormatText(s, cf, v.Files[v.File])

	if err != nil {
		return err
	}

	if co.SelectType == "full" {
		fu.Addstr(s, co.SelectStyle, co.XBuff, v.Y, formated + strings.Repeat(" ", v.Width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		fu.Addstr(s, co.SelectArrowStyle, co.XBuff, v.Y, co.SelectArrow)

		if co.SelectType == "arrow-default" {
			fu.Addstr(s, co.SelectStyle, co.XBuff+len(co.SelectArrow), v.Y, formated)
		} else {
			if _, err = ShowFile(s, cf, co.XBuff+len(co.SelectArrow), v.Y, v.Files[v.File]); err != nil {
				return err
			}
		}
	} else if co.SelectType == "default" {
		fu.Addstr(s, co.SelectStyle, co.XBuff, v.Y, formated)
	}

	return nil
}

func (v View) DSelFile(s tcell.Screen, cf Catfm) error {
	formated, err := ShowFile(s, cf, co.XBuff, v.Y, v.Files[v.File]) // format the text

	if err != nil {
		return err
	}

	if co.SelectType == "full" {
		// remove the full line

		fu.Addstr(s, tcell.StyleDefault, co.XBuff+len(formated), v.Y, strings.Repeat(" ", v.Width-(len(formated)+(co.XBuff*2))))
	} else if co.SelectType == "arrow" || co.SelectType == "arrow-default" {
		// remove the arrow

		fu.Addstr(s, tcell.StyleDefault, co.XBuff+len(formated), v.Y, strings.Repeat(" ", len(co.SelectArrow)+1))
	}

	return nil
}

func (v View) DispBar(s tcell.Screen, file string, curr int, vi int) error {
	if co.BarLocale != "" {
		var (
			x int = co.XBuff
			elemOutput string
			loc int
			err error

			keys []string
		)

		if co.BarLocale == "bottom" {
			loc = v.Height-(co.YBuffBottom)+1 // set the location to the bottom of the screen
		} else if co.BarLocale == "top" {
			loc = co.YBuffTop-2 // set the location ot the top of the screen
		}

		fu.Addstr(s, tcell.StyleDefault, co.XBuff, loc, strings.Repeat(" ", v.Width)) // clear the last bar

		for k, _ := range co.BarStyle { // add keys to the slice
			keys = append(keys, k)
		}

		sort.Strings(keys)

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
				elemOutput = fmt.Sprintf("%d/%d", curr, len(v.Files))
			} else if k[1] == '[' && k[len(k)-1] == ']' {
				replacedString := strings.Replace(k, "@", file, -1)
				cmdOutput, _ := exec.Command(co.Shell, "-c", replacedString[2:len(replacedString)-1]).Output()
				elemOutput = string(cmdOutput)
			} else if k[1:] == "tab" {
				elemOutput = fmt.Sprintf("[%d]", vi+1)
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
				elemOutput = strings.Replace(elemOutput, "$TAB", strconv.Itoa(vi+1), -1)
			}

			fu.Addstr(s, co.BarStyle[k], x, loc, elemOutput)
	
			if num, _ := strconv.Atoi(string(k[0])); num < len(keys) { // the element isn't at the end
				fu.Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), x+len(elemOutput), loc, co.BarDiv)
				x += len(elemOutput + co.BarDiv)
			} else {
				x += len(elemOutput)
			}
		}

		if x > v.Width {
			fu.Addstr(s, tcell.StyleDefault.Background(co.BarBg).Foreground(co.BarFg), v.Width-3, loc, "...")
		}
	}

	return nil
}

func (v View) DrawScreen(s tcell.Screen, cf Catfm) error {
	s.Clear()

	if err := v.DispFiles(s, cf); err != nil { // draw the files
		return err
	}

	go BorderPipes(s) // display the pipes

	if len(v.Files) > 0 {
		// if there are files to display, display the bar and select a file

		if err := v.DispBar(s, v.Files[v.File], v.File+1, cf.View); err != nil {
			return err
		}

		if err := v.SelFile(s, cf); err != nil {
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

		switch co.PipeType {
			case "thin": used = thin
			case "thick": used = thick
			case "hollow": used = hollow
			case "round": used = round
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
			user, eu := user.Current()
			host, eh := os.Hostname()

			var text string

			switch co.PipeText {
				case "user@host":
					if eu == nil && eh == nil {
						text = user.Name + "@" + host
					}
				case "catfm@host":
					if eh == nil {
						text = "catfm@" + host
					}
				case "user@catfm":
					if eu == nil {
						text = user.Name + "@catfm"
					}
				default: text = co.PipeText
			}

			fu.Addstr(s, co.PipeTextStyle, 1, 0, text)
			s.Show()
		}
	}
}

func (v *View) GoToLast(s tcell.Screen, cf Catfm) {
	v.File = len(v.Files) - 1 // set index to the last file
	
	if len(v.Files) + co.YBuffTop + co.YBuffBottom > v.Height { // should we scroll the line?
		v.Y = v.Height - (co.YBuffBottom + 1) // send the cursor to the bottom of the screens

		v.Buffer1 = (len(v.Files))-(v.Height-(co.YBuffTop+co.YBuffBottom)) // scroll
		v.Buffer2 = len(v.Files)
	} else {
		v.Y = v.File + co.YBuffTop // send the cursor to the last file
	}

	if err := v.DrawScreen(s, cf); err != nil {
		fu.Errout(s, "couldn't draw screen")
	}
}

func (v *View) GoToFirst(s tcell.Screen, cf Catfm) {
	v.ResetVals()

	if err := v.DrawScreen(s, cf); err != nil {
		fu.Errout(s, "couldn't draww screen")
	}
}

func (v *View) Move(s tcell.Screen, cf Catfm, n int) {
	if len(v.Files) > 0 {
		if v.File == 0 && n == -1 {
			v.GoToLast(s, cf) // go to last if moving up and the first file is selected
		} else if v.File == len(v.Files)-1 && n == 1 {
			v.GoToFirst(s, cf) // go to first if moving down and the last file is selected
		} else if (v.Y == (v.Height-1)-co.YBuffBottom && n == 1) || (v.Y == co.YBuffTop && n == -1) { // shall we scroll?
			v.Buffer1 += n // scroll
			v.Buffer2 += n

			v.File += n // move the file

			if err := v.DrawScreen(s, cf); err != nil {
				fu.Errout(s, "couldn't draw screen")
			}
		} else {
			c := make(chan error)

			go func(c chan error) {
				cff := 0

				if n == 1 {
					cff = 2
				}

				c <- v.DispBar(s, v.Files[v.File+n], v.File+cff, cf.View)
			}(c)

			err := <-c

			if err != nil {
				fu.Errout(s, "unable to display the infobar")
			}

			if err := v.DSelFile(s, cf); err != nil {
				fu.Errout(s, "unable to deselect file")
			}

			v.Y += n
			v.File += n

			if err := v.SelFile(s, cf); err != nil {
				fu.Errout(s, "unable to select file")
			}

			s.Show()

		}
	}
}

func (v *View) Resize(s tcell.Screen, cf Catfm) (out bool) {
	nw, nh := s.Size()
	var err error

	if v.Width != nw && v.Height == nh {
		v.Width, v.Height = nw, nh
		err = v.DrawScreen(s, cf)
		out = true
	} else if v.Height != nh {
		v.Width, v.Height = nw, nh

		if v.Buffer1 != 0 || v.File > v.Height-(co.YBuffTop+co.YBuffBottom) {
			v.ResetVals()
		}

		err = v.DrawScreen(s, cf)
		out = true
	}

	if err != nil {
		fu.Errout(s, "unable to draw screen")
	}

	return
}

func (v *View) Search(s tcell.Screen, cf Catfm) {
	var (
		original []string = v.Files
		searching bool = true
		search string
	)

	for searching {
		v.Resize(s, cf)

		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				if ke.MatchKey(input, co.KeyToggleSearch) { // stop search if the user types the toggle key
					searching = false
				} else if input.Key() == tcell.KeyDEL || input.Key() == tcell.KeyBS { // delete last character on backspace
					if len(search) > 0 {
						if len(search) == 1 {
							search = ""
						} else {
							search = search[:len(search)-1]
						}
					}
				} else {
					search += string(input.Rune()) // add character
				}

				v.Files = []string{}

				for _, e := range original { // add matching files
					if strings.Contains(e, search) {
						v.Files = append(v.Files, e)
					}
				}

				if v.File != 0 { // reset everything
					v.ResetVals()
				}

				if err := v.DrawScreen(s, cf); err != nil {
					fu.Errout(s, "unable to draw screen")
				}
		}
	}
}

// DotToggle
// toggle hidden files

func (v *View) DotToggle(s tcell.Screen, cf Catfm) {
	var err error

	v.Dot = !(v.Dot) // toggle the dot variable
	v.Files, err = fu.GetFiles(v.Cwd, v.Dot) // update the files

	if err == nil {
		// reset everything

		v.ResetVals()

		// draw screen

		if err := v.DrawScreen(s, cf); err != nil {
			fu.Errout(s, "couldn't draw screen")
		}
	}
		
}

// Rename
// Rename a file

func (v *View) Rename(s tcell.Screen, cf Catfm) {
	err := ioutil.WriteFile("/tmp/rename.catfm", []byte(v.Files[v.File]), 0644) // create a new temporary file containing the current file name

	if err == nil { // if that was successful...
		replaced := strings.Replace(co.FileOpen["*"][1], "@", "/tmp/rename.catfm", -1) // use the default editor and replace '@' with /tmp/rename.catfm
		s = v.ParseBinding(s, cf, []string{co.FileOpen["*"][0],  replaced}) // run that sucker

		if _, err := os.Stat("/tmp/rename.catfm"); err == nil { // make sure it still exists...
			file, err := os.Open("/tmp/rename.catfm") // open the file 

			defer file.Close()

			if err == nil {
				b, err := ioutil.ReadAll(file)

				if err == nil {
					newName := strings.Split(string(b), "\n")[0]
					_, in := fu.In(newName, v.Files)

					if newName != "" && !(in) { // if the new name isn't blank or already in existance, 
						os.Rename(v.Files[v.File], newName) // rename the file
						v.Files[v.File] = newName

						if err := v.DrawScreen(s, cf); err != nil { // draw the screen
							fu.Errout(s, "couldn't draw screen")
						}
					}
				}
			}
		}
	}
	
}

/// Right
// this will cd into a directory or open a file

func (v *View) Right(s tcell.Screen, cf Catfm) {
	if len(v.Files) > 0 { // is there anything here?
		if fu.Isd(v.Files[v.File]) { // is it a directory?
			v.ChangeDir(s, cf, v.Files[v.File])
		} else { // otherwise open the file
			var command []string = co.FileOpen["*"]
			splitFile := strings.Split(v.Files[v.File], ".")
			
			for k, v := range co.FileOpen {
				if k == splitFile[len(splitFile)-1] {
					command = v
					break
				}
			}

			s = v.ParseBinding(s, cf, command)
		}
	}
}

// Refresh
// refresh the files

func (v *View) Refresh(s tcell.Screen, cf Catfm) {
	old := v.Files
	var err error

	v.Files, err = fu.GetFiles(v.Cwd, v.Dot)

	if err != nil {
		fu.Errout(s, "unable to read files")
	}

	if len(v.Files) < len(old) {
		v.ResetVals()
	} 
	
	if err := v.DrawScreen(s, cf); err != nil {
		fu.Errout(s, "couldn't draw screen")
	}
}

// FormatText
// format a file so it is ready to be displayed

func (v View) FormatText(s tcell.Screen, cf Catfm, text string) (string, error) {
	var (
		tlen int
		buflen int

		wmin int = 0
		bufmin int = 3
		
		slash string = ""
		aster string = ""
		
		selected bool
	)

	if len(v.Files) > 0 {
		selected = v.Files[v.File] == text
	}
	
	if selected && (co.SelectType == "arrow" || co.SelectType == "arrow-default") { // check if the user is using arrow select mode
		tlen = len(text)+len(co.SelectArrow) // add arrow to text lengths
		buflen = len(co.SelectArrow) // ignore the arrow length
	} else {
		tlen = len(text)
		buflen = 0
	}

	if fu.Isd(text) {
		slash = "/"

		wmin = 1
		bufmin = 4

		if cf.IsSel(v.Cwd + "/" + text) {
			aster = "*"

			wmin = 2
			bufmin = 5
		}
	} else {
		if cf.IsSel(v.Cwd + "/" + text) {
			aster = "*"

			wmin = 1
			bufmin = 4
		}
	}

	if tlen > v.Width-(co.XBuff*2)-wmin {
		return aster+text[:v.Width-(co.XBuff*2)-buflen-bufmin] + "..." + slash, nil
	} else {
		return aster+text+slash, nil
	}
}

