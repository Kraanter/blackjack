package manager_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/manager"
)

func TestManagerGetJoinableGameReturnsNewGameIfEmpty(t *testing.T) {
	manager := manager.CreateManager(nil)

	game := manager.GetJoinableGame()

	if game == nil {
		t.Fatalf("GetJoinableGame() = %v, want pointer to newly created game", game)
	}

	gameCount := manager.GetGameCount()
	want := uint(1)
	if gameCount != want {
		t.Fatalf("GetGameCount() = %v, want gamecount to be %v", gameCount, want)
	}
}
