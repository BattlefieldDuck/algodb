package pkg

import (
	"sync"
)

// FindSolutionsParallel spawns one worker per first‐move.
func FindSolutionsParallel(
	initial *Cube,
	moves []string,
	check CheckFunc,
	maxDepth int,
	progress chan<- struct{},
) [][]string {
	var (
		wg        sync.WaitGroup
		solMu     sync.Mutex
		solutions [][]string
	)

	for _, first := range moves {
		// skip repeated‐face pruning on depth=1 if you like
		wg.Add(1)
		go func(firstMove string) {
			defer wg.Done()

			// each branch gets its own explicit stack of paths
			type frame struct{ path []string }
			stack := []frame{{path: []string{firstMove}}}

			for len(stack) > 0 {
				// pop
				f := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				// apply this branch's current sequence
				next := initial.Copy()

				for _, m := range f.path {
					if err := next.Move(m); err != nil {
						Printf("skipping invalid move %s: %v\n", m, err)
						continue
					}
				}

				// tick progress
				if progress != nil {
					progress <- struct{}{}
				}

				// record solution if solved
				if check(next) {
					solMu.Lock()
					solutions = append(solutions, f.path)
					solMu.Unlock()
				}

				// push deeper children
				if len(f.path) < maxDepth {
					lastFace := f.path[len(f.path)-1][0]
					for _, mv := range moves {
						if mv[0] == lastFace {
							continue
						}
						// copy path and append
						newPath := make([]string, len(f.path), len(f.path)+1)
						copy(newPath, f.path)
						newPath = append(newPath, mv)
						stack = append(stack, frame{path: newPath})
					}
				}
			}
		}(first)
	}

	wg.Wait()

	return solutions
}
