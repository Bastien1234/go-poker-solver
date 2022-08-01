package utils

import (
	"math"
)

// !!!!!!!!!!! It goes under 0 !!!!!!!!!

func UpdateFrenquencies(ev []int, frenquencies *[]int, delta int) {
	if len(ev) != len(*frenquencies) {
		panic("Error in inputs !")
	}

	maxEv := math.MinInt
	minEv := math.MaxInt

	for _, el := range ev {
		if el < minEv {
			minEv = el
		}

		if el > maxEv {
			maxEv = el
		}
	}

	middle := (maxEv + minEv) / 2

	pointsOffMiddle := make([]int, len(ev))
	totalOffPoints := 0

	for index, el := range ev {
		diff := el - middle
		pointsOffMiddle[index] = diff
		if diff > 0 {
			totalOffPoints += diff
		} else {
			totalOffPoints += (-diff)
		}
	}

	var pointValue float32 = float32(delta) / float32(totalOffPoints)

	newFreq := make([]int, len(ev))

	totalNew := 0
	for index, el := range *frenquencies {
		newFreq[index] = el + int(float32(pointsOffMiddle[index])*pointValue)
		totalNew += el + int(float32(pointsOffMiddle[index])*pointValue)
	}
	missingPoint := 100 - totalNew
	newFreq[len(ev)-1] = newFreq[len(ev)-1] + missingPoint
	*frenquencies = newFreq
}
