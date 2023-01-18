package solver 

import (
	"pokersolver/pkg/node"
	"sync"
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

var solvedHandsStruct = NewSolvedHandsStruct()

// ----------- divider !!!


func handleFold(curNode *node.Node, currentFrequency, currentHandFrenquency int, divider float32) float32 {
	var returnedValue float32 = 0.0

	/*
	if playerIsIP {
		if curNode.PlayersTurn == "oop" {
			// Loosing invested money before folding

			// ------------- One has to see if GTO just returns 0
			// -------------- Refactor to not have everything in float32 and avoid casting
			returnedValue -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			returnedValue += ((float64(curNode.PotSize)) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}
	} else {
		if curNode.PlayersTurn == "ip" {
			// Loosing invested money before folding
			returnedValue -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			returnedValue += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}
	}

	return returnedValue
	*/

	return 0.0
}

func handleCallAndCheckBack(curNode *node.Node, currentFrequency, currentHandFrenquency int, divider float32) float32 {
	var returnedValue float32 = 0.0

	// get who's winning...

	oopFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])
	ipFinalHand := append(constants.Board, subnode.Hand[0], subnode.Hand[1])

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
			returnedValue -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			returnedValue += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}

	} else if oopValue > ipValue {
		// OOP wins
		if playerIsIP {
			returnedValue -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			returnedValue += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}
	} else if oopValue == ipValue {
		// Equality, split dat pot babe
		returnedValue += ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
	}

	return returnedValue
}