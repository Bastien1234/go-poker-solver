package main

import (
	// "pokersolver/pkg/solver"
	"pokersolver/pkg/handSolver"

	// "pokersolver/pkg/utils"
	"fmt"
)

func main() {
	// solver.NashSolver()

	fmt.Println(handSolver.HandSolver([]string{"9h", "Kh", "Qh", "Jh", "Th", "2s", "2d"}))
	// ev := []int{500, 25, -1000, 100}
	// frequencies := []int{0, 0, 100, 0}
	// delta := 10

	// utils.UpdateFrenquencies(&ev, &frequencies, delta)

	// fmt.Println(frequencies)
}
