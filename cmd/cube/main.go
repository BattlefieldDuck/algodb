package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BattlefieldDuck/algodb/internal"
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

	// Display cube state
	pkg.Printf("ID: %s\n", targetID)
	pkg.Printf("MaxDepth: %d\n", maxDepth)
	pkg.Printf("MoveSet: %s\n", movesArg)

	fmt.Printf("\n%dx%dx%d Cube - %s\n\n", n, n, n, scramble)
	c.DisplayColorANSI()

	fmt.Printf("\nUp Face\n\n")
	c.DisplayColorANSIUFace()
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

	// Estimate time assuming 100000k nodes per second
	estSec := total / 100000000
	pkg.Printf("Estimated time at 100000k nodes/s: %d seconds (~%s)\n", estSec, time.Duration(estSec)*time.Second)

	// Measure start time
	start := time.Now()

	// Run parallel solver
	solutions := pkg.FindSolutionsParallelDFS(c, moves, isSolved, maxDepth, nil)

	// Measure elapsed time and throughput
	elapsed := time.Since(start)
	pkg.Printf("Elapsed time: %s\n", elapsed)
	throughput := float64(total) / elapsed.Seconds()
	pkg.Printf("Nodes per second: %.2f\n", throughput)

	// Print solutions
	pkg.Printf("Found %d solution(s):\n\n", len(solutions))
	for i, sol := range solutions {
		fmt.Printf("%2d [%d]: %s\n", i+1, len(sol), strings.Join(sol, " "))
	}
	fmt.Println()

	internal.CreateAlgorithms(name, targetID, solutions)
}
