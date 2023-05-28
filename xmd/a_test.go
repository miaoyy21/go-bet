package xmd

import (
	"log"
	"math"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	log.Println(2 << 24)
	log.Println(2 << 25)
	log.Println(2 << 26)
	log.Println(2 << 27)
	log.Println(math.Ceil(float64(26499671) / (2 << 25)))
	log.Println(math.Ceil(float64(26499671*3) / (2 << 25)))
	log.Println(26499671 / 10000 * 10000)
}
