package handSolver

import (
	"pokersolver/pkg/utils"
)

func HandSolver2(arr []string, safe bool) int {

	// Constants

	cards := []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	colors := []string{"h", "d", "c", "s"}

	cardsValueMap := map[string]int{
		"A": 14, // Value 1 as well !
		"K": 13,
		"Q": 12,
		"J": 11,
		"T": 10,
		"9": 9,
		"8": 8,
		"7": 7,
		"6": 6,
		"5": 5,
		"4": 4,
		"3": 3,
		"2": 2,
	}

	bitMaskStraightValues := map[int]uint16{
		14: 31744,
		13: 15872,
		12: 7936,
		11: 3968,
		10: 1984,
		9:  992,
		8:  496,
		7:  248,
		6:  124,
		5:  62,
	}

	// Checking input in safe mode

	if safe {
		// Correct length
		if len(arr) != 7 {
			return -1
		}

		// Duplicates
		newSlice := make([]string, 0)
		for _, el := range arr {
			if !utils.Contains(newSlice, el) {
				newSlice = append(newSlice, el)
			}
		}
		if len(newSlice) != 7 {
			return -1
		}

		// Correct inputs
		for _, el := range arr {
			if len(el) != 2 {
				return -1
			}

			var card string = el[0:1]
			var color string = el[1:]

			if !utils.Contains(cards, card) {
				return -1
			}

			if !utils.Contains(colors, color) {
				return -1
			}
		}
	}

	var stateStraight uint16 = 0
	var stateHearts uint16 = 0
	var stateDiamonds uint16 = 0
	var stateClubs uint16 = 0
	var stateSpades uint16 = 0

	fourOfAKind := -1
	threeOfAKind := -1
	bestPair := -1
	secondPair := -1

	possibleColor := false
	stateColor := ""

	straightValue := -1

	// 0: hearts, 1: diamonds, 2: clubs, 3: spades
	colorsCount := utils.FilledArrayInt(4, 0)
	valuesCount := utils.FilledArrayInt(14, 0)

	for _, element := range arr {

		card := element[0:1]
		color := element[1:2]

		switch color {
		case "h":
			colorsCount[0] += 1
			stateHearts |= 1 << cardsValueMap[card]
			if card == "A" {
				stateHearts |= 1 << 1
			}

		case "d":
			colorsCount[1] += 1
			stateDiamonds |= 1 << cardsValueMap[card]
			if card == "A" {
				stateDiamonds |= 1 << 1
			}

		case "c":
			colorsCount[2] += 1
			stateClubs |= 1 << cardsValueMap[card]
			if card == "A" {
				stateClubs |= 1 << 1
			}

		case "s":
			colorsCount[3] += 1
			stateSpades |= 1 << cardsValueMap[card]
			if card == "A" {
				stateSpades |= 1 << 1
			}

		default:
			panic("Nooooooooo")
		}

		stateStraight |= 1 << cardsValueMap[card]
		if card == "A" {
			stateStraight |= 1 << 1
		}

		valuesCount[cardsValueMap[card]-1] += 1
		if card == "A" {
			valuesCount[0] += 1
		}

	}

	for i := range colorsCount {
		if colorsCount[i] >= 5 {
			possibleColor = true

			switch i {
			case 0:
				stateColor = "h"
			case 1:
				stateColor = "d"
			case 2:
				stateColor = "c"
			case 3:
				stateColor = "s"

			default:
				panic("Bad colors assertion")
			}
		}
	}

	for index, el := range valuesCount {
		if el == 4 {
			fourOfAKind = index + 1
		}

		if el == 3 {
			if threeOfAKind != -1 {
				if el > threeOfAKind {
					threeOfAKind = index + 1
				}
			} else {
				threeOfAKind = index + 1
			}
		}

		if el == 2 {
			// If no current pair
			if bestPair == -1 && secondPair == -1 {
				bestPair = index + 1
			} else if bestPair != -1 && secondPair == -1 {
				currentBest := bestPair
				if el > currentBest {
					bestPair = index + 1
					secondPair = currentBest
				} else {
					secondPair = index + 1
				}
			} else {

				// If we already have two pairs
				currentBest := bestPair
				currentSecondBest := secondPair

				// Candidate pair is the best
				if el > currentBest {
					bestPair = index + 1
					secondPair = currentBest
				} else if el > currentSecondBest {
					secondPair = index + 1
				}

			}
		}

	}

	// Check if possible Straight
	for value, bitMask := range bitMaskStraightValues {
		if (stateStraight & bitMask) == bitMask {
			if value > straightValue {
				straightValue = value
			}
		}
	}

	/*
	   Check if possible straight flush
	   Return 9 billions then value of high card
	*/

	if possibleColor && straightValue > 0 {
		height := -1
		switch stateColor {
		case "h":
			for value, bitMask := range bitMaskStraightValues {
				if (stateHearts & bitMask) == bitMask {
					if value > height {
						height = value
					}
				}
			}

		case "d":
			for value, bitMask := range bitMaskStraightValues {
				if (stateDiamonds & bitMask) == bitMask {
					if value > height {
						height = value
					}
				}
			}

		case "c":
			for value, bitMask := range bitMaskStraightValues {
				if (stateClubs & bitMask) == bitMask {
					if value > height {
						height = value
					}
				}
			}

		case "s":
			for value, bitMask := range bitMaskStraightValues {
				if (stateSpades & bitMask) == bitMask {
					if value > height {
						height = value
					}
				}
			}
		}

		if height > 0 {
			return 900_000_000_000 + (height * 1e9)
		}
	}

	/*
	   Check is possible four of a kind
	   returns 8 billions then the FOAK value and finally the kicker
	*/

	if fourOfAKind > 0 {
		for index := 14; index > 0; index-- {
			if (stateStraight&(1<<index) != 0) && (fourOfAKind != index) {
				return 800_000_000_000 + (fourOfAKind * 1e9) + (index * 1e7)
			}
		}
	}

	/*
	   Check is possible full house
	*/

	if threeOfAKind > 0 && bestPair > 0 {
		return 700_000_000_000 + (threeOfAKind * 1e9) + (bestPair * 1e7)
	}

	/*
	   Check if possible color
	*/

	if possibleColor {
		values := []int{}
		totalValues := 0
		switch stateColor {
		case "h":
			for index := 14; index > 0; index-- {
				if stateHearts&(1<<index) != 0 {
					values = append(values, index)
					totalValues += 1
					if totalValues == 5 {
						break
					}
				}
			}

		case "d":
			for index := 14; index > 0; index-- {
				if stateDiamonds&(1<<index) != 0 {
					values = append(values, index)
					totalValues += 1
					if totalValues == 5 {
						break
					}
				}
			}

		case "c":
			for index := 14; index > 0; index-- {
				if stateClubs&(1<<index) != 0 {
					values = append(values, index)
					totalValues += 1
					if totalValues == 5 {
						break
					}
				}
			}

		case "s":
			for index := 14; index > 0; index-- {
				if stateSpades&(1<<index) != 0 {
					values = append(values, index)
					totalValues += 1
					if totalValues == 5 {
						break
					}
				}
			}
		}

		var multiplier int = 1e9
		returnValue := 600_000_000_000
		for i := 0; i < 5; i++ {
			returnValue += (values[i] * multiplier)
			multiplier /= 100
		}

		return returnValue
	}

	/*
		Check for possible straight
	*/

	if straightValue > 0 {
		return 500_000_000_000 + (straightValue * 1e9)
	}

	/*
		Check for three of a kind
	*/

	if threeOfAKind > 0 {
		values := []int{}
		totalValues := 0
		for index := 14; index > 0; index-- {
			if (stateStraight&(1<<index) != 0) && (threeOfAKind != index) {
				values = append(values, index)
				totalValues += 1
				if totalValues == 2 {
					break
				}
			}
		}

		var multiplier int = 1e7
		returnValue := 400_000_000_000 + (threeOfAKind * 1e9)
		for i := 0; i < 2; i++ {
			returnValue += (values[i] * multiplier)
			multiplier /= 100
		}

		return returnValue
	}

	/*
		Two pairs
	*/

	if secondPair > 0 {
		bestKicker := 0
		for index := 14; index > 0; index-- {
			if (stateStraight&(1<<index) != 0) && (bestPair != index) && (secondPair != index) {
				bestKicker = index
				break
			}
		}

		return 300_000_000_000 + (bestPair * 1e9) + (secondPair * 1e7) + (bestKicker * 1e5)
	}

	/*
		One pair
		Getting to the bad looking hands right ?
	*/

	if bestPair > 0 {
		values := []int{}
		totalValues := 0
		for index := 14; index > 0; index-- {
			if (stateStraight&(1<<index) != 0) && (bestPair != index) {
				values = append(values, index)
				totalValues += 1
				if totalValues == 3 {
					break
				}
			}
		}

		var multiplier int = 1e7
		returnValue := 200_000_000_000 + (bestPair * 1e9)
		for i := 0; i < 3; i++ {
			returnValue += (values[i] * multiplier)
			multiplier /= 100
		}

		return returnValue
	}

	/*
		High card
		Good enough to call though !!!
	*/

	values := []int{}
	totalValues := 0
	for index := 14; index > 0; index-- {
		if stateStraight&(1<<index) != 0 {
			values = append(values, index)
			totalValues += 1
			if totalValues == 5 {
				break
			}
		}
	}

	// FIX ME
	var multiplier int = 1e9
	returnValue := 100_000_000_000 + (bestPair * 1e9)
	for i := 0; i < 5; i++ {
		returnValue += (values[i] * multiplier)
		multiplier /= 100
	}

	return returnValue
}
