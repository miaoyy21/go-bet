package xmd

import (
	"log"
	"strconv"
	"testing"
)

func TestRiddleDetail(t *testing.T) {
	calc()

	origin := "..."
	url := "..."
	cookie := "..."
	unix := "..."
	keyCode := "..."
	deviceId := "..."
	userId := "..."
	token := "..."

	userBase := NewUserBase(true, 0, origin, url, cookie, unix, keyCode, deviceId, userId, token)

	for i := 0; i <= 20; i++ {
		issue := strconv.Itoa(1697574 - i)

		_, num, rate, err := RiddleDetail(userBase, issue)
		if err != nil {
			t.Fatalf("期数【%s】，出现错误【%s】\n", issue, err.Error())
		}

		log.Printf("期数【%s】：中奖数字【%02d】，实际赔率【%-7.2f】\n", issue, num, rate)
	}
}
