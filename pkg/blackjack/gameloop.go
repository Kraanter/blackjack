package blackjack

import (
	"time"
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

func (b *BlackjackGame) Start() {
	b.GameState = BettingState
	b.sendGameUpdate()

	collectionChannel := make(chan bool)
	go func() {
		for b.GameState == BettingState && len(b.GetPlayersWihoutBets()) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		collectionChannel <- true
	}()

	<-collectionChannel

	b.DealInitialCards()
	b.GameState = PlayingState
	b.sendGameUpdate()

	if isDealer, playerNum := b.nextPlayersTurn(); !isDealer {
		b.sendPlayerTurn(playerNum)
	}
}

func (b *BlackjackGame) DealInitialCards() {
	if b.GameState == PlayingState {
		return
	}

	b.Dealer = CreateHand(0)

	for i := 0; i < 2; i++ {
		for _, player := range b.playerMap {
			if player.Hand != nil {
				b.dealCard(player.Hand)
			}
		}

		b.dealCard(b.Dealer)
	}
}

// Returns true if player got a new card
// Returns false if game is not in state for player to receive card
// Returns error if any other reason like player could not be found
func (game *BlackjackGame) PlayerHit(playerNum PlayerId) (bool, error) {
	if game.GameState != PlayingState {
		return false, WrongGameStateError
	}

	if isDealer, num := game.nextPlayersTurn(); isDealer || num != playerNum {
		return false, WrongGameStateError
	}

	player, ok := game.GetPlayer(playerNum)
	if !ok {
		return false, PlayerNotFoundError
	}

	ok = game.dealCard(player.Hand)

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
	game.sendGameUpdate()
	if isDealer, nextNum := game.nextPlayersTurn(); !isDealer {
		game.sendPlayerTurn(nextNum)
	}

	return nil
}

func (game *BlackjackGame) DealerTurn() {
	if game.GameState != DealerState {
		return
	}

	for shouldDealerDrawCard(game.Dealer) {
		game.dealCard(game.Dealer)
	}

	game.Dealer.lock()
	time.Sleep(game.Settings.DealCardTime)
	game.sendGameUpdate()

	game.finishRound()
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

	game.sendGameUpdate()
	time.Sleep(game.Settings.DealCardTime)

	return ok
}

func (game *BlackjackGame) finishRound() {
	game.GameState = PayoutState
	// TODO: Maybe something with the payoutMap from this function
	game.sendGameUpdate()
	time.Sleep(game.Settings.TimeBetweenRounds)

	game.reset()
}
