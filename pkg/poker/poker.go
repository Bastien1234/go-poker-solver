package poker

import (
	"fmt"
	"math/rand"
	"pokersolver/pkg/constants"
	"pokersolver/pkg/deck"
	"pokersolver/pkg/handSolver"
	"pokersolver/pkg/ranges"
	"pokersolver/pkg/utils"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/timpalpant/go-cfr"
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

type memoMap struct {
	mm  map[string]int
	mtx sync.Mutex
}

var MemoMap = &memoMap{
	mm:  map[string]int{},
	mtx: sync.Mutex{},
}

func getMemoValue(array []string) int {
	MemoMap.mtx.Lock()
	defer MemoMap.mtx.Unlock()
	sorted := sort.StringSlice(array)
	final := ""
	for i := range sorted {
		final += sorted[i]
	}

	if val, ok := MemoMap.mm[final]; ok {
		return val
	} else {
		MemoMap.mm[final] = handSolver.HandSolver(sorted, true)
		return MemoMap.mm[final]
	}
}

var deckChannel = make(chan []string, 1e7)

func RunDeckChannel() {
	for {
		deck := deck.MakeDeck()
		deckChannel <- deck
	}
}

func GetLimitedRunOuts(nbRunouts int) [][]string {
	runouts := [][]string{}
	for {

		runout := []string{}

		deck := <-deckChannel

		for i := range deck {
			if !utils.Contains(constants.Board, deck[i]) {
				runout = append(runout, deck[i])
				if len(runout) == 2 {
					runouts = append(runouts, runout)
					continue
				}
			}
		}

		if len(runouts) == nbRunouts {
			return runouts
		}

	}
}

type PokerNode struct {
	parent        *PokerNode
	player        int
	children      []PokerNode
	probabilities []float64
	History       string

	// Hands held by either players
	P0Card, P1Card ranges.Hand

	// Perso
	RaiseLevel       int
	PotSize          int
	EffectiveSize    int
	CurrentFacingBet int
	Board            []string

	// Stage
	Stage NodeStage

	// From Kunh
	RegretSum   []float32
	StrategySum []float32
	Strategy    []float32
	ReachPr     float32
	ReachPrSum  float32

	LimitedRunOuts [][]string
	TurnIndex      int
}

func NewGame(handsOOP, handsIP []ranges.Hand, limitedRunOuts [][]string) *PokerNode {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsOOP), func(i, j int) {
		handsOOP[i], handsOOP[j] = handsOOP[j], handsOOP[i]
	})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(handsIP), func(i, j int) {
		handsIP[i], handsIP[j] = handsIP[j], handsIP[i]
	})

	var validHandOOP ranges.Hand
	var validHandIP ranges.Hand
	var currentForbiddenCards []string

	for _, hand := range handsOOP {
		if !utils.Contains(constants.Board, hand.Cards[0]) && !utils.Contains(constants.Board, hand.Cards[1]) {
			validHandOOP = hand
			currentForbiddenCards = append(constants.Board, validHandOOP.Cards...)
			continue
		}
	}

	for _, hand := range handsIP {
		if !utils.Contains(currentForbiddenCards, hand.Cards[0]) && !utils.Contains(currentForbiddenCards, hand.Cards[1]) {
			validHandIP = hand
			continue
		}
	}
	return &PokerNode{
		player:         chance,
		PotSize:        constants.Pot,
		Board:          constants.Board,
		RaiseLevel:     0,
		Stage:          Flop,
		History:        h_RootNode,
		P0Card:         validHandOOP,
		P1Card:         validHandIP,
		LimitedRunOuts: limitedRunOuts,
	}
}

func (n *PokerNode) Close() {
	n.children = nil
	n.probabilities = nil
}

// Player implements *PokerNode.
func (n *PokerNode) Player() int {
	return n.player
}

func (n *PokerNode) NumChildren() int {
	if n.children == nil {
		n.buildChildren()
	}

	return len(n.children)
}

