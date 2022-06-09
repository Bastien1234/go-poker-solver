package tree

import (
	"fmt"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/node"
	"pokersolver/pkg/ranges"
)

var cameOnBuildFromNode int = 0

func buildFromNode(tree *Tree, n *node.Node) {
	cameOnBuildFromNode++

	for _, action := range n.Actions {
		rangeToAdd := make([]ranges.Hand, 0)
		betsToAdd := make([]int, 0)
		raisesToAdd := make([]int, 0)
		newPlayerTurn := ""

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
		if action > 1000 && n.RaiseLevel == constants.MaxRaises {
			var raiseValue int = (action * n.CurrentFacingBet) / 1000
			threashold := int(constants.Threashold * float32(n.EffectiveSize))
			if raiseValue > threashold {
				raiseValue = n.EffectiveSize
			}

			newNode := node.NewNode(
				rangeToAdd,
				betsToAdd,
				raisesToAdd,
				n.RaiseLevel+1,
				n.PotSize+raiseValue,
				n.EffectiveSize-raiseValue,
				raiseValue,
				newPlayerTurn,
				"fc",
			)

			n.PostActionNodes[action] = &newNode
		} else if action == 0 || action == -2 || action == -3 {
			// Another base case : closing actions : fold, call, check back

			// We don't add any more node !
		} else if action == -1 {
			// Open check
			newNode := node.NewNode(
				rangeToAdd,
				betsToAdd,
				raisesToAdd,
				n.RaiseLevel,
				n.PotSize,
				n.EffectiveSize,
				0,
				newPlayerTurn,
				"xbb",
			)

			n.PostActionNodes[action] = &newNode
			buildFromNode(tree, &newNode)
		} else if action > 0 && action < 1000 {
			// We bet !
			betValue := (action * n.PotSize) / 100

			if betValue >= n.EffectiveSize {
				// Can't bet more than what we have, right !
				// Therefore it's an all in situation

				bet := (n.PotSize / n.EffectiveSize) * 100
				newNode := node.NewNode(
					rangeToAdd,
					betsToAdd,
					raisesToAdd,
					n.RaiseLevel+1,
					n.PotSize+bet,
					n.EffectiveSize-bet,
					bet,
					newPlayerTurn,
					"fc",
				)

				n.PostActionNodes[action] = &newNode
			} else {
				// Just raising, not all in

				newNode := node.NewNode(
					rangeToAdd,
					betsToAdd,
					raisesToAdd,
					n.RaiseLevel+1,
					n.PotSize+betValue,
					n.EffectiveSize-betValue,
					betValue,
					newPlayerTurn,
					"fcr",
				)

				n.PostActionNodes[action] = &newNode
				buildFromNode(tree, &newNode)

			}
		} else if action > 1000 {
			// Raising baby

			raiseValue := (action * n.CurrentFacingBet) / 1000
			threashold := int(constants.Threashold * float32(n.EffectiveSize))
			if raiseValue > threashold {
				// All in situation
				raiseValue = n.EffectiveSize
				newNode := node.NewNode(
					rangeToAdd,
					betsToAdd,
					raisesToAdd,
					n.RaiseLevel+1,
					n.PotSize+raiseValue,
					n.EffectiveSize-raiseValue,
					raiseValue,
					newPlayerTurn,
					"fc",
				)

				n.PostActionNodes[action] = &newNode
			} else {
				// Recursive situation as we reopen the action baby

				newNode := node.NewNode(
					rangeToAdd,
					betsToAdd,
					raisesToAdd,
					n.RaiseLevel+1,
					n.PotSize+raiseValue,
					n.EffectiveSize-raiseValue,
					raiseValue,
					newPlayerTurn,
					"fcr",
				)

				n.PostActionNodes[action] = &newNode
				buildFromNode(tree, &newNode)

			}
		}

	}
}

type Tree struct {
	PotSize       int
	EffectiveSize int
	OopRange      []ranges.Hand
	IpRange       []ranges.Hand
	OopBets       []int
	IpBets        []int
	OopRaises     []int
	IpRaises      []int

	// Init later
	MaxRaises int
	Root      *node.Node
}

func NewTree(potSize int, effectiveSize int, oopRange []ranges.Hand, ipRange []ranges.Hand, oopBets []int, ipBets []int, oopRaises []int, ipRaises []int) Tree {
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

	buildFromNode(t, t.Root)
	fmt.Printf("Came on recursive thing %v times\n", cameOnBuildFromNode)

}

func (t *Tree) PrintTree() {
	bfsLevel := 0
	currentNode := t.Root
	queue := []*node.Node{currentNode}

	for len(queue) > 0 {
		fmt.Printf("at level : %v\n", bfsLevel)

		currentNode = queue[0]
		queue = queue[1:]
		for _, el := range currentNode.Actions {
			fmt.Printf("Action : %v\n", el)
		}

		for _, value := range currentNode.PostActionNodes {
			if value != nil {
				queue = append(queue, value)
			}

		}

		bfsLevel++
	}

}
