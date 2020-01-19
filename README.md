<h1 align="center">catfm 🐟</h1>

<p align="center">a simple and programmable file manager written in Go</p>
<img align="right" src="media/catfm.png">

## About

catfm, or **C**ool **A**nd **T**echnical **F**ile **M**anager was written to emulate the pluto file manager that I started working on in late 2018. Pluto was slow and buggy so I decided that I needed to rewrite it from scratch in a new language. So voilà, you have catfm!

## Features

catfm is feature rich, but still maintains simiplicity by giving the user complete control over the program. catfm's most interesting feature is that it doesn't have a configuration file, all configuration is done by editing the source code! This may sound scary but it is quite simple and makes a lot of sense when you start to think about it. Most programs load large files that contain all of the settings that the user has set, catfm avoids this by having the user embed all of their options in an easy to read and understand source file located in ```config/```. The config file isn't just for enabling and disabling features, it allows you to utilize you favorite programs, scripts, and commands all inside of catfm. Do you want catfm to pull up an image of you father-in-law's 70th birthday party when you press the 'Y' key? well it's possible in catfm! Catfm also has the following features:

* Customizable file formating
* Opening files in your favorite programs
* Script intergration
* Bookmarks
* Tabs
* A customizable bar
* Overall aesthetic customizability

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
