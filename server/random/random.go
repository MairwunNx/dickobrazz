package random

import (
	"log/slog"
	"math/rand/v2"
	"sync"

	"github.com/sgade/randomorg"
)

type Random struct {
	randomOrg *randomorg.Random
	chacha8   *rand.ChaCha8
	mu        sync.Mutex
}

func New(log *slog.Logger, token string) *Random {
	rndorg := randomorg.NewRandom(token)
	log.Info("Random.org generator initialized")
	return &Random{
		randomOrg: rndorg,
		chacha8:   newChaCha8(),
	}
}

func (rnd *Random) IntN(log *slog.Logger, maxInclusive int) int {
	value, err := rnd.randomOrg.GenerateIntegers(1, 0, int64(maxInclusive+1))
	if err != nil {
		log.Warn("Random.org failed, using chacha8 fallback", "error", err)
		rnd.mu.Lock()
		result := int(rnd.chacha8.Uint64() % uint64(maxInclusive+1))
		rnd.mu.Unlock()
		return result
	}

	log.Info("Random.org value generated")
	return int(value[0])
}

func newChaCha8() *rand.ChaCha8 {
	return rand.NewChaCha8([32]byte{
		0x3d, 0x1a, 0x94, 0xde, 0xdd, 0x01, 0xc3, 0xa6, 0xb8, 0x09, 0x42, 0x18, 0xba, 0x90, 0xc5, 0x71,
		0x8b, 0x93, 0xf8, 0x0c, 0x66, 0xeb, 0x98, 0xba, 0x48, 0x0d, 0x05, 0x24, 0x3b, 0xfa, 0x0e, 0x9f,
	})
}
