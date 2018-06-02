package textbox

import (
	"errors"
	"fmt"
	"image"
	"os"
	"reflect"
	"unicode/utf8"

	"golang.org/x/sys/unix"
)

type Rect struct {
	P      image.Point
	Width  int
	Height int
}

type Textbox struct {
	canvas [][]rune
	pixels []rune
	width  int
	height int
	cursor image.Point
}

func NewTextbox(width, height int) *Textbox {
	t := new(Textbox)
	t.width = width
	t.height = height
	t.canvas = make([][]rune, height)
	t.pixels = make([]rune, width*height)
	for i := range t.canvas {
		start_index := i * width
		t.canvas[i] = t.pixels[start_index : start_index+width]
	}
	t.cursor.X, t.cursor.Y = 0, 0
	return t
}

func BoxOfStrings(strs ...string) *Textbox {
	width, height := 0, len(strs)
	for _, s := range strs {
		width = maxInt(len(s), width)
	}
	t := NewTextbox(width, height)
	for i, s := range strs {
		copy(t.canvas[i], []rune(s))
	}
	return t
}

func (t *Textbox) Runes() []rune {
	runes := make([]rune, len(t.pixels))
	copy(runes, t.pixels)
	return runes
}

func (t *Textbox) String() string {
	b := make([]byte, t.width*t.height+t.height)
	for i := range t.canvas {
		copy(b[t.width*i:], string(t.canvas[i]))
		b[t.width*(i+1)-1] = '\n'
	}
	return string(b)
}

func (t *Textbox) Bytes() []byte {
	return []byte(string(t.pixels))
}

func (t *Textbox) Size() (width, height int) {
	return t.width, t.height
}

func (t *Textbox) Runway() int {
	return cap(t.pixels) - (t.width * t.cursor.Y) - t.cursor.X
}

func (t *Textbox) Cursor() image.Point {
	return t.cursor
}

func (t *Textbox) SetCursor(p image.Point) error {
	// Test bounds
	if p.X < 0 || p.X > t.width {
		return errors.New("image.Point out of bounds.")
	}
	if p.Y < 0 || p.Y > t.height {
		return errors.New("image.Point out of bounds.")
	}

	t.cursor = p
	return nil
}

func (t *Textbox) incrementCursor() {
	if t.cursor.X < t.width-1 {
		t.cursor.X++
	} else {
		t.cursor.X = 0
		t.cursor.Y++
	}
}

func (t *Textbox) decrementCursor() {
	if t.cursor.X > 0 {
		t.cursor.X--
	} else {
		t.cursor.Y--
		t.cursor.X = t.width - 1
	}
}

func (t *Textbox) Write(i interface{}) (int, error) {
	var runes []rune

	switch i := i.(type) {
	case *Textbox:
		return t.Draw(i, t.cursor)

	case []byte:
		runes = []rune(string(i))

	case []rune:
		runes = i

	case string:
		runes = []rune(i)

	case fmt.Stringer:
		runes = []rune(i.String())

	default:
		return 0, fmt.Errorf("Textbox can not write type %s. Convert to string", reflect.TypeOf(i))
	}

	return t.writeRunes(runes)
}

func (t *Textbox) writeRunes(r []rune) (int, error) {
	count := 0
	freeSpace := t.Runway()
	if freeSpace == 0 {
		return count, errors.New("Textbox full.")
	}

	for ; count < minInt(freeSpace, len(r)); count++ {
		if r[count] != 0 {
			t.canvas[t.cursor.Y][t.cursor.X] = r[count]
		}
		t.incrementCursor()
	}
	return count, nil
}

func (t *Textbox) WriteWords(w ...string) (int, error) {
	count := 0
	freeSpace := t.Runway()
	if freeSpace == 0 {
		return count, errors.New("Textbox full.")
	}

	return count, nil
}

func (t *Textbox) Draw(tb *Textbox, p image.Point) (int, error) {
	count := 0
	// check for complete out if bounds
	if tb.width+p.X <= 0 || tb.height+p.Y <= 0 {
		return count, nil
	}

	// set crop point on foreground textbox
	cropPoint := image.Point{0, 0}
	if p.X < 0 {
		cropPoint.X = absInt(p.X)
		p.X = 0
	}
	if p.Y < 0 {
		cropPoint.Y = absInt(p.Y)
		p.Y = 0
	}

	// draw foreground textbox (tb) into background textbox (t)
	xIndex := p.X
	for i := range tb.canvas[cropPoint.Y:] {
		yIndex := p.Y + i
		if yIndex >= t.height {
			break
		}

		count += copy(t.canvas[yIndex][xIndex:], tb.canvas[cropPoint.Y+i][cropPoint.X:])
	}

	// set cursor
	t.SetCursor(image.Point{minInt(xIndex+tb.width-cropPoint.X-1, t.width-1), minInt(p.Y+tb.height-cropPoint.Y-1, t.height-1)})
	t.incrementCursor()

	return count, nil
}

