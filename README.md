<h1 align="center">catfm 🐟</h1>

<p align="center">catfm is a simple and programmable file manager that aims to improve on the now deprecated pluto file manager</p>
<img align="right" src="media/scrot.png">

## Dependencies 

* ```go```
* ```tcell```
* ```bytefmt```
* ```copy```

## Building

**Clone the repository:**

```git clone http://github.com/geremachek/catfm```

**Tnstall go from your distro's repositories: (Arch Linux for example)**

```sudo pacman -S go```

**Install tcell:**

```go get github.com/gdamore/tcell```

**Install bytefmt**

```go get code.cloudfoundry.org/bytefmt```

**Install copy**

```go get github.com/otiai10/copy```

**Build the program**

Move into the ```catfm/``` directory and type ```go build```

**Move the binary to somewhere in your path**

## Configuration

You can configure the program in the ```config/config.go``` file before compiling. This speeds up the program as it doesn't have to read and parse a giant config file everytime you start up the program

### ```cd``` on exit

shell function (put this in your ```.shellrc```):

```bash
fm() {
	catfm && cd "$(< /tmp/kitty)"
}
```

## Todo

- [X] Run program/script/command on keypress
- [X] Run custom commands in the bar
- [X] Add hotkeys for directories
- [ ] Add file searching
- [ ] Add file renaming
- [X] Add ability to customize the keys for delete, move, copy, movement, etc.
- [X] Move to top or bottom screen
