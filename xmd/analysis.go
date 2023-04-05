package xmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

var pw4s = make(map[int]struct{})
var xWins int

func analysis(cache *Cache) error {
	if err := cache.Sync(200); err != nil {
		return err
	}

	nextIssue := strconv.Itoa(cache.issue + 1)

	// 当前账户余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	// 输出
	if len(pw4s) == 0 {
		log.Printf("第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
	} else {
		if _, exists := pw4s[cache.result%10]; exists {
			xWins++
			log.Printf("第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, xWins, cache.result, surplus)
		} else {
			xWins = 0
			log.Printf("第【👀 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
		}
	}

	// 按照尾数取最热的8期
	w4s, sw4s := make(map[int]struct{}), make([]string, 0)
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if len(w4s) == 4 {
			break
		}

		w4s[cache.histories[i].result%10] = struct{}{}
		sw4s = append(sw4s, fmt.Sprintf("%d", cache.histories[i].result%10))
	}
	sort.Strings(sw4s)
	pw4s = w4s

	var total int
	for i := 0; i <= 27; i++ {
		if _, exists := w4s[i%10]; !exists {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注倍率【%-7.3f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]), 0.0)
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[i]) / 1000)
		if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), betGold)

		total = total + betGold
	}

	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, total, surplus)

	return nil
}
