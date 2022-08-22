package utils

import (
	"math"
	"fmt"
)

// !!!!!!!!!!! It goes under 0 !!!!!!!!!

/*
func UpdateFrenquencies(ev *[]int, frenquencies *[]int, delta int) {
	if len(*ev) != len(*frenquencies) {
		panic("Error in inputs !")
	}

	maxEv := math.MinInt
	minEv := math.MaxInt

	for _, el := range *ev {
		if el < minEv {
			minEv = el
		}

		if el > maxEv {
			maxEv = el
		}
	}

	middle := (maxEv + minEv) / 2

	pointsOffMiddle := make([]int, len(*ev))
	totalOffPoints := 0

	for index, el := range *ev {
		diff := el - middle
		pointsOffMiddle[index] = diff
		if diff > 0 {
			totalOffPoints += diff
		} else {
			totalOffPoints += (-diff)
		}
	}

	var pointValue float32 = float32(delta) / float32(totalOffPoints)

	newFreq := make([]int, len(*ev))

	totalNew := 0
	negatives := 0
	for index, el := range *frenquencies {
		newValue := el + int(float32(pointsOffMiddle[index])*pointValue)
		if newValue >= 0 {
		newFreq[index] = newValue
		} else {
			newFreq[index] = 0
			negatives -= newValue
		}
	}
	if negatives > 0 {
		for _, el := range newFreq {
			if el > 0 && (el + negatives) <= 100 {
				el += negatives
				negatives = 0
			}
		}
	}
	for _, el := range newFreq {
		totalNew += el
	}
	missingPoint := 100 - totalNew
	newFreq[len(*ev)-1] = newFreq[len(*ev)-1] + missingPoint
	*frenquencies = newFreq

	// Reset ev
	newEvArray := make([]int, len(*ev))
	for index, _ := range newEvArray {
		newEvArray[index] = 0
	}

	*ev = newEvArray
}
*/

func UpdateFrenquencies(ev *[]int, frenquencies *[]int, delta int) {
	if len(*ev) != len(*frenquencies) {
		panic("Error in inputs !")
	}

	fmt.Println(*ev)

	maxEv := math.MinInt
	minEv := math.MaxInt

	// Get min and max values

	for _, el := range *ev {
		if el < minEv {
			minEv = el
		}

		if el > maxEv {
			maxEv = el
		}
	}

	// Put everything over zero
	newPositiveEv := make([]int, len(*ev))


	if minEv < 0 {
		for index, el := range *ev {
			newPositiveEv[index] = el + (-minEv)
		}

		*ev = newPositiveEv

		//Recalculate min and max ev
		delta := -minEv
		minEv = 0
		maxEv = maxEv + delta
	}

	

	middle := (maxEv + minEv) / 2

	pointsOffMiddle := make([]int, len(*ev))
	totalOffPoints := 0

	for index, el := range *ev {
		diff := el - middle
		pointsOffMiddle[index] = diff
		if diff > 0 {
			totalOffPoints += diff
		} else {
			totalOffPoints += (-diff)
		}
	}

	var pointValue float32 = float32(delta) / float32(totalOffPoints)

	newFreq := make([]int, len(*ev))

	totalNew := 0
	for index, el := range *frenquencies {
		newValue := el + int(float32(pointsOffMiddle[index])*pointValue)
		newFreq[index] = newValue
	}

	for _, el := range newFreq {
		totalNew += el
	}
	missingPoint := 100 - totalNew
	newFreq[len(*ev)-1] = newFreq[len(*ev)-1] + missingPoint

	// Putting everything over zero
	totalNegatives := 0
	for _, el := range newFreq {
		if el < 0 {
			totalNegatives += (-el)
			el = 0
		}
	}

	// newAndPositiveFreq := make([]int, len(*ev))

	for totalNegatives > 0 {
		for i := range newFreq {
			if newFreq[i] > 0 && newFreq[i] < 100 && totalNegatives > 0 {
				newFreq[i] += 1
				totalNegatives -= 1
			} else if newFreq[i] < 0 {
				newFreq[i] = 0
			}
		}
	}

	totalPositives := 0
	for _, el := range newFreq {
		totalPositives += el
	}

	var multiplicationValue float64 = 100 / float64(totalPositives)
	floatArray := make([]float64, len(*ev))
	for index, el := range newFreq {
		floatArray[index] = float64(el) * multiplicationValue
	}

	for i := range floatArray {
		newFreq[i] = int(floatArray[i])
	}

	lastTotalNew := 0
	for _, el := range newFreq {
		lastTotalNew += el 
	}

	lastMissing := 100 - lastTotalNew
	newFreq[len(*ev) - 1] = newFreq[len(*ev) - 1] + lastMissing

	*frenquencies = newFreq

	// Reset ev
	newEvArray := make([]int, len(*ev))
	for index, _ := range newEvArray {
		newEvArray[index] = 0
	}

	*ev = newEvArray
}