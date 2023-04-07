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

	log.Printf("%.2f \n", 1.35+math.Pow(0.1, 0))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 1))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 2))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 3))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 4))
	fmt.Println(2.35 * 2.00 * 1.77)

	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 0))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 1))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 2))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 3))
	log.Printf("%.2f \n", 1.35+math.Pow(0.65, 4))
	fmt.Println()

	log.Printf("%.2f \n", 1.35+math.Pow(0.75, 0))
	log.Printf("%.2f \n", 1.35+math.Pow(0.75, 1))
	log.Printf("%.2f \n", 1.35+math.Pow(0.75, 2))
	log.Printf("%.2f \n", 1.35+math.Pow(0.75, 3))
	log.Printf("%.2f \n", 1.35+math.Pow(0.75, 4))
	fmt.Println()

	fmt.Println(math.Pow(1.25, float64(1)))
	fmt.Println(math.Pow(1.25, float64(2)))
	fmt.Println(math.Pow(1.25, float64(3)))
	fmt.Println(math.Pow(1.25, float64(4)))
	fmt.Println(math.Pow(1.25, float64(5)))
	fmt.Println(math.Pow(1.25, float64(6)))
	fmt.Println(math.Pow(1.25, float64(7)))
	fmt.Println(math.Pow(1.25, float64(8)))
	fmt.Printf("%.3f \n", 2.15*2.06*1.98*1.91*1.84*1.78)

	log.Println("标准赔率")
	for i := 0; i <= 27; i++ {
		log.Printf("%02d:  %.2f \n", i, 1000.0/float64(stds[i]))
	}
}
