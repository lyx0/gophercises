package main

import (
	"fmt"
	"strings"

	deck "github.com/lyx0/gophercises/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// Ace is currently worth 1 point, here we
			// change it so that it's worth 11
			return minScore + 10
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	score := 0

	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("===== FINAL HANDS =====")
	fmt.Printf("\n")
	fmt.Println("Player:", ret.Player, "\nScore:", pScore)
	fmt.Println("Dealer:", ret.Dealer, "\nScore:", dScore)
	switch {
	case pScore > 21:
		fmt.Printf("\n")
		fmt.Println("------------------")
		fmt.Println("--- You busted ---")
		fmt.Println("---- You lose ----")
		fmt.Println("------------------")
		fmt.Printf("\n\n\n")
	case dScore > 21:
		fmt.Printf("\n")
		fmt.Println("---------------------")
		fmt.Println("--- Dealer busted ---")
		fmt.Println("------ You win ------")
		fmt.Println("---------------------")
		fmt.Printf("\n\n\n")
	case pScore > dScore:
		fmt.Printf("\n")
		fmt.Println("---------------")
		fmt.Println("--- You win ---")
		fmt.Println("---------------")
		fmt.Printf("\n\n\n")
	case dScore > pScore:
		fmt.Printf("\n")
		fmt.Println("----------------")
		fmt.Println("--- You lose ---")
		fmt.Println("----------------")
		fmt.Printf("\n\n\n")
	case dScore == pScore:
		fmt.Printf("\n")
		fmt.Println("------------")
		fmt.Println("--- Draw ---")
		fmt.Println("------------")
		fmt.Printf("\n\n\n")
	}
	fmt.Println()

	ret.Player = nil
	ret.Dealer = nil

	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	for i := 0; i < 10; i++ {
		gs = Deal(gs)

		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Player: ", gs.Player)
			fmt.Println("Dealer: ", gs.Dealer.DealerString())
			fmt.Printf("Your cards are worth %d\n", gs.Player.Score())
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)
			fmt.Printf("\n")
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid option:", input)
			}
		}

		// If dealer score <= 16, we hit
		// If dealer has a soft 17, then we hit. (Ace as
		// 11 points and 16 points)
		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}

}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer

	default:
		panic("The gamestate got corrupted. Currently no players turn. Exiting")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}

	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)

	return ret
}
