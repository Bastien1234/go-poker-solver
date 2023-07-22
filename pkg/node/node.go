package node

import "pokersolver/pkg/utils"

type Node struct {
	Key         string
	NbActions   int
	RegretSum   []float32
	StrategySum []float32
	Strategy    []float32
	ReachPr     float32
	ReachPrSum  float32
	Actions     []string

	Stage  int
	Player int

	RaiseLevel       int
	PotSize          int
	EffectiveSize    int
	CurrentFacingBet int
	Board            []string
}

func NewNode(
	key string,
	actions []string,
	stage int,
	player int,
	board []string,
	potSize int) Node {
	n := Node{
		Key:         key,
		Actions:     actions,
		NbActions:   len(actions),
		RegretSum:   utils.FilledArrayFloat(len(actions), 0.0),
		StrategySum: utils.FilledArrayFloat(len(actions), 0.0),
		Strategy:    utils.FilledArrayFloat(len(actions), 1.0/float32(len(actions))),
		ReachPr:     0.0,
		ReachPrSum:  0.0,
		Stage:       stage,
		Player:      player,
		Board:       board,
		PotSize:     potSize,
	}

	return n
}

func (n *Node) UpdateStrategy() {
	for index := range n.StrategySum {
		n.StrategySum[index] += n.ReachPr * n.Strategy[index]
	}

	n.ReachPrSum += n.ReachPr

	n.Strategy = n.GetStrategy()

	n.ReachPr = 0.0
}

func (n *Node) GetStrategy() []float32 {
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

func (n *Node) GetAverageStrategy() []float32 {
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
