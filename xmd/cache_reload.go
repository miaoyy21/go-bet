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
		return false, nil
	}

	user := NewUserBase(
		conf.IsDebug, conf.IsBetMode, conf.Origin, conf.URL, conf.Cookie, conf.Agent,
		conf.Unix, conf.KeyCode, conf.DeviceId, conf.UserId, conf.Token,
	)

	o.md5 = h.Sum(nil)
	o.user = user

	log.Printf("当前是否启用设定投注模式【%t】 ... \n", o.user.isBetMode)

	return true, nil
}
