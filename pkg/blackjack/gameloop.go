package blackjack

import ()

type GameState = int

const (
	PlayerTurn GameState = iota
	DealerTurn
	BettingTurn
)

func (b *BlackjackGame) GetTurn() (GameState, interface{}) {
	if peopleToBet := b.getPeopleWihoutBets(); len(peopleToBet) > 0 {
		return BettingTurn, peopleToBet
	}
	if isDealer, nextPlayer := b.nextPlayersTurn(); isDealer {
		return DealerTurn, nil
	} else {
		return PlayerTurn, nextPlayer
	}
}

// Get list of all people that still need to bet
func (b *BlackjackGame) getPeopleWihoutBets() []playerId {
	peopleArr := make([]playerId, 0)
	for k, v := range b.playerMap {
		if v.Hand == nil && v.playing == false {
			peopleArr = append(peopleArr, k)
		}
	}

	return peopleArr
}

// Returns true if dealers turn is next
// Returns false, playerId of the player that is next
func (b *BlackjackGame) nextPlayersTurn() (isDealersTurn bool, turnPlayerId playerId) {
	for k, v := range b.playerMap {
		if v.Hand == nil {
			continue
		}
		if v.Hand.locked {
			continue
		}
		return false, k
	}

	return true, playerId(0)
}
