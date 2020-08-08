package structs

import (
	"github.com/gdamore/tcell"
	co "catfm/config"
	ke "catfm/keys"
)

func (cf *Catfm) Parse(s tcell.Screen, input *tcell.EventKey) {
	if ke.MatchKey(input, co.KeyQuit) {
		cf.Views[cf.View].Quit(s)	
	} else if ke.MatchKey(input, co.KeyDelete) {
		cf.RemoveFile(s)
	} else if ke.MatchKey(input, co.KeyRecycle) {
		cf.Recycle(s)	
	} else if ke.MatchKey(input, co.KeyCopy) {
		cf.OprSelected(s, 0)
	} else if ke.MatchKey(input, co.KeyMove) {
		cf.OprSelected(s, 1)
	} else if ke.MatchKey(input, co.KeyBulkDelete) {
		cf.OprSelected(s, 2)
	} else if ke.MatchKey(input, co.KeySelect) {
		cf.Select(s)
	} else if ke.MatchKey(input, co.KeySelectAll) {
		cf.SelectAll(s)
	} else if ke.MatchKey(input, co.KeyDeselectAll) {
		cf.DeselectAll(s)
	} else if ke.MatchKey(input, co.KeyDotToggle) {
		cf.Views[cf.View].DotToggle(s, *cf)
	} else if ke.MatchKey(input, co.KeyGoToFirst) {
		cf.Views[cf.View].GoToFirst(s, *cf)
	} else if ke.MatchKey(input, co.KeyGoToLast) {
		cf.Views[cf.View].GoToLast(s, *cf)
	} else if ke.MatchKey(input, co.KeyRename) {
		cf.Views[cf.View].Rename(s, *cf)
	} else if ke.MatchKey(input, co.KeyDown) || input.Key() == tcell.KeyDown {
		cf.Views[cf.View].Move(s, *cf, 1)
	} else if ke.MatchKey(input, co.KeyUp) || input.Key() == tcell.KeyUp {
		cf.Views[cf.View].Move(s, *cf, -1)
	} else if ke.MatchKey(input, co.KeyRight) || input.Key() == tcell.KeyRight {
		cf.Views[cf.View].Right(s, *cf)
	} else if ke.MatchKey(input, co.KeyLeft) || input.Key() == tcell.KeyLeft {
		cf.Views[cf.View].ChangeDir(s, *cf, "..")
	} else if ke.MatchKey(input, co.KeyRefresh) {
		cf.Views[cf.View].Refresh(s, *cf)
	} else if ke.MatchKey(input, co.KeyToggleSearch) {
		cf.Views[cf.View].Search(s, *cf)
	} else if input.Rune() >= 48 && input.Rune() <= 57 {
		cf.TabSwitch(s, input.Rune())
	} else {
		for k, v := range co.Bindings {
			if ke.MatchKey(input, k) {
				s = cf.Views[cf.View].ParseBinding(s, *cf, v)
				break
			}
		}
	}
}
