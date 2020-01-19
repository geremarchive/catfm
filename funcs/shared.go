package funcs

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

var (
	Selected []string
	ViewNumber int
)
