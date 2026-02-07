package utils

import "math/rand"

func RandomNumber(min, max int) int {
	return min + rand.Intn(max-min+1)
}
