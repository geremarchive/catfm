/*

Welcome To The Catfm Source Code

------------------------------------

Info:

∙ This project was created and is maintained by Jonah G. Rongstad (geremachek).
∙ Catfm is registered under the GPL v. 3 witch means you are free to modify and distribute it and its source code
∙ This project was based on the now deprecated pluto file-manager and shares a lot of the core design ideas of that file-manager

More Info:

∙ <http://github.com/geremachek/catfm/>
∙ <http://geremachek.io/software>

*/

package main

import (
	"github.com/gdamore/tcell"
	"os"
	"fmt"
	"catfm/funcs"
	"catfm/structs"
	"catfm/help"
	"catfm/config"
)

func main() {
	if len(os.Args) > 1 { // check if the user has supplied any arguments
		if os.Args[1] == "-h" || os.Args[1] == "--help" { // display help if they use the help flags
			help.Help()

			return
		} else { // otherwise move to the directory they supplied
			err := os.Chdir(os.Args[1])

			if err != nil {
				funcs.Perr(fmt.Sprintf("couldn't open '%s'", os.Args[1]))
			}
		}
	}

	config.Init()

	s, err := tcell.NewScreen()

	if err != nil {
		funcs.Perr("couldn't draw screen")
	}

	s.Init()

	defer s.Fini()

	session, err := structs.NewCatfm(s)

	if err != nil {
		funcs.Errout(s, "couldn't start catfm")
	}

	if err := session.Views[session.View].DrawScreen(s, session); err != nil {
		funcs.Errout(s, "couldn't draw screen")
	}

	for {
		session.Views[session.View].Resize(s, session)

		input := s.PollEvent()

		switch input := input.(type) {
			case *tcell.EventKey:
				s = session.Parse(s, input)
		}
	}
}
