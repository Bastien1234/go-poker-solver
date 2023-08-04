package poker

import (
	"math/rand"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/utils"
	"sort"
	"strconv"
	"sync"
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

	Visited int
}

func NewNodeStrategy(nbActions int) *NodeStrategy {
	n := &NodeStrategy{
		RegretSum:   utils.FilledArrayFloat(nbActions, 0.0),
		StrategySum: utils.FilledArrayFloat(nbActions, 0.0),
		Strategy:    utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions)),
		ReachPr:     0.0,
		ReachPrSum:  0.0,

		Visited: 0,
	}

	return n
}

type StrategyMap struct {
	SM map[string]*NodeStrategy

	mtx sync.Mutex
}

func (sm *StrategyMap) GetNodeStrategyKey(history string, nbActions int, playerCard ranges.Hand, player int, node PokerNode) string {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	str := playerCard.Cards
	sort.Strings(str)

	key := str[0] + str[1] + "-" + strconv.Itoa(player) + history

	_, ok := sm.SM[key]
	if ok {
		return key

	} else {
		sm.SM[key] = NewNodeStrategy(nbActions)
		return key
	}
}

func (sm *StrategyMap) SetChildren(nodes []PokerNode) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	for _, node := range nodes {
		str := []string{}
		if node.player == 0 {
			str = node.P0Card.Cards
		} else {
			str = node.P1Card.Cards
		}
		sort.Strings(str)

		key := str[0] + str[1] + "-" + strconv.Itoa(node.player) + node.History

		_, ok := sm.SM[key]
		if ok {
			return
		}

		n := &NodeStrategy{
			RegretSum:   node.RegretSum,
			StrategySum: node.StrategySum,
			Strategy:    node.Strategy,
			ReachPr:     node.ReachPr,
			ReachPrSum:  node.ReachPrSum,
		}

		sm.SM[key] = n

	}

}

func (sm *StrategyMap) _GetNodeStrategy(history string, nbActions int, playerCard ranges.Hand, player int, node PokerNode) *NodeStrategy {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	str := playerCard.Cards
	sort.Strings(str)

	key := str[0] + str[1] + "-" + strconv.Itoa(player) + history

	val, ok := sm.SM[key]
	if ok {
		return val

	} else {
		sm.SM[key] = NewNodeStrategy(nbActions)
		return sm.SM[key]
	}
}

func (sm *StrategyMap) GetStrategy(key string) []float32 {
	n := sm.SM[key]
	regrets := make([]float32, len(n.RegretSum))
	copy(regrets, n.RegretSum)
	var normalizingSum float32 = 0.0

	for _, el := range regrets {

		normalizingSum += el
	}

	if normalizingSum > 0.0 {
		for index := range regrets {
			value := regrets[index] / normalizingSum
			regrets[index] = value
		}

		return regrets
	}

	yeahDefaultList := utils.FilledArrayFloat(len(regrets), 1.0/float32(len(regrets)))

	return yeahDefaultList
}

func (sm *StrategyMap) UpdateStrategy() {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	for key := range sm.SM {

		n := sm.SM[key]

		for index := range n.StrategySum {
			n.StrategySum[index] += n.ReachPr * n.Strategy[index]
		}

		n.ReachPrSum += n.ReachPr

		n.Strategy = sm.GetStrategy(key)

		n.ReachPr = 0.0
	}
}

func (sm *StrategyMap) UpdatePolicy(key string, player int, reachP0, reachP1 float32, regrets []float32) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	n := sm.SM[key]

	if player == 0 {
		n.ReachPr += reachP0
		for i := range n.RegretSum {
			n.RegretSum[i] += reachP1 * regrets[i]
		}
	} else {
		n.ReachPr += reachP1
		for i := range n.RegretSum {
			n.RegretSum[i] += reachP0 * regrets[i]
		}
	}
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

func NewStrategyMap() *StrategyMap {
	return &StrategyMap{
		// Strategy: map[string]*NodeStrategy{},

		SM: make(map[string]*NodeStrategy, 800000),
	}
}

type CFR struct {
}

func New() *CFR {
	return &CFR{}
}

func (c *CFR) Run(node PokerNode, strategyMap *StrategyMap) float32 {
	return c.runHelper(node, node.Player(), 1.0, 1.0, 1.0, strategyMap)
}

func (c *CFR) runHelper(node PokerNode, lastPlayer int, reachP0, reachP1, reachChance float32, strategyMap *StrategyMap) float32 {
	var ev float32
	switch node.Type() {
	case TerminalNodeType:
		ev = float32(node.Utility(lastPlayer))
	case ChanceNodeType:
		ev = c.handleChanceNode(node, lastPlayer, reachP0, reachP1, reachChance, strategyMap)
	default:
		sgn := getSign(lastPlayer, node.Player())
		ev = sgn * c.handlePlayerNode(node, reachP0, reachP1, reachChance, strategyMap)
	}

	node.Close()
	return ev
}

func (c *CFR) handleChanceNode(node PokerNode, lastPlayer int, reachP0, reachP1, reachChance float32, strategyMap *StrategyMap) float32 {
	/*
		var expectedValue float32
			for i := 0; i < node.NumChildren(); i++ {
				child := node.GetChild(i)
				p := float32(node.GetChildProbability(i))
				expectedValue += p * c.runHelper(child, lastPlayer, reachP0, reachP1, reachChance*p, strategyMap)
			}
	*/

	// Try to randomize to see only one node
	nbChildren := node.NumChildren()
	childIndex := rand.Intn(nbChildren)
	child := node.GetChild(childIndex)
	// p := float32(node.GetChildProbability(i)) FIX ME OMG
	expectedValue := 1 * c.runHelper(child, lastPlayer, reachP0, reachP1, reachChance*1, strategyMap) // FIX ME: value p instead of 1

	return expectedValue
}

func (c *CFR) handlePlayerNode(node PokerNode, reachP0, reachP1, reachChance float32, strategyMap *StrategyMap) float32 {
	player := node.Player()
	nChildren := node.NumChildren()
	if nChildren == 1 {
		// Optimization to skip trivial nodes with no real choice.
		child := node.GetChild(0)
		return c.runHelper(child, player, reachP0, reachP1, reachChance, strategyMap)
	}

	var playerCard ranges.Hand
	if player == 0 {
		playerCard = node.P0Card
	} else if player == 1 {
		playerCard = node.P1Card
	}

	key := strategyMap.GetNodeStrategyKey(node.History, nChildren, playerCard, player, node)
	strategy := strategyMap.GetStrategy(key)
	actionUtils := utils.FilledArrayFloat(nChildren, 0.0)

	for i := 0; i < nChildren; i++ {
		child := node.GetChild(i)
		p := strategy[i]
		// var util float32
		if player == 0 {
			actionUtils[i] = c.runHelper(child, player, p*reachP0, reachP1, reachChance, strategyMap)
		} else {
			actionUtils[i] = c.runHelper(child, player, reachP0, p*reachP1, reachChance, strategyMap)
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

	strategyMap.UpdatePolicy(key, player, reachP0, reachP1, regrets)

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
