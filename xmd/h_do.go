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
	req.Header.Set("Cookie", "CLIENTKEY=6905-8779-3648; Hm_lvt_f8f6a0064a3e891522bdf044119d462a=1682001475,1682252733,1682319534,1682510129; CLIENTKEY_ShowLogin=1986-4594-9911; .ADWASPX7A5C561934E_PCEGGS=EE9648FCE0A5722A18204F85D91BCFCCEC8C4876F15652B439E625741D478858930E7B9DE52E7E4754F8FB5659B66C4C35B2DF0C975FB80C108CA1275A3F9856627F406EA9F3B88D2C1045A6DC5BFCF11577154F8CBF1D73384A411685B23741844279BE628A213FF52C07FCCF6F029ED5ADFD2769858B47AC9B253284611B3E3B0DA39D; ckurl.pceggs.com=ckurl=http://www.pceggs.com/game/gameindex/gameindex.aspx?gameid=4; Hm_lpvt_f8f6a0064a3e891522bdf044119d462a="+fmt.Sprintf("%d", time.Now().Unix()))
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
