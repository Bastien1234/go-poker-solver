package poker

import (
	"encoding"
	"io"
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/deck"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/utils"
)

// NodeType is the type of node in an extensive-form game tree.
type NodeType int

const (
	ChanceNodeType NodeType = iota
	TerminalNodeType
	PlayerNodeType
)

type NodeStage int

const (
	Flop NodeStage = iota
	Turn
	River
)

const (
	chance  = -1
	player0 = 0
	player1 = 1
)

// Notation history
const (
	h_RootNode  = "R"
	h_P0Deal    = "0"
	h_p1Deal    = "1"
	h_Chance    = "C"
	h_Check     = "x"
	h_CheckBack = "c"
	h_Bet1      = "q"
	h_Bet2      = "s"
	h_Bet3      = "d"
	h_Raise1    = "w"
	h_Raise2    = "v"
	h_AllIn     = "a"
	h_Fold      = "f"
	h_Call      = "k"
)

// InfoSet is the observable game history from the point of view of one player.
type InfoSet interface {
	// Key is an identifier used to uniquely look up this InfoSet
	// when accumulating probabilities in tabular CFR.
	//
	// It may be an arbitrary string of bytes and does not need to be
	// human-readable. For example, it could be a simplified abstraction
	// or hash of the full game history.
	Key() []byte
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// ChanceNode is a node that has a pre-defined probability distribution over its children.
type ChanceNode interface {
	// Get the probability of the ith child of this node.
	// May only be called for nodes with Type == Chance.
	GetChildProbability(i int) float64

	// Sample a single child from this Chance node according to the probability
	// distribution over children.
	//
	// Implementations may reuse sampling.SampleChanceNode to sample from the CDF,
	// (by scanning over GetChildProbability) or implement their own more efficient
	// sampling.
	SampleChild() (child GameTreeNode, p float64)
}

// PlayerNode is a node in which one of the player's acts.
type PlayerNode interface {
	// Player returns this current node's acting player.
	// It may only be called for nodes with IsChance() == false.
	Player() int
	// InfoSet returns the information set for this node for the given player.
	InfoSet(player int) InfoSet
	// InfoSetKey returns the equivalent of InfoSet(player).Key(),
	// but can be used to avoid allocations incurred by the InfoSet interface.
	InfoSetKey(player int) []byte
	// Utility returns this node's utility for the given player.
	// It must only be called for nodes with type == Terminal.
	Utility(player int) float64
}

// Tree node represents a node in a directed rooted tree.
type TreeNode interface {
	// The number of direct children of this node.
	NumChildren() int
	// Get the ith child of this node.
	GetChild(i int) GameTreeNode
	// Get the parent of this node.
	Parent() GameTreeNode
}

// GameTreeNode is the interface for a node in an extensive-form game tree.
type GameTreeNode interface {
	// NodeType returns the type of game node.
	Type() NodeType
	// Release resources held by this node (including any children).
	Close()

	TreeNode
	ChanceNode
	PlayerNode
}

// ----------------------------------------------------------

// Get ranges
var matrixOOP = constants.MatrixOOP
var matrixIP = constants.MatrixIp

var handsOOP = ranges.RangeToVector(matrixOOP)
var handsIP = ranges.RangeToVector(matrixIP)

type PokerNode struct {
	parent        *PokerNode
	player        int
	children      []PokerNode
	probabilities []float64
	history       string

	// Hands held by either players
	p0Card, p1Card ranges.Hand

	// Perso
	RaiseLevel       int
	PotSize          int
	EffectiveSize    int
	CurrentFacingBet int
	Board            []string

	// Stage
	Stage NodeStage
}

func NewGame() *PokerNode {
	return &PokerNode{
		player:     chance,
		PotSize:    constants.Pot,
		Board:      constants.Board,
		RaiseLevel: 0,
		Stage:      Flop,
		history:    h_Chance,
	}
}

// string ?

func (n *PokerNode) Close() {
	n.children = nil
	n.probabilities = nil
}

func (n *PokerNode) NumChildren() int {
	if n.children == nil {
		n.buildChildren()
	}

	return len(n.children)
}

func (n *PokerNode) GetChild(i int) *PokerNode {
	if n.children == nil {
		n.buildChildren()
	}

	return &n.children[i]
}

func (n *PokerNode) GetChildProbability(i int) float64 {
	if n.children == nil {
		n.buildChildren()
	}

	return n.probabilities[i]
}

// SampleChild implements cfr.GameTreeNode.
func (n *PokerNode) SampleChild() (*PokerNode, float64) {
	i := rand.Intn(n.NumChildren())
	return n.GetChild(i), n.GetChildProbability(i)
}

// Type implements cfr.GameTreeNode.
func (n *PokerNode) Type() NodeType {
	if n.IsTerminal() {
		return TerminalNodeType
	} else if n.player == chance {
		return ChanceNodeType
	}

	return PlayerNodeType
}

func (n *PokerNode) IsTerminal() bool {
	/*
		Possible cases of terminal node :

		Player folds
		Player is all in
		Player checks back or call on the river
	*/

	// FIX ME because I'm not sure
	if len(n.history) == 0 {
		return false
	}

	lastAction := n.history[len(n.history)-1]
	if lastAction == 'f' {
		return true
	}

	if len(n.history) > 1 {
		allInSituation := n.history[len(n.history)-2:]
		if allInSituation == "ac" || allInSituation == "af" {
			return true
		}
	}

	if n.Stage == River {
		lastAction := n.history[len(n.history)-1]
		if lastAction == 'f' || lastAction == 'c' {
			return true
		}
	}

	return false
}

func (n *PokerNode) Utility(player int) float64 {
	cardPlayer := n.playerCard(player)
	cardOpponent := n.playerCard(1 - player)

	// By convention, terminal nodes are labeled with the player whose
	// turn it would be (i.e. not the last acting player).

	var isShowdown bool
	if n.history[len(n.history)-1] == 'f' {
		isShowdown = false
	} else {
		isShowdown = true
	}

	if !isShowdown {
		if n.player == player {
			return float64(n.PotSize)
		} else {
			// FIX ME
			// return -float64(n.PotSize) * opponent or hero card frequency
			return -float64(n.PotSize)
		}
	}

	// Implement River and all ins situation with sampling

	// River situation
	if n.Stage == River {
		playerFinalhand := append(constants.Board, cardPlayer.Cards...)
		playerHandValue := handSolver.HandSolver2(playerFinalhand, false)

		opponentFinalHand := append(constants.Board, cardOpponent.Cards...)
		opponentHandValue := handSolver.HandSolver2(opponentFinalHand, false)

		if playerHandValue > opponentHandValue {
			return float64(n.PotSize)
		} else if playerHandValue < opponentHandValue {
			return -float64(n.PotSize)
		} else if playerHandValue == opponentHandValue {
			return 0.0
		}
	}

	/*
		All in situation
		To avoid running into all possible nodes, we sample the possibles outcomes
	*/

	var cumulativePlayerWinnings float64 = 0

	var averagePlayerWinnings float64

	// FIX ME : memoization
	for i := 0; i < constants.AllInSamplesize; i++ {
		fullBoard := getFullBoard(n.Board, cardPlayer.Cards, cardOpponent.Cards)

		playerFinalhand := append(fullBoard, cardPlayer.Cards...)
		playerHandValue := handSolver.HandSolver2(playerFinalhand, false)

		opponentFinalHand := append(fullBoard, cardOpponent.Cards...)
		opponentHandValue := handSolver.HandSolver2(opponentFinalHand, false)

		if playerHandValue > opponentHandValue {
			cumulativePlayerWinnings += float64(n.PotSize)
		} else if playerHandValue < opponentHandValue {
			cumulativePlayerWinnings -= float64(n.PotSize)
		} else if playerHandValue == opponentHandValue {
			// do nothing
		}

		averagePlayerWinnings = cumulativePlayerWinnings / float64(constants.AllInSamplesize)

	}

	return averagePlayerWinnings

}

func getFullBoard(currentBoard, player, opponent []string) []string {
	// FIX ME : optimize with get all 30? boards in one time ?
	deck := deck.MakeDeck()
	cardsOut := append(append(currentBoard, player...), opponent...) // Can be better ?

	availableCards := []string{}

	for _, card := range deck {
		if !utils.Contains(cardsOut, card) {
			availableCards = append(availableCards, card)
			if len(availableCards) > 1 {
				break
			}
		}
	}

	var fullBoard []string
	if len(currentBoard) == 3 {
		fullBoard = append(append(currentBoard, availableCards[0]), availableCards[1])
	}

	if len(currentBoard) == 4 {
		fullBoard = append(currentBoard, availableCards[0])
	}

	return fullBoard
}

func (n *PokerNode) playerCard(player int) ranges.Hand {
	if player == player0 {
		return n.p0Card
	}

	return n.p1Card
}

func uniformDist(n int) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = 1.0 / float64(n)
	}
	return result
}

