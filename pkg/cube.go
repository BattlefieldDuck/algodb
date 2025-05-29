package pkg

import (
	"fmt"
	"strings"
)

const (
	Uface = iota
	Rface
	Fface
	Dface
	Lface
	Bface
)

// Cube uses a fixed array of byte-slices for faces: 0=U,1=R,2=F,3=D,4=L,5=B
// Each byte stores the face index of that sticker.
type Cube struct {
	Size  int
	Faces [6][]byte
}

// colorChar maps face indices to display letters
var colorChar = []rune{'W', 'R', 'G', 'Y', 'O', 'B'}

// NewCube creates a solved n×n cube.
func NewCube(n int) *Cube {
	c := &Cube{Size: n}
	for f := 0; f < 6; f++ {
		c.Faces[f] = make([]byte, n*n)
		for i := range c.Faces[f] {
			c.Faces[f][i] = byte(f)
		}
	}
	return c
}

// Copy returns a deep copy of this cube.
func (c *Cube) Copy() *Cube {
	n := c.Size
	newC := &Cube{Size: n}
	for f := 0; f < 6; f++ {
		newC.Faces[f] = make([]byte, n*n)
		copy(newC.Faces[f], c.Faces[f])
	}
	return newC
}

// Display prints the cube in ASCII using face letters with correct indentation.
func (c *Cube) Display() {
	n := c.Size
	indent := strings.Repeat(" ", n*2)

	// print a single face row with optional indent
	printRow := func(f int, r int, withIndent bool) {
		if withIndent {
			fmt.Print(indent)
		}
		for x := 0; x < n; x++ {
			fmt.Printf("%c ", colorChar[c.Faces[f][r*n+x]])
		}
		fmt.Println()
	}

	// U face
	for r := 0; r < n; r++ {
		printRow(Uface, r, true)
	}

	// middle L-F-R-B
	for r := 0; r < n; r++ {
		for _, f := range []int{Lface, Fface, Rface, Bface} {
			for x := 0; x < n; x++ {
				fmt.Printf("%c ", colorChar[c.Faces[f][r*n+x]])
			}
		}
		fmt.Println()
	}

	// D face
	for r := 0; r < n; r++ {
		printRow(Dface, r, true)
	}
}

// ANSI background codes for your six face‐colors.
// We use two spaces “  ” as our “sticker” so it’s a square block.
var ansiBg = []string{
	"\x1b[47m",       // white
	"\x1b[41m",       // red
	"\x1b[42m",       // green
	"\x1b[43m",       // yellow
	"\x1b[48;5;208m", // orange (256‐color)
	"\x1b[44m",       // blue
}

const (
	sticker = "⠀⠀" // two spaces for a square sticker
	reset   = "\x1b[0m"
)

// colorChar maps face indices to display letters
// DisplayColor prints the cube with ANSI-colored stickers.
func (c *Cube) DisplayColor() {
	n := c.Size
	indent := strings.Repeat(" ", n*2)
	// helper to paint a row of stickers
	paintRow := func(cells []byte) {
		for _, idx := range cells {
			bg := ansiBg[int(idx)]
			fmt.Print(bg + sticker + reset)
		}
		fmt.Println()
	}

	// U face
	for r := 0; r < n; r++ {
		row := c.Faces[Uface][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}

	// middle L-F-R-B
	for r := 0; r < n; r++ {
		var row []byte
		for _, f := range []int{Lface, Fface, Rface, Bface} {
			row = append(row, c.Faces[f][r*n:r*n+n]...)
		}
		paintRow(row)
	}

	// D face
	for r := 0; r < n; r++ {
		row := c.Faces[Dface][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}
}

// rotateFaceCWBytes rotates an n×n byte-slice CW in-place.
func rotateFaceCWBytes(face []byte, n int) {
	tmp := make([]byte, len(face))
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			tmp[c*n+(n-1-r)] = face[r*n+c]
		}
	}
	copy(face, tmp)
}

// rotateFaceCCWBytes rotates CCW via three CWs.
func rotateFaceCCWBytes(face []byte, n int) {
	for i := 0; i < 3; i++ {
		rotateFaceCWBytes(face, n)
	}
}

func getRowBytes(face []byte, r, n int) []byte {
	out := make([]byte, n)
	copy(out, face[r*n:(r+1)*n])
	return out
}

func setRowBytes(face []byte, r, n int, data []byte) {
	for i := 0; i < n; i++ {
		face[r*n+i] = data[i]
	}
}

func getColBytes(face []byte, col, n int) []byte {
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = face[i*n+col]
	}
	return out
}

func setColBytes(face []byte, col, n int, data []byte) {
	for i := 0; i < n; i++ {
		face[i*n+col] = data[i]
	}
}