func (n *PokerNode) GetChild(i int) PokerNode {
	if n.children == nil {
		n.buildChildren()
	}

	return n.children[i]
}

func (n *PokerNode) Parent() *PokerNode {
	return n.parent
}

func (n *PokerNode) GetChildProbability(i int) float64 {
	if n.children == nil {
		n.buildChildren()
	}

	return n.probabilities[i]
}

// Type implements *PokerNode.
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
	if len(n.History) == 0 {
		return false
	}

	lastAction := n.History[len(n.History)-1:]
	if lastAction == h_Fold {
		return true
	}

	if len(n.History) > 2 {

		areWeAllIn := n.History[len(n.History)-2 : len(n.History)-1]
		if areWeAllIn == h_AllIn {
			return true
		}
	}

	if (n.Stage == River) && (lastAction == h_CheckBack) {
		return true
	}

	if (n.Stage == River) && (lastAction == h_Call) {
		return true
	}

	return false
}

func (n *PokerNode) Utility(player int) float64 {
	cardPlayer := n.playerCard(player)
	cardOpponent := n.playerCard(1 - player)

	// By convention, terminal nodes are labeled with the player whose
	// turn it would be (i.e. not the last acting player).

	var isShowdown bool
	lastAction := n.History[len(n.History)-1:]

	if lastAction == h_Fold {
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
			// Also increase pot size on bets ?
			return -float64(n.PotSize)
		}
	}

	// Implement River and all ins situation with sampling

	// River situation
	if n.Stage == River {
		playerFinalhand := append(n.Board, cardPlayer.Cards...)
		playerHandValue := getMemoValue(playerFinalhand)

		opponentFinalHand := append(n.Board, cardOpponent.Cards...)
		opponentHandValue := getMemoValue(opponentFinalHand)

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

	if len(n.Board) > 4 {
		panic("tooooo big board, something went wrong")
	}

	for i := 0; i < constants.AllInSamplesize; i++ {
		fullBoard := getFullBoard(n.Board, cardPlayer.Cards, cardOpponent.Cards)

		playerFinalhand := append(fullBoard, cardPlayer.Cards...)
		playerHandValue := getMemoValue(playerFinalhand)

		// FIX ME: bad inputs

		opponentFinalHand := append(fullBoard, cardOpponent.Cards...)

		opponentHandValue := getMemoValue(opponentFinalHand)

		if playerHandValue > opponentHandValue {
			cumulativePlayerWinnings += float64(n.PotSize)
		} else if playerHandValue < opponentHandValue {
			cumulativePlayerWinnings -= float64(n.PotSize)
		} else if playerHandValue == opponentHandValue {
			// do nothing
		}

	}

	averagePlayerWinnings = cumulativePlayerWinnings / float64(constants.AllInSamplesize)

	return averagePlayerWinnings

}

func getFullBoard(currentBoard, player, opponent []string) []string {
	// FIX ME : optimize with get all 30? boards in one time ?
	// deck := deck.MakeDeck()
	deck := <-deckChannel
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
		return n.P0Card
	}

	return n.P1Card
}

func uniformDist(n int) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = 1.0 / float64(n)
	}
	return result
}

