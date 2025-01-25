package blackjack

import (
	"fmt"
)

type Player struct {
	Balance uint
	// Nil if not playing in current round
	Hands     []*Hand
	PlayerNum uint

	// True if player has made a choice for the current round
	playing bool
}

func CreatePlayer(number uint, balance uint) *Player {
	return &Player{
		Balance:   balance,
		PlayerNum: number,
		Hands:     make([]*Hand, 0),
		playing:   false,
	}
}

var NotEnoughBalanceError = fmt.Errorf("Player wants to bet more than they have balance")
var WrongGameStateError = fmt.Errorf("Game is in the wrong state")
var NotHighEnoughBetError = fmt.Errorf("Bet needs to be higher to be valid")

func (p *Player) PlaceBet(bet uint) error {
	if p.playing || len(p.Hands) != 0 {
		return WrongGameStateError
	}
	if bet > p.Balance {
		return NotEnoughBalanceError
	} else if bet < 1 {
		return NotHighEnoughBetError
	}

	p.Balance -= bet
	p.Hands = append(p.Hands, CreateHand(bet))
	p.playing = true

	return nil
}

func (p *Player) Destroy() uint {
	p.reset()
	p.Hands = make([]*Hand, 0)
	p.PlayerNum = 0
	p.playing = false
	balance := p.Balance
	p.Balance = 0

	return balance
}

func (p *Player) GetActiveHand() *Hand {
	for _, hand := range p.Hands {
		if !hand.locked {
			return hand
		}
	}

	return nil
}

func (p *Player) stand() {
	p.GetActiveHand().lock()
}

func (p *Player) reset() {
	p.playing = false
	p.Hands = make([]*Hand, 0)
}

func (p *Player) String() string {
	return fmt.Sprintf("PlayerNr: %v | Balance: €%v | Hand: %v", p.PlayerNum, p.Balance, p.Hands)
}
