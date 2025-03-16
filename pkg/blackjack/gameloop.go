package blackjack

import "time"

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
	if b.GameState != NoState {
		return
	}

	b.GameState = BettingState
	b.sendGameUpdate()

	go func() {
		for b.GameState == BettingState && len(b.GetPlayersWihoutBets()) > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		if len(b.PlayerMap) == 0 {
			return
		}

		b.DealInitialCards()
		b.GameState = PlayingState
		b.sendGameUpdate()
	}()
}

func (b *BlackjackGame) DealInitialCards() {
	if b.GameState == PlayingState {
		return
	}

	b.Dealer = CreateHand(0)

	for i := 0; i < 2; i++ {
		for _, player := range b.PlayerMap {
			if hand := player.GetActiveHand(); hand != nil {
				b.dealCard(hand)
			}
		}

		if len(b.Dealer.Cards) == 1 {
			b.hiddenDealerCard = b.shoe.DrawCard()
		} else {
			b.dealCard(b.Dealer)
		}
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
	game.sendGameUpdate()

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

	game.sendGameUpdate()

	return nil
}

func (game *BlackjackGame) DealerTurn() {
	if game.GameState != DealerState {
		return
	}

	// Show the hidden dealer card
	game.Dealer.AddCard(game.hiddenDealerCard)
	game.hiddenDealerCard = nil
	game.sendGameUpdate()

	for shouldDealerDrawCard(game.Dealer) {
		game.dealCard(game.Dealer)
	}

	game.Dealer.lock()
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

	return ok
}

func (game *BlackjackGame) finishRound() {
	game.GameState = PayoutState
	game.sendGameUpdate()
	game.reset()
}
