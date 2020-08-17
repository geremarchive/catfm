<h1 align="center">catfm</h1>

<p align="center"><img src="media/catfm.png"></p>

<p align="center">An extensible interactive shell for the UNIX operating system</p>

## About

```catfm```, or **C**ompact **A**nd **T**weakable **F**ile **M**anager is an interactive shell for your UNIX system that runs right in your terminal. It not only allows the user to manage files, but to complete complex actions through simple user defined keyboard shortcuts

## Features

* Customizable file formating
* Opening files in your favorite programs
* Script intergration
* Bookmarks
* Tabs
* A customizable bar
* Overall aesthetic customizability

## Dependencies 

* ```go 1.13```
* a POSIX enviroment and shell
* a terminal that supports cursor adressing, etc.

## Building

**Clone the repository:**

```git clone http://github.com/geremachek/catfm```

Move into the ```catfm/``` directory and type ```go build```

**Move the binary to somewhere in your path**

## Configuration

You can configure the program in the ```config/config.go``` file before compiling. This speeds up the program as it doesn't have to read and parse a giant config file everytime you start up the program

I also recommend looking at the [tcell](https://godoc.org/github.com/gdamore/tcell) documentation on [color](https://godoc.org/github.com/gdamore/tcell#Color) and [styles](https://godoc.org/github.com/gdamore/tcell#Style)

### ```cd``` on exit

shell function (put this in your ```.shellrc```):

```bash
fm() {
	catfm && cd "$(< /tmp/kitty)"
}
```

## Tutorial

Coming soon...

## Todo

- [X] Run program/script/command on keypress
- [X] Run custom commands in the bar
- [X] Add hotkeys for directories
- [X] Add file searching
- [X] Add file renaming
- [X] Add ability to customize the keys for delete, move, copy, movement, etc.
- [X] Move to top or bottom screen
- [X] Add tabs
- [X] Improve config file and keybinds
