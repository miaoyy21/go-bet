package xmd

import (
	"log"
	"strconv"
)

var latest = make(map[int]struct{})
var times = 1

var wins int
var fails int

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

	var rate float64
	if len(latest) == 0 {
		wins, fails, rate = 0, 0, 0.1
		log.Printf("【%-4d】第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0

			if wins > 16 {
				rate = 2.0
			} else if wins > 8 {
				rate = 1.0
			} else {
				rate = 0.1
			}

			log.Printf("【%-4d】第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, cache.issue, wins, cache.result, surplus, rate)
		} else {
			fails++
			wins = 0
			rate = 0.1
			log.Printf("【%-4d】第【👀 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, cache.issue, fails, cache.result, surplus, rate)
		}
	}

	var total int

	latest = make(map[int]struct{})
	for i := 0; i <= 27; i++ {
		if i <= 5 || i >= 22 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]))
			continue
		}

		betGold := int(rate * float64(cache.user.gold) * float64(stds[i]) / 1000)
		if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
			return err
		}

		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), betGold)
		latest[i] = struct{}{}
		total = total + betGold
	}

	// 额外投注
	isExtra := true
	for i := len(cache.histories) - 1; i >= len(cache.histories)-12; i-- {
		result := cache.histories[i].result
		if result <= 5 || result >= 22 {
			isExtra = false
			break
		}
	}

	if isExtra {
		r1, r2 := cache.result, cache.histories[len(cache.histories)-2].result
		if (r1 >= 10 && r1 <= 17) && (r2 < 10 || r2 > 17) {
			for i := 0; i <= 27; i++ {
				if i == 4 || i == 23 || i == 5 || i == 22 {
					log.Printf("第【%s】期【额外投注】：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]))
					continue
				}

				delta := 1.0
				if i == 3 || i == 24 || i == 6 || i == 21 {
					delta = 0.5
				}

				betGold := int(5 * delta * float64(cache.user.gold) * float64(stds[i]) / 1000)
				if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
					return err
				}

				log.Printf("第【%s】期【额外投注】：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), betGold)
				total = total + betGold
			}
		}
	}

	times++
	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, total, surplus)

	return nil
}
