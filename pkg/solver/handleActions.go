package solver 

import (
	"pokersolver/pkg/node"
)

func handleFold(curNode *node.Node, currentFrequency, currentHandFrenquency int) float32 {
	if playerIsIP {
		if curNode.PlayersTurn == "oop" {
			// Loosing invested money before folding
			*valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			*valueOfAction += ((float64(curNode.PotSize)) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}
	} else {
		if curNode.PlayersTurn == "ip" {
			// Loosing invested money before folding
			*valueOfAction -= ((float64(curNode.PotSize) / 2) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		} else {
			*valueOfAction += (float64(curNode.PotSize) * (float64(currentFrequency) / 100) * (float64(currentHandFrenquency) / 100)) * divider
		}
	}
}