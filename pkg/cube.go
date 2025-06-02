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
// DisplayColorANSI prints the cube with ANSI-colored stickers.
func (c *Cube) DisplayColorANSI() {
	n := c.Size
	indent := strings.Repeat("⠀", n*2)
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

func (c *Cube) DisplayColorANSIUFace() {
	tmp := c.Copy()
	tmp.MoveRPrime(c.Size)

	n := tmp.Size
	indent := strings.Repeat("⠀", 4)
	// helper to paint a row of stickers
	paintRow := func(cells []byte) {
		for _, idx := range cells {
			bg := ansiBg[int(idx)]
			fmt.Print(bg + sticker + reset)
		}
		fmt.Println()
	}

	// U face
	for r := n - 1; r < n; r++ {
		row := tmp.Faces[Uface][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}

	// middle L-F-R-B
	for r := 0; r < n; r++ {
		var row []byte
		row = append(row, tmp.Faces[Lface][r*n+n-1:r*n+n]...)
		row = append(row, tmp.Faces[Fface][r*n:r*n+n]...)
		row = append(row, tmp.Faces[Rface][r*n:r*n+1]...)

		fmt.Print(sticker)
		paintRow(row)
	}

	// D face
	for r := 0; r < 1; r++ {
		row := tmp.Faces[Dface][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}
}

// colorEmoji maps face indices to colored square emojis
var colorEmoji = []string{
	"⬜", // Uface: white
	"🟥", // Rface: red
	"🟩", // Fface: green
	"🟨", // Dface: yellow
	"🟧", // Lface: orange
	"🟦", // Bface: blue
}

// DisplayColorUnicode prints the cube net using Unicode colored squares.
func (c *Cube) DisplayColorUnicode() {
	n := c.Size
	indent := strings.Repeat("  ", n) // two spaces per sticker width

	paintRow := func(cells []byte) {
		for _, idx := range cells {
			fmt.Print(colorEmoji[int(idx)])
		}
		fmt.Println()
	}

	// U face (centered above)
	for r := 0; r < n; r++ {
		fmt.Print(indent)
		paintRow(c.Faces[Uface][r*n : r*n+n])
	}

	// middle L-F-R-B
	for r := 0; r < n; r++ {
		var row []byte
		for _, f := range []int{Lface, Fface, Rface, Bface} {
			row = append(row, c.Faces[f][r*n:r*n+n]...)
		}
		paintRow(row)
	}

	// D face (centered below)
	for r := 0; r < n; r++ {
		fmt.Print(indent)
		paintRow(c.Faces[Dface][r*n : r*n+n])
	}
}

func (c *Cube) DisplayColorUnicodeUFace() {
	tmp := c.Copy()
	tmp.MoveRPrime(c.Size)

	n := tmp.Size
	indent := strings.Repeat("  ", 2) // two spaces per sticker width

	paintRow := func(cells []byte) {
		for _, idx := range cells {
			fmt.Print(colorEmoji[int(idx)])
		}
		fmt.Println()
	}

	// U face (centered above)
	for r := n - 1; r < n; r++ {
		fmt.Print(indent)
		paintRow(tmp.Faces[Uface][r*n : r*n+n])
	}

	// middle L-F-R
	for r := 0; r < n; r++ {
		var row []byte
		row = append(row, tmp.Faces[Lface][r*n+n-1:r*n+n]...)
		row = append(row, tmp.Faces[Fface][r*n:r*n+n]...)
		row = append(row, tmp.Faces[Rface][r*n:r*n+1]...)

		fmt.Print("  ")
		paintRow(row)
	}

	// D face (centered below)
	for r := 0; r < 1; r++ {
		fmt.Print(indent)
		paintRow(tmp.Faces[Dface][r*n : r*n+n])
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

func (c *Cube) MoveM() {
	c.MoveLPrime(1)
	c.MoveL(c.Size - 1)
}

func (c *Cube) MoveMPrime() {
	c.MoveL(1)
	c.MoveLPrime(c.Size - 1)
}

func (c *Cube) MoveE() {
	c.MoveDPrime(1)
	c.MoveD(c.Size - 1)
}

func (c *Cube) MoveEPrime() {
	c.MoveD(1)
	c.MoveDPrime(c.Size - 1)
}

func (c *Cube) MoveS() {
	c.MoveFPrime(1)
	c.MoveF(c.Size - 1)
}

func (c *Cube) MoveSPrime() {
	c.MoveF(1)
	c.MoveFPrime(c.Size - 1)
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

func (c *Cube) parseNotation(notation string) (face, count, width int, isPrime, isSlice bool, err error) {
	width = 1

	// leading digit as width
	if len(notation) > 0 && notation[0] >= '0' && notation[0] <= '9' {
		width = int(notation[0] - '0')
		notation = notation[1:]
	}
	// prime suffix
	isPrime = false
	if n := len(notation); n > 0 && notation[n-1] == '\'' {
		isPrime = true
		notation = notation[:n-1]
	}

	// double suffix
	count = 1
	if n := len(notation); n > 0 && notation[n-1] == '2' {
		count = 2
		notation = notation[:n-1]
	}

	// mapping
	isSlice = false
	if m, ok := moveMap[notation]; ok {
		face = int(m & faceMask)

		switch {
		case m&flagWide != 0:
			// if it’s a “wide” move and user didn’t specify width ≥ 2 already,
			// force width to 2
			if width < 2 {
				width = 2
			}

		case m&flagRot != 0:
			// a whole‐cube rotation (x, y, z) → width = cube size
			width = c.Size

		case m&flagSlice != 0:
			// slice move (E, M, or S)
			isSlice = true
		}

		return face, count, width, isPrime, isSlice, nil
	}

	// if we fall through, notation wasn’t in moveMap→ return zero values + an error
	return 0, 0, 0, false, false,
		fmt.Errorf("parseNotation: invalid move “%s”", notation)
}

// Move parses a notation (e.g. "R2'", "u"), then calls the appropriate face-turn.
func (c *Cube) Move(notation string) error {
	face, count, width, isPrime, isSlice, err := c.parseNotation(notation)
	if err != nil {
		return err
	}

	// FaceTurnFunc now returns func() error
	faceTurnFunc(face, count, width, isPrime, isSlice)(c)

	return nil
}

func (c *Cube) CreateOp(notation string) (int, func(c *Cube), func(c *Cube)) {
	face, count, width, isPrime, isSlice, _ := c.parseNotation(notation)
	apply := faceTurnFunc(face, count, width, isPrime, isSlice)
	undo := faceTurnFunc(face, count, width, !isPrime, isSlice)
	return face, apply, undo
}

// instead of executing immediately, return a closure that will perform the moves when called
func faceTurnFunc(face, count, width int, isPrime, isSlice bool) func(c *Cube) {
	switch face {
	case Uface:
		if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveUPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveU(width)
				}
			}
		}

	case Dface:
		if isSlice {
			// E‐slice move
			if isPrime {
				return func(c *Cube) {
					for range count {
						c.MoveEPrime()
					}
				}
			} else {
				return func(c *Cube) {
					for range count {
						c.MoveE()
					}
				}
			}
		} else if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveDPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveD(width)
				}
			}
		}

	case Rface:
		if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveRPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveR(width)
				}
			}
		}

	case Lface:
		if isSlice {
			// M‐slice move
			if isPrime {
				return func(c *Cube) {
					for range count {
						c.MoveMPrime()
					}
				}
			} else {
				return func(c *Cube) {
					for range count {
						c.MoveM()
					}
				}
			}
		} else if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveLPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveL(width)
				}
			}
		}

	case Fface:
		if isSlice {
			// S‐slice move
			if isPrime {
				return func(c *Cube) {
					for range count {
						c.MoveSPrime()
					}
				}
			} else {
				return func(c *Cube) {
					for range count {
						c.MoveS()
					}
				}
			}
		} else if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveFPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveF(width)
				}
			}
		}

	case Bface:
		if isPrime {
			return func(c *Cube) {
				for range count {
					c.MoveBPrime(width)
				}
			}
		} else {
			return func(c *Cube) {
				for range count {
					c.MoveB(width)
				}
			}
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
