package blackjack

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"
)

func TestGameLoop100_000Rounds(t *testing.T) {
	for range 10000 {
		game := CreateGame()
		players := make([]*Player, 0, 3)

		players = append(players, game.AddPlayerWithBalance(10))
		players = append(players, game.AddPlayerWithBalance(10))
		players = append(players, game.AddPlayerWithBalance(10))

		gameFinishedContext, testFinished := context.WithTimeout(context.Background(), time.Millisecond)
		defer testFinished()

		log := ""
		game.OnGameUpdate = func(game *BlackjackGame) {
			log += fmt.Sprintf("\n---\ngame_update: %v\n\nplayers:\n", game.GameState.Get())

			for _, player := range players {
				log += fmt.Sprintln(player.String())
			}

			log += fmt.Sprintln("Dealer: ", game.Dealer.String())
		}

		game.OnGameFinished = func(payout map[PlayerId]uint) { testFinished() }

		game.OnPlayerTurn = func(pi PlayerId) {
			log += fmt.Sprintf("Player (%v)'s turn", pi)
			game.PlayerHit(pi)
			game.PlayerStand(pi)
		}

		game.Initialize()

		err := game.SetPlayerBet(players[0].PlayerNum, 5)
		err = game.SetPlayerBet(players[2].PlayerNum, 2)
		if err != nil {
			t.Fatalf("Setting player bets went wrong: %v", err.Error())
		}
		err = game.SkipPlayerBet(players[1].PlayerNum)
		if err != nil {
			t.Fatalf("Skipping player bet went wrong: %v", err.Error())
		}

		err = game.Start(gameFinishedContext)
		if err != nil {
			t.Fatalf("Game failed to start after setting bets: %s\nlog:\n%s", err, log)
		}

		// Wait for all players to have their turn
		<-gameFinishedContext.Done()

		if game.GameState.Get() != NoState {
			t.Fatalf("Game should be in no state (%v) after a match is finished, was %v\nlog:\n%s", NoState, game.GameState.Get(), log)
		}
		if game.Dealer != nil {
			t.Fatalf("Dealer should not have a hand at the end of the game")
		}
		if slices.ContainsFunc(players, func(player *Player) bool { return len(player.Hands) != 0 }) {
			t.Fatalf("All players should not have a hand at the end of the game")
		}
	}
}

func playHand(t *testing.T, playerHand, dealerHand *Hand) *BlackjackGame {
	game := CreateGame()

	player := game.AddPlayerWithBalance(10)

	game.Initialize()

	err := game.SetPlayerBet(player.PlayerNum, 10)
	if err != nil {
		t.Fatalf("Setting player bet went wrong: %v", err.Error())
	}

	log := ""
	game.OnGameUpdate = func(game *BlackjackGame) {
		log += fmt.Sprintf("\n---\ngame_update: %v\n\nplayers:\n", game.GameState.Get())

		log += fmt.Sprintln(player.String())

		log += fmt.Sprintln("Dealer: ", game.Dealer.String())
	}

	player.Hands = []*Hand{playerHand}
	game.hiddenDealerCard = dealerHand.Cards[0]
	dealerHand.Cards = dealerHand.Cards[1:]
	game.Dealer = dealerHand
	game.DealInitialCards()
	game.GameState.Set(PlayingState)
	game.sendGameUpdate()

	game.OnPlayerTurn = func(pi PlayerId) {
		if player, ok := game.GetPlayer(pi); ok {
			fmt.Printf("\n--- On player turn \nplayer: (%s)\n---\n\n", player)
		} else {
			fmt.Printf("Not ok player turn")
		}
	}

	testFinishedContext, testsFinished := context.WithTimeout(context.Background(), time.Millisecond)
	defer testsFinished()
	game.OnGameFinished = func(m map[PlayerId]uint) { testsFinished() }

	<-testFinishedContext.Done()

	if game.GameState.Get() != NoState {
		t.Fatalf("Game should be in no state (%v) after a match is finished, was %v\nlog:%s", NoState, game.GameState.Get(), log)
	}
	if game.Dealer != nil {
		t.Fatalf("Dealer should not have a hand at the end of the game\nlog:%s", log)
	}
	if slices.ContainsFunc([]*Player{player}, func(player *Player) bool { return len(player.Hands) != 0 }) {
		t.Fatalf("All players should not have a hand at the end of the game\nlog:%s", log)
	}

	return game
}

func handWithCards(betSize uint, cards ...*Card) *Hand {
	hand := CreateHand(betSize)
	for _, v := range cards {
		hand.AddCard(v)
	}
	hand.lock()

	return hand
}

func TestGameLoopBlackjackHandPlayer(t *testing.T) {
	betSize := uint(10)
	playerHand := handWithCards(betSize, CreateCard(King, Hearts), CreateCard(Ace, Hearts))
	dealerHand := handWithCards(0, CreateCard(King, Spades), CreateCard(Seven, Spades))
	game := playHand(t, playerHand, dealerHand)

	expected := uint(25)
	game.forEachPlayer(func(_ PlayerId, player *Player) {
		actual := player.Balance
		if actual != expected {
			t.Fatalf("Player (%s) should have get payed out 2.5x (%v), was %v", player.String(), expected, actual)
		}
	})
}
