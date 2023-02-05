package solver

import (
	"pokersolver/pkg/constants"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/node"
	"sort"
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

func handleFold(curNode *node.Node, currentFrequency, currentHandFrenquency int, divider float64) float64 {
	var returnedValue float64 = 0.0

	/*
		if playerIsIP {
			if curNode.PlayersTurn == "oop" {
				// Loosing invested money before folding

				// ------------- One has to see if GTO just returns 0
				// -------------- Refactor to not have everything in float64 and avoid casting
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

	return returnedValue
}

func handleCallAndCheckBack(curNode *node.Node, currentSubnode *node.SubNode, action int) float64 {
	var returnedValue float64 = 0.0
	currentHandFrenquency := currentSubnode.Weight
	var actionIdx int
	for idx, el := range currentSubnode.Actions {
		if el == action {
			actionIdx = idx
		}
	}
	currentFrequency := currentSubnode.Frequencies[actionIdx]

	// get who's winning...

	oopFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])
	ipFinalHand := append(constants.Board, currentSubnode.Hand[0], currentSubnode.Hand[1])

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

	var playerIsIP bool
	if curNode.PlayersTurn == "ip" {
		playerIsIP = true
	} else {
		playerIsIP = false
	}

	if oopValue < ipValue {

		// IP wins
		if playerIsIP == false {
			returnedValue -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100))
		} else {
			returnedValue += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100))
		}

	} else if oopValue > ipValue {
		// OOP wins
		if playerIsIP {
			returnedValue -= (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100))
		} else {
			returnedValue += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100))
		}
	} else if oopValue == ipValue {
		// Equality, split dat pot babe
		returnedValue += ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100))
	}

	return returnedValue
}
