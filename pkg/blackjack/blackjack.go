package blackjack

import (
	"fmt"
)

type BlackjackGame struct {
	dealer      *Hand
	playerMap   map[uint]*Player
	playerCount uint
}

func CreateGame() *BlackjackGame {
	return &BlackjackGame{
		dealer:    CreateHand(0),
		playerMap: make(map[uint]*Player, 0),
	}
}

func (game *BlackjackGame) AddPlayerWithBalance(balance uint) *Player {
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)
	game.playerMap[game.playerCount] = newPlayer

	return newPlayer
}

var PlayerNotFoundError error = fmt.Errorf("Could not find player")

func (game *BlackjackGame) RemovePlayer(playerNum uint) (uint, error) {
	playerToDelete, ok := game.playerMap[playerNum]
	if !ok {
		return 0, PlayerNotFoundError
	}

	delete(game.playerMap, playerNum)
	playersBalance := playerToDelete.Destroy()

	return playersBalance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum uint) (*Player, bool) {
	player, ok := game.playerMap[playerNum]

	return player, ok
}
