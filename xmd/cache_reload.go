package xmd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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
		//if hrs := time.Now().Hour(); hrs == 13 || hrs == 18 || hrs == 20 {
		//	o.user.gold = conf.Gold / 5
		//} else {
		//	o.user.gold = conf.Gold
		//}

		return false, nil
	}

	user := NewUserBase(
		conf.IsDebug, conf.Gold, conf.Origin, conf.URL, conf.Cookie, conf.Agent,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)

	o.md5 = h.Sum(nil)
	o.user = user
	o.dev = conf.Dev

	log.Printf("当前设置仅赔率标准方差大于%.3f时，才进行投注 ...\n", o.dev)
	log.Printf("当前设置投注基数为 %d ...\n", o.user.gold)
	return true, nil
}
