package helper

import (
	"math/rand"
	"time"
)

const charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length int = 8

func GenerateUID() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
