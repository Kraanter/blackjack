package blackjack

import "fmt"

type Player struct {
	Balance uint
	Hand    *Hand
	playing bool
}

func CreatePlayer(balance uint) *Player {
	return &Player{
		Balance: balance,
		playing: false,
	}
}

var NotEnoughBalanceError = fmt.Errorf("Player wants to bet more than they have balance")
var WrongGameStateError = fmt.Errorf("Game is in the wrong state")
var NotHighEnoughBetError = fmt.Errorf("Game is in the wrong state")

func (p *Player) PlaceBet(bet uint) error {
	if p.playing {
		return WrongGameStateError
	}
	if bet > p.Balance {
		return NotEnoughBalanceError
	} else if bet < 1 {
		return NotHighEnoughBetError
	}

	p.Balance -= bet
	p.Hand = CreateHand(bet)

	return nil
}
