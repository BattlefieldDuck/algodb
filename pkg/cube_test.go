package pkg

import (
	"fmt"
	"strings"
	"testing"
)

func displayCube(c *Cube) {
	fmt.Printf("\nDisplay\n\n")
	c.Display()
	fmt.Printf("\nDisplayColorANSI\n\n")
	c.DisplayColorANSI()
	fmt.Printf("\nDisplayColorUnicode\n\n")
	c.DisplayColorUnicode()
	fmt.Printf("\nDisplayColorUnicodeFrontFaceOnly\n\n")
	c.DisplayColorUnicodeUFace()
	fmt.Println()
}

func Test222(t *testing.T) {
	c := NewCube(2)
	c.Moves("R2 U F2 U F' R' U2 R U2")
	displayCube(c)
	c.Moves("U2 R' U2 R F U' F2 U' R2")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func Test333(t *testing.T) {
	c := NewCube(3)
	c.Moves("U2 B L' R2 U2 B L B U2 R2 D' R' B' U B D2 F' U L B2")
	displayCube(c)
	c.Moves("B2 L' U' F D2 B' U' B R D R2 U2 B' L' B' U2 R2 L B' U2")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func Test444(t *testing.T) {
	c := NewCube(4)
	c.Moves("U2 Fw' L' F Bw U2 L2 D F Bw2 D' L2 Uw D' F' Lw' Uw' F R2 Lw2 Dw' Uw' Rw Lw Fw2 Rw2 Dw B2 D' F' U Dw' Lw Rw' Fw' D' R' D' U Bw")
	displayCube(c)
	c.Moves("Bw' U' D R D Fw Rw Lw' Dw U' F D B2 Dw' Rw2 Fw2 Lw' Rw' Uw Dw Lw2 R2 F' Uw Lw F D Uw' L2 D Bw2 F' D' L2 U2 Bw' F' L Fw U2")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func Test555(t *testing.T) {
	c := NewCube(5)
	c.Moves("L F B2 Dw Lw2 U Dw Fw' Bw' Rw Uw' Fw' Bw2 Uw2 Dw' R2 L' Dw Fw D2 Bw' Uw R2 Bw' F' Uw Lw' Bw2 R' D L Rw2 Bw D2 B2 F' Dw' Bw F2 U2 B R Bw' L2 F2 U' Fw2 B2 Dw2 Lw' D' U B' Fw2 L R Uw L' D2 R")
	displayCube(c)
	c.Moves("R' D2 L Uw' R' L' Fw2 B U' D Lw Dw2 B2 Fw2 U F2 L2 Bw R' B' U2 F2 Bw' Dw F B2 D2 Bw' Rw2 L' D' R Bw2 Lw Uw' F Bw R2 Uw' Bw D2 Fw' Dw' L R2 Dw Uw2 Bw2 Fw Uw Rw' Bw Fw Dw' U' Lw2 Dw' B2 F' L'")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func Test666(t *testing.T) {
	c := NewCube(6)
	c.Moves("R' Uw' Dw' 3Fw2 L Bw' Dw' R L 3Bw' Fw2 R2 B' Rw2 U2 F' Dw 3Rw2 3Dw2 3Bw2 Uw2 3Rw 3Lw2 3Uw' Dw2 Lw' Dw' L' U' 3Lw' Fw2 3Bw2 U2 L' Uw Bw F' Lw' Fw 3Bw2 Lw Uw2 3Rw2 Uw F2 Lw2 Bw2 Rw2 3Bw' 3Fw2 Lw' 3Fw' Bw 3Lw' 3Fw 3Rw2 Dw2 Uw' Bw Dw' Rw' Dw' B F2 3Lw' Rw' Uw2 3Lw' Rw2 Uw R2 Lw Bw' Dw2 U' R' Uw' D2 Fw2 3Uw")
	displayCube(c)
	c.Moves("3Uw' Fw2 D2 Uw R U Dw2 Bw Lw' R2 Uw' Rw2 3Lw Uw2 Rw 3Lw F2 B' Dw Rw Dw Bw' Uw Dw2 3Rw2 3Fw' 3Lw Bw' 3Fw Lw 3Fw2 3Bw Rw2 Bw2 Lw2 F2 Uw' 3Rw2 Uw2 Lw' 3Bw2 Fw' Lw F Bw' Uw' L U2 3Bw2 Fw2 3Lw U L Dw Lw Dw2 3Uw 3Lw2 3Rw' Uw2 3Bw2 3Dw2 3Rw2 Dw' F U2 Rw2 B R2 Fw2 3Bw L' R' Dw Bw L' 3Fw2 Dw Uw R")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func Test777(t *testing.T) {
	c := NewCube(7)
	c.Moves("3Bw2 3Uw' Fw2 3Bw' 3Lw' Bw' R U' 3Fw 3Dw2 Fw2 Bw Uw F2 Uw' Rw D2 3Uw2 F2 3Lw B2 3Uw' Rw' D2 3Uw2 R2 F2 Uw B' Uw2 Rw' D2 3Lw' U2 D2 B2 Dw2 U Rw2 3Dw 3Bw2 Fw' Dw 3Fw2 3Rw Uw2 3Lw2 D2 Rw2 3Dw' R2 Bw' D' R Lw Dw2 F2 Bw' R Fw2 Rw 3Bw' D' 3Lw Dw 3Bw2 Lw' 3Bw2 F L 3Rw2 Dw' U2 3Fw Lw' D2 3Uw' 3Bw2 3Fw' Uw 3Dw' Bw2 R2 3Uw' Bw2 Lw2 3Dw Uw2 R2 Bw")
	displayCube(c)
	c.Moves("Bw' R2 Uw2 3Dw' Lw2 Bw2 3Uw R2 Bw2 3Dw Uw' 3Fw 3Bw2 3Uw D2 Lw 3Fw' U2 Dw 3Rw2 L' F' 3Bw2 Lw 3Bw2 Dw' 3Lw' D 3Bw Rw' Fw2 R' Bw F2 Dw2 Lw' R' D Bw R2 3Dw Rw2 D2 3Lw2 Uw2 3Rw' 3Fw2 Dw' Fw 3Bw2 3Dw' Rw2 U' Dw2 B2 D2 U2 3Lw D2 Rw Uw2 B Uw' F2 R2 3Uw2 D2 Rw 3Uw B2 3Lw' F2 3Uw2 D2 Rw' Uw F2 Uw' Bw' Fw2 3Dw2 3Fw' U R' Bw 3Lw 3Bw Fw2 3Uw 3Bw2")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after applying inverse, but it was not")
	}
}

func TestCubeRotations(t *testing.T) {
	c := NewCube(3)
	c.Moves("R x U y R' z U'")
	displayCube(c)
	c.Moves("U z' R y' U' x' R'")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after 90 degree rotation, but it was not")
	}
}

func TestSliceMoves(t *testing.T) {
	c := NewCube(3)
	c.Moves("M E S M2 E2 S2")
	displayCube(c)
	c.Moves("S2 E2 M2 S' E' M'")

	if !c.IsSolved() {
		t.Error("expected cube to be solved after 90 degree rotation, but it was not")
	}
}

type FindSolutionsFunc func(initial *Cube, moveSet []string, check CheckFunc, maxDepth int, progress chan<- struct{}) []string

func testFindSolutions(t *testing.T, findSolutions FindSolutionsFunc, maxDepth int) {
	c := NewCube(3)
	c.Moves("R U2 R' U' R U' R'")
	displayCube(c)

	moves := []string{"R", "R'", "R2", "U", "U'", "U2", "F", "F'", "F2"}
	check := func(c *Cube) bool { return c.IsSolved() }
	solutions := findSolutions(c, moves, check, maxDepth, nil)

	// Print solutions
	Printf("Found %d solution(s):\n\n", len(solutions))
	for i, sol := range solutions {
		fmt.Printf("%2d [%d]: %s\n", i+1, len(strings.Fields(sol)), sol)

		next := c.Copy()
		next.Moves(sol)

		if !next.IsSolved() {
			t.Error("expected cube to be solved after 90 degree rotation, but it was not")
		}
	}
	fmt.Println()
}

func TestFindSolutions(t *testing.T) {
	testFindSolutions(t, FindSolutions, 7)
}

func TestFindSolutionsIter(t *testing.T) {
	testFindSolutions(t, FindSolutionsIter, 7)
}
func TestFindSolutionsParallel(t *testing.T) {
	testFindSolutions(t, FindSolutionsParallel, 7)
}

func TestFindSolutionsParallelDFS(t *testing.T) {
	testFindSolutions(t, FindSolutionsParallelDFS, 8)
}
