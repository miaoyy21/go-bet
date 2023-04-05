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
	req.Header.Set("Cookie", "Hm_lvt_f8f6a0064a3e891522bdf044119d462a=1678603948,1680138770,1680199323; CLIENTKEY_ShowLogin=4384-1807-1757; .ADWASPX7A5C561934E_PCEGGS=8019B2AB8884AA2912C1207FB8D59DD760652DE11925554557ABC32473E58CE3BB4D2BAA2279F66145D94E5E58569499439B5021928CBE7ECE87216B066384129ACE2081E8021B103F8CB3891DFE854A12F83C991385D903A696A00A32517680229A89270A11FA36AC0D314DFE5C0F663A4490A83793816C3497FF50B4A6128F943CB05F; Hm_lpvt_f8f6a0064a3e891522bdf044119d462a="+fmt.Sprintf("%d", time.Now().Unix()))
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
