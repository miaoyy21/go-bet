package xmd

import (
	"log"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	x := 0.991234723743
	log.Printf("%.8f \n", x)
	log.Printf("%.8f \n", (x-0.99)*100)
}
