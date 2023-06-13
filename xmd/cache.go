package xmd

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
)

type IssueResult struct {
	issue  int
	result int
	money  int
	member int
}

type HGold struct {
	Time string
	Gold int
}

type Cache struct {
	dir string
	md5 []byte

	db   *sql.DB
	user UserBase
	secs float64

	issue  int // 最新期数
	result int // 最新开奖结果
	money  int // 最新投注金额
	member int // 最新参与人数

	histories []IssueResult // 每期存在数据库的开奖记录
	hGolds    []HGold
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
		conf.IsDebug, conf.Gold, conf.Website, conf.Origin, conf.URL, conf.Cookie, conf.Agent,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)
	cache := &Cache{
		dir: dir,
		md5: h.Sum(nil),

		db:   db,
		user: user,
		secs: conf.Secs,

		issue:  -1,
		result: -1,
		money:  -1,
		member: -1,

		histories: make([]IssueResult, 0),
		hGolds:    make([]HGold, 0),
	}

	return cache, nil
}
