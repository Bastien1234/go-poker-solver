package main

import (
	"fmt"
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/poker"
	"pokersolver/pkg/ranges"
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

	lro := poker.GetLimitedRunOuts(100)

	handsOOP := ranges.RangeToVector(constants.MatrixOOP)
	handsIP := ranges.RangeToVector(constants.MatrixIp)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsOOP), func(i, j int) {
		handsOOP[i], handsOOP[j] = handsOOP[j], handsOOP[i]
	})
	rand.Shuffle(len(handsIP), func(i, j int) {
		handsIP[i], handsIP[j] = handsIP[j], handsIP[i]
	})

	handsOOP = handsOOP[0:constants.HandsToKeepFromRange]
	handsIP = handsIP[0:constants.HandsToKeepFromRange]

	vanillaCFR := poker.New()
	strategyMap := poker.NewStrategyMap()
	nIter := 100
	expectedValue := float32(0.0)

	for i := 0; i <= nIter; i++ {
		if i%10 == 0 {
			fmt.Printf("Starting iteration : %d\n", i)
			fmt.Println(time.Now())
		}
		poker := poker.NewGame(handsOOP, handsIP, lro)
		vanillaCFR.Run(*poker, strategyMap)

		strategyMap.UpdateStrategy()
	}

	expectedValue /= float32(nIter)
	fmt.Printf("Expected value is: %v\n", expectedValue)

	elapsed := time.Since(start)

	fmt.Printf("Did %d iterations in %s\n", nIter, elapsed)
	fmt.Println("The end, cross your fingers baby")
}
