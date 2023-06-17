package xmd

import (
	"fmt"
	"strconv"
	"strings"
)

type UserBaseRequest struct {
	Unix     string `json:"unix"`
	KeyCode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	UserId   string `json:"userid"`
	Token    string `json:"token"`
}

type UserBaseResponse struct {
	Status int `json:"status"`
	Data   struct {
		GoldEggs string `json:"goldeggs"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func hGetGold(user UserBase) (gold int, err error) {
	userBaseRequest := UserBaseRequest{
		Unix:     user.unix,
		KeyCode:  user.code,
		PType:    "3",
		DeviceId: user.device,
		UserId:   user.id,
		Token:    user.token,
	}

	var userBaseResponse UserBaseResponse

	// 执行查询开奖历史
	err = hDo(user, "POST", fmt.Sprintf("%s_UserBase.ashx", user.url), userBaseRequest, &userBaseResponse)
	if err != nil {
		return
	}

	// 开奖历史是否存在错误
	if userBaseResponse.Status != 0 {
		return gold, fmt.Errorf("查询用户信息存在错误返回：(%d) %s", userBaseResponse.Status, userBaseResponse.Msg)
	}

	sGold := strings.ReplaceAll(userBaseResponse.Data.GoldEggs, ",", "")
	iGold, err := strconv.Atoi(sGold)
	if err != nil {
		return gold, err
	}

	return iGold, nil
}
