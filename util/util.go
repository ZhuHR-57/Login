package util

import (
	"math/rand"
	"time"
)

//随机字符
func RandomName(n int) string {
	var letters = []byte("diabvciwnvonwmxwomonrvbuebivnwoineo")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
