package csp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkMapColoringWithoutOrdering(b *testing.B) {
	var neq = func(a CSPVal, b CSPVal) bool {
		return a != b
	}

	// Make all vars not equal to their next or next-next neighbor
	vars := []CSPVar{"WA", "NT", "SA", "Q", "NSW", "V", "T"}
	constraints := make(map[CSPPair]func(CSPVal, CSPVal) bool)
	constraints[CSPPair{a: "WA", b: "NT"}] = neq
	constraints[CSPPair{a: "WA", b: "SA"}] = neq
	constraints[CSPPair{a: "NT", b: "SA"}] = neq
	constraints[CSPPair{a: "Q", b: "NT"}] = neq
	constraints[CSPPair{a: "Q", b: "SA"}] = neq
	constraints[CSPPair{a: "Q", b: "NSW"}] = neq
	constraints[CSPPair{a: "V", b: "SA"}] = neq
	constraints[CSPPair{a: "V", b: "NSW"}] = neq
	constraints[CSPPair{a: "SA", b: "NSW"}] = neq
	sampleProblem := &BinaryCSP{
		vars:        vars,
		domain:      []CSPVal{"Red", "Blue", "Green"},
		constraints: constraints,
		options: map[string]bool{
			"checkViolation": true,
			"ordering":       false,
		},
	}

	fringe := &StackFringe{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchWithFringe(fringe, sampleProblem)
	}
}

func BenchmarkMapColoringWithOrdering(b *testing.B) {
	var neq = func(a CSPVal, b CSPVal) bool {
		return a != b
	}

	// Make all vars not equal to their next or next-next neighbor
	vars := []CSPVar{"WA", "NT", "SA", "Q", "NSW", "V", "T"}
	constraints := make(map[CSPPair]func(CSPVal, CSPVal) bool)
	constraints[CSPPair{a: "WA", b: "NT"}] = neq
	constraints[CSPPair{a: "WA", b: "SA"}] = neq
	constraints[CSPPair{a: "NT", b: "SA"}] = neq
	constraints[CSPPair{a: "Q", b: "NT"}] = neq
	constraints[CSPPair{a: "Q", b: "SA"}] = neq
	constraints[CSPPair{a: "Q", b: "NSW"}] = neq
	constraints[CSPPair{a: "V", b: "SA"}] = neq
	constraints[CSPPair{a: "V", b: "NSW"}] = neq
	constraints[CSPPair{a: "SA", b: "NSW"}] = neq
	sampleProblem := &BinaryCSP{
		vars:        vars,
		domain:      []CSPVal{"Red", "Blue", "Green"},
		constraints: constraints,
		options: map[string]bool{
			"checkViolation": true,
			"ordering":       true,
		},
	}

	fringe := &StackFringe{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchWithFringe(fringe, sampleProblem)
	}
}

func Benchmark8QueensWithOrdering(b *testing.B) {
	var notThreatening = func(a CSPVal, b CSPVal) bool {
		if a == "" || b == "" {
			return false
		}
		aCoords := strings.Split(string(a), ",")
		bCoords := strings.Split(string(b), ",")
		ax, _ := strconv.Atoi(aCoords[0])
		ay, _ := strconv.Atoi(aCoords[1])
		bx, _ := strconv.Atoi(bCoords[0])
		by, _ := strconv.Atoi(bCoords[1])

		// not same row
		if ax == bx {
			return false
		}

		// not same column
		if ay == by {
			return false
		}

		var abs = func(num int) int {
			if num < 0 {
				return num * -1
			}
			return num
		}
		// Not diagonal to each other
		if abs(ax-bx) == abs(ay-by) {
			return false
		}

		return true
	}

	// Make all vars not equal to their next or next-next neighbor
	vars := []CSPVar{"Q1", "Q2", "Q3", "Q4", "Q5", "Q6", "Q7", "Q8"}
	//vars := []CSPVar{"Q1", "Q2", "Q3", "Q4"}
	grid := []CSPVal{}
	for i := 0; i < len(vars); i++ {
		for j := 0; j < len(vars); j++ {
			grid = append(grid, CSPVal(fmt.Sprintf("%d,%d", i, j)))
		}
	}
	constraints := make(map[CSPPair]func(CSPVal, CSPVal) bool)
	for i := 1; i < len(vars)+1; i++ {
		for j := 1; j < len(vars)+1; j++ {
			a := fmt.Sprintf("Q%d", i)
			b := fmt.Sprintf("Q%d", j)
			if a != b {
				constraints[CSPPair{a: CSPVar(a), b: CSPVar(b)}] = notThreatening
			}
		}
	}
	sampleProblem := &BinaryCSP{
		vars:        vars,
		domain:      grid,
		constraints: constraints,
		options: map[string]bool{
			"checkViolation": true,
			"ordering":       true,
		},
	}

	fringe := &StackFringe{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		printQueens(searchWithFringe(fringe, sampleProblem))
	}
}

func printQueens(queensJson string) {
	queens := make(map[string]string)
	queensJson = strings.Split(queensJson, "found: ")[1]
	if err := json.Unmarshal([]byte(queensJson), &queens); err != nil {
		fmt.Println("Couldn't unmarshal jack shit")
		return
	}
	grid := make([][]bool, len(queens))
	for i := range grid {
		grid[i] = make([]bool, len(queens))
	}
	for _, v := range queens {
		coords := strings.Split(string(v), ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		grid[x][y] = true
	}

	outStr := ""
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] {
				outStr += "[Q]"
			} else {
				outStr += "[ ]"
			}
		}
		outStr += "\n"
	}
	fmt.Println(outStr)
}
