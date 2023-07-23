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

type NodeStrategy struct {
	RegretSum   []float32
	StrategySum []float32
	Strategy    []float32
	ReachPr     float32
	ReachPrSum  float32
}

func NewNodeStrategy(nbActions int) *NodeStrategy {
	n := &NodeStrategy{
		RegretSum:   utils.FilledArrayFloat(nbActions, 0.0),
		StrategySum: utils.FilledArrayFloat(nbActions, 0.0),
		Strategy:    utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions)),
		ReachPr:     0.0,
		ReachPrSum:  0.0,
	}

	return n
}

type CFR struct {
	Strategy map[string]*NodeStrategy
}

func (cfr *CFR) GetNodeStrategy(history string, nbActions int) *NodeStrategy {

	val, ok := cfr.Strategy[history]
	if ok {
		return val

	} else {
		cfr.Strategy[history] = NewNodeStrategy(nbActions)
		return cfr.Strategy[history]
	}
}

func (n *NodeStrategy) GetStrategy() []float32 {
	regrets := n.RegretSum
	var normalizingSum float32 = 0.0

	for _, el := range regrets {

		normalizingSum += el
	}

	if normalizingSum > 0.0 {
		for index := range regrets {
			regrets[index] = regrets[index] / normalizingSum
		}

		return regrets
	}

	yeahDefaultList := utils.FilledArrayFloat(len(regrets), 1.0/float32(len(regrets)))

	return yeahDefaultList
}

func (n *NodeStrategy) UpdateStrategy() {

	for index := range n.StrategySum {
		n.StrategySum[index] += n.ReachPr * n.Strategy[index]
	}

	n.ReachPrSum += n.ReachPr

	n.Strategy = n.GetStrategy()

	n.ReachPr = 0.0
}

func (n *NodeStrategy) GetAverageStrategy() []float32 {

	strategy := n.StrategySum
	var totalFloat float32 = 0.0
	for i := range strategy {
		strategy[i] = strategy[i] / n.ReachPrSum
		totalFloat += strategy[i]
	}

	for i := range strategy {
		strategy[i] /= totalFloat
	}

	return strategy
}

func New() *CFR {
	return &CFR{
		Strategy: map[string]*NodeStrategy{},
	}
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
	policy := c.GetNodeStrategy(node.History, nChildren)
	strategy := policy.GetStrategy()
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
		policy.ReachPr += reachP0
		for i := range policy.RegretSum {
			policy.RegretSum[i] += reachP1 * regrets[i]
		}
	} else {
		policy.ReachPr += reachP1
		for i := range policy.RegretSum {
			policy.RegretSum[i] += reachP0 * regrets[i]
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
