package regia

import (
	"math/rand"
	"strings"
	"time"
)

const randomStringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var maxRandomStringCharsLength = len(randomStringChars)

// Return a securely generated random string
func getRandomString(length int) string {
	if length > maxRandomStringCharsLength {
		length = maxRandomStringCharsLength
	}
	var builder strings.Builder
	builder.Grow(length)
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		index := rand.Intn(maxRandomStringCharsLength)
		builder.WriteString(string(randomStringChars[index]))
	}
	return builder.String()
}
