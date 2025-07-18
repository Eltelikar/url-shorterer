package random

import (
	"math/rand"
	"time"
)

func NewRandomString(aliasLenght int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPRSTUVWXYZ" +
		"abcdefghijklmnoprstuvwxyz" +
		"0123456789")

	b := make([]rune, aliasLenght)

	for i := range b {
		b[i] = rune(rnd.Intn(len(chars)))
	}

	return string(b)
}
