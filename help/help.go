package help

import (
	"fmt"
	co "catfm/config"
)

func Help() {
	var user string

	for k, v := range co.Bindings {
		if v[0] == "cd" {
			user += k + ": Jump to '" + v[1] + "'\n"
		} else {
			user += k + ": Execute '" + v[1] +"'\n"
		}
	}

	fmt.Printf(`Usage: catfm [OPTION] [DIRECTORY]
Extensible interactive shell for the UNIX operating system

--help, -h: Display this information
[DIRECTORY]: Jump to the specified directory

 Movement (arrow keys will always work)

%s: Move up
%s: Move down
%s: Move left
%s: Move right
%s: Go to first file
%s: Go to last file

 Manipulating files

%s: Select/deselect current file
%s: Select all
%s: Deselect all
%s: Copy selected file(s)
%s: Move selected file(s)
%s: Delete file
%s: Delete all selected files
%s: Rename file
%s: Move file to the trash directory

 Miscellaneous

%s: Refresh current directory
%s: Quit
%s: Toggle search mode
%s: Toggle hidden files and directories

 User defined

%s`, co.KeyUp, co.KeyDown, co.KeyLeft, co.KeyRight, co.KeyGoToFirst, co.KeyGoToLast, co.KeySelect, co.KeySelectAll, co.KeyDeselectAll, co.KeyCopy, co.KeyMove, co.KeyDelete, co.KeyBulkDelete, co.KeyRename, co.KeyRecycle, co.KeyRefresh, co.KeyQuit, co.KeyToggleSearch, co.KeyDotToggle, user)
}