func (t *Textbox) DrawWithTransparency(tb *Textbox, p image.Point, transparentChar rune) (int, error) {
	count := 0
	// check for complete out if bounds
	if tb.width+p.X <= 0 || tb.height+p.Y <= 0 {
		return count, nil
	}

	// set crop point on foreground textbox
	cropPoint := image.Point{0, 0}
	if p.X < 0 {
		cropPoint.X = absInt(p.X)
		p.X = 0
	}
	if p.Y < 0 {
		cropPoint.Y = absInt(p.Y)
		p.Y = 0
	}

	// draw foreground textbox (tb) into background textbox (t)
	xIndex := p.X
	for i := range tb.canvas[cropPoint.Y:] {
		yIndex := p.Y + i
		if yIndex >= t.height {
			break
		}

		for j := range tb.canvas[i][cropPoint.X:] {
			if xIndex+j >= t.width {
				break
			}

			if tb.canvas[cropPoint.Y+i][cropPoint.X+j] != transparentChar {
				t.canvas[yIndex][xIndex+j] = tb.canvas[cropPoint.Y+i][cropPoint.X+j]
			}
			count++
		}
	}

	// set cursor
	t.SetCursor(image.Point{minInt(xIndex+tb.width-cropPoint.X-1, t.width-1), minInt(p.Y+tb.height-cropPoint.Y-1, t.height-1)})
	t.incrementCursor()

	return count, nil
}

func (t *Textbox) Tile(tb *Textbox) (int, error) {
	count := 0
	var err error
	p := image.Point{0, 0}
	for {
		start := p
		n, err := t.Draw(tb, p)
		count += n
		if err != nil || t.Runway() == 0 {
			break
		}
		p = t.cursor
		if p.X != 0 {
			p.Y = start.Y
		}
	}
	return count, err
}

func (t *Textbox) Fill(u rune) error {
	if !utf8.ValidRune(u) {
		return errors.New("invalid rune.")
	}

	for i := range t.pixels {
		t.pixels[i] = u
	}

	return nil
}

func (t *Textbox) Replace(a, b rune) (int, error) {
	count := 0
	if !utf8.ValidRune(a) {
		return count, errors.New("invalid rune.")
	}

	for i := range t.pixels {
		if t.pixels[i] == a {
			t.pixels[i] = b
			count++
		}
	}
	return count, nil
}

func (t *Textbox) Crop(r Rect) (*Textbox, error) {
	// Check if point is in bounds
	err := t.SetCursor(r.P)
	if err != nil {
		return nil, err
	}
	// crop intersection
	cropWidth := minInt(t.width-r.P.X, r.Width)
	cropHeight := minInt(t.height-r.P.Y, r.Height)
	crop := NewTextbox(cropWidth, cropHeight)
	for i := range crop.canvas {
		copy(crop.canvas[i], t.canvas[r.P.Y+i][r.P.X:cropWidth+r.P.X])
	}
	return crop, err
}

func (t *Textbox) Copy() *Textbox {
	tb := NewTextbox(t.width, t.height)
	copy(tb.pixels, t.pixels)
	return tb
}

func TerminalSize() (width, height int, err error) {
	size, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(size.Col), int(size.Row), nil
}

func boundPoint(p image.Point, width, height int) (np, diff image.Point) {
	x, xd := bound(p.X, 0, width-1)
	y, yd := bound(p.Y, 0, height-1)
	np = image.Point{x, y}
	diff = image.Point{xd, yd}
	return np, diff
}

func bound(i, start, end int) (n, diff int) {
	n = i
	diff = 0
	if i < start {
		diff = start - i
		n = start
	} else if n > end {
		diff = end - i
		n = end
	}
	return n, diff
}

func absInt(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func maxInt(ints ...int) int {
	max := ints[0]
	for _, v := range ints {
		if v > max {
			max = v
		}
	}
	return max
}

func minInt(ints ...int) int {
	min := ints[0]
	for _, v := range ints {
		if v < min {
			min = v
		}
	}
	return min
}
