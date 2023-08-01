package url

import (
	"crypto/sha256"
	"math/big"
)

func base62Encode(number *big.Int) string {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	base := big.NewInt(int64(len(base62Chars)))
	zero := big.NewInt(0)
	encoded := ""

	for number.Cmp(zero) > 0 {
		mod := new(big.Int)
		number.DivMod(number, base, mod)

		encoded = string(base62Chars[mod.Int64()]) + encoded
	}

	return encoded
}

func encode(url string) string {
	hash := sha256.Sum256([]byte(url))
	number := new(big.Int).SetBytes(hash[:])

	return base62Encode(number)
}
