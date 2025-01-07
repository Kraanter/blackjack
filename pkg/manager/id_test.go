package manager_test

import (
	"testing"

	"github.com/kraanter/blackjack/pkg/manager"
)

var loopCount = 100000000

func TestRandomIdGeneratorLength1(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(1))
		if id > 9 || id < 0 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}

func TestRandomIdGeneratorLength2(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(2))
		if id > 99 || id < 10 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}

func TestRandomIdGeneratorLength3(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(3))
		if id > 999 || id < 100 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}

func TestRandomIdGeneratorLength4(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(4))
		if id > 9999 || id < 1000 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}

func TestRandomIdGeneratorLength5(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(5))
		if id > 99999 || id < 10000 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}

func TestRandomIdGeneratorLength6(t *testing.T) {
	for _ = range loopCount {
		id := manager.CreateRandomGameId(uint(6))
		if id > 999999 || id < 100000 {
			t.Fatalf("Generating random number with one digit should be one positive digit, was: %v", id)
			return
		}
	}
}