func (n *PokerNode) buildChildren() {
	if n.IsTerminal() {
		n.children = nil
		return
	}
	// Case chance node p0
	previousAction := n.History[len(n.History)-1 : len(n.History)]
	switch previousAction {
	case h_RootNode:
		n.children = buildRootDeals(n)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
		n.probabilities = uniformDist(len(n.children)) // everywhere !!!
		// push children onto the map !!!

	case h_P0Deal:
		n.children = buildP1Deals(n)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
		n.probabilities = uniformDist(len(n.children))

	// Case new chance node
	case h_p1Deal:
		n.children = buildOpenAction(n)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0

	case h_Chance:
		n.children = buildP0Deals(n)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
	case h_Check:
		n.children = buildCBAction(n)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0

	// Bet and raise have to take into account if we are all in...
	case h_Bet1:
		n.children = buildFCRAction(n, true)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
	case h_Bet2:

		n.children = buildFCRAction(n, true)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
	case h_Bet3:

		n.children = buildFCRAction(n, true)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0

	case h_Raise1:

		n.children = buildFCRAction(n, true)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0
	case h_Raise2:

		n.children = buildFCRAction(n, true)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0

	case h_AllIn:

		n.children = buildFCRAction(n, false)
		nbActions := len(n.children)
		n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
		n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
		n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
		n.ReachPr = 0.0
		n.ReachPrSum = 0.0

	case h_CheckBack:
		if n.parent.Stage == Flop || n.parent.Stage == Turn {

			n.children = buildChanceNode(n)
			nbActions := len(n.children)
			n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
			n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
			n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
			n.ReachPr = 0.0
			n.ReachPrSum = 0.0
			n.probabilities = uniformDist(len(n.children))

		} else {

			// FIX me !!!
			fmt.Println("wowowowowowow")
			// n.children = buildFCRAction(n, false)
			n.children = nil
		}

	case h_Call:
		if n.parent.Stage == Flop || n.parent.Stage == Turn {

			n.children = buildChanceNode(n)
			nbActions := len(n.children)
			n.RegretSum = utils.FilledArrayFloat(nbActions, 0.0)
			n.StrategySum = utils.FilledArrayFloat(nbActions, 0.0)
			n.Strategy = utils.FilledArrayFloat(nbActions, 1.0/float32(nbActions))
			n.ReachPr = 0.0
			n.ReachPrSum = 0.0
		} else {
			fmt.Println("wowowowowowow 2")

			// n.children = buildFCRAction(n, false)
			n.children = nil
		}
	}

	// Case check
	// Case check back
	// Case bet
	// Case raise
	// Case fold
	// case allin
}

func buildRootDeals(parent *PokerNode) []PokerNode {
	var results []PokerNode

	// hand := parent.P0Card

	child := PokerNode{
		parent:  parent,
		player:  chance,
		History: h_P0Deal,
		// P0Card:  hand,
		P0Card: parent.P0Card,
		P1Card: parent.P1Card,

		PotSize:        constants.Pot,
		EffectiveSize:  constants.EffectiveStack,
		RaiseLevel:     0,
		Board:          constants.Board,
		LimitedRunOuts: parent.LimitedRunOuts,
	}

	results = append(results, child)

	/*
		child := *parent
		child.parent = parent
		child.player = chance
		child.P0Card = hand
		child.history += h_P0Deal

		results = append(results, &child)
	*/

	return results
}

func buildP0Deals(parent *PokerNode) []PokerNode {
	var results []PokerNode

	// hand := parent.P0Card

	child := *parent
	child.parent = parent
	child.player = chance
	// child.P0Card = hand
	child.History += h_P0Deal

	results = append(results, child)

	return results
}

func buildP1Deals(parent *PokerNode) []PokerNode {
	var results []PokerNode

	// hand := parent.P1Card

	child := *parent
	child.parent = parent
	child.player = player0
	// child.P1Card = hand
	child.History += h_p1Deal

	results = append(results, child)

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
		child.History += choice

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
		child.History += choice

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

	// FC or FCR node
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
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
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
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {

							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise2)
							bets = append(bets, choice)
						}
					}
				}

			case River:
				for index, choice := range constants.IPRiverRaises {
					switch index {
					case 0:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise2)
							bets = append(bets, choice)
						}
					}
				}
			}
		} else {
			switch stage {
			case Flop:
				for index, choice := range constants.OOPFlopRaises {
					switch index {
					case 0:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
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
				for index, choice := range constants.OOPTurnRaises {
					switch index {
					case 0:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {

							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise2)
							bets = append(bets, choice)
						}
					}
				}

			case River:
				for index, choice := range constants.OOPRiverRaises {
					switch index {
					case 0:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise1)
							bets = append(bets, choice)
						}

					case 1:
						if isOverThreasholdRaise(parent, choice) || parent.RaiseLevel == constants.MaxRaises {
							choices = append(choices, h_AllIn)
							bets = append(bets, float64(parent.EffectiveSize))
							break
						} else {
							choices = append(choices, h_Raise2)
							bets = append(bets, choice)
						}
					}
				}
			}
		}
	}

	for index, choice := range choices {

		var addToPotSize float64 = float64(parent.PotSize) * bets[index]

		child := *parent
		child.parent = parent
		child.player = player
		child.History += choice

		child.PotSize += int(addToPotSize)

		result = append(result, child)

	}

	return result
}

