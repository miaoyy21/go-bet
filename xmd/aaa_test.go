package xmd

import (
	"log"
	"math"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	r0 := 13.70
	r1 := 13.40
	log.Printf("%.8f\n", r1/r0)
	log.Printf("%d\n", int(math.Floor(r1*10000/r0)))
	log.Printf("%.8f\n", 1.0-float64(10000-int(math.Floor(r1*10000/r0)))*0.005)

}
