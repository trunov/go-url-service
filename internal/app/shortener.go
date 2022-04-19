package app

import "math/rand"

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func GenerateShortLink() string {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(randInt(97, 122))
	}

	return string(bytes)
}
