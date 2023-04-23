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
	req.Header.Set("Cookie", "CLIENTKEY=2921-2461-4085; Hm_lvt_f8f6a0064a3e891522bdf044119d462a=1681225919,1681784522,1682001475,1682252733; CLIENTKEY_ShowLogin=5527-1676-4716; .ADWASPX7A5C561934E_PCEGGS=39386FA4BA9AFE3650BD534F01F68A526B4659708EF8675664504ECD1CB8FA1905573D645F2EB788913ED3CB21946951DA5BBD075FE5023DB779299747E7BE8089823A5CAEC8F7F057ED64FE3FE61EAF502C4B781E7D8098E4E757706A04BA79AA01698B4FF8DAEF9DCDAC6FADA955562579227B99049CB15190CBD251118C330D3E5288; ckurl.pceggs.com=ckurl=http://www.pceggs.com/game/gameindex/gameindex.aspx?gameid=4; Hm_lpvt_f8f6a0064a3e891522bdf044119d462a="+fmt.Sprintf("%d", time.Now().Unix()))
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
