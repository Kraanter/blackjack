package blackjack

import (
	"testing"
)

func playGameWithCards(playerCards, dealerCards []*Card) (*BlackjackGame, *Player) {
	game := CreateGame()
	player := game.AddPlayerWithBalance(10)

	game.GameState.Set(BettingState)
	game.SetPlayerBet(player.PlayerNum, 10)

	hand := player.GetActiveHand()
	hand.Cards = playerCards
	hand.lock()

	game.hiddenDealerCard.Set(dealerCards[0])
	game.Dealer.Cards = dealerCards[1:]

	game.GameState.Set(PlayingState)
	game.gameloopTick()

	game.finishRound()

	return game, player
}

func playGameWithSplit(playerCards, splitCards, dealerCards []*Card) (*BlackjackGame, *Player) {
	game := CreateGame()
	player := game.AddPlayerWithBalance(10)

	game.GameState.Set(BettingState)
	game.SetPlayerBet(player.PlayerNum, 5)

	hand := player.GetActiveHand()
	hand.Cards = playerCards
	game.GameState.Set(PlayingState)

	game.hiddenDealerCard.Set(dealerCards[0])
	game.Dealer.Cards = dealerCards[1:]

	game.PlayerSplit(player.PlayerNum)
	player.GetActiveHand().Cards[1] = splitCards[0]
	game.PlayerStand(player.PlayerNum)
	player.GetActiveHand().Cards[1] = splitCards[1]
	game.PlayerStand(player.PlayerNum)

	return game, player
}

func TestPayoutWithDealerAndPlayerBlackjackIsPush(t *testing.T) {
	wants := uint(10)
	playerCards := []*Card{CreateCard(Ace, Spades), CreateCard(Ten, Spades)}
	dealerCards := []*Card{CreateCard(Ace, Hearts), CreateCard(Ten, Hearts)}

	game, player := playGameWithCards(playerCards, dealerCards)

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState.Get(), NoState)
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

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState.Get(), NoState)
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

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState.Get(), NoState)
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

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState.Get(), NoState)
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

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout", game.GameState.Get(), NoState)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of losing a game", balance, wants)
	}
}

func TestPayoutWithSplitCardsPlayerWinningBoth(t *testing.T) {
	wants := uint(20)
	playerCards := []*Card{CreateCard(Ten, Spades), CreateCard(Ten, Hearts)}
	splitCards := []*Card{CreateCard(Ten, Clubs), CreateCard(Ten, Diamonds)}
	dealerCards := []*Card{CreateCard(Seven, Spades), CreateCard(Ten, Spades)}

	game, player := playGameWithSplit(playerCards, splitCards, dealerCards)

	if game.GameState.Get() != NoState {
		t.Fatalf("game.GameState = %v, expected state to be %v after payout,  game: (%v)", game.GameState.Get(), NoState, game)
	}

	balance := player.Balance
	if balance != wants {
		t.Fatalf("player.Balance = %v, wants balance to be %v after payout of losing a game", balance, wants)
	}
}
