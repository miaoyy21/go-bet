package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go-bet/xmd"
	"log"
	"os"
	"time"
)

func main() {
	log.Printf("当前版本 2023.05.13 23:55\n")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s \n", err.Error())
	}

	log.Println("延迟3秒启动...")
	time.Sleep(3 * time.Second)
	cache, err := xmd.NewCache(dir)
	if err != nil {
		log.Fatalf("xmd.NewCache() fail : %s\n", err.Error())
	}

	xmd.Run(cache)
}
