# Blackjack

A small terminal blackjack game written in Go.

## Requirements

- Go 1.20 or newer

## Run

```sh
go run blackjack.go
```

The game deals an initial hand to the player and dealer, then prompts the player to hit or stand:

```text
Hit or Stand? (h/s):
```

After each round, enter `y` to play again or any other value to quit.

## Build

```sh
go build -o blackjack blackjack.go
```

Then run the compiled binary:

```sh
./blackjack
```

## Gameplay Notes

- Aces count as 11 unless that would bust the hand, then they count as 1.
- The dealer draws until reaching at least 17.
- The game includes simple Perfect Pair and 21+3 side bet checks.
