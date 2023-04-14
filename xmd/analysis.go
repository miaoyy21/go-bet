package xmd

//
//import (
//	"fmt"
//	"log"
//	"strconv"
//	"strings"
//)
//
//var latest = make(map[int]struct{})
//var wins int
//var fails int
//var rate = 1.0
//var times = 1
//
//func analysis(cache *Cache) error {
//	if err := cache.Sync(200); err != nil {
//		return err
//	}
//
//	nextIssue := strconv.Itoa(cache.issue + 1)
//
//	// 当前账户余额
//	surplus, err := hGetGold(cache.user)
//	if err != nil {
//		return err
//	}
//
//	// 输出
//	if len(latest) == 0 {
//		log.Printf("【%-4d】第【✊ %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
//	} else {
//		if _, exists := latest[cache.result]; exists {
//			wins++
//
//			//rate = rate - 1.0
//			//if rate < 1.0 {
//			//	rate = 1.0
//			//}
//
//			log.Printf("【%-4d】第【👍 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
//		} else {
//			fails++
//			//rate = rate + 0.1
//
//			log.Printf("【%-4d】第【👀 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
//		}
//	}
//
//	latest = make(map[int]struct{})
//	bets, total := make([]string, 0), 0
//	for _, result := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27} {
//		if cache.result%10 != result%10 {
//			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注金额【    -】\n", nextIssue, result, 1000.0/float64(stds[result]))
//			continue
//		}
//
//		betGold := int(rate * float64(cache.user.gold) * float64(stds[result]) / 1000)
//		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
//			return err
//		}
//		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)
//
//		latest[result] = struct{}{}
//		bets = append(bets, fmt.Sprintf("%02d", result))
//
//		total = total + betGold
//	}
//
//	times++
//	surplus = surplus - total
//	log.Printf("第【%s】期：投注数字【%s】，投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, strings.Join(bets, ","), total, surplus)
//
//	return nil
//}
