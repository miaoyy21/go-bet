package xmd

import (
	"log"
	"math"
	"testing"
)

func TestSpaceFn(t *testing.T) {
	r1, r0 := 13.64, 13.70
	log.Printf("%.2f/%.2f = %.2f\n", r1, r0, math.Trunc(r1*100/r0))
	log.Printf("%t\n", int(math.Trunc(r1*100/r0)) < 100.0)
}
