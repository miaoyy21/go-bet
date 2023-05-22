package xmd

import (
	"errors"
	"fmt"
)

type QHistoryItem struct {
	Issue  string `json:"issue"`
	Result string `json:"lresult"`
	Money  string `json:"tmoney"`
	Member int    `json:"tmember"`
}

type QHistoryData struct {
	Items []QHistoryItem `json:"items"`
}

type QHistory struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`

	Data QHistoryData `json:"data"`
}

type QHistoryRequest struct {
	PageSize int    `json:"pagesize"`
	Unix     string `json:"unix"`
	KeyCode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	UserId   string `json:"userid"`
	Token    string `json:"token"`
}

func hGetHistories(pageSize int, user UserBase) ([]QHistoryItem, error) {

	// 查询近期历史
	hisRequest := QHistoryRequest{
		PageSize: pageSize,
		PType:    "3",
		Unix:     user.unix,
		KeyCode:  user.code,
		DeviceId: user.device,
		UserId:   user.id,
		Token:    user.token,
	}

	var hisResponse QHistory

	// 执行查询开奖历史
	err := hDo(user, "GET", fmt.Sprintf("%s_Analyse_History.ashx", user.url), hisRequest, &hisResponse)
	if err != nil {
		return nil, fmt.Errorf("查询开奖历史存在服务器错误：%s", err.Error())
	}

	// 开奖历史是否存在错误
	if hisResponse.Status != 0 {
		return nil, fmt.Errorf("查询开奖历史存在返回错误：(%d) %s", hisResponse.Status, hisResponse.Msg)
	}

	// 开奖历史为空
	if len(hisResponse.Data.Items) < 1 {
		return nil, errors.New("没有查询到开奖历史")
	}

	return hisResponse.Data.Items, nil
}
