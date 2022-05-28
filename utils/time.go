package utils

import (
	"math/rand"
	"time"
)

func WaitForRandomPeriod() {
	time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))
}
