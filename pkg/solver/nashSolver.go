package solver

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
	"time"
)

func NashSolver() {

	fmt.Println("Starting nash solver !")

	start := time.Now()

	// Ranges
	matrixOOP := constants.MatrixOOP
	matrixIP := constants.MatrixIp

	handsOOP := ranges.RangeToList(matrixOOP, 25)
	handsIP := ranges.RangeToList(matrixIP, 25)

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	board := constants.Board
	/*
		hero := constants.Hero
		heroPosition := constants.HeroPosition
	*/

	// Memoziation
	solvedHands := make(map[string]int)

	// ********** Solving here **********

	tree := tree.NewTree(
		constants.Pot,
		constants.EffectiveStack,
		handsOOP,
		handsIP,
		[]int{25, 60},
		[]int{25, 60},
		[]int{3000},
		[]int{3000},
	)

	tree.MakeRiverTree()

	fmt.Printf("Solving operation took %s\n", time.Since(start))

	fmt.Println("Checking inside root node")
	for _, value := range tree.Root.LocalActionMap {
		fmt.Println("hand : ", value)
		// for k, v := range tree.Root.LocalActionMap[k] {
		// 	v.ActionMap
		// }

	}
}
