package solver

import (
	"pokersolver/pkg/node"
)

func evRecursor(curNode *node.Node, subnode *node.SubNode, subnodeAction int, divider float64, localAction int) float64 {

	var currentFrequency int = subnode.Frequencies[subnodeActionIndex]
	var currentHandFrenquency int = subnode.Weight

	// Should return current action value
	// Entry point is subnode action

	switch localAction {

	// Case of closing actions
	case -3:
		// Handle fold

		// Maybe just return var float32 0.0
		return HandleFold(curNode, currentFrequency, currentHandFrenquency, divider)

	case 0, 2:
		// Handle call or check back
		return HandleCallAndCheckBack(curNode, currentFrequency, currentHandFrenquency, divider)

	// Case of non closing action, passing recursion
	default:
		// Handle open check, raise or bet
		nextNode := curNode.PostActionNodes[localAction].PostActionNodes[-1]
		return evRecursor(nextNode, 1, localAction)
	}
}

func GetSubnodeEv(currentNode node.Node, currentSubnode node.SubNode) float64 {

	// Iterate actions of the subnode
	for actionIdx, subnodeAction := range currentSubnode.Actions {

		currentSubnode.Ev[actionIdx] = evRecursor(currentNode, currentSubnode, subnodeAction, float64(1.0), action)
	}
}
