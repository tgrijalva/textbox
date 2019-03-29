package textbox

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

// Point is a 2d coordinate
type Point struct {
	X, Y int
}

// Rect defines a rectangle in a 2d-plane by its position and size
type Rect struct {
	Origin Point
	Width  int
	Height int
}

// Textbox is a canvas of runes of a defined dimension
type Textbox struct {
	canvas [][]rune
	pixels []rune
	width  int
	height int
	cursor int
}

// NewTextbox returns a textbox with dimentions widht and height
func NewTextbox(width, height int) *Textbox {
	tb := new(Textbox)
	tb.width = int(abs(int64(width)))
	tb.height = int(abs(int64(height)))
	tb.canvas = make([][]rune, height)
	tb.pixels = make([]rune, width*height)

	index := 0
	for i := range tb.canvas {
		tb.canvas[i] = tb.pixels[index : index+width]
		index += width
	}
	return tb
}

// FromStrings creates a new textbox exactly large enough to contain strs
func FromStrings(strs ...string) *Textbox {
	width, height := 0, len(strs)
	for _, s := range strs {
		if len(s) > width {
			width = len(s)
		}
	}
	tb := NewTextbox(width, height)
	for i, s := range strs {
		copy(tb.canvas[i], []rune(s))
	}
	return tb
}

// Runes returns the contents of the textbox as a slice of runes
func (tb *Textbox) Runes() []rune {
	runes := make([]rune, len(tb.pixels))
	copy(runes, tb.pixels)
	return runes
}

// String returns the contents of the textbox as a string
func (tb *Textbox) String() string {
	buf := new(bytes.Buffer)
	for i := range tb.canvas {
		buf.WriteString(string(tb.canvas[i]))
		buf.WriteByte('\n')
	}
	return buf.String()
}

// Bytes returns the contense of the textbox as a []byte
func (tb *Textbox) Bytes() []byte {
	return []byte(string(tb.pixels))
}

// Size returns the textbox dimensions
func (tb *Textbox) Size() (width, height int) {
	return tb.width, tb.height
}

// Cursor returns the current cursor location
func (tb *Textbox) Cursor() (x, y int) {
	return tb.cursor % len(tb.canvas[0]), tb.cursor / len(tb.canvas[0])
}

// SetCursor moves the cursor to the x, y location
func (tb *Textbox) SetCursor(x, y int) error {

	if x < 0 || x > tb.width ||
		y < 0 || y > tb.height {
		return errors.New("location out of bounds")
	}

	tb.cursor = (y * len(tb.canvas[0])) + x
	return nil
}

func (tb *Textbox) Write(b []byte) (int, error) {
	return tb.WriteRunes(bytes.Runes(b))
}

// WriteRunes writes r into the textbox at the cursor location
func (tb *Textbox) WriteRunes(r []rune) (int, error) {
	if tb.cursor >= len(tb.pixels) {
		return 0, errors.New("textbox full")
	}

	n := copy(tb.pixels[tb.cursor:], r)
	tb.cursor += n
	return n, nil
}

// WriteString writes s into the textbox at the cursor location
func (tb *Textbox) WriteString(s string) (int, error) {
	return tb.WriteRunes([]rune(s))
}

// Draw tbox into tb at the offset given by p
func (tb *Textbox) Draw(tbox *Textbox, p Point, transparents []rune) (int, error) {
	count := 0
	// check for complete out if bounds
	if tbox.width+p.X <= 0 || tbox.height+p.Y <= 0 {
		return count, nil
	}

	// set crop point on foreground textbox
	cropPoint := Point{0, 0}
	if p.X < 0 {
		cropPoint.X = int(abs(int64(p.X)))
		p.X = 0
	}
	if p.Y < 0 {
		cropPoint.Y = int(abs(int64(p.Y)))
		p.Y = 0
	}

	// draw foreground textbox (tb) into background textbox (t)
	xIndex := p.X
	if transparents == nil {
		// without transparency
		for i := range tbox.canvas[cropPoint.Y:] {
			yIndex := p.Y + i
			if yIndex >= tb.height {
				break
			}
			count += copy(tb.canvas[yIndex][xIndex:], tbox.canvas[cropPoint.Y+i][cropPoint.X:])
		}
	} else {
		// with transparency
		invisible := make(map[rune]bool)
		for _, r := range transparents {
			invisible[r] = true
		}

		for i := range tbox.canvas[cropPoint.Y:] {
			yIndex := p.Y + i
			if yIndex >= tb.height {
				break
			}

			for j := range tbox.canvas[i][cropPoint.X:] {
				if xIndex+j >= tb.width {
					break
				}

				if !invisible[tbox.canvas[cropPoint.Y+i][cropPoint.X+j]] {
					tb.canvas[yIndex][xIndex+j] = tbox.canvas[cropPoint.Y+i][cropPoint.X+j]
				}
				count++
			}
		}
	}

	// set cursor
	tb.SetCursor(
		min(xIndex+tbox.width-cropPoint.X-1, tb.width-1),
		min(p.Y+tbox.height-cropPoint.Y-1, tb.height-1))
	tb.cursor++

	return count, nil
}

// Tile fills tb by drawing adjacent copies of tbox into it
func (tb *Textbox) Tile(tbox *Textbox, transparents []rune) (int, error) {
	count := 0
	var err error
	p := Point{0, 0}
	for {
		start := p
		n, err := tb.Draw(tbox, p, transparents)
		count += n
		if err != nil || tb.cursor >= len(tb.pixels) {
			break
		}
		p.X, p.Y = tb.Cursor()
		if p.X != 0 {
			p.Y = start.Y
		}
	}
	return count, err
}

// Fill every rune in textbox with 'u'
func (tb *Textbox) Fill(u rune) error {
	if !utf8.ValidRune(u) {
		return errors.New("invalid rune")
	}

	for i := range tb.pixels {
		tb.pixels[i] = u
	}
	return nil
}

// Replace all 'a' runes in textbox with 'b'
func (tb *Textbox) Replace(a, b rune) (int, error) {
	if !utf8.ValidRune(a) {
		return 0, errors.New("invalid rune")
	}

	count := 0
	for i := range tb.pixels {
		if tb.pixels[i] == a {
			tb.pixels[i] = b
			count++
		}
	}
	return count, nil
}

// Crop returns a new textbox with the contense the crop region
func (tb *Textbox) Crop(r Rect) (*Textbox, error) {
	// Check if point is in bounds
	err := tb.SetCursor(r.Origin.X, r.Origin.Y)
	if err != nil {
		return nil, err
	}

	// crop intersection
	cropWidth := min(tb.width-r.Origin.X, r.Width)
	cropHeight := min(tb.height-r.Origin.Y, r.Height)
	crop := NewTextbox(cropWidth, cropHeight)
	for i := range crop.canvas {
		copy(crop.canvas[i], tb.canvas[r.Origin.Y+i][r.Origin.X:cropWidth+r.Origin.X])
	}
	return crop, err
}

// Helper functions
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
