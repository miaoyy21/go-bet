package xmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

var latest = make(map[int]struct{})
var wins int
var fails int
var rate = 2.25
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

	//for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
	//	result := cache.histories[i].result
	//	if result <= 5 || result >= 22 {
	//		log.Printf("【%-4d】第【🐔 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
	//		if len(latest) > 0 {
	//			if _, exists := latest[cache.result]; !exists {
	//				fails++
	//				rate = rate * 1.625
	//			}
	//		}
	//
	//		latest = make(map[int]struct{})
	//		total, coverage := 0, 0
	//		for _, n := range []int{0, 1, 2, 3, 4, 5, 22, 23, 24, 25, 26, 27} {
	//			betGold := int(float64(cache.user.gold) * float64(stds[n]) / 1000)
	//			if err := hPostBet(nextIssue, betGold, n, cache.user); err != nil {
	//				return err
	//			}
	//			log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, n, 1000.0/float64(stds[n]), betGold)
	//
	//			total = total + betGold
	//			coverage = coverage + stds[n]
	//		}
	//
	//		times++
	//		surplus = surplus - total
	//		log.Printf("第【%s】期：投注数字【0 ~ 5, 22 ~ 27】，投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)
	//		return nil
	//	}
	//}

	// 按照尾数取最热的8期
	w8s, w2s := make(map[int]struct{}), make(map[int]struct{})
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if len(w8s) == 8 {
			break
		}

		result := cache.histories[i].result
		mod := result % 10

		// 最近3期的尾数
		if len(w2s) < 2 {
			w2s[mod] = struct{}{}
		}

		// 最近2期的尾数
		w8s[mod] = struct{}{}
	}

	// String
	sw8s := make([]string, 0, len(w8s))
	for i := range w8s {
		if _, exists := w2s[i]; exists {
			sw8s = append(sw8s, fmt.Sprintf("%d★", i))
			continue
		}

		sw8s = append(sw8s, fmt.Sprintf("%d", i))
	}
	sort.Strings(sw8s)

	// 输出
	if len(latest) == 0 {
		for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
			result := cache.histories[i].result
			if result <= 5 || result >= 22 {
				latest = make(map[int]struct{})

				log.Printf("【%-4d】第【🐔 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，暂时停止投注【%d】 ...\n", times, cache.issue, wins, fails, cache.result, surplus, len(cache.histories)-i)
				return nil
			}
		}

		log.Printf("【%-4d】第【✊ %d %03d/%03d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus, strings.Join(sw8s, ","))
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			rate = rate * 0.825
			if rate < 1.0 {
				rate = 1.0
			}

			log.Printf("【%-4d】第【👍 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus, strings.Join(sw8s, ","))
		} else {
			fails++
			rate = rate * 1.675

			for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
				result := cache.histories[i].result
				if result <= 5 || result >= 22 {
					latest = make(map[int]struct{})

					log.Printf("【%-4d】第【🐔 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，暂时停止投注【%d】 ...\n", times, cache.issue, wins, fails, cache.result, surplus, len(cache.histories)-i)
					return nil
				}
			}

			log.Printf("【%-4d】第【👀 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，投注尾数【%s】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus, strings.Join(sw8s, ","))
		}
	}

	latest = make(map[int]struct{})
	bets, total := make([]string, 0), 0
	for _, result := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27} {
		if _, exists := w8s[result%10]; !exists {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，变化倍率【%.2f】，投注金额【    -】\n", nextIssue, result, 1000.0/float64(stds[result]), 0.0)
			continue
		}

		// 倍率变化率
		delta := 1.0
		if _, exists := w2s[result%10]; exists {
			delta = 1.1
		}

		betGold := int(rate * delta * float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，变化倍率【%.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), delta, betGold)

		latest[result] = struct{}{}
		bets = append(bets, fmt.Sprintf("%02d", result))

		total = total + betGold
	}

	times++
	surplus = surplus - total
	log.Printf("第【%s】期：投注数字【%s】，投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, strings.Join(bets, ","), total, surplus)

	return nil
}
