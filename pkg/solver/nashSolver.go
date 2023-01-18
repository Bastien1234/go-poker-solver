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
	"sort"
	"sync"
	"time"
)



func NashSolver() {
	var wg sync.WaitGroup

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

				// Calculation of ev logic comes here

				// multi threading
				wg.Add(1)

				// ************************************************************************************************************************************
				/*

					Refactoring from here

				*/
				// ********************************************************************************************************************************

				currentSubnode.Ev = getSubnodeEv(currentSubnode, &playerIsIP)

				// See how to handle multithreading

				go func() {
					
					defer wg.Done()
					for actionIdx, action := range currentSubnode.Actions {

						// Get action value
						
						if currentNode.PostActionNodes[action] == nil {
							continue
						}
						divider = 1 / float64((len(currentNode.PostActionNodes[action].LocalActionMap)))

						var evRecursor func(curNode *node.Node, divider float64, valueOfAction *float64, localAction int, plrPos *bool)

						evRecursor = func(curNode *node.Node, divider float64, valueOfAction *float64, localAction int, plrPos *bool) {

							if curNode.PostActionNodes[localAction] != nil {

								for _, subnode := range curNode.PostActionNodes[localAction].LocalActionMap {

									for subnodeActionIndex, subnodeAction := range subnode.Actions {
										var currentFrequency int = subnode.Frequencies[subnodeActionIndex]
										var currentHandFrenquency int = subnode.Weight

										// *** Base cases *** //

										// Open check

										// Check back !!!!!!

										// open check !!!!

										if subnodeAction == -3 { // fold
											if playerIsIP {
												if curNode.PlayersTurn == "oop" {
													// Loosing invested money before folding
													*valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												} else {
													*valueOfAction += ((float64(curNode.PotSize)) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												}
											} else {
												if curNode.PlayersTurn == "ip" {
													// Loosing invested money before folding
													*valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												}
											}

										} else if subnodeAction == -2 || subnodeAction == 0 { // call or check back
											// get who's winning...

											oopFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])
											ipFinalHand := append(constants.Board, subnode.Hand[0], subnode.Hand[1])
											// Check efficiency of sorting please...
											sort.Strings(oopFinalHand)
											sort.Strings(ipFinalHand)
											oopFinalHandString := ""
											ipFinalHandString := ""
											for i := 0; i < 7; i++ {
												oopFinalHandString += oopFinalHand[i]
												ipFinalHandString += ipFinalHand[i]
											}

											var oopValue int
											var ipValue int

											solvedHandsStruct.Lock()

											if val, ok := solvedHandsStruct.solvedHands[oopFinalHandString]; ok {
												oopValue = val
											} else {
												solvedHandsStruct.solvedHands[oopFinalHandString] = handSolver.HandSolver(oopFinalHand, false)
												oopValue = solvedHandsStruct.solvedHands[oopFinalHandString]
											}

											if val, ok := solvedHandsStruct.solvedHands[ipFinalHandString]; ok {
												ipValue = val
											} else {
												solvedHandsStruct.solvedHands[ipFinalHandString] = handSolver.HandSolver(ipFinalHand, false)
												ipValue = solvedHandsStruct.solvedHands[ipFinalHandString]
											}

											solvedHandsStruct.Unlock()

											// what if equality ?? lol

											if oopValue < ipValue {
												// IP wins
												if playerIsIP == false {
													*valueOfAction -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												}

											} else if oopValue > ipValue {
												// OOP wins
												if playerIsIP {
													*valueOfAction -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
												}

											}

										} else if subnodeAction == -1 {
											// Open check
											nodeCopy := curNode.PostActionNodes[action].PostActionNodes[-1]
											playerIsIP = !playerIsIP
											evRecursor(nodeCopy, divider, valueOfAction, -1, &playerIsIP)
										} else {
											nextActions := make([]int, 0)

											// Getting next actions
											// for _, possibleAction := range curNode
											nextActions = append(nextActions, subnodeAction)
											var newDivider float64 = divider / float64(len(nextActions))

											if len(nextActions) > 0 {
												for _, actToPass := range nextActions {
													nodeCopy := curNode.PostActionNodes[action].PostActionNodes[actToPass]
													playerIsIP = !playerIsIP
													evRecursor(nodeCopy, newDivider, valueOfAction, actToPass, &playerIsIP)
												}
											}
										}
									}
								}
							}
						}

						evRecursor(currentNode, divider, &valueOfAction, action, &playerIsIP)
						currentSubnode.Ev[actionIdx] = int(valueOfAction)
					}
				}()
			}

			wg.Wait()
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
