package xmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func hDo(method string, url string, s interface{}, t interface{}) error {
	buf := new(bytes.Buffer)

	// JSON Encode
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		return err
	}

	// Sync
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return err
	}

	// Sync Header
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", "CLIENTKEY=7082-2708-6480; Hm_lvt_f8f6a0064a3e891522bdf044119d462a=1680138770,1680199323,1680541450,1680785532; CLIENTKEY_ShowLogin=4052-3874-7509; .ADWASPX7A5C561934E_PCEGGS=6146ED75F3622D9085BCA889F181FFD7EEB638607F66D8AF8DC1DD7475B603F1071E6F178DA947DAB32F179D425F27BAF5349C4E56E00628E3859CF10B3F36783B80B0AC6BE000E69651026A4BB1AF24E65B8A3D481CA98A7C2FBF0EB1ED7DAFC89AB6D48C9D599C85D412E88E810B0EF9F51034ACAC34A8E30E69F6F923076E2DFE3CC3; ckurl.pceggs.com=ckurl=http://www.pceggs.com/game/gameindex/gameindex.aspx?gameid=4; Hm_lpvt_f8f6a0064a3e891522bdf044119d462a="+fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("Origin", "http://manorapp.pceggs.com")
	req.Header.Set("Pragma", "http://manorapp.pceggs.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	// Response
	http.DefaultClient.Timeout = 3 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// JSON Decode
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return err
	}

	return nil
}
