package utils

func PickIndexFromVector(distributionVector []int, n int) int {
	index := 0
	cumul := 0

	for _, number := range distributionVector {

		if n == 0 {
			return 0
		}

		cumul += number
		if n < cumul {
			return index
		}

		index++
	}

	return index
}
