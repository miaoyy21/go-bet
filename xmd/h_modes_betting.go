package xmd

import "fmt"

type XModesBettingRequest struct {
	Issue  string `json:"issue"`
	ModeId int    `json:"modeid"`

	Unix     string `json:"unix"`
	Keycode  string `json:"keycode"`
	PType    string `json:"ptype"`
	DeviceId string `json:"deviceid"`
	UserId   string `json:"userid"`
	Token    string `json:"token"`
}

type XModesBettingResponse struct {
	Status int      `json:"status"`
	Data   struct{} `json:"data"`
	Msg    string   `json:"msg"`
}

func hModesBetting(issue string, modeId int, user UserBase) error {
	if user.isDebug {
		return nil
	}

	betRequest := XModesBettingRequest{
		Issue:  issue,
		ModeId: modeId,

		Unix:     user.unix,
		Keycode:  user.code,
		PType:    "3",
		DeviceId: user.device,
		UserId:   user.id,
		Token:    user.token,
	}

	var betResponse XModesBettingResponse
	err := hDo(user, "POST", fmt.Sprintf("%s_ModesBetting.ashx", user.url), betRequest, &betResponse)
	if err != nil {
		return fmt.Errorf("下期开奖期数【%s】，执行押注模式ID[%d]，出现错误：%s", issue, modeId, err.Error())
	}

	if betResponse.Status != 0 {
		return fmt.Errorf("下期开奖期数【%s】，执行押注模式ID[%d]，服务器返回错误信息：%s", issue, modeId, betResponse.Msg)
	}

	return nil
}
