package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go-bet/xmd"
	"log"
	"os"
)

func main() {
	log.Printf("当前版本 2023.04.30 13:38\n")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s \n", err.Error())
	}

	cache, err := xmd.NewCache(dir)
	if err != nil {
		log.Fatalf("getCache() fail : %s\n", err.Error())
	}

	xmd.Run(cache)
}
