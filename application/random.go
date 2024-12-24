package application

import (
	"github.com/sgade/randomorg"
	"math/rand/v2"
	"os"
)

type Random struct {
	randomOrg *randomorg.Random
	chacha8   *rand.ChaCha8
}

func InitializeRandom(log *Logger) *Random {
	rndorg, chacha8 := InitializeRandomOrg(log), InitializeChaCha8(log)
	return &Random{randomOrg: rndorg, chacha8: chacha8}
}

func (rnd *Random) IntN(log *Logger, maxInclusive int) int {
	if value, err := rnd.randomOrg.GenerateIntegers(1, 0, int64(maxInclusive+1)); err != nil {
		log.E("Failed to generate random number", RndSource, "external/random.org", InnerError, err)
		log.I("Successfully generated random number", RndSource, "native/chacha8", InnerError, err)
		return int(rnd.chacha8.Uint64() % uint64(maxInclusive+1))
	} else {
		log.I("Successfully generated random number", RndSource, "external/random.org")
		return int(value[0])
	}
}

func InitializeRandomOrg(log *Logger) *randomorg.Random {
	token, exist := os.LookupEnv("RANDOMORG_TOKEN")
	if !exist || token == "" {
		log.F("Random.org token must be set and non-empty")
	}

	rnd := randomorg.NewRandom(token)
	log.I("Successfully initialized random.org random")
	return rnd
}

func InitializeChaCha8(log *Logger) *rand.ChaCha8 {
	log.I("Successfully initialized chacha8 random")
	return rand.NewChaCha8([32]byte{
		0x3d, 0x1a, 0x94, 0xde, 0xdd, 0x01, 0xc3, 0xa6, 0xb8, 0x09, 0x42, 0x18, 0xba, 0x90, 0xc5, 0x71,
		0x8b, 0x93, 0xf8, 0x0c, 0x66, 0xeb, 0x98, 0xba, 0x48, 0x0d, 0x05, 0x24, 0x3b, 0xfa, 0x0e, 0x9f,
	})
}
