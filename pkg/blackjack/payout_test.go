package blackjack

import (
	"testing"
)

func playGameWithCards(playerCards []*Card, dealerCards []*Card) (*BlackjackGame, *Player) {
	game := CreateGame()
	player := game.AddPlayerWithBalance(10)

	game.GameState = BettingState
	game.SetPlayerBet(player.PlayerNum, 10)

	player.Hand.Cards = playerCards
	player.Hand.lock()
	game.Dealer.Cards = dealerCards
	game.Dealer.lock()

	game.GameState = PlayingState
	game.sendGameUpdate()

	game.finishRound()

	return game, player
}

func TestPayoutWithDealerAndPlayerBlackjackIsPush(t *testing.T) {
	wants := uint(10)
	playerCards := []*Card{CreateCard(Ace, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Ace, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState, NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of blackjack push", balance, wants)
	}
}

func TestPayoutWithPlayerBlackjackIsPayed2To3(t *testing.T) {
	wants := uint(25)
	playerCards := []*Card{CreateCard(Ace, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Ten, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState, NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of blackjack win", balance, wants)
	}
}

func TestPayoutWithDealerBlackjack(t *testing.T) {
	wants := uint(0)
	playerCards := []*Card{CreateCard(Ten, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Ace, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState, NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of losing to blackjack", balance, wants)
	}
}

func TestPayoutWithNoBlackjackPlayerWinning(t *testing.T) {
	wants := uint(20)
	playerCards := []*Card{CreateCard(Ten, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Nine, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState, NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of winning a game", balance, wants)
	}
}

func TestPayoutWithNoBlackjackDealerWinning(t *testing.T) {
	wants := uint(0)
	playerCards := []*Card{CreateCard(Nine, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Ten, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState, NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of losing a game", balance, wants)
	}
}
