package xmd

import (
	"bytes"
	"crypto/md5"
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
	bs, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, err
	}

	// MD5
	h := md5.New()
	if _, err := h.Write(bs); err != nil {
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(bs, &conf); err != nil {
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
		dir: dir,
		md5: h.Sum(nil),

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

func (o *Cache) Reload() (bool, error) {
	bs, err := os.ReadFile(filepath.Join(o.dir, "config.json"))
	if err != nil {
		return false, err
	}

	// MD5
	h := md5.New()
	if _, err := h.Write(bs); err != nil {
		return false, err
	}

	if bytes.Equal(h.Sum(nil), o.md5) {
		return false, nil
	}

	var conf Config
	if err := json.Unmarshal(bs, &conf); err != nil {
		return false, err
	}

	user := NewUserBase(
		conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)

	o.md5 = h.Sum(nil)
	o.user = user
	o.isExtra = conf.IsExtra
	return true, nil
}
