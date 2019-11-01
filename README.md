i# lunae

<img src="scrot.png" alt="scrot"/>

## About

Lunae (pronounced loon-ay) is a simple file manager that aims to improve on the now deprecated pluto file manager

## Dependencies 

* ```go```
* ```tcell```
* ```dash```

## Building

**Clone the repository:**

```git clone http://github.com/geremachek/lunae```

**Tnstall go from your distro's repositories: (Arch Linux for example)**

```sudo pacman -S go```

**Install dash from your distro's repositories: (Arch Linux for example)**

```sudo pacman -S dash```

**Install tcell:**

```go get github.com/gdamore/tcell```

**Build the program**

Move into the ```lunae/``` directory and type ```go build```

**Move the binary to somewhere in your path**

## Configuration

You can configure the program in the ```config/config.go``` file before compiling
