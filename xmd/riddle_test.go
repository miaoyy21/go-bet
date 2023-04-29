package xmd

import (
	"log"
	"strconv"
	"testing"
)

func TestRiddleDetail(t *testing.T) {
	calc()

	unix := "1682729340"
	keyCode := "81cabab20f55ccf81bcf170c91e264a3"
	deviceId := "0E6EE3CC-8184-4CD7-B163-50AE8AD4516F"
	userId := "31591499"
	token := "o787oj82qisa3d7fu2ur9r1uct9bo5s8fpfacw41"

	userBase := NewUserBase(true, 0, unix, keyCode, deviceId, userId, token)

	for i := 0; i <= 20; i++ {
		issue := strconv.Itoa(1697574 - i)

		_, num, rate, err := RiddleDetail(userBase, issue)
		if err != nil {
			t.Fatalf("期数【%s】，出现错误【%s】\n", issue, err.Error())
		}

		log.Printf("期数【%s】：中奖数字【%02d】，实际赔率【%-7.2f】\n", issue, num, rate)
	}
}