var rundowns = constants.MaxRundowns

func (n *PokerNode) buildChildren() {
	// Case chance node p0
	previousAction := n.history[len(n.history)-1 : len(n.history)]
	switch previousAction {
	case h_Chance:
		n.children = buildP0Deals(n)
		n.probabilities = uniformDist(len(n.children))

	case h_P0Deal:
		n.children = buildP1Deals(n)
		n.probabilities = uniformDist(len(n.children))

	// Case new chance node
	case h_p1Deal:
		n.children = buildOpenAction(n)
	}

	// Case check
	// Case check back
	// Case bet
	// Case raise
	// Case fold
	// case allin
}

func buildP0Deals(parent *PokerNode) []PokerNode {
	var results []PokerNode

	for _, hand := range handsOOP {
		// Card not on the board
		if !utils.Contains(constants.Board, hand.Cards[0]) && !utils.Contains(constants.Board, hand.Cards[1]) {

			child := PokerNode{
				parent:  parent,
				player:  chance,
				history: h_P0Deal,
				p0Card:  hand,

				PotSize:       constants.Pot,
				EffectiveSize: constants.EffectiveStack,
				RaiseLevel:    0,
				Board:         constants.Board,
			}

			results = append(results, child)
		}
	}

	return results
}

func buildP1Deals(parent *PokerNode) []PokerNode {
	var results []PokerNode

	for _, hand := range handsIP {
		// Card not on the board

		// FIX ME: remove properly duplicates between two players
		if !utils.Contains(constants.Board, hand.Cards[0]) && !utils.Contains(constants.Board, hand.Cards[1]) {

			child := *parent
			child.parent = parent
			child.player = player0
			child.p1Card = hand
			child.history += h_p1Deal

			results = append(results, child)
		}
	}

	return results
}

