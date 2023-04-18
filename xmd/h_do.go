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
	req.Header.Set("Cookie", "CLIENTKEY=8492-0163-9389; Hm_lvt_f8f6a0064a3e891522bdf044119d462a=1680541450,1680785532,1681225919,1681784522; CLIENTKEY_ShowLogin=0626-0363-8896; .ADWASPX7A5C561934E_PCEGGS=87475763F77893F08AC283A232F75D8896E1F116182C1C5D06DBD6B77BDB3627D57F5A954BCED4B50A240A8452F91D662B4A1D403963256830DE464FE94CDACF27BD3961D424E20A1DFCB82E5F35157A76A9EA8A75EDA2D0EB2EF15719DEFD9512FD92681D10A8C4140C5DA1DA41BBFF26BA95F9EFE20EED5E78B4C5013E69833923C9EA; ckurl.pceggs.com=ckurl=http://www.pceggs.com/game/gameindex/gameindex.aspx?gameid=4; Hm_lpvt_f8f6a0064a3e891522bdf044119d462a="+fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("Origin", "http://manorapp.pceggs.com")
	req.Header.Set("Pragma", "http://manorapp.pceggs.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	// Response
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
