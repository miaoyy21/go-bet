package xmd

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	IsDebug    bool   `json:"is_debug"`
	DataSource string `json:"datasource"`
	Gold       int    `json:"gold"`
	Origin     string `json:"origin"`
	URL        string `json:"url"`
	Cookie     string `json:"cookie"`
	UserId     string `json:"user_id"`
	Token      string `json:"token"`
	Unix       string `json:"unix"`
	KeyCode    string `json:"key_code"`
	DeviceId   string `json:"device_id"`
}

func NewCache(dir string) (*Cache, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/bet?charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true")
	if err != nil {
		log.Fatalf("sql.Open() fail : %s \n", err.Error())
	}

	file, err := os.Open(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conf Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return nil, err
	}

	user := NewUserBase(
		conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)
	cache := &Cache{
		db:   db,
		user: user,

		issue:  -1,
		result: -1,

		histories: make([]IssueResult, 0),
		hGolds:    make([]HGold, 0),
	}

	return cache, nil
}
