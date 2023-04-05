package xmd

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

var latest = make(map[int]struct{})
var rate = 1.0
var sigma float64
var wins int
var fails int
var times = 1

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

	// 按照尾数取最热的8期
	w8s := make(map[int]struct{})
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if len(w8s) == 8 {
			break
		}

		result := cache.histories[i].result
		mod := result % 10

		// 最近2期的尾数
		w8s[mod] = struct{}{}
	}

	// String
	sw8s := make([]string, 0, len(w8s))
	for i := range w8s {
		sw8s = append(sw8s, fmt.Sprintf("%d", i))
	}
	sort.Strings(sw8s)

	// 输出
	if len(latest) == 0 {
		log.Printf("【%-4d】第【%d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，开始执行分析 ...\n", times, cache.issue, cache.result, surplus, strings.Join(sw8s, ","))
	} else {
		// L:	2 			-1
		// W:	2 			+0.4
		// W:	2*0.75		+0.3
		// W:	2*0.75^2	+0.225
		// W:	2*0.75^3	+0.16875

		// L:	2 			-1
		// W:	2 			+0.4
		// W:	2*0.8		+0.32
		// W:	2*0.8^2		+0.256
		// W:	2*0.8^3		+0.2048
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0
			rate = rate * 0.825

			// 连续赢4次，那么再降1次
			if wins >= 4 {
				rate = rate * 0.9
			}

			// 连续开奖
			if rate < 1.0 {
				rate = 1.0
			}

			// 倍率变化量
			if rate <= 3.0 && sigma > 0 {
				if rate <= 1.5 {
					log.Printf("当前倍率【%.3f】，重置倍率变化量【%.2f】为零\n", rate, sigma)
					sigma = 0
				} else {
					log.Printf("当前倍率【%.3f】，将倍率变化量由【%.2f】调整为【%.2f】\n", rate, sigma, sigma*0.825)
					sigma = sigma * 0.825
				}
			} else if sigma >= 1 {
				sigma = sigma - 1
			}

			// 防止引高倍率
			if wins > 5 {
				rate = rate + 0.125*float64(wins-5)
			}

			log.Printf("【%-4d %.2f】第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，投注倍率【%.3f】，开始执行分析 ...\n", times, sigma, cache.issue, wins, cache.result, surplus, strings.Join(sw8s, ","), rate)
		} else {
			wins = 0
			fails++
			sigma = sigma + 4.0

			// 0.90: 2.0 * 1.90 * 1.81 * 1.73 * 1.66 * 1.59 * 1.53 = 48
			// 0.88: 2.0 * 1.88 * 1.77 * 1.68 * 1.60 * 1.53 * 1.46 = 40
			rate = rate * (1.0 + math.Pow(0.9, (float64(sigma)/4.0)-1))

			log.Printf("【%-4d %.2f】第【👀 %d %02d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，投注倍率【%.3f】，开始执行分析 ...\n", times, sigma, cache.issue, fails, cache.result, surplus, strings.Join(sw8s, ","), rate)
		}
	}

	latest = make(map[int]struct{})
	bets, total := make([]string, 0), 0
	for i := 0; i <= 27; i++ {
		_, exists := w8s[i%10]
		if !exists {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，变化倍率【%.2f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]), 0.0)
			continue
		}

		betGold := int(rate * float64(cache.user.gold) * float64(stds[i]) / 1000)
		if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), betGold)

		latest[i] = struct{}{}
		bets = append(bets, fmt.Sprintf("%02d", i))

		total = total + betGold
	}

	times++
	surplus = surplus - total
	log.Printf("第【%s】期：投注数字【%s】，投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, strings.Join(bets, ","), total, surplus)

	return nil
}
