package million

import (
	"fmt"
	"pokersolver/pkg/deck"
	"pokersolver/pkg/handSolver"
	"time"
)

// var wg sync.WaitGroup

func SolveOneMillion(n int) {
	fmt.Println("Starting to solver 1 million hands")
	start := time.Now()

	for i := 0; i < n; i++ {
		deck := deck.MakeDeck()
		handToSolve := deck[0:7]
		handSolver.HandSolver(handToSolve)
	}

	elapsed := time.Since(start)
	fmt.Println("Took %S to solver 1M hands baby", elapsed)

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
