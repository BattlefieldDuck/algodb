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
	Size   int
	Faces  [6][]byte
	buffer []byte // reusable temp buffer
}

// colorChar maps face indices to display letters
var colorChar = []rune{'W', 'R', 'G', 'Y', 'O', 'B'}

// NewCube creates a solved n×n cube.
func NewCube(n int) *Cube {
	c := &Cube{Size: n, buffer: make([]byte, n*n)}
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
	newC := &Cube{Size: n, buffer: make([]byte, n*n)}
	for f := range 6 {
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
func (c *Cube) rotateFaceCWBytes(face int) {
	for r := range c.Size {
		for col := range c.Size {
			c.buffer[col*c.Size+(c.Size-1-r)] = c.Faces[face][r*c.Size+col]
		}
	}
	copy(c.Faces[face], c.buffer)
}

func (c *Cube) rotateFaceCCWBytes(face int) {
	for r := range c.Size {
		for col := range c.Size {
			c.buffer[(c.Size-1-col)*c.Size+r] = c.Faces[face][r*c.Size+col]
		}
	}
	copy(c.Faces[face], c.buffer)
}

// Face turns
func (c *Cube) MoveU(width int) {
	n := c.Size
	faces := c.Faces
	f, r, b, l := faces[Fface], faces[Rface], faces[Bface], faces[Lface]

	// rotate the up face CW, and if it's a full slice also rotate down CCW
	c.rotateFaceCWBytes(Uface)
	if width == n {
		c.rotateFaceCCWBytes(Dface)
	}

	var tmp byte
	for layer := 0; layer < width; layer++ {
		base := layer * n
		for x := 0; x < n; x++ {
			i := base + x
			// cycle F → R → B → L → F
			tmp = f[i]
			f[i] = r[i]
			r[i] = b[i]
			b[i] = l[i]
			l[i] = tmp
		}
	}
}

func (c *Cube) MoveUPrime(width int) {
	n := c.Size
	faces := c.Faces
	f, r, b, l := faces[Fface], faces[Rface], faces[Bface], faces[Lface]

	// rotate the up face CCW, and if it's a full slice also rotate down CW
	c.rotateFaceCCWBytes(Uface)
	if width == n {
		c.rotateFaceCWBytes(Dface)
	}

	var tmp byte
	for layer := 0; layer < width; layer++ {
		base := layer * n
		for x := 0; x < n; x++ {
			i := base + x
			// inverse cycle F ← R ← B ← L ← F
			tmp = f[i]
			f[i] = l[i]
			l[i] = b[i]
			b[i] = r[i]
			r[i] = tmp
		}
	}
}

func (c *Cube) MoveD(width int) {
	n := c.Size
	faces := c.Faces
	f, r, b, l := faces[Fface], faces[Rface], faces[Bface], faces[Lface]

	// rotate the down face CW, and if it's a full slice also rotate up CCW
	c.rotateFaceCWBytes(Dface)
	if width == n {
		c.rotateFaceCCWBytes(Uface)
	}

	var tmp byte
	for layer := 0; layer < width; layer++ {
		row := n - 1 - layer
		base := row * n
		for x := 0; x < n; x++ {
			i := base + x
			// cycle F ← L ← B ← R ← F
			tmp = f[i]
			f[i] = l[i]
			l[i] = b[i]
			b[i] = r[i]
			r[i] = tmp
		}
	}
}

func (c *Cube) MoveDPrime(width int) {
	n := c.Size
	faces := c.Faces
	f, r, b, l := faces[Fface], faces[Rface], faces[Bface], faces[Lface]

	// rotate the down face CCW, and if it's a full slice also rotate up CW
	c.rotateFaceCCWBytes(Dface)
	if width == n {
		c.rotateFaceCWBytes(Uface)
	}

	var tmp byte
	for layer := 0; layer < width; layer++ {
		row := n - 1 - layer
		base := row * n
		for x := 0; x < n; x++ {
			i := base + x
			// inverse cycle F → R → B → L → F
			tmp = f[i]
			f[i] = r[i]
			r[i] = b[i]
			b[i] = l[i]
			l[i] = tmp
		}
	}
}

func (c *Cube) MoveR(width int) {
	n := c.Size
	faces := c.Faces
	u, f, d, b := faces[Uface], faces[Fface], faces[Dface], faces[Bface]

	// rotate the right face CW, and if it's a full slice also rotate left CCW
	c.rotateFaceCWBytes(Rface)
	if width == n {
		c.rotateFaceCCWBytes(Lface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		col := n1 - layer
		for i := 0; i < n; i++ {
			j := i*n + col        // index on U/F/D
			k := (n1-i)*n + layer // index on B

			// cycle U → F → D → B → U
			tmp = u[j]
			u[j] = f[j]
			f[j] = d[j]
			d[j] = b[k]
			b[k] = tmp
		}
	}
}

func (c *Cube) MoveRPrime(width int) {
	n := c.Size
	faces := c.Faces
	u, f, d, b := faces[Uface], faces[Fface], faces[Dface], faces[Bface]

	// rotate the right face CCW, and if it's a full slice also rotate left CW
	c.rotateFaceCCWBytes(Rface)
	if width == n {
		c.rotateFaceCWBytes(Lface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		col := n1 - layer
		for i := 0; i < n; i++ {
			j := i*n + col        // index on U/F/D
			k := (n1-i)*n + layer // index on B

			// inverse cycle U ← F ← D ← B ← U
			tmp = u[j]
			u[j] = b[k]
			b[k] = d[j]
			d[j] = f[j]
			f[j] = tmp
		}
	}
}

func (c *Cube) MoveL(width int) {
	n := c.Size
	faces := c.Faces
	u, f, d, b := faces[Uface], faces[Fface], faces[Dface], faces[Bface]

	// rotate the left face CW, and if full‐cube slice also rotate right CCW
	c.rotateFaceCWBytes(Lface)
	if width == n {
		c.rotateFaceCCWBytes(Rface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		col := layer
		for i := 0; i < n; i++ {
			j := i*n + col               // index on U/F/D
			k := (n1-i)*n + (n1 - layer) // index on B

			// cycle U ← B ← D ← F ← U
			tmp = u[j]
			u[j] = b[k]
			b[k] = d[j]
			d[j] = f[j]
			f[j] = tmp
		}
	}
}

func (c *Cube) MoveLPrime(width int) {
	n := c.Size
	faces := c.Faces
	u, f, d, b := faces[Uface], faces[Fface], faces[Dface], faces[Bface]

	// rotate the left face CCW, and if full‐cube slice also rotate right CW
	c.rotateFaceCCWBytes(Lface)
	if width == n {
		c.rotateFaceCWBytes(Rface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		col := layer
		for i := 0; i < n; i++ {
			j := i*n + col               // index on U/F/D
			k := (n1-i)*n + (n1 - layer) // index on B

			// inverse cycle U → F → D → B → U
			tmp = u[j]
			u[j] = f[j]
			f[j] = d[j]
			d[j] = b[k]
			b[k] = tmp
		}
	}
}

func (c *Cube) MoveF(width int) {
	n := c.Size
	faces := c.Faces
	u, r, d, l := faces[Uface], faces[Rface], faces[Dface], faces[Lface]

	// rotate the front face CW, and if it’s a full slice also rotate back CCW
	c.rotateFaceCWBytes(Fface)
	if width == n {
		c.rotateFaceCCWBytes(Bface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		rowU := (n1 - layer) * n
		rowD := layer * n
		colL, colR := layer, n1-layer

		for i := 0; i < n; i++ {
			uIdx := rowU + i
			rIdx := i*n + colL
			dIdx := rowD + (n1 - i)
			lIdx := (n1-i)*n + colR

			// cycle U ← L ← D ← R ← U
			tmp = u[uIdx]
			u[uIdx] = l[lIdx]
			l[lIdx] = d[dIdx]
			d[dIdx] = r[rIdx]
			r[rIdx] = tmp
		}
	}
}

func (c *Cube) MoveFPrime(width int) {
	n := c.Size
	faces := c.Faces
	u, r, d, l := faces[Uface], faces[Rface], faces[Dface], faces[Lface]

	// rotate the front face CCW, and if it’s a full slice also rotate back CW
	c.rotateFaceCCWBytes(Fface)
	if width == n {
		c.rotateFaceCWBytes(Bface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		rowU := (n1 - layer) * n
		rowD := layer * n
		colL, colR := layer, n1-layer

		for i := 0; i < n; i++ {
			uIdx := rowU + i
			rIdx := i*n + colL
			dIdx := rowD + (n1 - i)
			lIdx := (n1-i)*n + colR

			// inverse cycle U → R → D → L → U
			tmp = u[uIdx]
			u[uIdx] = r[rIdx]
			r[rIdx] = d[dIdx]
			d[dIdx] = l[lIdx]
			l[lIdx] = tmp
		}
	}
}

func (c *Cube) MoveB(width int) {
	n := c.Size
	faces := c.Faces
	u, r, d, l := faces[Uface], faces[Rface], faces[Dface], faces[Lface]

	// rotate the back face CW, and if it's a full slice also rotate front CCW
	c.rotateFaceCWBytes(Bface)
	if width == n {
		c.rotateFaceCCWBytes(Fface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		rowU := layer
		rowD := n1 - layer
		for i := 0; i < n; i++ {
			ui := rowU*n + i
			ri := i*n + (n1 - layer)
			di := rowD*n + (n1 - i)
			li := (n1-i)*n + layer

			// cycle U → R → D → L → U
			tmp = u[ui]
			u[ui] = r[ri]
			r[ri] = d[di]
			d[di] = l[li]
			l[li] = tmp
		}
	}
}

func (c *Cube) MoveBPrime(width int) {
	n := c.Size
	faces := c.Faces
	u, r, d, l := faces[Uface], faces[Rface], faces[Dface], faces[Lface]

	// rotate the back face CCW, and if it's a full slice also rotate front CW
	c.rotateFaceCCWBytes(Bface)
	if width == n {
		c.rotateFaceCWBytes(Fface)
	}

	n1 := n - 1
	var tmp byte
	for layer := 0; layer < width; layer++ {
		rowU := layer
		rowD := n1 - layer
		for i := 0; i < n; i++ {
			ui := rowU*n + i
			ri := i*n + (n1 - layer)
			di := rowD*n + (n1 - i)
			li := (n1-i)*n + layer

			// inverse cycle U ← R ← D ← L ← U
			tmp = u[ui]
			u[ui] = l[li]
			l[li] = d[di]
			d[di] = r[ri]
			r[ri] = tmp
		}
	}
}

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

const (
	flagWide  = 1 << 3
	flagRot   = 1 << 4
	flagSlice = 1 << 5
	faceMask  = (1 << 3) - 1 // lower 3 bits for face id
)

var moveMap = map[string]uint8{
	// normals
	"U": Uface,
	"D": Dface,
	"R": Rface,
	"L": Lface,
	"F": Fface,
	"B": Bface,

	// wides
	"Uw": Uface | flagWide,
	"u":  Uface | flagWide,
	"Dw": Dface | flagWide,
	"d":  Dface | flagWide,
	"Rw": Rface | flagWide,
	"r":  Rface | flagWide,
	"Lw": Lface | flagWide,
	"l":  Lface | flagWide,
	"Fw": Fface | flagWide,
	"f":  Fface | flagWide,
	"Bw": Bface | flagWide,
	"b":  Bface | flagWide,

	// rotations
	"x": Rface | flagRot,
	"y": Uface | flagRot,
	"z": Fface | flagRot,

	// slices
	"M": Lface | flagSlice,
	"E": Dface | flagSlice,
	"S": Fface | flagSlice,
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
	var face int
	isSlice := false
	if m, ok := moveMap[notation]; ok {
		face = int(m & faceMask)

		switch {
		case m&flagWide != 0:
			// set width=2 if isWide && width<=1
			if width < 2 {
				width = 2
			}

		case m&flagRot != 0:
			width = c.Size

		case m&flagSlice != 0:
			isSlice = true
		}
	}

	// apply move times times
	for range times {
		switch face {
		case Uface:
			if isPrime {
				c.MoveUPrime(width)
			} else {
				c.MoveU(width)
			}
		case Dface:
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
		case Rface:
			if isPrime {
				c.MoveRPrime(width)
			} else {
				c.MoveR(width)
			}
		case Lface:
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
		case Fface:
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
		case Bface:
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
