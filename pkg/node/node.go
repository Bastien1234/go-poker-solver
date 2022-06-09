package node

import (
	"math"
	"pokersolver/pkg/ranges"
)

// Subnodes

type SubNode struct {
	Hand        []string
	Actions     []int
	Ev          []int
	Frequencies []int
}

func NewSubNode(hand []string, actions []int) SubNode {
	sn := SubNode{}
	sn.Hand = hand
	sn.Actions = actions

	numOfActions := len(actions)

	sn.Ev = make([]int, numOfActions)
	sn.Frequencies = make([]int, numOfActions)

	var defaultValue int = 100 / len(sn.Actions)
	for i := 0; i < len(actions); i++ {
		sn.Frequencies[i] = defaultValue
	}

	return sn
}

// Nodes

type Node struct {
	HandRange        []ranges.Hand
	OpponentRange    []ranges.Hand
	Actions          []int
	Raises           []int
	RaiseLevel       int
	PotSize          int
	EffectiveSize    int
	CurrentFacingBet int
	PlayersTurn      string
	NodeType         string

	// To be init later
	PostActionNodes map[int]*Node
	GlobalActionMap map[int]int // Pas computée ...
	LocalActionMap  map[string]*SubNode
	GlobalBestScore int // Pas computé non plus ...
}

func NewNode(handRange []ranges.Hand, actions []int, raises []int, raiseLevel int, potSize int, effectiveSize int, currentFacingBet int, playersTurn string, nodeType string) Node {
	n := Node{}
	n.HandRange = handRange
	n.Actions = actions
	n.Raises = raises
	n.RaiseLevel = raiseLevel
	n.PotSize = potSize
	n.EffectiveSize = effectiveSize
	n.CurrentFacingBet = currentFacingBet
	n.PlayersTurn = playersTurn
	n.NodeType = nodeType

	n.GlobalBestScore = math.MinInt

	n.PostActionNodes = make(map[int]*Node)
	// n.GlobalActionMap = make(map[int]int)
	n.LocalActionMap = make(map[string]*SubNode)

	switch n.NodeType {
	case "fcr":
		n.Actions = n.Raises
		n.Actions = append(n.Actions, -2, -3)
		for _, act := range n.Actions {
			n.PostActionNodes[act] = nil
		}
		break

	case "root":
		n.Actions = append(n.Actions, -1)
		for _, act := range n.Actions {
			n.PostActionNodes[act] = nil
		}
		break

	case "xbb":
		n.Actions = append(n.Actions, 0)
		for _, act := range n.Actions {
			n.PostActionNodes[act] = nil
		}
		break

	case "fc":
		n.Actions = nil
		n.Actions = append(n.Actions, -2, -3)
		n.PostActionNodes[-2] = nil
		n.PostActionNodes[-3] = nil
		break

	}

	return n
}
