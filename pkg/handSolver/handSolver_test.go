package handSolver

import (
	"testing"
)

func TestHandSolver(t *testing.T) {

	// Bad inputs

	shortArray := HandSolver([]string{"Kh", "Qh", "Jh", "Th", "2s", "2d"}, true)
	if shortArray != -1 {
		t.Errorf("Problem accepting too short array")
	}

	longArray := HandSolver([]string{"Ah", "Ad", "Kh", "Qh", "Jh", "Th", "2s", "2d"}, true)
	if longArray != -1 {
		t.Errorf("Problem acceptign too long array")
	}

	badInputs := HandSolver([]string{"Ah", "Kh", "Qh", "Jh", "Th", "2s", "2x"}, true)
	if badInputs != -1 {
		t.Errorf("Problem acceppting bad inputs")
	}

	badInputs2 := HandSolver([]string{"Ah", "Kh", "Qh", "Jh", "Th", "2s", "Xs"}, true)
	if badInputs2 != -1 {
		t.Errorf("Problem acceppting bad inputs")
	}

	badInputs3 := HandSolver([]string{"Ahh", "Kh", "Qh", "Jh", "Th", "2s", "Xs"}, true)
	if badInputs3 != -1 {
		t.Errorf("Problem acceppting bad inputs")
	}

	duplicates := HandSolver([]string{"Ah", "Ah", "Kh", "Jh", "Th", "2s", "2d"}, true)
	if duplicates != -1 {
		t.Errorf("Problem accepting duplicates")
	}

	// Correct inputs

	straightFlush1 := HandSolver([]string{"Ah", "Kh", "Qh", "Jh", "Th", "2s", "2d"}, true)
	if straightFlush1 != 914000000000 {
		t.Errorf("Problem on straight flush ace to high")
	}

	straightFlush2 := HandSolver([]string{"Ah", "2h", "3h", "4h", "5h", "2s", "2d"}, true)
	if straightFlush2 != 905000000000 {
		t.Errorf("Problem on straight flush ace to low")
	}

	quads1 := HandSolver([]string{"Ah", "Ad", "Ac", "As", "Th", "2s", "2d"}, true)
	if quads1 != 814_100_000_000 {
		t.Errorf("Problem with quads 1\nValue found : %v", quads1)
	}
}
