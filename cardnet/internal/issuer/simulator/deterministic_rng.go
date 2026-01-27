package simulator

import (
	"hash/fnv"
	"math/rand"
)

func rngFromSeed(seed string) *rand.Rand {
	h := fnv.New64a()
	h.Write([]byte(seed))
	return rand.New(rand.NewSource(int64(h.Sum64())))
}
