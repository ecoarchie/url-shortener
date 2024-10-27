package random

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// random string not less than 4 letter
func RandomString(length int) string {
	if length < 4 {
		length = 4
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alias := make([]rune, length)
	for i := 0; i < length; i++ {
		idx := rnd.Intn(len(letters))
		alias[i] = letters[idx]
	}
	return string(alias)
}