func buildOpenAction(parent *PokerNode) []PokerNode {
	var result []PokerNode

	// First to act, action is check and different betsizes
	choices := []string{h_Check}
	bets := []float64{0}

	stage := parent.Stage

	switch stage {
	case Flop:
		for index, choice := range constants.OOPFlopBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}

	case Turn:
		for index, choice := range constants.OOPTurnBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}

	case River:
		for index, choice := range constants.OOPRiverBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}
	}

	for index, choice := range choices {

		var addToPotSize float64 = float64(parent.PotSize) * bets[index]

		child := *parent
		child.parent = parent
		child.player = player1
		child.history += choice

		child.PotSize += int(addToPotSize)

		result = append(result, child)

	}

	return result
}

// FIX ME: Bet size with threashold etc
func buildCBAction(parent *PokerNode) []PokerNode {
	var result []PokerNode

	// This is only after open check
	choices := []string{h_CheckBack}
	bets := []float64{0}

	stage := parent.Stage

	switch stage {
	case Flop:
		for index, choice := range constants.IPFlopBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}

	case Turn:
		for index, choice := range constants.IPTurnBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}

	case River:
		for index, choice := range constants.IPRiverBets {
			switch index {
			case 0:
				choices = append(choices, h_Bet1)
				bets = append(bets, choice)

			case 1:
				choices = append(choices, h_Bet2)
				bets = append(bets, choice)
			case 2:
				choices = append(choices, h_Bet3)
				bets = append(bets, choice)
			}
		}
	}

	for index, choice := range choices {

		var addToPotSize float64 = float64(parent.PotSize) * bets[index]

		child := *parent
		child.parent = parent
		child.player = player0
		child.history += choice

		child.PotSize += int(addToPotSize)

		result = append(result, child)

	}

	return result
}