func buildChanceNode(parent *PokerNode) []PokerNode {
	var results []PokerNode

	// allPossibleCards := deck.MakeDeck()
	// allPossibleCards := <-deckChannel
	validCards := []string{}

	// FIX ME: not sure
	var player int

	if parent.player == player0 {
		player = player1
	} else {
		player = player0
	}

	/*
		for _, c := range allPossibleCards {
			if !utils.Contains(parent.Board, c) {
				validCards = append(validCards, c)
			}
		}

		if len(validCards) > constants.MaxChanceNodes {
			validCards = validCards[0:constants.MaxChanceNodes]
		}
	*/

	// ---------------
	hasTurnIndex := false
	var turnIndexes = []int{}

	if parent.Stage == Flop {
		for index, cards := range parent.LimitedRunOuts {
			card := cards[0]

			if !utils.Contains(parent.Board, card) {
				validCards = append(validCards, card)
				turnIndexes = append(turnIndexes, index)
				hasTurnIndex = true
			}
		}
	} else if parent.Stage == Turn {
		card := parent.LimitedRunOuts[parent.TurnIndex][1]
		validCards = append(validCards, card)
	}

	// ------------

	for index, newCard := range validCards {
		// Card not on the board

		var newNodeStage NodeStage
		if parent.Stage == Flop {
			newNodeStage = Turn
		} else if parent.Stage == Turn {
			newNodeStage = River
		} else {
			panic("Wrong chance node I think... I am sure !")
		}

		/*

			child := PokerNode{
				parent:  parent,
				player:  player,
				history: parent.history + h_Chance,
				P0Card:  hand,

				PotSize:       parent.PotSize,
				EffectiveSize: parent.EffectiveSize,
				RaiseLevel:    parent.RaiseLevel,
				Board:         append(parent.Board, newCard),
				Stage:         newNodeStage,
			}

		*/
		finalTurnIndex := -1
		if hasTurnIndex {
			finalTurnIndex = turnIndexes[index]
		}

		child := *parent
		child.player = player
		child.History = parent.History + "*" + newCard + "*" + h_Chance
		child.Board = append(parent.Board, newCard)
		child.Stage = newNodeStage
		child.TurnIndex = finalTurnIndex

		results = append(results, child)
	}

	return results
}

// --------------------------

type pokerInfoSet struct {
	history string
	card    string
}

func (p pokerInfoSet) Key() []byte {
	return []byte(p.history + "-" + p.card)
}

func (p pokerInfoSet) MarshalBinary() ([]byte, error) {
	return p.Key(), nil
}

func (p *pokerInfoSet) UnmarshalBinary(buf []byte) error {
	parts := strings.SplitN(string(buf), "-", 1)
	if len(parts) != 2 {
		return fmt.Errorf("invalid binary poker info set: %v", parts)
	}

	p.history = parts[0]
	p.card = parts[1]
	return nil
}

// InfoSet implements *PokerNode.
func (n *PokerNode) InfoSet(player int) cfr.InfoSet {
	cardString := n.playerCard(player).Cards[0] + n.playerCard(player).Cards[1]
	return &pokerInfoSet{
		history: n.History,
		card:    cardString,
	}
}

func (n *PokerNode) InfoSetKey(player int) []byte {
	return n.InfoSet(player).Key()
}
