package solver

import (
	"pokersolver/pkg/constants"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/node"
	"sort"
)

func GetActionValue(divider float64, action int, subnodeActionIndex int, curNode *node.Node, currentSubnode *node.SubNode, playerIsIP bool, solvedHandsStruct *SolvedHandsStruct) float64 {
	var valueOfAction float64 = 0.0
	var futurDevider = 1 / float64((len(curNode.PostActionNodes[action].LocalActionMap)))
	var currentFrequency int = currentSubnode.Frequencies[subnodeActionIndex]
	var currentHandFrenquency int = currentSubnode.Weight

	if action == -3 {
		// Case Fold
		if playerIsIP {
			if curNode.PlayersTurn == "oop" {
				// Loosing invested money before folding
				valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			} else {
				valueOfAction += ((float64(curNode.PotSize)) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			}
		} else {
			if curNode.PlayersTurn == "ip" {
				// Loosing invested money before folding
				valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			} else {
				valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			}
		}

	} else if action == -2 || action == 0 { // call or check back
		// get who's winning...

		oopFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])
		ipFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])
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

		if oopValue < ipValue {
			// IP wins
			if playerIsIP == false {
				valueOfAction -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			} else {
				valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			}

		} else if oopValue > ipValue {
			// OOP wins
			if playerIsIP {
				valueOfAction -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			} else {
				valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
			}

		} else {
			// Equality
			return 0.0
		}

	} else if action == -1 {
		// Open check
		// start recursion stuff now
	}

	return valueOfAction

	// ---------------------------------------------------------------------------------------

	for _, nextAction := range curNode.Actions {
		// Check all possibilities

		

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

	return valueOfAction
}
