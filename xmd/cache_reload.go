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
	o.wx = conf.Wx
	o.rx = conf.Rx
	o.ex = conf.Ex

	log.Printf("当前设置活动状态（%t） ...\n", o.isExtra)
	log.Printf("当前设置投注基数为 %d ...\n", o.user.gold)
	log.Printf("当前设置当不存在超过实际赔率%.2f%%的数字时，仅进行全部投注 ...\n", o.wx*100-100)
	log.Printf("当前设置当返奖率不超过%.2f%%时，仅进行全部投注 ...\n", o.rx*100)
	log.Printf("当前设置本时内中奖比率变化量达到设定标准时，基础投注变为%.2f...\n", o.ex)
	log.Printf("当前前期投注失败时，本期投注倍率变为%.2f倍...\n", o.dx)
	return true, nil
}
