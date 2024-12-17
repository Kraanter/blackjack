package blackjack

import (
	"fmt"
)

type playerId = uint

type BlackjackGame struct {
	dealer      *Hand
	playerMap   map[playerId]*Player
	playerCount playerId
	shoe        Shoe
}

func CreateGame() *BlackjackGame {
	return &BlackjackGame{
		dealer:    CreateHand(0),
		playerMap: make(map[playerId]*Player, 0),
		shoe:      *CreateShoe(1),
	}
}

func (game *BlackjackGame) AddPlayerWithBalance(balance playerId) *Player {
	game.playerCount++
	newPlayer := CreatePlayer(game.playerCount, balance)
	game.playerMap[game.playerCount] = newPlayer

	return newPlayer
}

var PlayerNotFoundError error = fmt.Errorf("Could not find player")

func (game *BlackjackGame) RemovePlayer(playerNum playerId) (playerId, error) {
	playerToDelete, ok := game.playerMap[playerNum]
	if !ok {
		return 0, PlayerNotFoundError
	}

	delete(game.playerMap, playerNum)
	playersBalance := playerToDelete.Destroy()

	return playersBalance, nil
}

func (game *BlackjackGame) GetPlayer(playerNum playerId) (*Player, bool) {
	player, ok := game.playerMap[playerNum]

	return player, ok
}
