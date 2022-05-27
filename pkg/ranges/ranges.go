package ranges

func RangeToList(matrix [][]int, pctToKeep int) [][]string {
	var vectorToReturn [][]string = make([][]string)
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

	for i := 0; i< 13; i++ {
		for j := 0; j<13; j++ {
			// case pair : i = j
			if i == j {
				if matrix[i][j] > 0 {
					for k := 0; k<(matrix[i][j] / 10); k++ {
						for _, colorCombo := range allColorsCombo {
							card1 := mapMatrixToCards[i] + colorCombo[0]
							card2 := mapMatrixToCards[i] + colorCombo[1]
							handToAdd := []string {card1, card2}
							vectorToReturn = append(vectorToReturn, handToAdd)
						}
					}
				}
			}

			// case suited, i < j
			if i < j {
				if matrix [i][jk] > 0 {
					for k :=0; k<(matrix[i][j] / 10); k++ {
						for _, color := range colors {
							card1 := mapMatrixToCards[i] + color
							card2 := mapMatrixToCards[j] + color 
							handToAdd := []string {card1, card2}
							vectorToReturn = append(vectorToReturn, handToAdd)
						}
					}
				}
			}

			// case offsuited, i > j
			if i > j {
				if matrix[i][j] > 0 {
					for k := 0 ; k< (matrix[i][j] / 100); k++ {
						for _, colorCombo := range allColorsCombo {
							card1 := mapMatrixToCards[i] + colorCombo[0]
							card2 := mapMatrixToCards[j] + colorCombo[1]
							handToAdd1 := []string {card1, card2}

							card3 := mapMatrixToCards[j] + colorCombo[0]
							card4 := mapMatrixToCards[i] + colorCombo[1]
							handToAdd2 := []string {card3, card4}

							vectorToReturn = append(vectorToReturn, handToAdd1, handToAdd2)

						}
					})
				}
			}
		}
	}

	return vectorToReturn
}
