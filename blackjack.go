package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Suit string
	Rank string
}

func (c Card) String() string {
	return c.Rank + c.Suit
}

type Hand []Card

func (h Hand) String() string {
	cards := make([]string, len(h))
	for i, c := range h {
		cards[i] = c.String()
	}
	return strings.Join(cards, " ")
}

func (h Hand) Value() int {
	total, aces := 0, 0

	for _, c := range h {
		switch c.Rank {
		case "A":
			total += 11
			aces++
		case "K", "Q", "J":
			total += 10
		default:
			n, _ := strconv.Atoi(c.Rank)
			total += n
		}
	}

	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}

	return total
}

func (h Hand) Blackjack() bool {
	return len(h) == 2 && h.Value() == 21
}

func (h Hand) Bust() bool {
	return h.Value() > 21
}

type Stats struct {
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

type Game struct {
	deck  []Card
	input *bufio.Scanner
	rng   *rand.Rand
	stats Stats
}

var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

func main() {
	game := NewGame()
	game.Shuffle()

	if !game.AskBankroll() {
		return
	}

	for game.stats.Bankroll > 0 {
		if !game.PlayRound() {
			break
		}
	}

	game.PrintStats()
}

func NewGame() *Game {
	return &Game{
		input: bufio.NewScanner(os.Stdin),
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewDeck() []Card {
	deck := make([]Card, 0, len(suits)*len(ranks))

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}

	return deck
}

func (g *Game) Shuffle() {
	g.deck = NewDeck()

	g.rng.Shuffle(len(g.deck), func(i, j int) {
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	})
}

func (g *Game) EnsureDeck(cards int) {
	if len(g.deck) >= cards {
		return
	}

	g.Shuffle()
	fmt.Println("\nDeck reshuffled.")
}

func (g *Game) Draw() Card {
	if len(g.deck) == 0 {
		g.Shuffle()
	}

	card := g.deck[0]
	g.deck = g.deck[1:]
	return card
}

func (g *Game) Ask(prompt string) (string, bool) {
	fmt.Print(prompt)

	if !g.input.Scan() {
		return "", false
	}

	return strings.ToLower(strings.TrimSpace(g.input.Text())), true
}

func (g *Game) AskBankroll() bool {
	for {
		text, ok := g.Ask("Starting bankroll: ")
		if !ok {
			return false
		}

		amount, err := strconv.Atoi(text)
		if err == nil && amount > 0 {
			g.stats.StartingBankroll = amount
			g.stats.Bankroll = amount
			return true
		}

		fmt.Println("Enter a positive whole number.")
	}
}

func (g *Game) AskBet() (int, bool) {
	for {
		fmt.Printf("\nBankroll: %d\n", g.stats.Bankroll)

		text, ok := g.Ask("Bet amount, [b]ankroll, [stats], [q]uit: ")
		if !ok {
			return 0, false
		}

		switch text {
		case "q", "quit":
			return 0, false

		case "b", "bankroll":
			continue

		case "stats":
			g.PrintStats()

		default:
			bet, err := strconv.Atoi(text)
			if err != nil || bet <= 0 {
				fmt.Println("Enter a positive whole-number bet.")
				continue
			}

			if bet > g.stats.Bankroll {
				fmt.Println("Bet cannot exceed bankroll.")
				continue
			}

			return bet, true
		}
	}
}

func (g *Game) PlayRound() bool {
	bet, ok := g.AskBet()
	if !ok {
		return false
	}

	g.EnsureDeck(15)

	player := Hand{g.Draw()}
	dealer := Hand{g.Draw()}

	player = append(player, g.Draw())
	dealer = append(dealer, g.Draw())

	fmt.Printf("\nRound %d | Bet %d\n", g.stats.Matches+1, bet)
	PrintTable(player, dealer, false)

	PrintSideChecks(player, dealer[0])

	if player.Blackjack() || dealer.Blackjack() {
		PrintTable(player, dealer, true)
		g.Settle(player, dealer, bet)
		return g.AskNextRound()
	}

	if !g.PlayerTurn(&player, bet) {
		return false
	}

	if !player.Bust() {
		g.DealerTurn(&dealer)
	}

	PrintTable(player, dealer, true)
	g.Settle(player, dealer, bet)

	return g.AskNextRound()
}

func (g *Game) PlayerTurn(player *Hand, bet int) bool {
	for {
		value := player.Value()

		if value == 21 {
			fmt.Println("Player has 21. Standing.")
			return true
		}

		if value > 21 {
			fmt.Println("Player busts.")
			return true
		}

		text, ok := g.Ask("Action: [h]it, [s]tand, [b]ankroll, [stats], [q]uit. Enter = stand: ")
		if !ok {
			return false
		}

		switch text {
		case "", "s", "stand":
			return true

		case "h", "hit":
			card := g.Draw()
			*player = append(*player, card)
			fmt.Printf("Draw: %s\n", card)
			PrintHand("Player", *player)

		case "b", "bankroll":
			fmt.Printf("Bankroll: %d | Active bet: %d\n", g.stats.Bankroll, bet)

		case "stats":
			g.PrintStats()

		case "q", "quit":
			return false

		default:
			fmt.Println("Invalid action.")
		}
	}
}

func (g *Game) DealerTurn(dealer *Hand) {
	fmt.Println("\nDealer turn.")

	for dealer.Value() < 17 {
		card := g.Draw()
		*dealer = append(*dealer, card)
		fmt.Printf("Dealer draws: %s\n", card)
	}
}

func (g *Game) Settle(player, dealer Hand, bet int) {
	pv, dv := player.Value(), dealer.Value()

	g.stats.Matches++
	g.stats.TotalWagered += bet

	switch {
	case player.Blackjack() && dealer.Blackjack():
		g.stats.Pushes++
		fmt.Println("Both have blackjack. Push.")

	case player.Blackjack():
		g.stats.PlayerWins++
		g.stats.Bankroll += bet
		fmt.Println("Blackjack. Player wins.")

	case dealer.Blackjack():
		g.stats.DealerWins++
		g.stats.Bankroll -= bet
		fmt.Println("Dealer has blackjack. Dealer wins.")

	case player.Bust():
		g.stats.PlayerBusts++
		g.stats.DealerWins++
		g.stats.Bankroll -= bet
		fmt.Println("Dealer wins.")

	case dealer.Bust():
		g.stats.DealerBusts++
		g.stats.PlayerWins++
		g.stats.Bankroll += bet
		fmt.Println("Dealer busts. Player wins.")

	case pv > dv:
		g.stats.PlayerWins++
		g.stats.Bankroll += bet
		fmt.Println("Player wins.")

	case pv < dv:
		g.stats.DealerWins++
		g.stats.Bankroll -= bet
		fmt.Println("Dealer wins.")

	default:
		g.stats.Pushes++
		fmt.Println("Push.")
	}

	fmt.Printf("Bankroll: %d\n", g.stats.Bankroll)
}

func (g *Game) AskNextRound() bool {
	if g.stats.Bankroll <= 0 {
		fmt.Println("Bankroll depleted. Game over.")
		return false
	}

	for {
		text, ok := g.Ask("\nNext round? [y]es, [n]o, [b]ankroll, [stats]: ")
		if !ok {
			return false
		}

		switch text {
		case "", "y", "yes":
			return true

		case "n", "no", "q", "quit":
			return false

		case "b", "bankroll":
			fmt.Printf("Bankroll: %d\n", g.stats.Bankroll)

		case "stats":
			g.PrintStats()

		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func PrintTable(player, dealer Hand, revealDealer bool) {
	fmt.Println()
	PrintHand("Player", player)

	if revealDealer {
		PrintHand("Dealer", dealer)
		return
	}

	fmt.Printf("Dealer: %s X\n", dealer[0])
}

func PrintHand(label string, hand Hand) {
	fmt.Printf("%s: %s (%d)\n", label, hand, hand.Value())
}

func PrintSideChecks(player Hand, dealerUp Card) {
	if len(player) < 2 {
		return
	}

	if player[0].Rank == player[1].Rank {
		fmt.Println("Side check: Perfect Pair.")
	}

	if Value(player[0])+Value(player[1])+Value(dealerUp) == 21 {
		fmt.Println("Side check: 21+3.")
	}
}

func Value(card Card) int {
	switch card.Rank {
	case "A":
		return 11
	case "K", "Q", "J":
		return 10
	default:
		n, _ := strconv.Atoi(card.Rank)
		return n
	}
}

func (g *Game) PrintStats() {
	s := g.stats

	fmt.Println("\nSession statistics")
	fmt.Println("------------------")
	fmt.Printf("Matches played: %d\n", s.Matches)
	fmt.Printf("Player wins: %d (%.1f%%)\n", s.PlayerWins, Percent(s.PlayerWins, s.Matches))
	fmt.Printf("Dealer wins: %d (%.1f%%)\n", s.DealerWins, Percent(s.DealerWins, s.Matches))
	fmt.Printf("Pushes: %d (%.1f%%)\n", s.Pushes, Percent(s.Pushes, s.Matches))
	fmt.Printf("Player busts: %d\n", s.PlayerBusts)
	fmt.Printf("Dealer busts: %d\n", s.DealerBusts)
	fmt.Printf("Starting bankroll: %d\n", s.StartingBankroll)
	fmt.Printf("Current bankroll: %d\n", s.Bankroll)
	fmt.Printf("Net change: %+d\n", s.Bankroll-s.StartingBankroll)
	fmt.Printf("Total wagered: %d\n", s.TotalWagered)
}

func Percent(part, total int) float64 {
	if total == 0 {
		return 0
	}

	return float64(part) * 100 / float64(total)
}
