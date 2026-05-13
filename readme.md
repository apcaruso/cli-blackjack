# Blackjack

A small terminal blackjack game written in Go.

## Requirements

- Go 1.20 or newer

## Run

```sh
go run blackjack.go
```

The game first asks for a starting bankroll:

```text
Set starting bankroll:
```

Before each match, enter a whole-number bet that does not exceed your current bankroll:

```text
Enter bet, (b)ankroll, (stats), or (q)uit:
```

The game then deals an initial hand to the player and dealer, then prompts the player to hit or stand:

```text
Hit, Stand, (b)ankroll, or (stats)? (h/s/b/stats):
```

After each match, enter `y` to play again, `n` or `q` to quit, `b` to view your bankroll, or `stats` to view statistics for all completed matches in the current session.

## Build

```sh
go build -o blackjack blackjack.go
```

Then run the compiled binary:

```sh
./blackjack
```

## Gameplay Notes

- The game includes simple Perfect Pair and 21+3 side bet checks.
- Bankroll increases by the bet amount on a player win, decreases by the bet amount on a dealer win, and stays unchanged on a push.
- Session statistics track matches played, wins, losses, pushes, busts, bankroll change, and total wagered.
