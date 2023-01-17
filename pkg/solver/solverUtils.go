package solver

import (
	"pokersolver/pkg/node"
)

func evRecursor(curNode *node.Node, devider float64, localAction int) float32 {

	/*
		**********************************************************************************************************
			Refactoring here
		**********************************************************************************************************

	*/

	var currentFrequency int = subnode.Frequencies[subnodeActionIndex]
	var currentHandFrenquency int = subnode.Weight


	// Should return current action value
	// Entry point is subnode action

	switch localAction {

	// Case of closing actions
	case -3:
		// Handle fold
		return handleFold(curNode, currentFrequency, currentHandFrenquency, devider)

	case 0, 2:
		// Handle call or check back
		return handleCallAndCheckBack(curNode, currentFrequency, currentHandFrenquency, devider)
	
	// Case of non closing action, passing recursion
	case -1:
		// Handle open check
		nextNode := 
		return evRecursor()
		

	default:
		// Handle raise or bet
		

	}






	return 0.0

	



	/*
		**********************************************************************************************************
			Refactoring here
		**********************************************************************************************************

	*/

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

func getSubnodeEv(currentSubnode node.SubNode, playerIsIP *bool) float32 {

	// Set up hero's position
	// See if necessary
	if currentNode.PlayersTurn == "root" || currentNode.PlayersTurn == "oop" {
		playerIsIP = false
	} else if currentNode.PlayersTurn == "ip" {
		playerIsIP = true
	} else {
		panic("WTF !!!")
	}

	// Iterate actions of the subnode
	for actionIdx, action := range currentSubnode.Actions {
		var valueOfAction float64 = 0.0
		var divider float64

		currentSubnode.Ev[actionIdx] = evRecursor()

		

		
	}

	
	return 0.0
}
