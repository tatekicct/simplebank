package util

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	if min >= max {
		panic("invalid range: min should be less than max")
	}
	return min + rnd.Int63n(max-min+1)
}

// RandomString generates a random string of length n using lowercase letters.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "AUD", "JPY"}
	return currencies[rnd.Intn(len(currencies))]
}
