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
			lenstv := 0

			subnodesToVisit := []*node.SubNode{}
			for _, sn := range currentNode.LocalActionMap {
				subnodesToVisit = append(subnodesToVisit, sn)
				lenstv++
			}

			for len(subnodesToVisit) > 0 {
				currentSubnode := subnodesToVisit[0]
				subnodesToVisit = subnodesToVisit[1:]

				// Calculation of ev logic comes here
				for _, action := range currentSubnode.Actions {
					var valueOfAction float64 = 0.0
					var divider float64
					if currentNode.PostActionNodes[action] != nil {
						divider = 1 / float64((len(currentNode.PostActionNodes[action].LocalActionMap)))

						for _, subnode := range currentNode.PostActionNodes[action].LocalActionMap {
							for subnodeActionIndex, subnodeAction := range subnode.Actions {
								var currentFrequency int = subnode.Frequencies[subnodeActionIndex]
								var currentHandFrenquency int = subnode.Weight

								currentFrequency++
								currentHandFrenquency++

								// for idx, n := range subnode.Actions {
								// 	if n == subnodeAction {
								// 		currentFrequency = subnode.Frequencies[idx]
								// 	}
								// }

								// for _, n := range currentNode.HandRange {
								// 	card := n.Cards[0] + n.Cards[1]
								// 	subNodeCard := subnode.Hand[0] + subnode.Hand[1]
								// 	if card == subNodeCard {
								// 		currentHandFrenquency = n.Frequency
								// 	}
								// }

								if subnodeAction == -3 { // fold
									valueOfAction += (float64(currentNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
								}
							}
						}
						divider += 1
						valueOfAction += 1
					}

				}

			}

		}

		// 2 - Update frequencies

		// 3 - Pass hands to next nodes

		// fmt.Println("Iter : ", iter)
	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))
}
