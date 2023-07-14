package million

import (
	"fmt"
	"pokersolver/pkg/deck"
	"pokersolver/pkg/handSolver"
	"sync"
	"time"
)

type ThreadCounter struct {
	threads int
	mtx     sync.Mutex
}

// var wg sync.WaitGroup

func SolveOneMillion(n int, bitwiseVersion bool) {
	fmt.Printf("Starting to solve %d hands\n", n)

	deckStart := time.Now()
	// hands := [][]string{}
	hands := make([][]string, n)
	for i := 0; i < n; i++ {
		deck := deck.MakeDeck()
		handToSolve := deck[0:7]
		// hands = append(hands, handToSolve) // Oh la caca !!!!!!!
		hands[i] = handToSolve
	}
	deckMakingDuration := time.Since(deckStart)
	fmt.Printf("Took %s to build decks\n", deckMakingDuration)

	start := time.Now()

	if !bitwiseVersion {

		for i := 0; i < n; i++ {
			handSolver.HandSolver(hands[i], false)
		}
	} else {
		// GO GO GO
		for i := 0; i < n; i++ {
			handSolver.HandSolver2(hands[i], false)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)

	// wg.Done()
}

/*
func SolveOneMillion8Threads() {

	wg.Add(8)
	for i := 0; i < 8; i++ {
		go SolveOneMillion(100)
	}

	wg.Wait()
}
*/
