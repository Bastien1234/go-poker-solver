package million

import (
	"fmt"
	"pokersolver/pkg/deck"
	"pokersolver/pkg/handSolver"
	"time"
)

func SolveOneMillion() {
	fmt.Println("Starting to solver 1 million hands")
	start := time.Now()

	for i := 0; i < 1000000; i++ {
		deck := deck.MakeDeck()
		handToSolve := deck[0:7]
		handSolver.HandSolver(handToSolve)
	}

	elapsed := time.Since(start)
	fmt.Println("Took %S to solver 1M hands baby", elapsed)
}
