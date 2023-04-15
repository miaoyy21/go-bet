package xmd

import (
	"log"
	"sort"
	"strconv"
)

var latest = make(map[int]struct{})
var wins int
var fails int
var rate = 1.0
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

	// 输出
	if len(latest) == 0 {
		log.Printf("【%-4d】第【✊ %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			log.Printf("【%-4d】第【👍 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
		} else {
			fails++
			log.Printf("【%-4d】第【👀 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, cache.issue, wins, fails, cache.result, surplus)
		}
	}

	target := getTarget(cache)

	latest = make(map[int]struct{})
	total, coverage := 0, 0
	for _, result := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27} {
		if _, exists := target[result]; !exists {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注金额【    -】\n", nextIssue, result, 1000.0/float64(stds[result]))
			continue
		}

		betGold := int(rate * float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	times++
	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	return nil
}

func getTarget(cache *Cache) map[int]struct{} {
	type Space struct {
		Result int
		Space  int
	}

	spaces := make(map[int]int)
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if _, ok := spaces[cache.histories[i].result]; ok {
			continue
		}

		spaces[cache.histories[i].result] = len(cache.histories) - i
	}

	newSpaces := make([]Space, 0, len(spaces))
	for result, space := range spaces {
		newSpaces = append(newSpaces, Space{Result: result, Space: space})
	}
	sort.Slice(newSpaces, func(i, j int) bool {
		return float64(newSpaces[i].Space)/(1000/float64(stds[newSpaces[i].Result])) < float64(newSpaces[j].Space)/(1000/float64(stds[newSpaces[j].Result]))
	})

	target, price := make(map[int]struct{}), 0
	for _, newSpace := range newSpaces {
		price = price + stds[newSpace.Result]
		if price > 800 {
			break
		}

		target[newSpace.Result] = struct{}{}
	}

	return target
}
