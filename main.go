package main

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

func main() {

	fmt.Println("Program started !")

	start := time.Now()

	// Ranges
	matrixOOP := constants.MatrixOOP
	matrixIP := constants.MatrixIp

	handsOOP := ranges.RangeToList(matrixOOP, 10)
	handsIP := ranges.RangeToList(matrixIP, 10)

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	board := constants.Board
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
		[]int{23, 60, 120},
		[]int{2500, 3500},
		[]int{33, 60, 125},
		[]int{2700, 3200},
	)

	tree.MakeRiverTree()

	for iter1 := 0; iter1 < constants.Iterations1; iter1++ {
		/*
			Choice of hands for both players
			Excluding the cards that are already on the board and onto oponent's hand
		*/

		randomIdxOOP := rand.Intn(len(handsOOP))
		randomIdxIP := rand.Intn(len(handsIP))
		OOPHand := handsOOP[randomIdxOOP]
		IPHand := handsIP[randomIdxIP]

		duplicatesOOP1 := utils.Contains(board, OOPHand[0])
		duplicatesOOP2 := utils.Contains(board, OOPHand[1])
		duplicatesOOP3 := utils.Contains(IPHand, OOPHand[0])
		duplicatesOOP4 := utils.Contains(IPHand, OOPHand[1])
		duplicatesIP1 := utils.Contains(board, IPHand[0])
		duplicatesIP2 := utils.Contains(board, IPHand[1])
		duplicatesIP3 := utils.Contains(OOPHand, IPHand[0])
		duplicatesIP4 := utils.Contains(OOPHand, IPHand[1])

		// Check if no duplicates, otherwise chose new cards
		if duplicatesOOP1 == false && duplicatesOOP2 == false && duplicatesOOP3 == false && duplicatesOOP4 == false && duplicatesIP1 == false && duplicatesIP2 == false && duplicatesIP3 == false && duplicatesIP4 == false {

			// Get hand value of players ... when I'll code that
			OOPValue := 9
			IPValue := 5
			stringHandOOP := OOPHand[0] + OOPHand[1]
			stringHandIP := IPHand[0] + IPHand[1]

			traversalOOPScore := 0
			traversalIPScore := 0
			/*
				random2 := utils.Generate2()
				random3 := utils.Generate3()
				random4 := utils.Generate4()
			*/

			for iter2 := 0; iter2 < constants.Iterations2; iter2++ {
				vectorActions := []int{}
				nodesToVisit := []node.Node{*tree.Root}

				// You like BFS, do you ? =)
				for len(nodesToVisit) > 0 {
					// Let's call this one Traversal

					currentNode := nodesToVisit[0]
					currentHand := make([]string, 2)

					if currentNode.PlayersTurn == "oop" {
						currentHand = OOPHand
					} else if currentNode.PlayersTurn == "ip" {
						currentHand = IPHand
					}

					currentSubnode := node.NewSubNode(currentHand, currentNode.Actions)

					// Populate subnodes hashmap
					if currentNode.PlayersTurn == "ip" {
						if !utils.ContainsLocalActionMap(currentNode.LocalActionMap, stringHandIP) {
							currentNode.LocalActionMap[stringHandIP] = &currentSubnode
						}
					} else if currentNode.PlayersTurn == "oop" {
						if !utils.ContainsLocalActionMap(currentNode.LocalActionMap, stringHandOOP) {
							currentNode.LocalActionMap[stringHandOOP] = &currentSubnode
						}
					}

					// Poppingaments
					nodesToVisit = nodesToVisit[1:]
					actionDistribution := []int{}

					var numberOfPossibleActions = len(currentSubnode.Actions)

					if numberOfPossibleActions == 2 {
						actionDistribution = utils.Generate2()
					} else if numberOfPossibleActions == 3 {
						actionDistribution = utils.Generate3()
					} else if numberOfPossibleActions == 4 {
						actionDistribution = utils.Generate4()
					}

					fmt.Println(numberOfPossibleActions, " ", actionDistribution)

					// Random choice of action
					rdmIdx := rand.Intn(numberOfPossibleActions)
					action := currentNode.Actions[rdmIdx]

					// Open check
					if action == -1 {
						// Going to the next node
						vectorActions = append(vectorActions, action)
						nodesToVisit = append(nodesToVisit, *currentNode.PostActionNodes[action])
					} else if action == -2 {
						// Call
						vectorActions = append(vectorActions, action)

						// Showdown baby
						if OOPValue > IPValue {
							traversalOOPScore += currentNode.CurrentFacingBet
							traversalOOPScore += currentNode.PotSize
							traversalIPScore -= currentNode.CurrentFacingBet
						} else if OOPValue < IPValue {
							traversalIPScore += currentNode.CurrentFacingBet
							traversalIPScore += currentNode.PotSize
							traversalOOPScore -= currentNode.CurrentFacingBet
						} else {
							// Nobody wins, don't move scores
						}
					} else if action == 3 {
						// Foldaments
						vectorActions = append(vectorActions, action)

						if currentNode.PlayersTurn == "oop" {
							traversalIPScore += currentNode.PotSize
						} else {
							traversalOOPScore += currentNode.PotSize
						}
					} else if action == 0 {
						vectorActions = append(vectorActions, action)

						if OOPValue > IPValue {
							// OOP wins
							// But you guessed it, right ? =)
							traversalOOPScore += currentNode.PotSize
						} else {
							traversalIPScore += currentNode.PotSize
						}
					} else if action > 0 {
						// Bet or raise, just go to next node
						vectorActions = append(vectorActions, action)
						nodesToVisit = append(nodesToVisit, *currentNode.PostActionNodes[action])
					}

				} // Traversal

				// BFS Updater
				// Yes we fucking love BFSs !!! <3 <3 <3
				currentNode := tree.Root
				data := [][]int{}
				queue := []*node.Node{currentNode}

				for len(queue) > 0 {
					currentNode = queue[0]
					queue = queue[1:]
					data = append(data, currentNode.Actions)

					for _, value := range currentNode.PostActionNodes {
						if value != nil {
							queue = append(queue, value)
						}
					}

					// Compare scores
					if currentNode.PlayersTurn == "ip" {
						if currentNode.LocalActionMap[stringHandIP] != nil {
							previousBestScore := currentNode.LocalActionMap[stringHandIP].BestScore
							if traversalIPScore > previousBestScore {
								// Update
								currentNode.LocalActionMap[stringHandIP].BestScore = traversalIPScore
							}
						}

					} else {
						if currentNode.LocalActionMap[stringHandOOP] != nil {
							previousBestScore := currentNode.LocalActionMap[stringHandOOP].BestScore
							if traversalOOPScore > previousBestScore {
								// Update
								currentNode.LocalActionMap[stringHandOOP].BestScore = traversalOOPScore
							}
						}

					}
				}

			} // iter2

		} else {
			iter1++
		}

	}

	fmt.Printf("Solving operation took %s", time.Since(start))

	fmt.Println("Checking inside root node")
	for _, value := range tree.Root.LocalActionMap {
		fmt.Println("hand : ", value)
		// for k, v := range tree.Root.LocalActionMap[k] {
		// 	v.ActionMap
		// }

	}

	fmt.Println("Program finished with code 0")
}
