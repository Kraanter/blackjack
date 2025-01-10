package blackjack

import (
	"fmt"
)

type Player struct {
	Balance uint
	// Nil if not playing in current round
	Hand      *Hand
	PlayerNum uint

	// True if player has made a choice for the current round
	playing bool
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
	if p.playing || p.Hand != nil {
		return WrongGameStateError
	}
	if bet > p.Balance {
		return NotEnoughBalanceError
	} else if bet < 1 {
		return NotHighEnoughBetError
	}

	p.Balance -= bet
	p.Hand = CreateHand(bet)
	p.playing = true

	return nil
}

func (p *Player) Destroy() uint {
	p.reset()
	p.Hand = nil
	p.PlayerNum = 0
	p.playing = false
	balance := p.Balance
	p.Balance = 0

	return balance
}

func (p *Player) stand() {
	p.Hand.lock()
}

func (p *Player) reset() {
	p.playing = false
	p.Hand = nil
}

func (p *Player) String() string {
	return fmt.Sprintf("PlayerNr: %v | Balance: €%v | Hand: %v", p.PlayerNum, p.Balance, p.Hand.String())
}
