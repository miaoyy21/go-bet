package xmd

import (
	"fmt"
)

type XBet struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type XBetRequest struct {
	Issue    string `json:"issue"`
	GoldEggs int    `json:"totalgoldeggs"`
	Numbers  int    `json:"numbers"`
	Unix     string `json:"unix"`
	Keycode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	Userid   string `json:"userid"`
	Token    string `json:"token"`
}

func hPostBet(nextIssue string, betGold int, result int, user UserBase) error {
	if user.isDebug {
		return nil
	}

	betRequest := XBetRequest{
		Issue:    nextIssue,
		GoldEggs: betGold,
		Numbers:  result,
		Unix:     user.unix,
		Keycode:  user.code,
		PType:    "3",
		DeviceId: user.device,
		Userid:   user.id,
		Token:    user.token,
	}

	var betResponse XBet
	err := hDo("GET", fmt.Sprintf("%s_Betting_1.ashx", user.url), betRequest, &betResponse)
	if err != nil {
		return fmt.Errorf("下期开奖期数【%s】，执行押注[% 5d]，出现错误：%s", nextIssue, result, err.Error())
	}

	if betResponse.Status != 0 {
		return fmt.Errorf("下期开奖期数【%s】，执行押注[%d]，服务器返回错误信息：%s", nextIssue, result, betResponse.Msg)
	}

	return nil
}
