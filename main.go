package main

import (
	"fmt"
	"pokersolver/pkg/cfr"
	"pokersolver/pkg/poker"
	"time"
	// "pokersolver/pkg/utils"
)

func main() {
	// solver.NashSolver()

	// million.SolveOneMillion(1000000, true)

	// ev := []int{500, 25, -1000, 100}
	// frequencies := []int{0, 0, 100, 0}
	// delta := 10

	// utils.UpdateFrenquencies(&ev, &frequencies, delta)

	// fmt.Println(frequencies)

	start := time.Now()
	go poker.RunDeckChannel()

	poker := poker.NewGame()
	vanillaCFR := cfr.New()
	nIter := 100000
	expectedValue := float32(0.0)
	for i := 1; i <= nIter; i++ {
		if i%10000 == 0 {
			fmt.Printf("Starting iteration : %d\n", i)
			expectedValue += vanillaCFR.Run(poker)
			for _, node := range vanillaCFR.Strategy {
				node.UpdateStrategy()
			}
		}
	}

	expectedValue /= float32(nIter)
	fmt.Printf("Expected value is: %v\n", expectedValue)

	elapsed := time.Since(start)

	fmt.Printf("Did %d iterations in %s\n", nIter, elapsed)
	/*
		seen := make(map[string]struct{})
		tree.Visit(poker, func(node cfr.GameTreeNode) {
			if node.Type() != cfr.PlayerNodeType {
				return
			}

			key := node.InfoSet(node.Player()).Key()
			if _, ok := seen[string(key)]; ok {
				return
			}

			actionProbs := policy.GetPolicy(node).GetAverageStrategy()
			if actionProbs != nil {
				fmt.Printf("[player %d] %6s: \n", node.Player(), key)
				for _, a := range actionProbs {
					fmt.Printf("%f\n", a)
				}
			}

			seen[string(key)] = struct{}{}
		})
	*/
}
