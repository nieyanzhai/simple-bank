package util

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate a random string of length n
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Generate a random integer between min and max (inclusive)
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Generate a random currency between USD EUR CNY JPY
func RandCurrency() string {
	currencies := []string{"USD", "EUR", "CNY", "JPY"}
	return currencies[RandInt(0, len(currencies)-1)]
}
