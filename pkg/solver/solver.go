package solver

import (
	"fmt"
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/node"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
	"pokersolver/pkg/utils"
	"time"
)

func Solve() {

	fmt.Println("Program started !")

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

	// After a while, players get better, and rewards must be more importants
	var multiplicator float32 = 1.0

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

			// Get hand value of players
			var OOPValue int = -1
			var IPValue int = -1
			stringHandOOP := OOPHand[0] + OOPHand[1]
			stringHandIP := IPHand[0] + IPHand[1]
			boardOOP := append(board, OOPHand...)
			boardIP := append(board, IPHand...)

			if val, ok := solvedHands[stringHandOOP]; ok {
				OOPValue = val
			} else {
				OOPValue = handSolver.HandSolver(boardOOP)
				solvedHands[stringHandOOP] = OOPValue
			}

			if val, ok := solvedHands[stringHandIP]; ok {
				IPValue = val
			} else {
				IPValue = handSolver.HandSolver(boardIP)
				solvedHands[stringHandIP] = IPValue
			}

			traversalOOPScore := 0
			traversalIPScore := 0

			random2 := utils.Generate2()
			random3 := utils.Generate3()
			random4 := utils.Generate4()
			random5 := utils.Generate5()

			var optimalPlayer string

			for iter2 := 0; iter2 < constants.Iterations2; iter2++ {
				vectorActions := []int{}
				nodesToVisit := []node.Node{*tree.Root}

				if iter1%2 == 0 {
					optimalPlayer = "oop"
				} else {
					optimalPlayer = "ip"
				}

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

					// Use this new subnode only if it doesn't exist already, see below
					currentSubnode := node.NewSubNode(currentHand, currentNode.Actions)

					// Populate subnodes hashmap
					if currentNode.PlayersTurn == "ip" {
						if !utils.ContainsLocalActionMap(currentNode.LocalActionMap, stringHandIP) {
							currentNode.LocalActionMap[stringHandIP] = &currentSubnode
						} else {
							currentSubnode = *currentNode.LocalActionMap[stringHandIP]
						}
					} else if currentNode.PlayersTurn == "oop" {
						if !utils.ContainsLocalActionMap(currentNode.LocalActionMap, stringHandOOP) {
							currentNode.LocalActionMap[stringHandOOP] = &currentSubnode
						} else {
							currentSubnode = *currentNode.LocalActionMap[stringHandOOP]
						}
					}

					// Poppingaments
					nodesToVisit = nodesToVisit[1:]
					actionDistribution := []int{}

					var numberOfPossibleActions = len(currentSubnode.Actions)

					// If not playing optimal : try new frequencies, otherwise we'll use current best frequencies
					if currentNode.PlayersTurn != optimalPlayer {
						if numberOfPossibleActions == 2 {
							actionDistribution = random2
						} else if numberOfPossibleActions == 3 {
							actionDistribution = random3
						} else if numberOfPossibleActions == 4 {
							actionDistribution = random4
						} else if numberOfPossibleActions == 5 {
							actionDistribution = random5
						}
					} else {
						// actionDistribution = currentSubnode.Actions
						for _, el := range currentSubnode.Actions {
							actionDistribution = append(actionDistribution, currentSubnode.ActionMap[el])
						}
					}

					// Random choice of action
					randomChoice := rand.Intn(99)
					action := currentSubnode.Actions[utils.PickIndexFromVector(actionDistribution, randomChoice)]

					// Open check
					if action == -1 {
						// Going to the next node
						vectorActions = append(vectorActions, action)
						nodesToVisit = append(nodesToVisit, *currentNode.PostActionNodes[action])
					} else if action == -2 {
						// Call
						vectorActions = append(vectorActions, action)

						// We need to populate ???????

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
				bfsCurrentNode := tree.Root
				data := [][]int{}
				bfsQueue := []*node.Node{bfsCurrentNode}

				for len(bfsQueue) > 0 {
					bfsCurrentNode = bfsQueue[0]
					bfsQueue = bfsQueue[1:]
					data = append(data, bfsCurrentNode.Actions)

					for _, value := range bfsCurrentNode.PostActionNodes {
						if value != nil {
							bfsQueue = append(bfsQueue, value)
						}
					}

					actionDistribution := []int{}

					var numberOfPossibleActions = len(bfsCurrentNode.Actions)

					if bfsCurrentNode.PlayersTurn != optimalPlayer {
						if numberOfPossibleActions == 2 {
							actionDistribution = random2
						} else if numberOfPossibleActions == 3 {
							actionDistribution = random3
						} else if numberOfPossibleActions == 4 {
							actionDistribution = random4
						} else if numberOfPossibleActions == 5 {
							actionDistribution = random5
						}
					}

					// Compare scores
					if bfsCurrentNode.PlayersTurn == "ip" && optimalPlayer != "ip" {
						if bfsCurrentNode.LocalActionMap[stringHandIP] != nil {
							previousBestScore := bfsCurrentNode.LocalActionMap[stringHandIP].BestScore
							if traversalIPScore > previousBestScore {
								// Update
								bfsCurrentNode.LocalActionMap[stringHandIP].BestScore = int(float32(traversalIPScore) * multiplicator)
								// Updating actions
								for idx, el := range actionDistribution {
									bfsCurrentNode.LocalActionMap[stringHandIP].ActionMap[bfsCurrentNode.LocalActionMap[stringHandIP].Actions[idx]] = el
								}

							}
						}

					} else if bfsCurrentNode.PlayersTurn == "oop" && optimalPlayer != "oop" {
						if bfsCurrentNode.LocalActionMap[stringHandOOP] != nil {
							previousBestScore := bfsCurrentNode.LocalActionMap[stringHandOOP].BestScore
							if traversalOOPScore > previousBestScore {
								// Update
								bfsCurrentNode.LocalActionMap[stringHandOOP].BestScore = int(float32(traversalOOPScore) * multiplicator)
								// Updating actions
								for idx, el := range actionDistribution {
									bfsCurrentNode.LocalActionMap[stringHandOOP].ActionMap[bfsCurrentNode.LocalActionMap[stringHandOOP].Actions[idx]] = el
								}
							}
						}

					}
				}

			} // iter2

		} else {
			iter1++
		}

		if iter1%1000 == 0 {
			fmt.Printf("Did %v iterations\n", iter1)

			multiplicator *= constants.MultiplicatorImprovement
		}

	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))

	fmt.Println("Checking inside root node")
	for _, value := range tree.Root.LocalActionMap {
		fmt.Println("hand : ", value)
		// for k, v := range tree.Root.LocalActionMap[k] {
		// 	v.ActionMap
		// }

	}
}
