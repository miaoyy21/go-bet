package main

import (
	"encoding/json"
	"go-bet/xmd"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	IsDebug  bool   `json:"is_debug"`
	Gold     int    `json:"gold"`
	Origin   string `json:"origin"`
	URL      string `json:"url"`
	Cookie   string `json:"cookie"`
	UserId   string `json:"user_id"`
	Token    string `json:"token"`
	Unix     string `json:"unix"`
	KeyCode  string `json:"key_code"`
	DeviceId string `json:"device_id"`
}

func main() {
	log.Printf("当前版本 2023.04.29 20:05\n")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() fail : %s\n", err.Error())
	}

	file, err := os.Open(filepath.Join(dir, "config.json"))
	if err != nil {
		log.Fatalf("os.Open() fail : %s\n", err.Error())
	}
	defer file.Close()

	var conf Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatalf("json.Decode() fail : %s\n", err.Error())
	}

	user := xmd.NewUserBase(conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie, conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token)
	cache, err := xmd.NewCache(user)
	if err != nil {
		log.Fatalf("xmd.NewCache() fail :: %s\n", err.Error())
	}

	xmd.Run(cache)
}
