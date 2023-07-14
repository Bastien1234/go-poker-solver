package utils

func FilledArrayInt(lenght int, values int) (a []int) {
	a = make([]int, lenght)
	for i := range a {
		a[i] = values
	}

	return
}