func reverseBytes(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Face turns
func (c *Cube) MoveU(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Uface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Dface], n)
	}
	for layer := 0; layer < width; layer++ {
		f := getRowBytes(c.Faces[Fface], layer, n)
		r := getRowBytes(c.Faces[Rface], layer, n)
		b := getRowBytes(c.Faces[Bface], layer, n)
		l := getRowBytes(c.Faces[Lface], layer, n)
		setRowBytes(c.Faces[Fface], layer, n, r)
		setRowBytes(c.Faces[Rface], layer, n, b)
		setRowBytes(c.Faces[Bface], layer, n, l)
		setRowBytes(c.Faces[Lface], layer, n, f)
	}
}

func (c *Cube) MoveUPrime(width int) { c.MoveU(width); c.MoveU(width); c.MoveU(width) }

func (c *Cube) MoveD(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Dface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Uface], n)
	}
	for layer := 0; layer < width; layer++ {
		row := n - 1 - layer
		f := getRowBytes(c.Faces[Fface], row, n)
		r := getRowBytes(c.Faces[Rface], row, n)
		b := getRowBytes(c.Faces[Bface], row, n)
		l := getRowBytes(c.Faces[Lface], row, n)
		setRowBytes(c.Faces[Fface], row, n, l)
		setRowBytes(c.Faces[Lface], row, n, b)
		setRowBytes(c.Faces[Bface], row, n, r)
		setRowBytes(c.Faces[Rface], row, n, f)
	}
}

func (c *Cube) MoveDPrime(width int) { c.MoveD(width); c.MoveD(width); c.MoveD(width) }

func (c *Cube) MoveR(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Rface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Lface], n)
	}
	for layer := 0; layer < width; layer++ {
		colU := getColBytes(c.Faces[Uface], n-1-layer, n)
		colF := getColBytes(c.Faces[Fface], n-1-layer, n)
		colD := getColBytes(c.Faces[Dface], n-1-layer, n)
		colB := getColBytes(c.Faces[Bface], layer, n)
		setColBytes(c.Faces[Uface], n-1-layer, n, colF)
		setColBytes(c.Faces[Fface], n-1-layer, n, colD)
		reverseBytes(colB)
		setColBytes(c.Faces[Dface], n-1-layer, n, colB)
		reverseBytes(colU)
		setColBytes(c.Faces[Bface], layer, n, colU)
	}
}

func (c *Cube) MoveRPrime(width int) { c.MoveR(width); c.MoveR(width); c.MoveR(width) }

func (c *Cube) MoveL(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Lface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Rface], n)
	}
	for layer := 0; layer < width; layer++ {
		colU := getColBytes(c.Faces[Uface], layer, n)
		colF := getColBytes(c.Faces[Fface], layer, n)
		colD := getColBytes(c.Faces[Dface], layer, n)
		colB := getColBytes(c.Faces[Bface], n-1-layer, n)
		reverseBytes(colB)
		setColBytes(c.Faces[Uface], layer, n, colB)
		reverseBytes(colD)
		setColBytes(c.Faces[Bface], n-1-layer, n, colD)
		setColBytes(c.Faces[Dface], layer, n, colF)
		setColBytes(c.Faces[Fface], layer, n, colU)
	}
}

func (c *Cube) MoveLPrime(width int) { c.MoveL(width); c.MoveL(width); c.MoveL(width) }

func (c *Cube) MoveF(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Fface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Bface], n)
	}
	for layer := 0; layer < width; layer++ {
		rowU := getRowBytes(c.Faces[Uface], n-1-layer, n)
		colR := getColBytes(c.Faces[Rface], layer, n)
		rowD := getRowBytes(c.Faces[Dface], layer, n)
		colL := getColBytes(c.Faces[Lface], n-1-layer, n)
		setColBytes(c.Faces[Rface], layer, n, rowU)
		reverseBytes(colR)
		setRowBytes(c.Faces[Dface], layer, n, colR)
		setColBytes(c.Faces[Lface], n-1-layer, n, rowD)
		reverseBytes(colL)
		setRowBytes(c.Faces[Uface], n-1-layer, n, colL)
	}
}

func (c *Cube) MoveFPrime(width int) { c.MoveF(width); c.MoveF(width); c.MoveF(width) }

func (c *Cube) MoveB(width int) {
	n := c.Size
	rotateFaceCWBytes(c.Faces[Bface], n)
	if width == n {
		rotateFaceCCWBytes(c.Faces[Fface], n)
	}
	for layer := 0; layer < width; layer++ {
		rowU := getRowBytes(c.Faces[Uface], layer, n)
		colL := getColBytes(c.Faces[Lface], layer, n)
		rowD := getRowBytes(c.Faces[Dface], n-1-layer, n)
		colR := getColBytes(c.Faces[Rface], n-1-layer, n)
		reverseBytes(rowU)
		setColBytes(c.Faces[Lface], layer, n, rowU)
		setRowBytes(c.Faces[Dface], n-1-layer, n, colL)
		reverseBytes(rowD)
		setColBytes(c.Faces[Rface], n-1-layer, n, rowD)
		setRowBytes(c.Faces[Uface], layer, n, colR)
	}
}

