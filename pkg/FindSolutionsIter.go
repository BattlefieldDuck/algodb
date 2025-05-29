package pkg

import (
	"slices"
	"strings"
)

func FindSolutionsIter(
	initial *Cube,
	moveSet []string,
	check CheckFunc,
	maxDepth int,
	progress chan<- struct{},
) []string {
	type frame struct {
		path []string // moves so far
	}

	solutions := make([]string, 0)
	// our explicit stack of frames
	stack := []frame{{path: []string{}}}

	for len(stack) > 0 {
		// pop
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, m := range moveSet {
			// skip same-axis repeats
			if len(f.path) > 0 && f.path[len(f.path)-1][0] == m[0] {
				continue
			}

			// build next path
			newPath := append(slices.Clone(f.path), m)

			// apply move onto a fresh copy
			next := initial.Copy()
			for _, m := range newPath {
				if err := next.Move(m); err != nil {
					Printf("skipping invalid move %s: %v\n", m, err)
					continue
				}
			}

			// progress tick
			if progress != nil {
				progress <- struct{}{}
			}

			// record solution
			if check(next) {
				seq := strings.Join(newPath, " ")
				solutions = append(solutions, seq)
				Printf("Solution found: %s\n", seq)
			}

			// if we can go deeper, push onto stack
			if len(newPath) < maxDepth {
				stack = append(stack, frame{path: newPath})
			}
		}
	}

	return solutions
}
