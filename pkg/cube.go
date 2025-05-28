package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

type Cube struct {
	Size  int
	Faces map[string][]string
}

func NewCube(size int) *Cube {
	colors := map[string]string{"U": "W", "D": "Y", "F": "G", "B": "B", "L": "O", "R": "R"}
	faces := make(map[string][]string, len(colors))
	for face, col := range colors {
		stickers := make([]string, size*size)
		for i := range stickers {
			stickers[i] = col
		}
		faces[face] = stickers
	}
	return &Cube{Size: size, Faces: faces}
}

func (c *Cube) Display() {
	n := c.Size
	indent := strings.Repeat(" ", n*2)

	// U face
	for r := 0; r < n; r++ {
		row := c.Faces["U"][r*n : r*n+n]
		fmt.Println(indent + strings.Join(row, " "))
	}

	// Middle L-F-R-B
	for r := 0; r < n; r++ {
		var row []string
		for _, face := range []string{"L", "F", "R", "B"} {
			row = append(row, c.Faces[face][r*n:r*n+n]...)
		}
		fmt.Println(strings.Join(row, " "))
	}

	// D face
	for r := 0; r < n; r++ {
		row := c.Faces["D"][r*n : r*n+n]
		fmt.Println(indent + strings.Join(row, " "))
	}
}

// ANSI background codes for your six face‐colors.
// We use two spaces “  ” as our “sticker” so it’s a square block.
var ansiBg = map[string]string{
	"W": "\x1b[47m",       // white
	"Y": "\x1b[43m",       // yellow
	"R": "\x1b[41m",       // red
	"O": "\x1b[48;5;208m", // orange (256‐color)
	"G": "\x1b[42m",       // green
	"B": "\x1b[44m",       // blue
}

const (
	sticker = "  "
	reset   = "\x1b[0m"
)

func (c *Cube) DisplayColor() {
	n := c.Size
	indent := strings.Repeat(" ", n*2)

	// helper to paint a row of stickers
	paintRow := func(cells []string) {
		for _, s := range cells {
			bg, ok := ansiBg[s]
			if !ok {
				bg = "\x1b[40m" // fallback: black
			}
			fmt.Print(bg + sticker + reset)
		}
		fmt.Println()
	}

	// U face
	for r := 0; r < n; r++ {
		row := c.Faces["U"][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}

	// middle L-F-R-B
	for r := 0; r < n; r++ {
		var row []string
		for _, face := range []string{"L", "F", "R", "B"} {
			row = append(row, c.Faces[face][r*n:r*n+n]...)
		}
		paintRow(row)
	}

	// D face
	for r := 0; r < n; r++ {
		row := c.Faces["D"][r*n : r*n+n]
		fmt.Print(indent)
		paintRow(row)
	}
}

// rotate an N×N face CW
func rotateFaceCW(face []string, n int) []string {
	out := make([]string, len(face))
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			out[c*n+(n-1-r)] = face[r*n+c]
		}
	}
	return out
}

// rotate an N×N face CCW (one CCW = three CWs)
func rotateFaceCCW(face []string, n int) []string {
	f := face
	for i := 0; i < 3; i++ {
		f = rotateFaceCW(f, n)
	}
	return f
}

// helper: extract row ‘r’ from an n×n face
func getRow(face []string, r, n int) []string {
	out := make([]string, n)
	copy(out, face[r*n:(r+1)*n])
	return out
}

// helper: write data into row ‘r’ of an n×n face
func setRow(face []string, r, n int, data []string) {
	for i := 0; i < n; i++ {
		face[r*n+i] = data[i]
	}
}

// helper: extract column col from an n×n face
func getCol(face []string, col, n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = face[i*n+col]
	}
	return out
}

// helper: write data into column col of an n×n face
func setCol(face []string, col, n int, data []string) {
	for i := 0; i < n; i++ {
		face[i*n+col] = data[i]
	}
}

