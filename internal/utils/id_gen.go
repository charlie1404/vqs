package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"sync"
)

// const ID_LENGTH = 16
// const SERIAL_LENGTH = 5
// const RANDOM_LENGTH = ID_LENGTH - SERIAL_LENGTH

// var defaultAlphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ07") // 07 are chosen random

// func serial(desc []byte) {
// 	now := time.Now().Unix() - 1514764800 // 1546300800 is 2018-01-01 00:00:00

// 	for i := SERIAL_LENGTH - 1; i >= 0; i-- {
// 		desc[i] = defaultAlphabet[now%62]
// 		now /= 62
// 	}
// }

// func random(desc []byte) {
// 	bytes := make([]byte, RANDOM_LENGTH)

// 	_, err := rand.Read(bytes)
// 	if err != nil {
// 		panic("WTF is happening")
// 	}
// 	for i := 0; i < RANDOM_LENGTH; i++ {
// 		desc[i] = defaultAlphabet[bytes[i]&63]
// 	}
// }

// func GenUniqueId() string {
// 	id := make([]byte, ID_LENGTH)

// 	serial(id)
// 	random(id[SERIAL_LENGTH:])

// 	return string(id)
// }

var bytesPool = sync.Pool{
	New: func() interface{} { return make([]byte, 10) },
}

func GenerateUUIDLikeId() string {
	b := bytesPool.Get().([]byte)
	defer bytesPool.Put(b)

	rand.Read(b)
	randomHash := md5.Sum(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", randomHash[0:4], randomHash[4:6], randomHash[6:8], randomHash[8:10], randomHash[10:])
}
