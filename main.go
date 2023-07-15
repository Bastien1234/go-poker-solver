package main

import (
	"fmt"
	"pokersolver/pkg/million"
	"pokersolver/pkg/poker"

	"github.com/timpalpant/go-cfr"
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

	poker := poker.NewGame()
	policy := cfr.NewPolicyTable(cfr.DiscountParams{UseRegretMatchingPlus: true})
	vanillaCFR := cfr.New(policy)
	nIter := 50000
	expectedValue := float32(0.0)
	for i := 1; i <= nIter; i++ {
		expectedValue += vanillaCFR.Run(poker)
	}

	expectedValue /= float32(nIter)
	fmt.Printf("Expected value is: %v\n", expectedValue)
}
