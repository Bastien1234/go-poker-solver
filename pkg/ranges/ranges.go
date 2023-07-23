package ranges

import (
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/utils"
	"time"
)

type Hand struct {
	Cards     []string
	Frequency float64
}

/*
func RangeToList(matrix [][]float64, pctToKeep int) [][]string {
	var vectorToReturn [][]string = make([][]string, 0)
	mapMatrixToCards := map[int]string{
		0:  "A",
		1:  "K",
		2:  "Q",
		3:  "J",
		4:  "T",
		5:  "9",
		6:  "8",
		7:  "7",
		8:  "6",
		9:  "5",
		10: "4",
		11: "3",
		12: "2",
	}

	colors := []string{"h", "d", "c", "s"}

	allColorsCombo := [][]string{
		{"h", "d"},
		{"h", "c"},
		{"h", "s"},
		{"d", "c"},
		{"d", "s"},
		{"c", "s"},
	}

	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			// case pair : i = j
			if i == j {
				if matrix[i][j] > 0 {
					for k := 0; k < (matrix[i][j] / 10); k++ {
						for _, colorCombo := range allColorsCombo {
							card1 := mapMatrixToCards[i] + colorCombo[0]
							card2 := mapMatrixToCards[i] + colorCombo[1]
							handToAdd := []string{card1, card2}
							vectorToReturn = append(vectorToReturn, handToAdd)
						}
					}
				}
			}

			// case suited, i < j
			if i < j {
				if matrix[i][j] > 0 {
					for k := 0; k < (matrix[i][j] / 10); k++ {
						for _, color := range colors {
							card1 := mapMatrixToCards[i] + color
							card2 := mapMatrixToCards[j] + color
							handToAdd := []string{card1, card2}
							vectorToReturn = append(vectorToReturn, handToAdd)
						}
					}
				}
			}

			// case offsuited, i > j
			if i > j {
				if matrix[i][j] > 0 {
					for k := 0; k < (matrix[i][j] / 100); k++ {
						for _, colorCombo := range allColorsCombo {
							card1 := mapMatrixToCards[i] + colorCombo[0]
							card2 := mapMatrixToCards[j] + colorCombo[1]
							handToAdd1 := []string{card1, card2}

							card3 := mapMatrixToCards[j] + colorCombo[0]
							card4 := mapMatrixToCards[i] + colorCombo[1]
							handToAdd2 := []string{card3, card4}

							vectorToReturn = append(vectorToReturn, handToAdd1, handToAdd2)

						}
					}
				}
			}
		}
	}

	// shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vectorToReturn), func(i, j int) {
		vectorToReturn[i], vectorToReturn[j] = vectorToReturn[j], vectorToReturn[i]
	})

	// Keep only a percentage of hands, as per second parameter of the function
	var handsToKeep int = (len(vectorToReturn) * pctToKeep) / 100
	finalVector := make([][]string, 0)
	for i := 0; i < handsToKeep; i++ {
		finalVector = append(finalVector, vectorToReturn[i])
	}

	return finalVector
}
*/

func RangeToVector(matrix [][]float64) []Hand {
	var vectorToReturn []Hand = make([]Hand, 0)
	mapMatrixToCards := map[int]string{
		0:  "A",
		1:  "K",
		2:  "Q",
		3:  "J",
		4:  "T",
		5:  "9",
		6:  "8",
		7:  "7",
		8:  "6",
		9:  "5",
		10: "4",
		11: "3",
		12: "2",
	}

	colors := []string{"h", "d", "c", "s"}

	allColorsCombo := [][]string{
		{"h", "d"},
		{"h", "c"},
		{"h", "s"},
		{"d", "c"},
		{"d", "s"},
		{"c", "s"},
	}

	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			// case pair : i = j
			if i == j {
				if matrix[i][j] > 0 {
					for _, colorCombo := range allColorsCombo {
						card1 := mapMatrixToCards[i] + colorCombo[0]
						card2 := mapMatrixToCards[i] + colorCombo[1]
						handToAdd := []string{card1, card2}
						h := Hand{}
						h.Cards = handToAdd
						h.Frequency = matrix[i][j]
						vectorToReturn = append(vectorToReturn, h)
					}
				}
			}

			// case suited, i < j
			if i < j {
				if matrix[i][j] > 0 {
					for _, color := range colors {
						card1 := mapMatrixToCards[i] + color
						card2 := mapMatrixToCards[j] + color
						handToAdd := []string{card1, card2}
						h := Hand{}
						h.Cards = handToAdd
						h.Frequency = matrix[i][j]
						vectorToReturn = append(vectorToReturn, h)
					}
				}
			}

			// case offsuited, i > j
			if i > j {
				if matrix[i][j] > 0 {
					for _, colorCombo := range allColorsCombo {
						card1 := mapMatrixToCards[i] + colorCombo[0]
						card2 := mapMatrixToCards[j] + colorCombo[1]
						handToAdd1 := []string{card1, card2}
						h1 := Hand{}
						h1.Cards = handToAdd1
						h1.Frequency = matrix[i][j]

						card3 := mapMatrixToCards[j] + colorCombo[0]
						card4 := mapMatrixToCards[i] + colorCombo[1]
						handToAdd2 := []string{card3, card4}
						h2 := Hand{}
						h2.Cards = handToAdd2
						h2.Frequency = matrix[i][j]

						vectorToReturn = append(vectorToReturn, h1, h2)

					}
				}
			}
		}
	}

	return vectorToReturn
}

// Get ranges
var matrixOOP = constants.MatrixOOP
var matrixIP = constants.MatrixIp

// var handsOOP = ranges.RangeToVector(matrixOOP)
// var handsIP = ranges.RangeToVector(matrixIP)

func GetHands(board []string) (Hand, Hand) {
	handsOOP := RangeToVector(matrixOOP)
	handsIP := RangeToVector(matrixIP)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsOOP), func(i, j int) {
		handsOOP[i], handsOOP[j] = handsOOP[j], handsOOP[i]
	})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsIP), func(i, j int) {
		handsIP[i], handsIP[j] = handsIP[j], handsIP[i]
	})

	var validHandOOP Hand
	var validHandIP Hand
	var currentForbiddenCards []string

	for _, hand := range handsOOP {
		if !utils.Contains(board, hand.Cards[0]) && !utils.Contains(board, hand.Cards[1]) {
			validHandOOP = hand
			currentForbiddenCards = append(board, validHandOOP.Cards...)
			continue
		}
	}

	for _, hand := range handsIP {
		if !utils.Contains(currentForbiddenCards, hand.Cards[0]) && !utils.Contains(currentForbiddenCards, hand.Cards[1]) {
			validHandIP = hand
			continue
		}
	}

	return validHandIP, validHandOOP
}
