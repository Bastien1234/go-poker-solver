package utils

func FilledArrayInt(lenght int, values int) (a []int) {
	a = make([]int, lenght)
	for i := range a {
		a[i] = values
	}

	return
}

func FilledArrayFloat(length int, values float32) []float32 {
	a := make([]float32, length)
	for i := range a {
		a[i] = values
	}

	return a
}
