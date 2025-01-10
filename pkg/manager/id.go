package manager

import (
	"math"
	"math/rand/v2"
)

type GameId = uint64

func CreateRandomGameId(length uint) GameId {
	minNum := int(math.Pow10(int(length) - 1))
	id := rand.IntN(int(math.Pow10(int(length))) - minNum)

	return GameId(minNum + id)
}
