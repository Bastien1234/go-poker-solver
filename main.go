package main

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
)

func main() {

	fmt.Println("Program started !")

	// Ranges
	matrixOOP := constants.MatrixOOP
	matrixIP := constants.MatrixIp

	handsOOP := ranges.RangeToList(matrixOOP, 10)
	handsIP := ranges.RangeToList(matrixIP, 10)

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	// Memoziation
	// solvedHands := make(map[string]int)

	// ********** Solving here **********

	tree := tree.NewTree(
		constants.Pot,
		constants.EffectiveStack,
		handsOOP,
		handsIP,
		[]int{23, 60, 120},
		[]int{2500, 3500},
		[]int{33, 60, 125},
		[]int{2700, 3200},
	)

	tree.MakeRiverTree()

	fmt.Println("program finished with code 0")

	fmt.Print(tree.Root.PostActionNodes)
}
