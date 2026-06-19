# Blackjack

A small terminal blackjack game written in Go.

The game runs entirely in the terminal and provides a simple blackjack flow: set a bankroll, place a bet, play each round, and review session statistics at any time.

## Requirements

* Go 1.20 or newer

## Run

```sh
go run blackjack.go
```

## Build

```sh
go build -o blackjack blackjack.go
```

Then run the compiled binary:

```sh
./blackjack
```

## Gameplay Flow

At startup, enter a positive whole-number bankroll:

```text
Starting bankroll:
```

Before each round, the game shows your current bankroll and asks for a bet:

```text
Bankroll: 10000
Bet amount, [b]ankroll, [stats], [q]uit:
```

Available commands:

* Enter a number to place a bet
* `b` or `bankroll` to view bankroll
* `stats` to view session statistics
* `q` or `quit` to exit

After the initial deal, the player sees both player cards and only the dealer's visible card:

```text
Round 1 | Bet 100

Player: J♥ Q♣ (20)
Dealer: 7♥ X
```

During the player's turn:

```text
Action: [h]it, [s]tand, [b]ankroll, [stats], [q]uit. Enter = stand:
```

Available commands:

* `h` or `hit` to draw another card
* `s` or `stand` to stand
* Press `Enter` to stand
* `b` or `bankroll` to view bankroll and active bet
* `stats` to view session statistics
* `q` or `quit` to exit

After each round:

```text
Next round? [y]es, [n]o, [b]ankroll, [stats]:
```

Available commands:

* `y` or `yes` to play another round
* Press `Enter` to play another round
* `n`, `no`, `q`, or `quit` to stop
* `b` or `bankroll` to view bankroll
* `stats` to view session statistics

## Gameplay Rules

* The player and dealer each start with two cards.
* The dealer shows one card until the round is resolved.
* The player may hit or stand.
* Pressing `Enter` during the player's turn stands.
* If the player reaches 21, the player automatically stands.
* If the player busts, the dealer does not draw.
* If either side has a natural blackjack, the round is resolved immediately.
* The dealer draws until reaching at least 17.
* Aces count as 11 unless that would bust the hand; then they count as 1.
* The deck is reshuffled automatically when there are not enough cards reserved for the next round.

## Side Checks

The game includes two simple side checks:

* **Perfect Pair**: the player's first two cards have the same rank.
* **21+3**: the player's first two cards plus the dealer's visible card total 21.

These checks only print a message. They do not currently affect bankroll payout.

## Bankroll Rules

* Player win: bankroll increases by the bet amount.
* Dealer win: bankroll decreases by the bet amount.
* Push: bankroll stays unchanged.
* If bankroll reaches zero, the session ends.

## Session Statistics

The game tracks:

* Matches played
* Player wins
* Dealer wins
* Pushes
* Player busts
* Dealer busts
* Starting bankroll
* Current bankroll
* Net bankroll change
* Total wagered
