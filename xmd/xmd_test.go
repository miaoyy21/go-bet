package xmd

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"testing"
	"time"
)

func TestCache_Sync2(t *testing.T) {

	log.Println(time.Now().Format("2006-01-02 15:04:05.999"))

	log.Println(math.Log2(float64(245644532)))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(245644532))))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(145644532))))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(75644532))))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(55644532))))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(35644532))))
	fmt.Println(math.Pow(0.5, 27.0-math.Log2(float64(15644532))))
}
