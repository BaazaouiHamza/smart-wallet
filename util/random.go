package util

import (
	"math/rand"
	"strings"
	"time"
)

const randomNym = "123456789adcdefghijklmnopqrstuvwxyz"
const alphabet = "adcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomString generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[(rand.Intn(k))]
		sb.WriteByte(c)
	}
	return sb.String()
}

//RandomOwner generate a random  description/text
func RandomPostDescriptionOrText() string {
	return RandomString(20)
}

//RandomOwner generate a random  name
func RandomName() string {
	return RandomString(6)
}

//RandomString generate a random nym of length n
func RandomNym(n int) string {
	var sb strings.Builder
	k := len(randomNym)

	for i := 0; i < n; i++ {
		c := randomNym[(rand.Intn(k))]
		sb.WriteByte(c)
	}
	return sb.String()
}
