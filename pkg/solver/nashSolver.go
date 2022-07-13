package solver

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/node"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/tree"
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

				// multi threading
				wg.Add(1)
				go func() {
					defer wg.Done()
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

									if subnodeAction == -3 { // fold
										valueOfAction += (float64(currentNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
									}

									if subnodeAction == -2 { // call
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

										if oopValue < ipValue {
											valueOfAction -= (float64(currentNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
										} else {
											valueOfAction += (float64(currentNode.PotSize) * float64(currentFrequency) * float64(currentHandFrenquency)) * divider
										}
									}
								}
							}
							divider += 1
							valueOfAction += 1
						}

					}
				}()

			}

			wg.Wait()

		}

		// 2 - Update frequencies

		// 3 - Pass hands to next nodes

		if iter%20 == 0 {
			fmt.Println("Iter : ", iter)
		}
	}

	fmt.Printf("Solving operation took %s\n", time.Since(start))
}
