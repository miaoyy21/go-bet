package xmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type XSignInRequest struct {
	Action string `json:"action"`
	T      string `json:"t"`
}

type XSignInResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Data    struct {
		Item []struct {
			Img      string `json:"img"`
			Title    string `json:"title"`
			Memo     string `json:"memo"`
			Gtype    int    `json:"gtype"`
			Winnum   int    `json:"winnum"`
			Btntitle string `json:"btntitle"`
			Btnurl   string `json:"btnurl"`
			Btntype  int    `json:"btntype"`
		} `json:"item"`
	} `json:"data"`
}

func signIn(cache *Cache) {
	signInRequest := XSignInRequest{
		Action: "gopay",
		T:      fmt.Sprintf("%s GMT+0800 (中国标准时间)", time.Now().Format("Mon Jan 02 2006 15:04:05")),
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(signInRequest); err != nil {
		log.Printf("%s 【签到抽奖】失败1 >>> %s \n", cache.user.id, err.Error())
		return
	}

	url := fmt.Sprintf("%s/activity2022/ttqdcj/ashx/ttqdcj.ashx", cache.user.origin)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("%s 【签到抽奖】失败2 >>> %s \n", cache.user.id, err.Error())
		return
	}

	// Sync Header
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", fmt.Sprintf("%s=%d", cache.user.cookie, time.Now().Unix()))
	req.Header.Set("Origin", cache.user.origin)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", cache.user.agent)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%s 【签到抽奖】失败3 >>> %s \n", cache.user.id, err.Error())
		return
	}
	defer resp.Body.Close()

	var signInResponse XSignInResponse
	if err := json.NewDecoder(resp.Body).Decode(&signInResponse); err != nil {
		io.Copy(os.Stdout, resp.Body)
		log.Printf("%s 【签到抽奖】失败4 >>> %s \n", cache.user.id, err.Error())
		return
	}

	if strings.EqualFold(signInResponse.Result, "0") {
		log.Printf("%s 【签到抽奖】成功 ... \n", cache.user.id)
	} else {
		log.Printf("%s 【签到抽奖】失败5 >>> %s \n", cache.user.id, signInResponse.Message)
	}
}
