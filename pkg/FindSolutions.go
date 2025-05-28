package pkg

import (
	"strings"
)

// CheckFunc defines a predicate to test cube state (e.g., solved).
type CheckFunc func(c *Cube) bool

// FindSolutions searches for move sequences (up to maxDepth) that satisfy the check function.
// moveSet is the list of moves to try (e.g., []string{"R", "R'", "R2", ...}).
// Returns a slice of solutions, each solution being a slice of move tokens.
func FindSolutions(initial *Cube, moveSet []string, check CheckFunc, maxDepth int, progress chan<- struct{}) []string {
	solutions := make([]string, 0)

	var dfs func(path []string)
	dfs = func(path []string) {
		// Explore each move
		for _, m := range moveSet {
			// optional: skip consecutive moves on same face
			if len(path) > 0 && path[len(path)-1][0] == m[0] {
				continue
			}

			next := initial.Copy()
			path := append(path, m)

			// apply move
			for _, m := range path {
				if err := next.Move(m); err != nil {
					Printf("skipping invalid move %s: %v\n", m, err)
					continue
				}
			}

			// if err := next.Moves(path...); err != nil {
			// 	Printf("skipping invalid move %s: %v\n", m, err)
			// 	continue
			// }

			// tick the progress bar
			if progress != nil {
				progress <- struct{}{}
			}

			// If check passes, record solution
			if check(next) {
				seq := strings.Join(path, " ")
				solutions = append(solutions, seq)
				Printf("Solution found: %s\n", seq)
			}

			// recurse
			if len(path) < maxDepth {
				dfs(path)
			}
		}
	}

	// start DFS
	dfs([]string{})

	return solutions
}
