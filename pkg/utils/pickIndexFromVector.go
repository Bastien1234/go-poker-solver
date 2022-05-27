package utils

func PickIndexFromVector(distributionVector []int, n int) int {
	index := 0
	cumul := 0

	for _, number := range distributionVector {
		cumul += number
		if n < cumul {
			return index
		}

		index++
	}

	return index
}
