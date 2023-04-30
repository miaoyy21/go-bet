package xmd

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	IsDebug    bool   `json:"is_debug"`
	IsExtra    bool   `json:"is_extra"`
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

	file, err := os.Open(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conf Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return nil, err
	}

	// MySQL
	db, err := sql.Open("mysql", conf.DataSource)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	user := NewUserBase(
		conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)
	cache := &Cache{
		db:      db,
		user:    user,
		isExtra: conf.IsExtra,

		issue:  -1,
		result: -1,

		histories: make([]IssueResult, 0),
		hGolds:    make([]HGold, 0),
	}

	return cache, nil
}
