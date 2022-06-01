package utils

import (
	"math/rand"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

// r1.Intn(100)

func Generate2() []int {
	random1 := r1.Intn(10)
	random2 := 10 - random1

	return []int{random1 * 10, random2 * 10}
}

func Generate3() []int {
	random1 := r1.Intn(10)
	random2 := r1.Intn(10)
	random3 := r1.Intn(10)

	total := random1 + random2 + random3

	if total == 0 {
		return []int{40, 30, 30}
	}

	n1 := random1 * 10 / total
	n2 := random2 * 10 / total
	n3 := 10 - n1 - n2

	return []int{n1 * 10, n2 * 10, n3 * 10}
}

func Generate4() []int {
	random1 := r1.Intn(10)
	random2 := r1.Intn(10)
	random3 := r1.Intn(10)
	random4 := r1.Intn(10)

	total := random1 + random2 + random3 + random4

	if total == 0 {
		return []int{30, 30, 20, 20}
	}

	n1 := random1 * 10 / total
	n2 := random2 * 10 / total
	n3 := random3 * 10 / total
	n4 := 10 - n1 - n2 - n3

	return []int{n1 * 10, n2 * 10, n3 * 10, n4 * 10}
}

func Generate5() []int {
	random1 := r1.Intn(10)
	random2 := r1.Intn(10)
	random3 := r1.Intn(10)
	random4 := r1.Intn(10)
	random5 := r1.Intn(10)

	total := random1 + random2 + random3 + random4 + random5

	if total == 0 {
		return []int{20, 20, 20, 20, 20}
	}
	n1 := random1 * 10 / total
	n2 := random2 * 10 / total
	n3 := random3 * 10 / total
	n4 := random4 * 10 / total
	n5 := 10 - n1 - n2 - n3 - n4

	return []int{n1 * 10, n2 * 10, n3 * 10, n4 * 10, n5 * 10}
}
