// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

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
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		index := rand.Intn(maxRandomStringCharsLength)
		builder.WriteString(string(randomStringChars[index]))
	}
	return builder.String()
}
