package solver

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/node"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
	"time"
)

func NashSolver() {

	// CHECK POT SIZE L'ENKULAY !!!

	// CHECK TO REMOVE DUPLICATES

	fmt.Println("Starting nash solver !")

	start := time.Now()

	// Ranges
	matrixOOP := constants.MatrixOOP
	matrixIP := constants.MatrixIp

	handsOOP := ranges.RangeToVector(matrixOOP)
	handsIP := ranges.RangeToVector(matrixIP)

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	// board := constants.Board
	/*
		hero := constants.Hero
		heroPosition := constants.HeroPosition
	*/

	// Memoziation
	// solvedHands := make(map[string]int)

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

	for iter := 0; iter < constants.NashIterations; iter++ {
		// 1 - Compute ev

		// BFS of nodes
		currentNode := tree.Root
		queue := []*node.Node{currentNode}
		for len(queue) > 0 {
			currentNode = queue[0]
			queue = queue[1:]

			for _, value := range currentNode.PostActionNodes {
				if value != nil {
					queue = append(queue, value)
				}
			}

			// Now get every value off the node

			// Create subnodes
			if iter == 0 {
				for _, hand := range currentNode.HandRange {
					name := hand.Cards[0] + hand.Cards[1]
					if currentNode.LocalActionMap[name] == nil {
						n := node.NewSubNode(hand.Cards, currentNode.Actions)
						currentNode.LocalActionMap[name] = &n
					}
				}
			}

			subnodesToVisit := []*node.SubNode{}
			for _, sn := range currentNode.LocalActionMap {
				subnodesToVisit = append(subnodesToVisit, sn)
			}

			for len(subnodesToVisit) > 0 {
				currentSubnode := subnodesToVisit[0]
				subnodesToVisit = subnodesToVisit[1:]

				fmt.Println(currentSubnode)

				// Calculation of ev logic comes here
				for index, action := range currentSubnode.Actions {
					// Keep fractions... or anyways call that some better way lol
					var ev float64 = 0.0

					// More trees traversal because that's what we do here
					



				}

			}

		}

		// 2 - Update frequencies

		// 3 - Pass hands to next nodes
	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))
}
