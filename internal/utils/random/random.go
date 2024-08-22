package random

import (
	"math/rand"
	"time"
)

func New(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789")

	alias := make([]rune, length)
	for i := range alias {
		alias[i] = chars[rnd.Intn(len(chars))]
	}
	return string(alias)
}
