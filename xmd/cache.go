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
}

type HGold struct {
	Time string
	Gold int
}

type Cache struct {
	dir string
	md5 []byte

	db      *sql.DB
	user    UserBase
	isExtra bool
	wx      float64
	rx      float64
	ex      float64
	dx      float64

	issue  int // 最新期数
	result int // 最新开奖结果

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
		conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)
	cache := &Cache{
		dir: dir,
		md5: h.Sum(nil),

		db:      db,
		user:    user,
		isExtra: conf.IsExtra,
		wx:      conf.Wx,
		rx:      conf.Rx,
		ex:      conf.Ex,
		dx:      conf.Dx,

		issue:  -1,
		result: -1,

		histories: make([]IssueResult, 0),
		hGolds:    make([]HGold, 0),
	}

	return cache, nil
}

func (o *Cache) IsExtra() bool {
	return o.isExtra
}
