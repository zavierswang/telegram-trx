package utils

import (
	"math"
	"math/rand"
	"time"
)

func RandPoint() float64 {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)
	return float64(r) / 1000
}

func Trunc(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
