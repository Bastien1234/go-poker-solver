package solver

import (
	"fmt"
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

// Memoziation
type SolvedHandsStruct struct {
	solvedHands map[string]int
	sync.Mutex
}

func NewSolvedHandsStruct() SolvedHandsStruct {
	s := SolvedHandsStruct{}
	s.solvedHands = make(map[string]int)

	return s
}

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

	handsOOP = handsOOP[0:constants.HandsToKeepFromRange]
	handsIP = handsIP[0:constants.HandsToKeepFromRange]

	fmt.Printf("IP player has : %v hands in his range\n", len(handsIP))

	/*
		hero := constants.Hero
		heroPosition := constants.HeroPosition
	*/

	solvedHandsStruct := NewSolvedHandsStruct()

	// ********** Solving here **********

	tree := tree.NewTree(
		constants.Pot,
		constants.EffectiveStack,
		handsOOP,
		handsIP,
		[]int{25, 60, 120},
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
				go func() {
					if currentNode.PlayersTurn == "root" || currentNode.PlayersTurn == "oop" {
						playerIsIP = false
					} else if currentNode.PlayersTurn == "ip" {
						playerIsIP = true
					} else {
						panic("WTF !!!")
					}
					defer wg.Done()
					for actionIdx, action := range currentSubnode.Actions {
						var valueOfAction float64 = 0.0
						var divider float64
						if currentNode.PostActionNodes[action] == nil {
							return
						}
						divider = 1 / float64((len(currentNode.PostActionNodes[action].LocalActionMap)))

						var evRecursor func(curNode *node.Node, devider float64, valueOfAction *float64, localAction int, plrPos *bool)

						evRecursor = func(curNode *node.Node, devider float64, valueOfAction *float64, localAction int, plrPos *bool) {

							if curNode.PostActionNodes[localAction] != nil {

								for _, subnode := range curNode.PostActionNodes[localAction].LocalActionMap {

									for subnodeActionIndex, subnodeAction := range subnode.Actions {
										// fmt.Println("Action : ", subnodeAction)
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
													*valueOfAction -= ((float64(curNode.PotSize) / 2) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												} else {
													*valueOfAction += ((float64(curNode.PotSize)) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												}
											} else {
												if curNode.PlayersTurn == "ip" {
													// Loosing invested money before folding
													*valueOfAction -= ((float64(curNode.PotSize) / 2) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												}
											}

										} else if subnodeAction == -2 { // call
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
												solvedHandsStruct.solvedHands[oopFinalHandString] = handSolver.HandSolver(oopFinalHand)
												oopValue = solvedHandsStruct.solvedHands[oopFinalHandString]
											}

											if val, ok := solvedHandsStruct.solvedHands[ipFinalHandString]; ok {
												ipValue = val
											} else {
												solvedHandsStruct.solvedHands[ipFinalHandString] = handSolver.HandSolver(ipFinalHand)
												ipValue = solvedHandsStruct.solvedHands[ipFinalHandString]
											}

											solvedHandsStruct.Unlock()

											// what if equality ?? lol

											if oopValue < ipValue {
												// IP wins
												if playerIsIP == false {
													*valueOfAction -= (float64(curNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												}

											} else if oopValue > ipValue {
												// OOP wins
												if playerIsIP {
													*valueOfAction -= (float64(curNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												} else {
													*valueOfAction += (float64(curNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
												}

											}

										} else if subnodeAction == -1 {
											// Open check
											nodeCopy := curNode.PostActionNodes[action].PostActionNodes[-1]
											playerIsIP = !playerIsIP
											evRecursor(nodeCopy, divider, valueOfAction, -1, &playerIsIP)
										} else if subnodeAction == 0 {
											// get winner from check back
										}

									}

									nextActions := make([]int, 0)

									for _, subnodeAction := range subnode.Actions {
										if subnodeAction != 0 && subnodeAction != -2 && subnodeAction != -3 {
											nextActions = append(nextActions, subnodeAction)
										}
									}

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

						evRecursor(currentNode, divider, &valueOfAction, action, &playerIsIP)
						currentSubnode.Ev[actionIdx] = int(valueOfAction)

					}

				}()

			}

			wg.Wait()

		}

		// 2 - Update frequencies

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

			subnodesToVisit := []*node.SubNode{}
			for _, sn := range currentNodeUpdate.LocalActionMap {
				subnodesToVisit = append(subnodesToVisit, sn)
			}

			for len(subnodesToVisit) > 0 {
				currentSubnode := subnodesToVisit[0]
				subnodesToVisit = subnodesToVisit[1:]

				// Calculation of ev logic comes here
				// fmt.Printf("In subnode : ", currentSubnode.Hand)
				utils.UpdateFrenquencies(currentSubnode.Ev, &currentSubnode.Frequencies, constants.Delta)
			}
		}

		// 3 - Pass hands to next nodes

		fmt.Println("Iter : ", iter)

	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))
}
