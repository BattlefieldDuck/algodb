package pkg

import (
	"strings"
	"sync"
)

// op bundles a move notation and its inverse for fast backtracking
type op struct {
	notation string
	inv      string
}

// inverseOf returns the inverse of a move notation (e.g. "R"→"R'", "R'"→"R", "R2"→"R2")
func inverseOf(m string) string {
	if strings.HasSuffix(m, "'") {
		return m[:len(m)-1]
	}
	if strings.HasSuffix(m, "2") {
		return m
	}
	return m + "'"
}

// FindSolutionsParallelDFS launches one goroutine per first move and performs
// in-place DFS with backtracking to find all sequences up to maxDepth.
func FindSolutionsParallelDFS(
	initial *Cube,
	moves []string,
	check CheckFunc,
	maxDepth int,
	progress chan<- struct{},
) []string {
	// precompute ops
	ops := make([]op, len(moves))
	for i, m := range moves {
		ops[i] = op{notation: m, inv: inverseOf(m)}
	}

	var (
		wg        sync.WaitGroup
		solMu     sync.Mutex
		solutions []string
	)

	// recursive DFS closure
	var dfs func(c *Cube, path []string, depth int)
	dfs = func(c *Cube, path []string, depth int) {
		// tick progress
		if progress != nil {
			progress <- struct{}{}
		}
		// record solution
		if check(c) {
			seq := strings.Join(path, " ")
			solMu.Lock()
			solutions = append(solutions, seq)
			solMu.Unlock()
		}
		if depth == maxDepth {
			return
		}

		lastFace := path[len(path)-1][0]
		for _, op := range ops {
			if op.notation[0] == lastFace {
				continue
			}
			// apply move
			c.Move(op.notation)
			path = append(path, op.notation)

			dfs(c, path, depth+1)

			// backtrack
			path = path[:len(path)-1]
			c.Move(op.inv)
		}
	}

	// spawn one goroutine per first move
	for _, root := range ops {
		wg.Add(1)
		go func(root op) {
			defer wg.Done()
			// one copy per branch
			c := initial.Copy()
			c.Move(root.notation)

			// start path
			path := []string{root.notation}
			dfs(c, path, 1)
		}(root)
	}

	wg.Wait()
	return solutions
}
