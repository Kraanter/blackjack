package blackjack

import (
	"context"
	"fmt"
)

type GameState = int

const (
	NoState GameState = iota
	BettingState
	DealingState
	PlayingState
	DealerState
	PayoutState
)

var NoPlayersInGameError = fmt.Errorf("Game does not have any players")

func (b *BlackjackGame) Initialize() {
	if b.GameState.Get() != NoState {
		return
	}

	b.bettingStateChannel = make(chan struct{})
	b.GameState.Set(BettingState)
	b.gameloopTick()
}

// This starts the game until all bets are in
// Then the gameplay loop is managed through the players
// This ends when the game broadcasts that it goes into PlayingState
func (b *BlackjackGame) Start(ctx context.Context) error {
	if b.GameState.Get() == NoState {
		b.Initialize()
	}

	if b.GameState.Get() != BettingState {
		return b.createGameStateError(BettingState)
	}

	if err := b.WaitUntilBettingFinished(ctx); err != nil {
		return err
	}

	if len(b.PlayerMap) == 0 {
		return NoPlayersInGameError
	}

	b.DealInitialCards()
	b.GameState.Set(PlayingState)
	b.gameloopTick()

	return nil
}

func (b *BlackjackGame) WaitUntilBettingFinished(ctx context.Context) error {
	select {
	case <-b.bettingStateChannel:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *BlackjackGame) DealInitialCards() {
	if b.GameState.Get() == PlayingState {
		return
	}

	b.Dealer.Set(CreateHand(0))

	for i := 0; i < 2; i++ {
		for _, player := range b.PlayerMap {
			if hand := player.GetActiveHand(); hand != nil {
				b.dealCard(hand)
			}
		}

		if len(b.Dealer.Get().Cards) == 1 {
			b.hiddenDealerCard.Set(b.shoe.DrawCard())
		} else {
			b.dealCard(b.Dealer.Get())
		}
	}
}

// Returns true if player got a new card
// Returns false if game is not in state for player to receive card
// Returns error if any other reason like player could not be found
func (game *BlackjackGame) PlayerHit(playerNum PlayerId) (bool, error) {
	if game.GameState.Get() != PlayingState {
		return false, game.createGameStateError(PlayingState)
	}

	if isDealer, num := game.nextPlayersTurn(); isDealer || num != playerNum {
		return false, WrongGameStateError
	}

	player, ok := game.GetPlayer(playerNum)
	if !ok {
		return false, PlayerNotFoundError
	}

	hand := player.GetActiveHand()
	if hand == nil {
		return false, nil
	}
	ok = game.dealCard(hand)

	return ok, nil
}

func (game *BlackjackGame) PlayerStand(playerNum PlayerId) error {
	if isDealer, num := game.nextPlayersTurn(); isDealer || num != playerNum {
		return WrongGameStateError
	}

	player, ok := game.GetPlayer(playerNum)
	if !ok {
		return PlayerNotFoundError
	}

	player.stand()
	game.gameloopTick()

	return nil
}

func (game *BlackjackGame) PlayerSplit(playerNum PlayerId) error {
	if isDealer, num := game.nextPlayersTurn(); isDealer || num != playerNum {
		return WrongGameStateError
	}

	player, ok := game.GetPlayer(playerNum)
	if !ok {
		return PlayerNotFoundError
	}

	hand := player.GetActiveHand()
	if hand == nil || hand.locked {
		return WrongGameStateError
	}

	canSplit := len(player.GetActiveHand().Cards) == 2 && hand.Cards[0].Face == hand.Cards[1].Face && player.Balance >= hand.Bet
	if !canSplit {
		return WrongGameStateError
	}

	secondCard := hand.Cards[1]
	newHand := CreateHand(hand.Bet)
	player.Balance -= hand.Bet

	newHand.Cards = append(newHand.Cards, secondCard)
	hand.Cards = hand.Cards[:1]
	player.Hands = append(player.Hands, newHand)

	game.sendGameUpdate()

	game.dealCard(hand)
	game.dealCard(newHand)

	game.gameloopTick()

	return nil
}

func (game *BlackjackGame) DealerTurn() error {
	if game.GameState.Get() != DealerState {
		return game.createGameStateError(DealerState)
	}

	hiddenCard := game.hiddenDealerCard.Get()
	if hiddenCard == nil {
		return fmt.Errorf("No dealer card assigned: %w", WrongGameStateError)
	}

	// Show the hidden dealer card
	game.Dealer.Get().AddCard(hiddenCard)
	game.hiddenDealerCard.Set(nil)
	game.sendGameUpdate()

	for shouldDealerDrawCard(game.Dealer.Get()) {
		game.dealCard(game.Dealer.Get())
	}

	game.Dealer.Get().lock()
	game.gameloopTick()

	game.finishRound()

	return nil
}

func shouldDealerDrawCard(hand *Hand) bool {
	if hand.IsLocked() || hand == nil {
		return false
	}

	total := hand.Total()
	if total > 16 {
		return false
	}

	return true
}

func (game *BlackjackGame) dealCard(hand *Hand) bool {
	if hand == nil {
		return false
	}

	card := game.shoe.DrawCard()
	ok := hand.AddCard(card)

	game.gameloopTick()

	return ok
}

func (game *BlackjackGame) finishRound() {
	game.GameState.Set(PayoutState)
	game.sendGameUpdate()
	game.reset()
}
