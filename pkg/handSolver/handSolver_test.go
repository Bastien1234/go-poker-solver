package handSolver

import (
	"testing"
)

func TestHandSolver(t *testing.T) {
	straightFlush1 := HandSolver([]string{"Ah", "Kh", "Qh", "Jh", "Th", "2s", "2d"})
	if straightFlush1 != 914000000000 {
		t.Errorf("Problem on straight flush ace to high")
	}

	straightFlush2 := HandSolver([]string{"Ah", "2h", "3h", "4h", "5h", "2s", "2d"})
	if straightFlush2 != 905000000000 {
		t.Errorf("Problem on straight flush ace to low")
	}

	quads1 := HandSolver([]string{"Ah", "Ad", "Ac", "As", "Th", "2s", "2d"})
	if quads1 != 81410000000 {
		t.Errorf("Problem with quads 1\nValue found : %v", quads1)
	}
}
