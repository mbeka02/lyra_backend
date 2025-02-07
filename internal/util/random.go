package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // return random int btwn min and max
}

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		character := alphabet[rand.Intn(k)]
		sb.WriteByte(character)

	}
	return sb.String()
}

func RandName() string {
	name := RandString(10)
	return name
}

func RandEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandString(10))
}
