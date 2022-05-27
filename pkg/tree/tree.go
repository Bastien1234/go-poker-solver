package tree

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/node"
)

func buildFromNode(tree *Tree, n *node.Node) {
	for _, action := range n.Actions {
		fmt.Println(action)
		building := true
		rangeToAdd := make([][]string, 0)
		betsToAdd := make([]int, 0)
		raisesToAdd := make([]int, 0)
		newPlayerTurn := ""
		newNodeToAdd := node.Node{}

		if n.PlayersTurn == "oop" {
			rangeToAdd = tree.IpRange
			betsToAdd = tree.IpBets
			raisesToAdd = tree.IpRaises
			newPlayerTurn = "ip"
		} else if n.PlayersTurn == "ip" {
			rangeToAdd = tree.OopRange
			betsToAdd = tree.OopBets
			raisesToAdd = tree.OopRaises
			newPlayerTurn = "oop"
		} else {
			fmt.Println("Error !!!!!!")
		}

		// Building according to action

		// Base case : maximum raises are reached
		fmt.Println("action : ", action)

	}
}

type Tree struct {
	PotSize       int
	EffectiveSize int
	OopRange      [][]string
	IpRange       [][]string
	OopBets       []int
	IpBets        []int
	OopRaises     []int
	IpRaises      []int

	// Init later
	MaxRaises int
	Root      *node.Node
}

func NewTree(potSize int, effectiveSize int, oopRange [][]string, ipRange [][]string, oopBets []int, ipBets []int, oopRaises []int, ipRaises []int) Tree {
	tree := Tree{}
	tree.PotSize = potSize
	tree.EffectiveSize = effectiveSize
	tree.OopRange = oopRange
	tree.IpRange = ipRange
	tree.OopBets = oopBets
	tree.IpBets = ipBets
	tree.OopRaises = oopRaises
	tree.IpRaises = ipRaises

	tree.MaxRaises = constants.MaxRaises
	tree.Root = nil

	return tree
}

func (t *Tree) MakeRiverTree() {
	fmt.Println("starting to build river tree")

	rootNode := node.NewNode(
		t.OopRange,
		t.OopBets,
		t.OopRaises,
		0,
		t.PotSize,
		t.EffectiveSize,
		0,
		"oop",
		"root",
	)

	t.Root = &rootNode

}
