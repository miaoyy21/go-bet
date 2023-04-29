package xmd

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRiddleDetail(t *testing.T) {
	calc()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() fail : %s \n", err.Error())
	}

	dir0, _ := filepath.Split(dir)
	cache, err := NewCache(dir0)
	if err != nil {
		log.Fatalf("NewCache() fail : %s \n", err.Error())
	}

	user := cache.user
	for i := 0; i <= 20; i++ {
		issue := strconv.Itoa(1698694 - i)

		_, num, rate, err := RiddleDetail(user, issue)
		if err != nil {
			t.Fatalf("期数【%s】，出现错误【%s】\n", issue, err.Error())
		}

		log.Printf("期数【%s】：中奖数字【%02d】，实际赔率【%-7.2f】\n", issue, num, rate)
	}
}
