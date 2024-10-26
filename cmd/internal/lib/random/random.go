package random

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func RandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alias := make([]rune, length)
	for i := 0; i < length; i++ {
		idx := rnd.Intn(len(letters))
		alias[i] = letters[idx]
	}
	return string(alias)
}