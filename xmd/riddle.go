package xmd

import (
	"fmt"
	"math"
	"strconv"
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
		Riddle []struct {
			Num    string `json:"num"`
			Rate   string `json:"rate"`
			Tmoney string `json:"tmoney"`
			Gmoney string `json:"gmoney"`
		} `json:"myriddle"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func RiddleDetail(user UserBase, issue string) (map[int]float64, float64, float64, error) {
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

	err := hDo(user, "GET", fmt.Sprintf("%s_MyRiddleDetail.ashx", user.url), riddleRequest, &riddleResponse)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("查询开奖明细存在服务器错误：%s", err.Error())
	}

	if riddleResponse.Status != 0 {
		return nil, 0, 0, fmt.Errorf("查询开奖明细存在返回错误：(%d) %s", riddleResponse.Status, riddleResponse.Msg)
	}

	var exp float64
	rts := make(map[int]float64)
	for _, riddle := range riddleResponse.Data.Riddle {
		n, err := strconv.Atoi(riddle.Num)
		if err != nil {
			return nil, 0, 0, err
		}

		r0, err := strconv.ParseFloat(riddle.Rate, 64)
		if err != nil {
			return nil, 0, 0, err
		}

		rts[n] = r0
		exp = exp + (float64(stds[n])/1000)*(r0/(1000.0/float64(stds[n])))
	}

	var dev float64
	for n, r0 := range rts {
		r1 := (float64(stds[n]) / 1000) * (r0 / (1000.0 / float64(stds[n])))
		dev = dev + (r1-exp)*(r1-exp)*(float64(stds[n])/1000)
	}

	return rts, exp, math.Sqrt(dev), nil
}
