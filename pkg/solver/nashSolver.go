package solver

import (
	"fmt"
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/node"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
	"pokersolver/pkg/utils"
	"time"
)

func NashSolver() {
	// var wg sync.WaitGroup

	// CHECK POT SIZE L'ENKULAY !!!

	// CHECK TO REMOVE DUPLICATES

	fmt.Println("Starting nash solver !")

	start := time.Now()

	// Ranges
	matrixOOP := constants.MatrixOOP
	matrixIP := constants.MatrixIp

	handsOOP := ranges.RangeToVector(matrixOOP)
	handsIP := ranges.RangeToVector(matrixIP)

	// Shuffle arrays
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsOOP), func(i, j int) { handsOOP[i], handsOOP[j] = handsOOP[j], handsOOP[i] })
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsIP), func(i, j int) { handsIP[i], handsIP[j] = handsIP[j], handsIP[i] })

	// Eliminate hands already on the board

	filteredOOP := make([]ranges.Hand, 0)
	filteredIP := make([]ranges.Hand, 0)

	for _, hand := range handsOOP {
		if hand.Cards[0] != constants.Board[0] && hand.Cards[1] != constants.Board[0] && hand.Cards[0] != constants.Board[1] && hand.Cards[1] != constants.Board[1] && hand.Cards[0] != constants.Board[2] && hand.Cards[1] != constants.Board[2] && hand.Cards[0] != constants.Board[3] && hand.Cards[1] != constants.Board[3] && hand.Cards[0] != constants.Board[4] && hand.Cards[1] != constants.Board[4] {
			filteredOOP = append(filteredOOP, hand)
		}
	}

	for _, hand := range handsIP {
		if hand.Cards[0] != constants.Board[0] && hand.Cards[1] != constants.Board[0] && hand.Cards[0] != constants.Board[1] && hand.Cards[1] != constants.Board[1] && hand.Cards[0] != constants.Board[2] && hand.Cards[1] != constants.Board[2] && hand.Cards[0] != constants.Board[3] && hand.Cards[1] != constants.Board[3] && hand.Cards[0] != constants.Board[4] && hand.Cards[1] != constants.Board[4] {
			filteredIP = append(filteredIP, hand)
		}
	}

	handsOOP = filteredOOP
	handsIP = filteredIP

	handsOOP = handsOOP[0:constants.HandsToKeepFromRange]
	handsIP = handsIP[0:constants.HandsToKeepFromRange]

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	// Making sure hero's hand is in the game

	// hero := constants.Hero
	// heroPosition := constants.HeroPosition

	// if heroPosition == "oop" {
	// 	handsOOP = append(handsOOP, hero)
	// }

	// ********** Solving here **********

	tree := tree.NewTree(
		constants.Pot,
		constants.EffectiveStack,
		handsOOP,
		handsIP,
		[]int{25, 60, 200},
		[]int{25, 60, 130},
		[]int{3000, 4500},
		[]int{3000, 4500},
	)

	tree.MakeRiverTree()

	var playerIsIP bool

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

				currentSubnode.Ev = GetSubnodeEv(currentSubnode, &playerIsIP)

			}
		}

		// 2 - Update frequencies

		// Multithreading please !

		currentNodeUpdate := tree.Root
		queueUpdate := []*node.Node{currentNodeUpdate}
		for len(queueUpdate) > 0 {
			currentNodeUpdate = queueUpdate[0]
			queueUpdate = queueUpdate[1:]

			for _, value := range currentNodeUpdate.PostActionNodes {
				if value != nil {
					queueUpdate = append(queueUpdate, value)
				}
			}

			for _, currentSubnode := range currentNodeUpdate.LocalActionMap {
				// debug

				/*
					fmt.Println("\n\nCurrent node : ", currentNodeUpdate.RaiseLevel)
					fmt.Println("Ev : ", currentSubnode.Ev)
					fmt.Println("Frequencies : ", currentSubnode.Frequencies)
					fmt.Println("Actions : ", currentSubnode.Actions)
					fmt.Println("Hand : ", currentSubnode.Hand)
				*/

				utils.UpdateFrenquencies(&currentSubnode.Ev, &currentSubnode.Frequencies, constants.Delta)
				/*
					fmt.Println("New freq : ", currentSubnode.Frequencies)
				*/

				// Debug
				currentSubnode.Updated = true

			}

			/*
				subnodesToVisit := []*node.SubNode{}
				for _, sn := range currentNodeUpdate.LocalActionMap {
					subnodesToVisit = append(subnodesToVisit, sn)
				}

				for len(subnodesToVisit) > 0 {
					currentSubnode := subnodesToVisit[0]
					subnodesToVisit = subnodesToVisit[1:]

					// Calculation of ev logic comes here
					// fmt.Printf("In subnode : ", currentSubnode.Hand)
					utils.UpdateFrenquencies(&currentSubnode.Ev, &currentSubnode.Frequencies, constants.Delta)
				}
			*/
		}

		fmt.Println("Iter : ", iter)

	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))
}
