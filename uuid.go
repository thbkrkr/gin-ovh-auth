package ovhauth

import (
	"math/rand"
	"time"
)

// https://github.com/germ/go-bits/blob/master/puuid/puuid.go

// GenerateUUID generates an UUID
func GenerateUUID() string {
	var uuid string
	rand.Seed(time.Now().UnixNano())

	for index := 1; index <= 25; index++ {
		if index%5 == 0 && index != 25 {
			uuid += "-"
		}
		uuid += string('A' + (rand.Int() % 26))
	}

	return uuid
}
