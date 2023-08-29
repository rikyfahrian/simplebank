package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {

	var sb strings.Builder
	length := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(length)]
		sb.WriteByte(c)
	}

	return sb.String()

}

func RandomName() string {
	return RandomString(8)
}

func RandomBalance() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	cur := []string{"RP", "USD", "PSTR"}
	return cur[rand.Intn(len(cur))]
}
