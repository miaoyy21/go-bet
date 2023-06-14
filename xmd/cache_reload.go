package xmd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

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

	var conf Config
	if err := json.Unmarshal(bs, &conf); err != nil {
		return false, err
	}

	if bytes.Equal(h.Sum(nil), o.md5) {
		// 动态调整投注基数
		if hrs := time.Now().Hour(); (hrs >= 8 && hrs <= 11) || (hrs >= 14 && hrs <= 17) {
			o.user.gold = 0
		} else {
			o.user.gold = conf.Gold
		}

		return false, nil
	}

	user := NewUserBase(
		conf.IsDebug, conf.Gold, conf.Website, conf.Origin, conf.URL, conf.Cookie, conf.Agent,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)

	o.md5 = h.Sum(nil)
	o.user = user

	log.Printf("当前设置投注基数为 %d ...\n", o.user.gold)
	return true, nil
}
