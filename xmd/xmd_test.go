package xmd

import (
	"fmt"
	"log"
	"math"
	"testing"
)

func TestInit(t *testing.T) {
	calc()

	var i float64

	i = 1
	fmt.Printf("%.2f\n", i/4)

	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 0))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 1))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 2))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 3))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 4))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 5))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 6))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 7))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 8))
	log.Printf("%.3f \n", 1.375+math.Pow(0.75, 9))
	fmt.Printf("%.3f \n", 2.15*2.06*1.98*1.91*1.84*1.78)

	log.Println("标准赔率")
	for i := 0; i <= 27; i++ {
		log.Printf("%02d:  %.2f \n", i, 1000.0/float64(stds[i]))
	}
}
