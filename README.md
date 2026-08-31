# Blackjack

Concurrent blackjack engine with a small HTTP and Server-Sent Events (SSE) API.

## Run

```sh
go run ./cmd/restAPI
```

The server listens on port `42069` by default. Use `-port` to change it:

```sh
go run ./cmd/restAPI -port 8080
```

`cmd/restAPI` spaces outbound SSE updates by 150ms. The game engine itself has no tick timer: a game advances when a player places a bet or performs an action.

## Core Usage

Create a game, add players, collect bets, then start the round once everyone has either bet or skipped:

```go
game := blackjack.CreateGame()
alice := game.AddPlayerWithBalance(100)
bob := game.AddPlayerWithBalance(100)

game.OnGameUpdate = func(snapshot *blackjack.GameSnapshot) {
	// Safe to retain, serialize, or send from another goroutine.
	fmt.Println(snapshot.GameState, snapshot.CurrentTurn)
}

if err := game.SetPlayerBet(alice.PlayerNum, 10); err != nil {
	panic(err)
}
if err := game.SetPlayerBet(bob.PlayerNum, 10); err != nil {
	panic(err)
}
if err := game.Start(context.Background()); err != nil {
	panic(err)
}
```

Set `OnPlayerTurn` before starting if a program controls the players:

```go
game.OnPlayerTurn = func(id blackjack.PlayerId) {
	_, _ = game.PlayerHit(id)
	_ = game.PlayerStand(id)
}
```

For interactive players, call `PlayerHit`, `PlayerStand`, or `PlayerSplit` for the player whose ID equals `GameSnapshot.CurrentTurn`. Invalid, duplicate, or out-of-turn actions return an error.

`Start` begins a round only while the game is in betting state and every current player has bet or skipped. It returns without starting when more bets are needed. `manager.ManagedPlayer.Bet` and `SkipBet` start the round automatically after the final decision.

## States And Rules

Game state values are `NoState` (0), `BettingState` (1), `DealingState` (2), `PlayingState` (3), `DealerState` (4), and `PayoutState` (5). The current implementation transitions directly from betting to playing while dealing the initial cards.

Players act in ascending player ID order. Split hands are played in their slice order. Once every active player hand is locked, the dealer reveals the hole card and draws through 16. A blackjack pays 3:2, a win pays 1:1, and a push returns the bet.

## Concurrency

Use game methods to inspect or mutate a live game. They serialize core mutations. `OnGameUpdate` receives an immutable `GameSnapshot`; use that rather than retaining the live `BlackjackGame`, `Player`, or `Hand` pointers. This also prevents network delivery from delaying gameplay.

## HTTP API

Start with `POST /join`, retain the `PlayerId` cookie, then subscribe to `GET /gamestate?sse`. The stream sends `initial` and `update` events containing game snapshots. Use `POST /hand` to bet and `PATCH /hit`, `PATCH /stand`, or `PATCH /split` on your turn.

The full machine-readable contract is in [openapi.yaml](openapi.yaml). The server sets a `Secure`, HTTP-only cookie, so browsers must use HTTPS for authenticated requests.

## Checks

```sh
go test ./...
go test -race ./...
```