// helper: in-place reverse a slice
func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// wide R move: rotate the outer face CW, then cycle `width` right‐most layers
func (c *Cube) MoveR(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid R width %d for size %d", width, n))
	}

	// always rotate the outer R face CW
	c.Faces["R"] = rotateFaceCW(c.Faces["R"], n)

	if n == width {
		c.Faces["L"] = rotateFaceCCW(c.Faces["L"], n)
	}

	// cycle each of the `width` right‐most slices
	for layer := 0; layer < width; layer++ {
		outerCol := n - 1 - layer // column index on U, F, D
		innerCol := layer         // matching column on B (reversed)

		u := getCol(c.Faces["U"], outerCol, n)
		f := getCol(c.Faces["F"], outerCol, n)
		d := getCol(c.Faces["D"], outerCol, n)
		b := getCol(c.Faces["B"], innerCol, n)

		// U <- F
		setCol(c.Faces["U"], outerCol, n, f)
		// F <- D
		setCol(c.Faces["F"], outerCol, n, d)
		// D <- reversed B
		reverse(b)
		setCol(c.Faces["D"], outerCol, n, b)
		// B <- reversed old U
		reverse(u)
		setCol(c.Faces["B"], innerCol, n, u)
	}
}

// wide R' move: rotate the outer face CCW, then cycle the same `width` slices in reverse
func (c *Cube) MoveRPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid R' width %d for size %d", width, n))
	}

	// rotate the outer R face CCW
	c.Faces["R"] = rotateFaceCCW(c.Faces["R"], n)

	if n == width {
		c.Faces["L"] = rotateFaceCW(c.Faces["L"], n)
	}

	// cycle each of the `width` right‐most slices in the opposite direction
	for layer := 0; layer < width; layer++ {
		outerCol := n - 1 - layer
		innerCol := layer

		u := getCol(c.Faces["U"], outerCol, n)
		f := getCol(c.Faces["F"], outerCol, n)
		d := getCol(c.Faces["D"], outerCol, n)
		b := getCol(c.Faces["B"], innerCol, n)

		// U <- reversed B
		reverse(b)
		setCol(c.Faces["U"], outerCol, n, b)
		// B <- reversed D
		reverse(d)
		setCol(c.Faces["B"], innerCol, n, d)
		// D <- F
		setCol(c.Faces["D"], outerCol, n, f)
		// F <- U
		setCol(c.Faces["F"], outerCol, n, u)
	}
}

// wide U move: rotate the U face CW, then cycle `width` top‐most layers
func (c *Cube) MoveU(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid U width %d for size %d", width, n))
	}

	// rotate the U face CW
	c.Faces["U"] = rotateFaceCW(c.Faces["U"], n)

	if n == width {
		c.Faces["D"] = rotateFaceCCW(c.Faces["D"], n)
	}

	// cycle each of the `width` top‐most slices
	for layer := 0; layer < width; layer++ {
		// extract layer-th row on F, R, B, L
		fRow := getRow(c.Faces["F"], layer, n)
		rRow := getRow(c.Faces["R"], layer, n)
		bRow := getRow(c.Faces["B"], layer, n)
		lRow := getRow(c.Faces["L"], layer, n)

		// F ← R, R ← B, B ← L, L ← F
		setRow(c.Faces["F"], layer, n, rRow)
		setRow(c.Faces["R"], layer, n, bRow)
		setRow(c.Faces["B"], layer, n, lRow)
		setRow(c.Faces["L"], layer, n, fRow)
	}
}

// wide U′ move: rotate the U face CCW, then cycle `width` top‐most layers in reverse
func (c *Cube) MoveUPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid U' width %d for size %d", width, n))
	}

	// rotate the U face CCW
	c.Faces["U"] = rotateFaceCCW(c.Faces["U"], n)

	if n == width {
		c.Faces["D"] = rotateFaceCW(c.Faces["D"], n)
	}

	// cycle each of the `width` top‐most slices in inverse order
	for layer := 0; layer < width; layer++ {
		fRow := getRow(c.Faces["F"], layer, n)
		rRow := getRow(c.Faces["R"], layer, n)
		bRow := getRow(c.Faces["B"], layer, n)
		lRow := getRow(c.Faces["L"], layer, n)

		// F ← L, L ← B, B ← R, R ← F
		setRow(c.Faces["F"], layer, n, lRow)
		setRow(c.Faces["L"], layer, n, bRow)
		setRow(c.Faces["B"], layer, n, rRow)
		setRow(c.Faces["R"], layer, n, fRow)
	}
}

