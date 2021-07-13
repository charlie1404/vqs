package utils

import (
	"crypto/rand"
	"fmt"
	"hash/adler32"
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

func Hash32(str string) string {
	return fmt.Sprintf("%x", adler32.Checksum([]byte(str)))
}
