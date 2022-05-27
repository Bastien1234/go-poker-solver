package utils

import (
	"math/rand"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

// r1.Intn(100)

func Generate2() []int {
	random1 := r1.Intn(100)
	random2 := 100 - random1

	return []int{random1, random2}
}

func Generate3() []int {
	random1 := r1.Intn(100)
	random2 := r1.Intn(100)
	random3 := r1.Intn(100)

	total := random1 + random2 + random3
	n1 := random1 * 100 / total
	n2 := random2 * 100 / total
	n3 := 100 - n1 - n2

	return []int{n1, n2, n3}
}

func Generate4() []int {
	random1 := r1.Intn(100)
	random2 := r1.Intn(100)
	random3 := r1.Intn(100)
	random4 := r1.Intn(100)

	total := random1 + random2 + random3 + random4
	n1 := random1 * 100 / total
	n2 := random2 * 100 / total
	n3 := random3 * 100 / total
	n4 := 100 - n1 - n2 - n3

	return []int{n1, n2, n3, n4}
}

func Generate5() []int {
	random1 := r1.Intn(100)
	random2 := r1.Intn(100)
	random3 := r1.Intn(100)
	random4 := r1.Intn(100)
	random5 := r1.Intn(100)

	total := random1 + random2 + random3 + random4 + random5
	n1 := random1 * 100 / total
	n2 := random2 * 100 / total
	n3 := random3 * 100 / total
	n4 := random4 * 100 / total
	n5 := 100 - n1 - n2 - n3 - n4

	return []int{n1, n2, n3, n4, n5}
}