// F move: rotate the front face CW, then cycle `width` front‐most slices
func (c *Cube) MoveF(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid F width %d for size %d", width, n))
	}

	// rotate the front face CW
	c.Faces["F"] = rotateFaceCW(c.Faces["F"], n)

	if n == width {
		c.Faces["B"] = rotateFaceCCW(c.Faces["B"], n)
	}

	for layer := 0; layer < width; layer++ {
		uRow := n - 1 - layer // bottom row of U
		dRow := layer         // top row of D
		rCol := layer         // left column of R
		lCol := n - 1 - layer // right column of L

		u := getRow(c.Faces["U"], uRow, n)
		r := getCol(c.Faces["R"], rCol, n)
		d := getRow(c.Faces["D"], dRow, n)
		l := getCol(c.Faces["L"], lCol, n)

		// U → R
		setCol(c.Faces["R"], rCol, n, u)

		// R → D (reversed)
		reverse(r)
		setRow(c.Faces["D"], dRow, n, r)

		// D → L
		setCol(c.Faces["L"], lCol, n, d)

		// L → U (reversed)
		reverse(l)
		setRow(c.Faces["U"], uRow, n, l)
	}
}

// F′ move: rotate the front face CCW, then cycle `width` front‐most slices in reverse
func (c *Cube) MoveFPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid F' width %d for size %d", width, n))
	}

	// rotate the front face CCW
	c.Faces["F"] = rotateFaceCCW(c.Faces["F"], n)

	if n == width {
		c.Faces["B"] = rotateFaceCW(c.Faces["B"], n)
	}

	for layer := 0; layer < width; layer++ {
		uRow := n - 1 - layer
		dRow := layer
		rCol := layer
		lCol := n - 1 - layer

		u := getRow(c.Faces["U"], uRow, n)
		r := getCol(c.Faces["R"], rCol, n)
		d := getRow(c.Faces["D"], dRow, n)
		l := getCol(c.Faces["L"], lCol, n)

		// inverse cycle:
		// U ← R
		setRow(c.Faces["U"], uRow, n, r)

		// R ← D (reversed)
		reverse(d)
		setCol(c.Faces["R"], rCol, n, d)

		// D ← L
		setRow(c.Faces["D"], dRow, n, l)

		// L ← U (reversed)
		reverse(u)
		setCol(c.Faces["L"], lCol, n, u)
	}
}

// L move: rotate the left face CW, then cycle `width` left‐most slices
func (c *Cube) MoveL(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid L width %d for size %d", width, n))
	}

	// rotate the L face CW
	c.Faces["L"] = rotateFaceCW(c.Faces["L"], n)

	if n == width {
		c.Faces["R"] = rotateFaceCCW(c.Faces["R"], n)
	}

	for layer := 0; layer < width; layer++ {
		outerCol := layer
		innerCol := n - 1 - layer

		u := getCol(c.Faces["U"], outerCol, n)
		f := getCol(c.Faces["F"], outerCol, n)
		d := getCol(c.Faces["D"], outerCol, n)
		b := getCol(c.Faces["B"], innerCol, n)

		// U ← reversed B
		reverse(b)
		setCol(c.Faces["U"], outerCol, n, b)
		// B ← reversed D
		reverse(d)
		setCol(c.Faces["B"], innerCol, n, d)
		// D ← F
		setCol(c.Faces["D"], outerCol, n, f)
		// F ← U
		setCol(c.Faces["F"], outerCol, n, u)
	}
}

