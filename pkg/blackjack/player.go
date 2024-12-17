package blackjack

import "fmt"

type Player struct {
	Balance   uint
	Hand      *Hand
	PlayerNum uint
	playing   bool
}

func CreatePlayer(number uint, balance uint) *Player {
	return &Player{
		Balance:   balance,
		PlayerNum: number,
		playing:   false,
	}
}

var NotEnoughBalanceError = fmt.Errorf("Player wants to bet more than they have balance")
var WrongGameStateError = fmt.Errorf("Game is in the wrong state")
var NotHighEnoughBetError = fmt.Errorf("Bet needs to be higher to be valid")

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

func (p *Player) Destroy() uint {
	p.Hand = nil
	p.PlayerNum = 0
	p.playing = false
	balance := p.Balance
	p.Balance = 0

	return balance
}
