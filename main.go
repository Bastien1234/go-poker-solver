package main

import (
	"pokersolver/pkg/million"
	// "pokersolver/pkg/utils"
)

func main() {
	// solver.NashSolver()

	million.SolveOneMillion(1000000, true)

	// ev := []int{500, 25, -1000, 100}
	// frequencies := []int{0, 0, 100, 0}
	// delta := 10

	// utils.UpdateFrenquencies(&ev, &frequencies, delta)

	// fmt.Println(frequencies)
}
