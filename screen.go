package termloop

import "github.com/nsf/termbox-go"

// A Canvas is a 2D array of Cells, used for drawing
type Canvas [][]Cell

func newCanvas(width, height int) Canvas {
	canvas := make(Canvas, width)
	for i := range canvas {
		canvas[i] = make([]Cell, height)
	}
	return canvas
}

// A Screen represents the current state of the display.
// To draw on the screen, create Drawables and set their positions.
// Then, add them to the Screen's Level, or to the Screen directly (e.g. a HUD).
type Screen struct {
	canvas   Canvas
	level    Level
	entities []Drawable
	width    int
	height   int
	delta    float64
}

// NewScreen creates a new Screen, with no entities or level.
// Returns a pointer to the new Screen.
func NewScreen() *Screen {
	s := Screen{entities: make([]Drawable, 0)}
	s.canvas = newCanvas(10, 10)
	return &s
}

// Tick is used to process events such as input. It is called
// on every frame by the Game.
func (s *Screen) Tick(ev Event) {
	// TODO implement ticks using worker pools
	if s.level != nil {
		s.level.Tick(ev)
	}
	if ev.Type != EventNone {
		for _, e := range s.entities {
			e.Tick(ev)
		}
	}
}

// Draw is called every frame by the Game to render the current
// state of the screen.
func (s *Screen) Draw() {
	// Update termloop canvas
	s.canvas = newCanvas(s.width, s.height)
	if s.level != nil {
		s.level.DrawBackground(s)
		s.level.Draw(s)
	}
	for _, e := range s.entities {
		e.Draw(s)
	}
	// Draw to terminal
	for i, row := range s.canvas {
		for j, cell := range row {
			termbox.SetCell(i, j, cell.Ch,
				termbox.Attribute(cell.Fg),
				termbox.Attribute(cell.Bg))
		}
	}
	termbox.Flush()
}

func (s *Screen) resize(w, h int) {
	s.width = w
	s.height = h
	c := newCanvas(w, h)
	// Copy old data that fits
	for i := 0; i < min(w, len(s.canvas)); i++ {
		for j := 0; j < min(h, len(s.canvas[0])); j++ {
			c[i][j] = s.canvas[i][j]
		}
	}
	s.canvas = c
}

// Size returns the width and height of the Screen
// in characters.
func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

// TimeDelta returns the number of seconds since the previous
// frame was rendered. Can be used for timings and animation.
func (s *Screen) TimeDelta() float64 {
	return s.delta
}

// RenderCell updates the Cell at a given position on the Screen
// with the attributes in Cell c.
func (s *Screen) RenderCell(x, y int, c *Cell) {
	renderCell(&s.canvas[x][y], c)
}

func renderCell(old, new_ *Cell) {
	if new_.Ch != 0 {
		old.Ch = new_.Ch
	}
	if new_.Bg != 0 {
		old.Bg = new_.Bg
	}
	if new_.Fg != 0 {
		old.Fg = new_.Fg
	}
}
