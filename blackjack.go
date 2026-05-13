package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Suit  string
	Value string
}

var suits = []string{"♠", "♥", "♦", "♣"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type GameStats struct {
	StartingBankroll int
	Bankroll         int
	Matches          int
	PlayerWins       int
	DealerWins       int
	Pushes           int
	PlayerBusts      int
	DealerBusts      int
	TotalWagered     int
}

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

func readInput() (string, bool) {
	var input string
	if _, err := fmt.Scan(&input); err != nil {
		return "", false
	}
	return strings.ToLower(strings.TrimSpace(input)), true
}

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func printBankroll(bankroll int) {
	fmt.Printf("Bankroll: %d\n", bankroll)
}

func printStats(stats GameStats) {
	fmt.Println("\nSession statistics")
	fmt.Println("------------------")
	fmt.Printf("Matches played: %d\n", stats.Matches)
	fmt.Printf("Player wins: %d (%.1f%%)\n", stats.PlayerWins, percentage(stats.PlayerWins, stats.Matches))
	fmt.Printf("Dealer wins: %d (%.1f%%)\n", stats.DealerWins, percentage(stats.DealerWins, stats.Matches))
	fmt.Printf("Pushes: %d (%.1f%%)\n", stats.Pushes, percentage(stats.Pushes, stats.Matches))
	fmt.Printf("Player busts: %d\n", stats.PlayerBusts)
	fmt.Printf("Dealer busts: %d\n", stats.DealerBusts)
	fmt.Printf("Starting bankroll: %d\n", stats.StartingBankroll)
	fmt.Printf("Current bankroll: %d\n", stats.Bankroll)
	fmt.Printf("Net change: %+d\n", stats.Bankroll-stats.StartingBankroll)
	fmt.Printf("Total wagered: %d\n", stats.TotalWagered)
}

func promptStartingBankroll() (int, bool) {
	for {
		fmt.Print("Set starting bankroll: ")
		input, ok := readInput()
		if !ok {
			return 0, false
		}

		bankroll, err := strconv.Atoi(input)
		if err != nil || bankroll <= 0 {
			fmt.Println("Enter a positive whole number.")
			continue
		}

		return bankroll, true
	}
}

func promptBet(stats GameStats) (int, bool) {
	for {
		fmt.Print("Enter bet, (b)ankroll, (stats), or (q)uit: ")
		input, ok := readInput()
		if !ok {
			return 0, false
		}

		switch input {
		case "b", "bankroll":
			printBankroll(stats.Bankroll)
			continue
		case "stats":
			printStats(stats)
			continue
		case "q", "quit":
			return 0, false
		}

		bet, err := strconv.Atoi(input)
		if err != nil || bet <= 0 {
			fmt.Println("Enter a positive whole-number bet.")
			continue
		}
		if bet > stats.Bankroll {
			fmt.Println("Bet cannot exceed your current bankroll.")
			continue
		}

		return bet, true
	}
}

func promptPlayerAction(stats GameStats, activeBet int) (string, bool) {
	for {
		fmt.Print("Hit, Stand, (b)ankroll, or (stats)? (h/s/b/stats): ")
		input, ok := readInput()
		if !ok {
			return "", false
		}

		switch input {
		case "h", "hit":
			return "h", true
		case "s", "stand":
			return "s", true
		case "b", "bankroll":
			fmt.Printf("Bankroll: %d (active bet: %d)\n", stats.Bankroll, activeBet)
		case "stats":
			printStats(stats)
		default:
			fmt.Println("Enter h, s, b, or stats.")
		}
	}
}

func promptPlayAgain(stats GameStats) bool {
	for {
		fmt.Print("Play again, (b)ankroll, (stats), or (q)uit? (y/n/b/stats): ")
		input, ok := readInput()
		if !ok {
			return false
		}

		switch input {
		case "y", "yes":
			return true
		case "n", "no", "q", "quit":
			return false
		case "b", "bankroll":
			printBankroll(stats.Bankroll)
		case "stats":
			printStats(stats)
		default:
			fmt.Println("Enter y, n, b, or stats.")
		}
	}
}

func main() {
	deck := newDeck()
	shuffle(deck)

	startingBankroll, ok := promptStartingBankroll()
	if !ok {
		return
	}
	stats := GameStats{
		StartingBankroll: startingBankroll,
		Bankroll:         startingBankroll,
	}

	for {
		if stats.Bankroll <= 0 {
			fmt.Println("Bankroll depleted. Game over.")
			break
		}

		bet, ok := promptBet(stats)
		if !ok {
			break
		}

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
			action, ok := promptPlayerAction(stats, bet)
			if !ok {
				printStats(stats)
				return
			}
			if action == "s" {
				break
			}

			player = append(player, deck[0])
			deck = deck[1:]
			fmt.Println("Player hits:", player[len(player)-1])
			fmt.Println("Hand Value:", handValue(player))
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
		stats.Matches++
		stats.TotalWagered += bet
		if pv > 21 {
			stats.PlayerBusts++
		}
		if dv > 21 {
			stats.DealerBusts++
		}
		if pv > 21 {
			stats.DealerWins++
			stats.Bankroll -= bet
			fmt.Println("Player busts. Dealer wins.")
		} else if dv > 21 || pv > dv {
			stats.PlayerWins++
			stats.Bankroll += bet
			fmt.Println("Player wins!")
		} else if pv < dv {
			stats.DealerWins++
			stats.Bankroll -= bet
			fmt.Println("Dealer wins.")
		} else {
			stats.Pushes++
			fmt.Println("Push (tie).")
		}
		printBankroll(stats.Bankroll)

		if stats.Bankroll <= 0 {
			fmt.Println("Bankroll depleted. Game over.")
			break
		}

		// Chiedi se continuare
		if !promptPlayAgain(stats) {
			break
		}
	}

	printStats(stats)
}