// L′ move: rotate the left face CCW, then cycle `width` left‐most slices in reverse
func (c *Cube) MoveLPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid L' width %d for size %d", width, n))
	}

	// rotate the L face CCW
	c.Faces["L"] = rotateFaceCCW(c.Faces["L"], n)

	if n == width {
		c.Faces["R"] = rotateFaceCW(c.Faces["R"], n)
	}

	for layer := 0; layer < width; layer++ {
		outerCol := layer         // column on U, F, D
		innerCol := n - 1 - layer // matching column on B (reversed)

		u := getCol(c.Faces["U"], outerCol, n)
		f := getCol(c.Faces["F"], outerCol, n)
		d := getCol(c.Faces["D"], outerCol, n)
		b := getCol(c.Faces["B"], innerCol, n)

		// U ← F
		setCol(c.Faces["U"], outerCol, n, f)
		// F ← D
		setCol(c.Faces["F"], outerCol, n, d)
		// D ← reversed B
		reverse(b)
		setCol(c.Faces["D"], outerCol, n, b)
		// B ← reversed old U
		reverse(u)
		setCol(c.Faces["B"], innerCol, n, u)
	}
}

// D move: rotate the down face CW, then cycle `width` bottom‐most slices
func (c *Cube) MoveD(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid D width %d for size %d", width, n))
	}

	// rotate the D face CW
	c.Faces["D"] = rotateFaceCW(c.Faces["D"], n)

	if n == width {
		c.Faces["U"] = rotateFaceCCW(c.Faces["U"], n)
	}

	for layer := 0; layer < width; layer++ {
		row := n - 1 - layer

		fRow := getRow(c.Faces["F"], row, n)
		rRow := getRow(c.Faces["R"], row, n)
		bRow := getRow(c.Faces["B"], row, n)
		lRow := getRow(c.Faces["L"], row, n)

		// cycle: F ← L, L ← B, B ← R, R ← F
		setRow(c.Faces["F"], row, n, lRow)
		setRow(c.Faces["L"], row, n, bRow)
		setRow(c.Faces["B"], row, n, rRow)
		setRow(c.Faces["R"], row, n, fRow)
	}
}

// D′ move: rotate the down face CCW, then cycle `width` bottom‐most slices in reverse
func (c *Cube) MoveDPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid D' width %d for size %d", width, n))
	}

	// rotate the D face CCW
	c.Faces["D"] = rotateFaceCCW(c.Faces["D"], n)

	if n == width {
		c.Faces["U"] = rotateFaceCW(c.Faces["U"], n)
	}

	for layer := 0; layer < width; layer++ {
		row := n - 1 - layer

		fRow := getRow(c.Faces["F"], row, n)
		rRow := getRow(c.Faces["R"], row, n)
		bRow := getRow(c.Faces["B"], row, n)
		lRow := getRow(c.Faces["L"], row, n)

		// inverse cycle: F ← R, R ← B, B ← L, L ← F
		setRow(c.Faces["F"], row, n, rRow)
		setRow(c.Faces["R"], row, n, bRow)
		setRow(c.Faces["B"], row, n, lRow)
		setRow(c.Faces["L"], row, n, fRow)
	}
}

// B move: rotate the back face CW, then cycle `width` back‐most slices
func (c *Cube) MoveB(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid B width %d for size %d", width, n))
	}

	// rotate the B face CW
	c.Faces["B"] = rotateFaceCW(c.Faces["B"], n)

	if n == width {
		c.Faces["F"] = rotateFaceCCW(c.Faces["F"], n)
	}

	for layer := 0; layer < width; layer++ {
		// layer‐th back slice touches:
		//   U row `layer`
		//   L col `layer`
		//   D row `n-1-layer`
		//   R col `n-1-layer`
		uRow, lCol := layer, layer
		dRow, rCol := n-1-layer, n-1-layer

		u := getRow(c.Faces["U"], uRow, n)
		l := getCol(c.Faces["L"], lCol, n)
		d := getRow(c.Faces["D"], dRow, n)
		r := getCol(c.Faces["R"], rCol, n)

		// cycle: U → L (reversed), L → D, D → R (reversed), R → U

		// U → L (reversed)
		reverse(u)
		setCol(c.Faces["L"], lCol, n, u)

		// L → D
		setRow(c.Faces["D"], dRow, n, l)

		// D → R (reversed)
		reverse(d)
		setCol(c.Faces["R"], rCol, n, d)

		// R → U
		setRow(c.Faces["U"], uRow, n, r)
	}
}

