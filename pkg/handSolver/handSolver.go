package handSolver

import (
	"pokersolver/pkg/utils"
	"sort"
)

func HandSolver(arr []string) int {

	// Constants

	cards := [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	cardsLow := [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	colors := [4]string{"h", "d", "c", "s"}

	// Hashmaps that we gonna populate later on
	hashValues := make(map[string]int)
	hashValuesLow := make(map[string]int)
	valuesArray := make([]int, 0)
	valuesArrayLow := make([]int, 0)

	possibleColor := false
	possibleStraightHigh := false
	possibleStraightLow := false
	straightValueHigh := -1
	straightValueLow := -1

	colorHearts := 0
	listHeartsHigh := make([]int, 0)
	listHeartsLow := make([]int, 0)

	colorDiamonds := 0
	listDiamondsHigh := make([]int, 0)
	listDiamondsLow := make([]int, 0)

	colorClubs := 0
	listClubsHigh := make([]int, 0)
	listClubsLow := make([]int, 0)

	colorSpades := 0
	listSpadesHigh := make([]int, 0)
	listSpadesLow := make([]int, 0)

	// Populate value arrays
	for i := 0; i < 7; i++ {
		valueOfCard := arr[i][0:1]
		valuesArray[i] = hashValues[valueOfCard]
		valuesArrayLow[i] = hashValuesLow[valueOfCard]
	}

	for _, el := range arr {
		color := el[1:2]
		switch color {
		case "h":
			colorHearts++
			listHeartsHigh = append(listHeartsHigh, hashValues[el[0:1]])
			listHeartsLow = append(listHeartsLow, hashValuesLow[el[0:1]])
			break

		case "d":
			colorDiamonds++
			listDiamondsHigh = append(listDiamondsHigh, hashValues[el[0:1]])
			listDiamondsLow = append(listDiamondsLow, hashValuesLow[el[0:1]])
			break
		case "c":
			colorClubs++
			listClubsHigh = append(listClubsHigh, hashValues[el[0:1]])
			listClubsLow = append(listClubsLow, hashValuesLow[el[0:1]])
			break
		case "s":
			colorSpades++
			listSpadesHigh = append(listSpadesHigh, hashValues[el[0:1]])
			listSpadesLow = append(listSpadesLow, hashValuesLow[el[0:1]])
			break
		}
	}

	// Sort arrays
	sort.Ints(valuesArray)
	sort.Ints(valuesArrayLow)
	sort.Ints(listHeartsHigh)
	sort.Ints(listHeartsLow)
	sort.Ints(listDiamondsHigh)
	sort.Ints(listDiamondsLow)
	sort.Ints(listClubsHigh)
	sort.Ints(listClubsLow)
	sort.Ints(listSpadesHigh)
	sort.Ints(listSpadesLow)

	// Check if possible color
	if colorHearts >= 5 || colorDiamonds >= 5 || colorClubs >= 5 || colorSpades >= 5 {
		possibleColor = true
	}

	// Check if possible straight

	/*
	   The set_of_array variable has different size options :
	   If it's under 5, there can't be any straight
	   If it's 5 or over, we can check if now the straight is possible
	*/

	setOfValues := make([]int, 0)
	setOfValuesLow := make([]int, 0)

	for _, el := range valuesArray {
		if !utils.ContainsInt(setOfValues, el) {
			setOfValues = append(setOfValues, el)
		}
	}

	for _, el := range valuesArrayLow {
		if !utils.ContainsInt(setOfValuesLow, el) {
			setOfValuesLow = append(setOfValuesLow, el)
		}
	}

	if len(setOfValues) >= 5 {
		iterations := len(setOfValues) - 4

		for i := 0; i < iterations; i++ {
			if (setOfValues[i]+1 == setOfValues[i+1]) && (setOfValues[i+1] +1 == setOfValues[i+2]) && (setOfValues[i+2] +1 == setOfValues[i+3]) && (setOfValues[i+3] +1 == setOfValues[i+4]) {
					possibleStraightHigh = true 
					straightValueHigh = setOfValues[i+4]
			}
		}

		for i := 0; i < iterations; i++ {
			if (setOfValuesLow[i]+1 == setOfValuesLow[i+1]) && (setOfValuesLow[i+1] +1 == setOfValuesLow[i+2]) && (setOfValuesLow[i+2] +1 == setOfValuesLow[i+3]) && (setOfValuesLow[i+3] +1 == setOfValuesLow[i+4]) {
					possibleStraightLow = true 
					straightValueLow = setOfValuesLow[i+4]
			}
		}

	}

	/*
        Check if possible straight flush
        Return 9 billions then value of high card
    */

	if possibleColor==true && straightValueHigh > 0 {
		listOfSuits := make([]int, 0)
		if colorHearts >= 0 {
			for _, el := range listHeartsHigh {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorDiamonds >= 0 {
			for _, el := range listDiamondsHigh {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorClubs >= 0 {
			for _, el := range listClubsHigh {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorSpades >= 0 {
			for _, el := range listSpadesHigh {
				listOfSuits = append(listOfSuits, el)
			}
		}

		valueToReturn := 900_000_000_000 
		for i:= len(listOfSuits)-5; i>=0; i-- {
			if listOfSuits[i] == listOfSuits[i+4] - 4 {
				return valueToReturn + (listOfSuits[i+4] * 1e7)
			}
		}
	}

	if possibleColor==true && straightValueLow > 0 {
		listOfSuits := make([]int, 0)
		if colorHearts >= 0 {
			for _, el := range listHeartsLow {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorDiamonds >= 0 {
			for _, el := range listDiamondsLow {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorClubs >= 0 {
			for _, el := range listClubsLow {
				listOfSuits = append(listOfSuits, el)
			}
		} else if colorSpades >= 0 {
			for _, el := range listSpadesLow {
				listOfSuits = append(listOfSuits, el)
			}
		}

		valueToReturn := 900_000_000_000 
		for i:= len(listOfSuits)-5; i>=0; i-- {
			if listOfSuits[i] == listOfSuits[i+4] - 4 {
				return valueToReturn + (listOfSuits[i+4] * 1e7)
			}
		}
	}

	/*
        Check is possible four of a kind
        returns 8 billions then the FOAK value and finally the kicker
    */

	if len(setOfValues) < 5 {
		quadValue := -1
		bestKicker := -1
		counter := make(map[int]int)

		for _, el := range valuesArray {
			if _, ok :=  counter[el]; ok {
				counter[el] += 1
				if counter[el] == 4 {
					quadValue = el 
					for i:= 0; i<len(setOfValues); i++ {
						if setOfValues[i] != el {
							if setOfValues[i] > bestKicker {
								bestKicker = setOfValues[i]
							}
						}
					}

					return 800_000_000_000 + (quadValue * 1e9) + (bestKicker * 1e7)
				}
			} else {
				counter[el] = 1
			}
		}
	}

	/*
        Check is possible full house
    */

	if len(setOfValues) <= 4 {
		counter := make(map[int]int)
		bestSet := -1
		bestPair := -1

		for _, el := range valuesArray {
			if _, ok := counter[el]; ok {
				counter[el] += 1
				if counter[el] == 3 {
					if el > bestSet {
						bestSet = el 
					}
				}
			}
		}

		for _, el := range valuesArray {
			if counter[el] >= 2 && bestSet != el {
				if el > bestPair {
					bestPair = el 
				}
			}
		}

		if bestSet > 0 && bestPair > 0 {
			return 700_000_000_000 + (bestSet * 1e9) + (bestPair * 1e7)
		}
	}

	/*
        Check if possible color
    */

	if (possibleColor) {
		
	}




	return 0
}
