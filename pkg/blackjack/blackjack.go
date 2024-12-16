package blackjack

import ()

type BlackjackGame struct {
	dealer  *Hand
	players []*Player
}

func createGame() *BlackjackGame {
	return &BlackjackGame{
		dealer:  CreateHand(0),
		players: make([]*Player, 0),
	}
}

func (game *BlackjackGame) AddPlayerWithBalance(balance uint) error {
	game.players = append(game.players, CreatePlayer(balance))

	return nil
}