// B′ move: rotate the back face CCW, then cycle `width` back‐most slices in reverse
func (c *Cube) MoveBPrime(width int) {
	n := c.Size
	if width < 1 || width > n {
		panic(fmt.Sprintf("invalid B' width %d for size %d", width, n))
	}

	// rotate the B face CCW
	c.Faces["B"] = rotateFaceCCW(c.Faces["B"], n)

	if n == width {
		c.Faces["F"] = rotateFaceCW(c.Faces["F"], n)
	}

	for layer := 0; layer < width; layer++ {
		uRow, lCol := layer, layer
		dRow, rCol := n-1-layer, n-1-layer

		// grab original edges
		u := getRow(c.Faces["U"], uRow, n)
		l := getCol(c.Faces["L"], lCol, n)
		d := getRow(c.Faces["D"], dRow, n)
		r := getCol(c.Faces["R"], rCol, n)

		// cycle CCW around the back face with reversals on every transfer:
		// U → R
		setCol(c.Faces["R"], rCol, n, u)

		// R → D (reversed)
		reverse(r)
		setRow(c.Faces["D"], dRow, n, r)

		// D → L
		setCol(c.Faces["L"], lCol, n, d)

		// L → U (reversed)
		reverse(l)
		setRow(c.Faces["U"], uRow, n, l)
	}
}

// M: slice between L and R, turning like an inner L
func (c *Cube) MoveM() {
	c.MoveLPrime(1)
	c.MoveL(c.Size - 1)
}

func (c *Cube) MoveMPrime() {
	c.MoveL(1)
	c.MoveLPrime(c.Size - 1)
}

// E: slice between U and D, turning like an inner D
func (c *Cube) MoveE() {
	c.MoveDPrime(1)
	c.MoveD(c.Size - 1)
}

func (c *Cube) MoveEPrime() {
	c.MoveD(1)
	c.MoveDPrime(c.Size - 1)
}

// S: slice between F and B, turning like an inner F
func (c *Cube) MoveS() {
	c.MoveFPrime(1)
	c.MoveF(c.Size - 1)
}

func (c *Cube) MoveSPrime() {
	c.MoveF(1)
	c.MoveFPrime(c.Size - 1)
}

// RotateX  = turn the entire cube around the R-axis (like Rw on every layer),
func (c *Cube) RotateX() {
	c.MoveR(c.Size)
}

func (c *Cube) RotateXPrime() {
	c.MoveRPrime(c.Size)
}

// RotateY  = turn the entire cube around the U-axis (like Uw on every layer)
func (c *Cube) RotateY() {
	c.MoveU(c.Size)
}

func (c *Cube) RotateYPrime() {
	c.MoveUPrime(c.Size)
}

// RotateZ  = turn the entire cube around the F-axis (like Fw on every layer)
func (c *Cube) RotateZ() {
	c.MoveF(c.Size)
}

func (c *Cube) RotateZPrime() {
	c.MoveFPrime(c.Size)
}

