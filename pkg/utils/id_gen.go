package utils

import (
	"crypto/rand"
)

var defaultAlphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ07") // 07 are chosen random

func GenRandomId() string {
	size := 32

	id := make([]rune, size)
	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		panic("WTF is happening")
	}
	for i := 0; i < size; i++ {
		id[i] = defaultAlphabet[bytes[i]&63]
	}

	return string(id)
}