func (c *Cube) MoveBPrime(width int) { c.MoveB(width); c.MoveB(width); c.MoveB(width) }

// Moves applies a sequence of moves, stripping parentheses.
func (c *Cube) Moves(seqs ...string) error {
	seq := strings.Join(seqs, " ")
	seq = strings.NewReplacer("(", "", ")", "").Replace(seq)
	for _, m := range strings.Fields(seq) {
		if err := c.Move(m); err != nil {
			return err
		}
	}
	return nil
}

// mapping for wide vs normal moves
var mapping = map[string]struct {
	face       int
	isWide     bool
	isRotation bool
	isSlice    bool
}{
	"U": {Uface, false, false, false}, "D": {Dface, false, false, false}, "R": {Rface, false, false, false},
	"L": {Lface, false, false, false}, "F": {Fface, false, false, false}, "B": {Bface, false, false, false},
	"Uw": {Uface, true, false, false}, "u": {Uface, true, false, false},
	"Dw": {Dface, true, false, false}, "d": {Dface, true, false, false},
	"Rw": {Rface, true, false, false}, "r": {Rface, true, false, false},
	"Lw": {Lface, true, false, false}, "l": {Lface, true, false, false},
	"Fw": {Fface, true, false, false}, "f": {Fface, true, false, false},
	"Bw": {Bface, true, false, false}, "b": {Bface, true, false, false},
	"x": {Rface, false, true, false}, "y": {Uface, false, true, false}, "z": {Fface, false, true, false},
	"M": {Lface, false, false, true}, "E": {Dface, false, false, true}, "S": {Fface, false, false, true},
}

// Move parses a notation (e.g. "R2'", "u"), then calls the appropriate face-turn.
func (c *Cube) Move(notation string) error {
	width := 1
	// leading digit as width
	if len(notation) > 0 && notation[0] >= '0' && notation[0] <= '9' {
		width = int(notation[0] - '0')
		notation = notation[1:]
	}
	// prime suffix
	isPrime := strings.HasSuffix(notation, "'")
	if isPrime {
		notation = notation[:len(notation)-1]
	}
	// double suffix
	times := 1
	if strings.HasSuffix(notation, "2") {
		times = 2
		notation = notation[:len(notation)-1]
	}

	// mapping
	isSlice := false
	if m, ok := mapping[notation]; ok {
		notation = ""
		// use face index from m.face
		// set width=2 if isWide && width<=1
		if m.isWide && width < 2 {
			width = 2
		}

		if m.isRotation {
			width = c.Size
		}

		isSlice = m.isSlice

		notation = fmt.Sprint(m.face) // replaced below
	}
	// apply move times times
	for i := 0; i < times; i++ {
		switch notation {
		case fmt.Sprint(Uface):
			if isPrime {
				c.MoveUPrime(width)
			} else {
				c.MoveU(width)
			}
		case fmt.Sprint(Dface):
			if isSlice {
				// E slice move
				if isPrime {
					c.MoveD(1)
					c.MoveDPrime(c.Size - 1)
				} else {
					c.MoveDPrime(1)
					c.MoveD(c.Size - 1)
				}
			} else if isPrime {
				c.MoveDPrime(width)
			} else {
				c.MoveD(width)
			}
		case fmt.Sprint(Rface):
			if isPrime {
				c.MoveRPrime(width)
			} else {
				c.MoveR(width)
			}
		case fmt.Sprint(Lface):
			if isSlice {
				// M slice move
				if isPrime {
					c.MoveL(1)
					c.MoveLPrime(c.Size - 1)
				} else {
					c.MoveLPrime(1)
					c.MoveL(c.Size - 1)
				}
			} else if isPrime {
				c.MoveLPrime(width)
			} else {
				c.MoveL(width)
			}
		case fmt.Sprint(Fface):
			if isSlice {
				// S slice move
				if isPrime {
					c.MoveF(1)
					c.MoveFPrime(c.Size - 1)
				} else {
					c.MoveFPrime(1)
					c.MoveF(c.Size - 1)
				}
			} else if isPrime {
				c.MoveFPrime(width)
			} else {
				c.MoveF(width)
			}
		case fmt.Sprint(Bface):
			if isPrime {
				c.MoveBPrime(width)
			} else {
				c.MoveB(width)
			}
		default:
			return fmt.Errorf("invalid move: %s", notation)
		}
	}
	return nil
}

// IsSolved returns true if every face of the cube is uniform (all stickers match the face index).
func (c *Cube) IsSolved() bool {
	for f := 0; f < 6; f++ {
		for _, v := range c.Faces[f] {
			if v != byte(f) {
				return false
			}
		}
	}
	return true
}
