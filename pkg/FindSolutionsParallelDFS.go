package pkg

import (
	"sync"
)

// op bundles a move notation and its inverse for fast backtracking
type op struct {
	notation string
	face     int
	apply    func(c *Cube)
	undo     func(c *Cube)
}

// FindSolutionsParallelDFS launches one goroutine per first move and performs
// in-place DFS with backtracking to find all sequences up to maxDepth.
func FindSolutionsParallelDFS(
	initial *Cube,
	moves []string,
	check CheckFunc,
	maxDepth int,
	progress chan<- struct{},
) [][]string {
	// precompute ops
	ops := make([]op, len(moves))
	for i, m := range moves {
		face, apply, undo := initial.CreateOp(m)
		ops[i] = op{m, face, apply, undo}
	}

	var (
		wg        sync.WaitGroup
		solMu     sync.Mutex
		solutions [][]string
	)

	// spawn one goroutine per first move
	for _, root := range ops {
		wg.Add(1)
		go func(root op) {
			defer wg.Done()

			// === per-goroutine local buffer ===
			var local [][]op

			// one copy per branch
			c := initial.Copy()
			root.apply(c)

			// recursive DFS closure
			var dfs func(c *Cube, path []op)
			dfs = func(c *Cube, path []op) {
				// tick progress
				if progress != nil {
					progress <- struct{}{}
				}

				l := len(path)

				// record solution
				if check(c) {
					// deep-copy the path before recording
					cp := make([]op, l)
					copy(cp, path)
					local = append(local, cp)
					return
				}
				if l == maxDepth {
					return
				}

				lastFace := path[l-1].face
				for _, op := range ops {
					if op.face == lastFace {
						continue
					}

					// apply move
					path = append(path, op)
					op.apply(c)

					dfs(c, path)

					// backtrack
					path = path[:l]
					op.undo(c)
				}
			}

			// start path
			path := make([]op, 0, maxDepth+1)
			path = append(path, root)
			dfs(c, path)

			// merge once
			solMu.Lock()

			for _, ops := range local {
				// make a []string the same length as this []op
				seq := make([]string, len(ops))
				// fill it with the notation of each op
				for i, op := range ops {
					seq[i] = op.notation
				}
				// now append that []string to your solutions
				solutions = append(solutions, seq)
			}

			solMu.Unlock()
		}(root)
	}

	wg.Wait()
	return solutions
}