func isOverThreasholdBet(parent *PokerNode, choice float64) bool {
	return float64(parent.PotSize)+(float64(parent.PotSize)*choice) >= float64(parent.EffectiveSize)*constants.Threashold
}

func isOverThreasholdRaise(parent *PokerNode, choice float64) bool {
	potentialRaise := float64(parent.CurrentFacingBet) * choice
	return potentialRaise+float64(parent.EffectiveSize) >= float64(parent.EffectiveSize)*constants.Threashold
}

func buildFCRAction(parent *PokerNode, includeRaise bool) []PokerNode {
	var result []PokerNode

	// This is only after open check
	choices := []string{h_Fold, h_Call}
	bets := []float64{0, 0}

	stage := parent.Stage
	var player int

	if parent.player == player0 {
		player = player1
	} else {
		player = player0
	}

	if includeRaise {

		if player == player0 {

			switch stage {
			case Flop:
				for index, choice := range constants.IPFlopRaises {
					switch index {
					case 0:
						if isOverThreasholdRaise(parent, choice) {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise2)
							bets = append(bets, choice)
						}
					}
				}

			case Turn:
				for index, choice := range constants.IPTurnRaises {
					switch index {
					case 0:
						choices = append(choices, h_Bet1)
						bets = append(bets, choice)

					case 1:
						choices = append(choices, h_Bet2)
						bets = append(bets, choice)
					}
				}

			case River:
				for index, choice := range constants.IPRiverRaises {
					switch index {
					case 0:
						choices = append(choices, h_Bet1)
						bets = append(bets, choice)

					case 1:
						choices = append(choices, h_Bet2)
						bets = append(bets, choice)
					case 2:
						choices = append(choices, h_Bet3)
						bets = append(bets, choice)
					}
				}
			}
		}
	}

	for index, choice := range choices {

		var addToPotSize float64 = float64(parent.PotSize) * bets[index]

		child := *parent
		child.parent = parent
		child.player = player0
		child.history += choice

		child.PotSize += int(addToPotSize)

		result = append(result, child)

	}

	return result
}

func buildChanceNode(parent *PokerNode) (result []PokerNode) {
	// ...

	return
}

// ----------------------------------------------------------

// StrategyProfile maintains a collection of regret-matching policies for each
// player node in the game tree.
//
// The policytable and deepcfr packages provide implementations of StrategyProfile.
type StrategyProfile interface {
	// GetPolicy returns the NodePolicy for the given node.
	GetPolicy(node GameTreeNode) NodePolicy

	// Calculate the next strategy profile for all visited nodes.
	Update()
	// Get the current iteration (number of times update has been called).
	Iter() int

	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	io.Closer
}

// NodePolicy maintains the action policy for a single Player node.
type NodePolicy interface {
	// AddRegret provides new observed instantaneous regrets
	// to add to the total accumulated regret with the given weight.
	AddRegret(w float32, samplingQ, instantaneousRegrets []float32)
	// GetStrategy gets the current vector of probabilities with which the ith
	// available action should be played.
	GetStrategy() []float32

	// GetBaseline gets the current vector of action-dependend baseline values,
	// used in VR-MCCFR.
	GetBaseline() []float32
	// UpdateBaseline updates the current vector of baseline values.
	UpdateBaseline(w float32, action int, value float32)

	// AddStrategyWeight adds the current strategy with weight w to the average.
	AddStrategyWeight(w float32)
	// GetAverageStrategy returns the average strategy over all iterations.
	GetAverageStrategy() []float32

	// IsEmpty returns true if the NodePolicy is new and has no accumulated regret.
	IsEmpty() bool
}
