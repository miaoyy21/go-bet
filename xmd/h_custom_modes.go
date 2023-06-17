package xmd

import (
	"fmt"
	"strconv"
	"strings"
)

type XCustomModesRequest struct {
	Unix     string `json:"unix"`
	Keycode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	UserId   string `json:"userid"`
	Token    string `json:"token"`
}

type XCustomModesResponse struct {
	Status int `json:"status"`
	Data   struct {
		Items []struct {
			Modeid     int    `json:"modeid"`
			Name       string `json:"name"`
			Eggs       string `json:"eggs"`
			Winmodeid  int    `json:"winmodeid"`
			Losemodeid int    `json:"losemodeid"`
		} `json:"items"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func hCustomModes(user UserBase) (int, error) {
	betRequest := XCustomModesRequest{
		Unix:     user.unix,
		Keycode:  user.code,
		PType:    "3",
		DeviceId: user.device,
		UserId:   user.id,
		Token:    user.token,
	}

	var betResponse XCustomModesResponse
	err := hDo(user, "POST", fmt.Sprintf("%s_Custom_Modes.ashx", user.url), betRequest, &betResponse)
	if err != nil {
		return 0, fmt.Errorf("查询自动投注模式，出现错误：%s", err.Error())
	}

	if betResponse.Status != 0 {
		return 0, fmt.Errorf("查询自动投注模式，服务器返回错误信息：%s", betResponse.Msg)
	}

	m1 := betResponse.Data.Items[0]
	eggs := strings.Split(m1.Eggs, ",")

	var total int
	for _, egg := range eggs {
		i, err := strconv.Atoi(egg)
		if err != nil {
			return 0, err
		}

		total = total + i
	}

	return total, nil
}
