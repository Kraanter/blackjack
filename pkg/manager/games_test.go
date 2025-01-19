package manager_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/blackjack"
	"github.com/kraanter/blackjack/pkg/manager"
)

func TestManagerGetJoinableGameReturnsNewGameIfEmpty(t *testing.T) {
	manager := manager.CreateManager(nil)

	id, game := manager.GetJoinableGame()

	if game == nil {
		t.Fatalf("GetJoinableGame() = %v, want pointer to newly created game", game)
	}

	if id == 0 {
		t.Fatalf("GetJoinableGame() = 0, want game ID not zero Id")
	}

	gameCount := manager.GetGameCount()
	want := uint(1)
	if gameCount != want {
		t.Fatalf("GetGameCount() = %v, want gamecount to be %v", gameCount, want)
	}
}

func TestManagerGetJoinableGameCreatesNewGameIfAllGamesFull(t *testing.T) {
	want := uint(5)
	settings := manager.CreateSettings()
	settings.MinPlayerCount = 1
	manager := manager.CreateManager(settings)

	for range want {
		_, game := manager.GetJoinableGame()
		if game == nil {
			t.Fatalf("GetJoinableGame() = %v, want pointer to newly created game", game)
		}
		game.AddPlayerWithBalance(0)
	}

	gameCount := manager.GetGameCount()
	if gameCount != want {
		t.Fatalf("GetGameCount() = %v, want gamecount to be %v", gameCount, want)
	}
}

func TestManagerGetsCorrectGameIfGivenGameCode(t *testing.T) {
	settings := manager.CreateSettings()
	settings.IdLength = 1
	settings.MinPlayerCount = 1
	man := manager.CreateManager(settings)
	numGames := 9

	games := make([]*blackjack.BlackjackGame, 0)
	for range numGames {
		_, game := man.GetJoinableGame()
		if game == nil {
			t.Fatalf("GetJoinableGame() = %v, want pointer to newly created game", game)
		}
		game.AddPlayerWithBalance(0)
		games = append(games, game)
	}

	for i := range numGames {
		i++
		want := manager.GameId(i)
		game, err := man.GetGameWithId(want)
		if game == nil {
			t.Fatalf("GetJoinableGame() = %v, want pointer to previously created game with gameId = %v.\n Error = %v", game, want, err.Error())
		}
	}
}

func TestManagerReturnsErrorIfGameNotFound(t *testing.T) {
	man := manager.CreateManager(nil)
	game, err := man.GetGameWithId(manager.CreateRandomGameId(2))
	if game != nil {
		t.Fatalf("manager.GetGameWithId(invalidId) = %v, want to return %v", game, nil)
	}

	want := manager.GameNotFoundError
	if err != want {
		t.Fatalf("manager.GetGameWithId(invalidId) should return game not found error but got: %v", err)
	}
}
