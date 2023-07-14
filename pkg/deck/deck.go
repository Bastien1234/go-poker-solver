package deck

import (
	"math/rand"
	"time"
)

// var deck []string

// func init() {
// 	d := make([]string, 0)
// 	colors := []string{"h", "d", "c", "s"}
// 	cards := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}

// 	for _, card := range cards {
// 		for _, color := range colors {
// 			d = append(d, (card + color))
// 		}
// 	}

// 	deck = d
// }

func MakeDeck() []string {
	/*
		// shuffle
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})

		return deck
	*/

	deck := make([]string, 0)
	colors := []string{"h", "d", "c", "s"}
	cards := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}

	for _, card := range cards {
		for _, color := range colors {
			deck = append(deck, (card + color))
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}
