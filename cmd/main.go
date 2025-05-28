package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BattlefieldDuck/algodb/pkg"
	"github.com/schollz/progressbar/v3"
)

// isSolved wraps your cube’s solved‐state check.
func isSolved(c *pkg.Cube) bool { return c.IsSolved() }

func main() {
	// Expect exactly 5 args: program, config, id, depth, moves
	if len(os.Args) != 5 {
		log.Fatalf("Usage: %s <config.csv> <id> <maxDepth> <move_set>\n", os.Args[0])
	}
	configPath := os.Args[1]
	targetID := os.Args[2]
	depthArg := os.Args[3]
	movesArg := os.Args[4]

	// Parse maxDepth and moves from CLI
	maxDepth, err := strconv.Atoi(depthArg)
	if err != nil {
		log.Fatalf("Invalid maxDepth %q: %v", depthArg, err)
	}
	moves := strings.Fields(movesArg)

	// Derive cube size (n) and config base name
	base := filepath.Base(configPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	parts := strings.Split(name, "-")
	sizeChar := parts[0][0]
	n, err := strconv.Atoi(string(sizeChar))
	if err != nil {
		log.Fatalf("Cannot parse cube size from %s: %v", parts[0], err)
	}

	// Read CSV for scramble
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening %s: %v", configPath, err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	// Find scramble by ID
	var scramble string
	found := false
	for i, rec := range records {
		if i == 0 {
			continue // header
		}
		if len(rec) < 2 {
			continue
		}
		if rec[0] == targetID {
			scramble = rec[1]
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("ID %s not found in %s", targetID, configPath)
	}

	// Create and scramble the cube
	c := pkg.NewCube(n)
	if err := c.Moves(scramble); err != nil {
		log.Fatalf("Error scrambling cube: %v", err)
	}

	// Measure start time
	start := time.Now()

	// Display cube state
	pkg.Printf("ID: %s\n", targetID)
	pkg.Printf("MaxDepth: %d\n", maxDepth)
	pkg.Printf("MoveSet: %s\n", movesArg)

	fmt.Printf("\n%dx%dx%d Cube - %s\n\n", n, n, n, scramble)
	c.DisplayColor()
	fmt.Println()

	// Compute branching parameters
	faceSet := make(map[rune]struct{})
	for _, m := range moves {
		faceSet[rune(m[0])] = struct{}{}
	}
	distinctFaces := len(faceSet)
	branchingFactor := len(moves) - distinctFaces

	// Compute total DFS nodes
	total := len(moves) *
		(int(math.Pow(float64(branchingFactor), float64(maxDepth))) - 1) /
		(branchingFactor - 1)

	pkg.Printf("Total nodes to explore: %d\n", total)

	// Setup progress bar
	progress := make(chan struct{}, 1000)
	bar := progressbar.NewOptions(
		total,
		progressbar.OptionSetDescription("Searching"),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
	)
	go func() {
		for range progress {
			bar.Add(1)
		}
		bar.Close()
	}()

	// Estimate time assuming 90k nodes per second
	estSec := total / 90000
	pkg.Printf("Estimated time at 90k nodes/s: %d seconds (~%s)\n", estSec, time.Duration(estSec)*time.Second)

	// Run parallel solver
	solutions := pkg.FindSolutionsParallel(c, moves, isSolved, maxDepth, nil)

	// Measure elapsed time and throughput
	elapsed := time.Since(start)
	pkg.Printf("Elapsed time: %s\n", elapsed)
	throughput := float64(total) / elapsed.Seconds()
	pkg.Printf("Nodes per second: %.2f\n", throughput)

	// Sort solutions by length (fewest moves first)
	sort.Slice(solutions, func(i, j int) bool {
		return len(strings.Fields(solutions[i])) < len(strings.Fields(solutions[j]))
	})

	// Print solutions
	pkg.Printf("Found %d solution(s):\n\n", len(solutions))
	for i, sol := range solutions {
		fmt.Printf("%2d [%d]: %s\n", i+1, len(strings.Fields(sol)), sol)
	}
	fmt.Println()

	// Prepare result struct
	result := struct {
		ConfigPath string        `json:"config_path"`
		ID         string        `json:"id"`
		Scramble   string        `json:"scramble"`
		MaxDepth   int           `json:"max_depth"`
		MoveSet    []string      `json:"move_set"`
		TotalNodes int           `json:"total_nodes"`
		ElapsedNS  time.Duration `json:"elapsed_ns"`
		Elapsed    string        `json:"elapsed"`
		Throughput float64       `json:"throughput"`
		Solutions  []string      `json:"solutions"`
		CreatedAt  time.Time     `json:"created_at"`
	}{
		ConfigPath: configPath,
		ID:         targetID,
		Scramble:   scramble,
		MaxDepth:   maxDepth,
		MoveSet:    moves,
		TotalNodes: total,
		ElapsedNS:  elapsed,
		Elapsed:    elapsed.String(),
		Throughput: throughput,
		Solutions:  solutions,
		CreatedAt:  time.Now(),
	}

	// Save results to JSON file in ./db/<configName>
	outDir := filepath.Join("db", name)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Error creating results dir: %v", err)
	}
	fname := filepath.Join(outDir, fmt.Sprintf("%s.json", targetID))
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling result: %v", err)
	}
	if err := os.WriteFile(fname, data, 0644); err != nil {
		log.Fatalf("Error writing result file: %v", err)
	}
	pkg.Printf("Saved results to %s\n", fname)
}
