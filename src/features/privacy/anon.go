package privacy

import (
	"fmt"
	"math/big"
)

// LCG multiplier (Knuth), used for deterministic anonymous number from user ID.
const anonNumberMultiplier = 6364136223846793005

// generateAnonymousNumber produces a deterministic 4-digit anonymous number from user ID.
func generateAnonymousNumber(userID int64) string {
	id := big.NewInt(userID)
	multiplier := big.NewInt(anonNumberMultiplier)
	mask := big.NewInt(0xffffffff)
	mod := big.NewInt(10000)

	var n big.Int
	n.Mul(id, multiplier)
	n.And(&n, mask)
	n.Mod(&n, mod)

	v := n.Int64()
	if v < 0 {
		v = -v
	}
	return fmt.Sprintf("%04d", v%10000)
}
