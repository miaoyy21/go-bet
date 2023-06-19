package xmd

import (
	"fmt"
	"log"
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

func hBetting1(issue string, betGold int, result int, user UserBase) error {
	if user.isDebug {
		log.Printf("第【%s】期：<<<执行>>> 押注数字【%02d】，押注金额【% 6d】\n", issue, result, betGold)
		return nil
	}

	betRequest := XBetRequest{
		Issue:    issue,
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
	err := hDo(user, "POST", fmt.Sprintf("%s_Betting_1.ashx", user.url), betRequest, &betResponse)
	if err != nil {
		return fmt.Errorf("下期开奖期数【%s】，执行押注[% 5d]，出现错误：%s", issue, result, err.Error())
	}

	if betResponse.Status != 0 {
		return fmt.Errorf("下期开奖期数【%s】，执行押注[%d]，服务器返回错误信息：%s", issue, result, betResponse.Msg)
	}

	return nil
}
