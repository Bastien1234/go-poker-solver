package ranges

func RangeToList(matrix [][]int, pctToKeep int) [][]string {
	var vectorToReturn [][]string = make([][]string)
	mapMatrixToCards = map[int]string{
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

	return vectorToReturn
}
