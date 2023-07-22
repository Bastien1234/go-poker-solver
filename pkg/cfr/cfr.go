package cfr

import (
	"pokersolver/pkg/poker"
	"pokersolver/pkg/utils"
)

/*
type CFR struct {
	strategyProfile StrategyProfile
	slicePool       *floatSlicePool
}

func New(strategyProfile StrategyProfile) *CFR {
	return &CFR{
		strategyProfile: strategyProfile,
		slicePool:       &floatSlicePool{},
	}
}
*/

type CFR struct{}

func New() *CFR {
	return &CFR{}
}

func (c *CFR) Run(node *poker.PokerNode) float32 {
	return c.runHelper(node, node.Player(), 1.0, 1.0, 1.0)
}

func (c *CFR) runHelper(node *poker.PokerNode, lastPlayer int, reachP0, reachP1, reachChance float32) float32 {
	var ev float32
	switch node.Type() {
	case poker.TerminalNodeType:
		ev = float32(node.Utility(lastPlayer))
	case poker.ChanceNodeType:
		ev = c.handleChanceNode(node, lastPlayer, reachP0, reachP1, reachChance)
	default:
		sgn := getSign(lastPlayer, node.Player())
		ev = sgn * c.handlePlayerNode(node, reachP0, reachP1, reachChance)
	}

	node.Close()
	return ev
}

func (c *CFR) handleChanceNode(node *poker.PokerNode, lastPlayer int, reachP0, reachP1, reachChance float32) float32 {
	var expectedValue float32
	for i := 0; i < node.NumChildren(); i++ {
		child := node.GetChild(i)
		p := float32(node.GetChildProbability(i))
		expectedValue += p * c.runHelper(child, lastPlayer, reachP0, reachP1, reachChance*p)
	}

	return expectedValue
}

func (c *CFR) handlePlayerNode(node *poker.PokerNode, reachP0, reachP1, reachChance float32) float32 {
	player := node.Player()
	nChildren := node.NumChildren()
	if nChildren == 1 {
		// Optimization to skip trivial nodes with no real choice.
		child := node.GetChild(0)
		return c.runHelper(child, player, reachP0, reachP1, reachChance)
	}

	// policy := c.strategyProfile.GetPolicy(node)
	strategy := node.Strategy
	actionUtils := utils.FilledArrayFloat(nChildren, 0.0)

	for i := 0; i < nChildren; i++ {
		child := node.GetChild(i)
		p := strategy[i]
		// var util float32
		if player == 0 {
			actionUtils[i] = c.runHelper(child, player, p*reachP0, reachP1, reachChance)
		} else {
			actionUtils[i] = c.runHelper(child, player, reachP0, p*reachP1, reachChance)
		}

		// regrets[i] = util
		// cfValue += p * util
	}

	// Utility of information set
	var util float32 = 0.0
	for i := range actionUtils {
		util += actionUtils[i] * strategy[i]
	}

	// Compute regrets
	regrets := make([]float32, nChildren)
	for i := range regrets {
		// Only positive regrets baby
		candidateRegret := actionUtils[i] - util
		if candidateRegret > 0 {
			regrets[i] = candidateRegret
		} else {
			regrets[i] = 0
		}
	}

	if player == 0 {
		node.ReachPr += reachP0
		for i := range node.RegretSum {
			node.RegretSum[i] += reachP1 * regrets[i]
		}
	} else {
		node.ReachPr += reachP1
		for i := range node.RegretSum {
			node.RegretSum[i] += reachP0 * regrets[i]
		}
	}

	return util

	/*

		// Transform action utilities into instantaneous regrets by
		// subtracting out the expected utility over all possible actions.
		f32.AddConst(-cfValue, regrets)
		counterFactualP := counterFactualProb(player, reachP0, reachP1, reachChance)
		ones := c.slicePool.alloc(nChildren)
		defer c.slicePool.free(ones)
		for i := range ones {
			ones[i] = 1.0
		}
		policy.AddRegret(counterFactualP, ones, regrets)
		reachP := reachProb(player, reachP0, reachP1, reachChance)
		policy.AddStrategyWeight(reachP)
		return cfValue
	*/

}

func getSign(player1, player2 int) float32 {
	if player1 == player2 {
		return 1.0
	}

	return -1.0
}

func reachProb(player int, reachP0, reachP1, reachChance float32) float32 {
	if player == 0 {
		return reachP0 * reachChance
	} else {
		return reachP1 * reachChance
	}
}

// The probability of reaching this node, assuming that the current player
// tried to reach it.
func counterFactualProb(player int, reachP0, reachP1, reachChance float32) float32 {
	if player == 0 {
		return reachP1 * reachChance
	} else {
		return reachP0 * reachChance
	}
}

/*
func (k *CFR) Train(nbIterations int) {
	var expectedGameValue float32 = 0.0
	for i := 0; i < nbIterations; i++ {
		k.Deck = utils.Shuffle(k.Deck)
		expectedGameValue = k.CFR("", 1, 1)
		for _, node := range k.NodeMap {
			node.UpdateStrategy()
		}
	}

	expectedGameValue /= float32(nbIterations)
	fmt.Println("Finished training, print results here")

	for k, v := range k.NodeMap {
		fmt.Printf("key : %s\n", k)
		finalStrat := v.GetAverageStrategy()
		fmt.Println(finalStrat)
	}
}
*/