// Moves applies a sequence of face turns to the cube.
// You can pass either:
//   - A single string containing moves and optional parentheses, e.g. c.Moves("(U) R U2 R' F")
//   - Multiple string arguments, e.g.     c.Moves("U", "R", "U2", "R'", "F")
//
// Parentheses in any string will be stripped before parsing.
func (c *Cube) Moves(seqs ...string) error {
	// 1) flatten into one space-separated string
	seq := strings.Join(seqs, " ")

	// 2) strip parentheses
	seq = strings.NewReplacer("(", "", ")", "").Replace(seq)

	// 3) split into individual move tokens
	for _, m := range strings.Fields(seq) {
		if err := c.Move(m); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cube) Move(notation string) error {
	token := notation
	times := 1
	isPrime := false

	// 1) strip a trailing "'" → prime (unless it was a 2-turn)
	if strings.HasSuffix(token, "'") {
		isPrime = true
		token = token[:len(token)-1]
	}

	// 2) strip a trailing "2" → double‐turn
	if strings.HasSuffix(token, "2") {
		times = 2
		token = token[:len(token)-1]
	} else if strings.HasSuffix(token, "3") {
		times = 3
		token = token[:len(token)-1]
	}

	if token == "" {
		return fmt.Errorf("invalid move: %s", notation)
	}

	// determine face/axis and width
	var face byte
	width := 1

	replacer := strings.NewReplacer(
		"u", "Uw",
		"d", "Dw",
		"r", "Rw",
		"l", "Lw",
		"f", "Fw",
		"b", "Bw",
	)
	token = replacer.Replace(token)

	// wide move notation ends in 'w' or 'W'
	if last := token[len(token)-1]; last == 'w' || last == 'W' {
		prefix := token[:len(token)-1]
		if len(prefix) < 1 {
			return fmt.Errorf("invalid wide move: %s", notation)
		}
		face = prefix[len(prefix)-1]
		num := prefix[:len(prefix)-1]
		if num == "" {
			width = 2
		} else {
			w, err := strconv.Atoi(num)
			if err != nil {
				return fmt.Errorf("invalid width in move %s: %v", notation, err)
			}
			width = w
		}
	} else if len(token) == 1 {
		// normal single‐layer
		face = token[0]
		width = 1
	} else {
		// maybe a numeric-prefix wide without 'w', e.g. "3R"
		// parse leading number
		i := 0
		for ; i < len(token) && token[i] >= '0' && token[i] <= '9'; i++ {
		}
		if i > 0 && i < len(token) {
			num := token[:i]
			w, err := strconv.Atoi(num)
			if err == nil {
				width = w
				face = token[i]
			} else {
				return fmt.Errorf("invalid move: %s", notation)
			}
		} else {
			return fmt.Errorf("invalid move: %s", notation)
		}
	}

	// apply the move the required number of times
	for range times {
		switch face {
		// face-turns
		case 'R':
			if isPrime {
				c.MoveRPrime(width)
			} else {
				c.MoveR(width)
			}
		case 'L':
			if isPrime {
				c.MoveLPrime(width)
			} else {
				c.MoveL(width)
			}
		case 'U':
			if isPrime {
				c.MoveUPrime(width)
			} else {
				c.MoveU(width)
			}
		case 'D':
			if isPrime {
				c.MoveDPrime(width)
			} else {
				c.MoveD(width)
			}
		case 'F':
			if isPrime {
				c.MoveFPrime(width)
			} else {
				c.MoveF(width)
			}
		case 'B':
			if isPrime {
				c.MoveBPrime(width)
			} else {
				c.MoveB(width)
			}
		case 'M':
			if isPrime {
				c.MoveMPrime()
			} else {
				c.MoveM()
			}
		case 'E':
			if isPrime {
				c.MoveEPrime()
			} else {
				c.MoveE()
			}
		case 'S':
			if isPrime {
				c.MoveSPrime()
			} else {
				c.MoveS()
			}

		// whole-cube rotations
		case 'x':
			if isPrime {
				c.RotateXPrime()
			} else {
				c.RotateX()
			}
		case 'y':
			if isPrime {
				c.RotateYPrime()
			} else {
				c.RotateY()
			}
		case 'z':
			if isPrime {
				c.RotateZPrime()
			} else {
				c.RotateZ()
			}

		default:
			return fmt.Errorf("invalid face or axis in move %s", notation)
		}
	}

	return nil
}

func (c *Cube) IsSolved() bool {
	for _, stickers := range c.Faces {
		for _, s := range stickers {
			if s != stickers[0] {
				return false
			}
		}
	}
	return true
}

// Copy creates a deep copy of the given cube.
func (c *Cube) Copy() *Cube {
	faces := make(map[string][]string, len(c.Faces))
	for face, stickers := range c.Faces {
		copyStickers := make([]string, len(stickers))
		copy(copyStickers, stickers)
		faces[face] = copyStickers
	}
	return &Cube{Size: c.Size, Faces: faces}
}

func main() {
	for _, sz := range []int{2, 3, 4, 5, 6, 7} {
		fmt.Printf("\n%dx%dx%d Cube:\n", sz, sz, sz)
		NewCube(sz).DisplayColor()
	}

	c := NewCube(3)
	c.Moves("(U2) x R2 D2 R U R' D2 R U' R")
	c.DisplayColor()
}
