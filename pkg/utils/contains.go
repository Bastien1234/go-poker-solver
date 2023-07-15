package utils

func Contains(array []string, card string) bool {
	for _, el := range array {
		if el == card {
			return true
		}
	}

	return false
}

func ContainsInt(array []int, card int) bool {
	for _, el := range array {
		if el == card {
			return true
		}
	}

	return false
}
