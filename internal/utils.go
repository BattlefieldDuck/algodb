package internal

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// CreateAlgorithms will:
// 1. Normalize any U-layer first moves into a y-rotation.
// 2. Sort by move-count (ignoring any x/y/z rotations) then lexicographically.
// 3. Write out a CSV at /db/<name>/<targetID>.csv with columns: length,prefix,algorithm
func CreateAlgorithms(name, targetID string, solutions [][]string) error {
	type entry struct {
		prefix   string   // the x/y/z rotation (if any)
		algMoves []string // the face-turns only (no x/y/z)
		fullAlg  []string // the full move list after normalization
	}

	var list []entry
	for _, sol := range solutions {
		// copy so we don’t clobber callers’ slice
		moves := append([]string(nil), sol...)

		// 1) if first move is U, U' or U2 → turn it into a cube-rotation on y
		if len(moves) > 0 {
			switch moves[0] {
			case "U":
				moves[0] = "y"
			case "U'":
				moves[0] = "y'"
			case "U2":
				moves[0] = "y2"
			}
		}

		// extract a prefix if it’s an x/y/z rotation
		prefix := ""
		if len(moves) > 0 {
			p := moves[0]
			if strings.HasPrefix(p, "x") || strings.HasPrefix(p, "y") || strings.HasPrefix(p, "z") {
				prefix = p
			}
		}

		// gather only the face-turns (drop any x/y/z)
		var faceTurns []string
		for _, m := range moves {
			if strings.HasPrefix(m, "x") || strings.HasPrefix(m, "y") || strings.HasPrefix(m, "z") {
				continue
			}
			faceTurns = append(faceTurns, m)
		}

		list = append(list, entry{
			prefix:   prefix,
			algMoves: faceTurns,
			fullAlg:  moves,
		})
	}

	// 2) sort by length, then prefix, then moves, then by the lexicographic join
	sort.Slice(list, func(i, j int) bool {
		if len(list[i].algMoves) != len(list[j].algMoves) {
			return len(list[i].algMoves) < len(list[j].algMoves)
		}
		if list[i].prefix != list[j].prefix {
			return list[i].prefix < list[j].prefix
		}
		// join here to compare the sequences
		ai := strings.Join(list[i].algMoves, " ")
		aj := strings.Join(list[j].algMoves, " ")
		return ai < aj
	})

	// ensure the output directory exists
	outDir := filepath.Join("db", name)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	outFile := filepath.Join(outDir, targetID+".csv")
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header row
	if err := w.Write([]string{"length", "prefix", "algorithm"}); err != nil {
		return err
	}

	// 3) write each sorted entry
	for _, e := range list {
		lengthStr := strconv.Itoa(len(e.algMoves))
		// algorithm is the trimmed, face-turn sequence
		algStr := strings.Join(e.algMoves, " ")
		if err := w.Write([]string{lengthStr, e.prefix, algStr}); err != nil {
			return err
		}
	}

	return nil
}
