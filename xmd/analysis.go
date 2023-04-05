package xmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

var pw8s = make(map[int]struct{})
var isBet bool
var sigma int
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
	if len(pw8s) == 0 {
		log.Printf("第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
	} else {
		if _, exists := pw8s[cache.result%10]; exists {
			xWins++
			if isBet {
				log.Printf("【%d】第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", sigma, cache.issue, xWins, cache.result, surplus)
			} else {
				log.Printf("【%d】第【🧠 %d %02d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", sigma, cache.issue, xWins, cache.result, surplus)
			}

			// 只有连续成功2次后，才进行投注
			if xWins >= 2 {
				isBet = true
				if sigma > 1 {
					sigma = sigma - 1
				}
			} else {
				isBet = false
			}
		} else {
			xWins = 0
			if isBet {
				sigma = sigma + 4
			} else {
				sigma = sigma + 2
			}

			isBet = false
			log.Printf("【%d】第【👀 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", sigma, cache.issue, cache.result, surplus)
		}
	}

	// 按照尾数取最热的8期
	w8s, sw8s := make(map[int]struct{}), make([]string, 0)
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if len(w8s) == 8 {
			break
		}

		w8s[cache.histories[i].result%10] = struct{}{}
		sw8s = append(sw8s, fmt.Sprintf("%d", cache.histories[i].result%10))
	}
	sort.Strings(sw8s)
	pw8s = w8s

	if !isBet {
		log.Printf("第【%s】期：本期没有进行投注 >>>>>>>>>> \n", nextIssue)
		return nil
	}

	var total int
	for i := 0; i <= 27; i++ {
		if _, exists := w8s[i%10]; !exists {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注倍率【%-7.3f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]), 0.0)
			continue
		}

		rate := 2.0
		//rate := 0.725 + 0.75*(float64(sigma)+3)/4
		betGold := int(rate * float64(cache.user.gold) * float64(stds[i]) / 1000)
		if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注倍率【%-7.3f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), rate, betGold)

		total = total + betGold
	}

	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, total, surplus)

	return nil
}
