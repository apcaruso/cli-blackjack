package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Card struct {
	Suit  string
	Value string
}

var suits = []string{"♠", "♥", "♦", "♣"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

func newDeck() []Card {
	deck := []Card{}
	for _, s := range suits {
		for _, v := range values {
			deck = append(deck, Card{Suit: s, Value: v})
		}
	}
	return deck
}

func shuffle(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	for i := range deck {
		j := rand.Intn(len(deck))
		deck[i], deck[j] = deck[j], deck[i]
	}
}

func cardValue(c Card) int {
	switch c.Value {
	case "A":
		return 11
	case "K", "Q", "J":
		return 10
	default:
		var v int
		fmt.Sscanf(c.Value, "%d", &v)
		return v
	}
}

func handValue(hand []Card) int {
	total, aces := 0, 0
	for _, c := range hand {
		v := cardValue(c)
		if c.Value == "A" {
			aces++
		}
		total += v
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func perfectPair(hand []Card) bool {
	return len(hand) >= 2 && hand[0].Value == hand[1].Value
}

func twentyOneThree(hand []Card) bool {
	if len(hand) < 3 {
		return false
	}
	v1 := cardValue(hand[0])
	v2 := cardValue(hand[1])
	v3 := cardValue(hand[2])
	return v1+v2+v3 == 21
}

func main() {
	deck := newDeck()
	shuffle(deck)

	for {
		if len(deck) < 10 { // ricrea e mischia se il mazzo è quasi finito
			deck = newDeck()
			shuffle(deck)
		}

		// inizializza mani
		player := []Card{deck[0], deck[1]}
		dealer := []Card{deck[2], deck[3]}
		deck = deck[4:]

		// mostra carte iniziali
		fmt.Println("\nPlayer:", player, "Value:", handValue(player))
		fmt.Println("Dealer:", dealer[0], "X")

		// Side bet semplice
		if perfectPair(player) {
			fmt.Println("Perfect Pair Side Bet Wins!")
		}
		if len(player) >= 3 && twentyOneThree(player) {
			fmt.Println("21+3 Side Bet Wins!")
		}

		// Loop del giocatore
		for handValue(player) < 21 {
			fmt.Print("Hit or Stand? (h/s): ")
			var action string
			fmt.Scan(&action)
			if action == "s" {
				break
			} else if action == "h" {
				player = append(player, deck[0])
				deck = deck[1:]
				fmt.Println("Player hits:", player[len(player)-1])
				fmt.Println("Hand Value:", handValue(player))
			}
		}

		// Logica dealer
		for handValue(dealer) < 17 {
			dealer = append(dealer, deck[0])
			deck = deck[1:]
		}

		// Mostra esito
		fmt.Println("Player final:", player, "Value:", handValue(player))
		fmt.Println("Dealer final:", dealer, "Value:", handValue(dealer))

		// Determina vincitore
		pv := handValue(player)
		dv := handValue(dealer)
		if pv > 21 {
			fmt.Println("Player busts. Dealer wins.")
		} else if dv > 21 || pv > dv {
			fmt.Println("Player wins!")
		} else if pv < dv {
			fmt.Println("Dealer wins.")
		} else {
			fmt.Println("Push (tie).")
		}

		// Chiedi se continuare
		fmt.Print("Play again? (y/n): ")
		var again string
		fmt.Scan(&again)
		if again != "y" {
			break
		}
	}
}
