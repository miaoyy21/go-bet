package xmd

import (
	"fmt"
	"strconv"
	"strings"
)

type QRiddleDetailRequest struct {
	Issue    string `json:"issue"`
	Unix     string `json:"unix"`
	Keycode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	Userid   string `json:"userid"`
	Token    string `json:"token"`
}

type QRiddleDetail struct {
	Status int `json:"status"`
	Data   struct {
		MyRiddle []struct {
			Num    string `json:"num"`
			Rate   string `json:"rate"`
			Tmoney string `json:"tmoney"`
			Gmoney string `json:"gmoney"`
		} `json:"myriddle"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func RiddleDetail(user UserBase, issue string) (map[int]float64, int, float64, error) {
	riddleRequest := &QRiddleDetailRequest{
		Issue:    issue,
		Unix:     user.unix,
		Keycode:  user.code,
		PType:    "3",
		DeviceId: user.device,
		Userid:   user.id,
		Token:    user.token,
	}

	var riddleResponse QRiddleDetail

	err := hDo("GET", fmt.Sprintf("%s_MyRiddleDetail.ashx", user.url), riddleRequest, &riddleResponse)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("查询开奖明细存在服务器错误：%s", err.Error())
	}

	if riddleResponse.Status != 0 {
		return nil, 0, 0, fmt.Errorf("查询开奖明细存在返回错误：(%d) %s", riddleResponse.Status, riddleResponse.Msg)
	}

	var rate float64
	var num int

	rts := make(map[int]float64)
	for _, riddle := range riddleResponse.Data.MyRiddle {
		n, err := strconv.Atoi(riddle.Num)
		if err != nil {
			return nil, 0, 0, err
		}

		r0, err := strconv.ParseFloat(riddle.Rate, 64)
		if err != nil {
			return nil, 0, 0, err
		}

		if !strings.EqualFold(riddle.Gmoney, "0") {
			num = n
		}

		rts[n] = r0
		rate = rate + r0/(1000.0/float64(stds[n]))
	}

	return rts, num, rate / 28.0, nil
}
