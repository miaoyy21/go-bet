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

	log.Println(1.0 + math.Pow(0.9, (1.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (2.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (3.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (4.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (5.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (6.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (7.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (8.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (15.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (16.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (17.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (18.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (19.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (20.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (24.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (28.0/4.0)-1))
	log.Println(1.0 + math.Pow(0.9, (32.0/4.0)-1))

	log.Println("标准赔率")
	for i := 0; i <= 27; i++ {
		log.Printf("%02d:  %.2f \n", i, 1000.0/float64(stds[i]))
	}

	log.Println(math.Round(3 / 4))
}
